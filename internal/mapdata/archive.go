// Package mapdata reads bounded PMTiles v3 ranges during an explicit map install.
// Browsing never accesses the source archive. Format: https://pmtiles.org/spec/v3/
package mapdata

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type Entry struct{ ID, Offset, Length, Run uint64 }
type Tile struct {
	Z, X, Y        uint32
	Offset, Length uint64
}
type Archive struct {
	Client                                                                 *http.Client
	URL, ETag                                                              string
	Size                                                                   uint64
	RootOffset, RootLength, LeafOffset, LeafLength, DataOffset, DataLength uint64
	Compression, TileCompression, MaxZoom                                  byte
	directories                                                            map[uint64][]Entry
	directoryEntries                                                       int
}

func Open(ctx context.Context, client *http.Client, source string) (*Archive, error) {
	a := &Archive{Client: client, URL: source, directories: map[uint64][]Entry{}}
	body, err := a.Range(ctx, 0, 127)
	if err != nil {
		return nil, err
	}
	if len(body) != 127 || string(body[:7]) != "PMTiles" || body[7] != 3 || body[99] != 1 {
		return nil, errors.New("unsupported map archive")
	}
	a.RootOffset = binary.LittleEndian.Uint64(body[8:16])
	a.RootLength = binary.LittleEndian.Uint64(body[16:24])
	a.LeafOffset = binary.LittleEndian.Uint64(body[40:48])
	a.LeafLength = binary.LittleEndian.Uint64(body[48:56])
	a.DataOffset = binary.LittleEndian.Uint64(body[56:64])
	a.DataLength = binary.LittleEndian.Uint64(body[64:72])
	a.Compression = body[97]
	a.TileCompression = body[98]
	a.MaxZoom = body[101]
	if (a.Compression != 1 && a.Compression != 2) || (a.TileCompression != 1 && a.TileCompression != 2) || a.MaxZoom > 16 || !within(a.RootOffset, a.RootLength, a.Size) || !within(a.LeafOffset, a.LeafLength, a.Size) || !within(a.DataOffset, a.DataLength, a.Size) {
		return nil, errors.New("unsupported or invalid map archive layout")
	}
	return a, nil
}

func within(offset, length, total uint64) bool { return offset <= total && length <= total-offset }

func (a *Archive) Range(ctx context.Context, offset, length uint64) ([]byte, error) {
	if length == 0 || length > 16<<20 || offset > math.MaxInt64-length {
		return nil, errors.New("map range exceeds limit")
	}
	req, err := http.NewRequestWithContext(ctx, "GET", a.URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", offset, offset+length-1))
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("User-Agent", "OneWorkspace-MyMaps map installer")
	if a.ETag != "" {
		req.Header.Set("If-Match", a.ETag)
	}
	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 206 {
		return nil, fmt.Errorf("map download requires byte ranges (status %d)", resp.StatusCode)
	}
	prefix := fmt.Sprintf("bytes %d-%d/", offset, offset+length-1)
	contentRange := resp.Header.Get("Content-Range")
	if !strings.HasPrefix(contentRange, prefix) {
		return nil, errors.New("map download returned the wrong byte range")
	}
	size, err := strconv.ParseUint(strings.TrimPrefix(contentRange, prefix), 10, 64)
	if err != nil || !within(offset, length, size) {
		return nil, errors.New("invalid map archive size")
	}
	tag := resp.Header.Get("ETag")
	if tag == "" || strings.HasPrefix(tag, "W/") {
		return nil, errors.New("map download needs a stable source version")
	}
	if a.ETag != "" && (a.ETag != tag || a.Size != size) {
		return nil, errors.New("map source changed; create a new download plan")
	}
	a.ETag = tag
	a.Size = size
	body, err := io.ReadAll(io.LimitReader(resp.Body, int64(length)+1))
	if err != nil || uint64(len(body)) != length {
		return nil, errors.New("incomplete map download")
	}
	return body, nil
}

func Inflate(body []byte, compression byte, limit int64) ([]byte, error) {
	var reader io.Reader = bytes.NewReader(body)
	if compression == 2 {
		g, err := gzip.NewReader(reader)
		if err != nil {
			return nil, err
		}
		defer g.Close()
		reader = g
	} else if compression != 1 {
		return nil, errors.New("unsupported map compression")
	}
	result, err := io.ReadAll(io.LimitReader(reader, limit+1))
	if err != nil || int64(len(result)) > limit {
		return nil, errors.New("map payload exceeds limit")
	}
	return result, nil
}

func DecodeDirectory(body []byte) ([]Entry, error) {
	reader := bytes.NewReader(body)
	count, err := binary.ReadUvarint(reader)
	if err != nil || count > 1000000 {
		return nil, errors.New("invalid map directory count")
	}
	entries := make([]Entry, int(count))
	var last uint64
	for i := range entries {
		delta, e := binary.ReadUvarint(reader)
		if e != nil || (i > 0 && delta == 0) || delta > math.MaxUint64-last {
			return nil, errors.New("invalid map tile IDs")
		}
		last += delta
		entries[i].ID = last
	}
	for i := range entries {
		n, e := binary.ReadUvarint(reader)
		if e != nil || n > 1<<32 {
			return nil, errors.New("invalid map tile run")
		}
		entries[i].Run = n
	}
	for i := range entries {
		n, e := binary.ReadUvarint(reader)
		if e != nil || n == 0 || n > 16<<20 {
			return nil, errors.New("invalid map entry size")
		}
		entries[i].Length = n
	}
	for i := range entries {
		n, e := binary.ReadUvarint(reader)
		if e != nil {
			return nil, e
		}
		if n == 0 {
			if i == 0 {
				return nil, errors.New("invalid initial map offset")
			}
			if entries[i-1].Offset > math.MaxUint64-entries[i-1].Length {
				return nil, errors.New("map offset overflow")
			}
			entries[i].Offset = entries[i-1].Offset + entries[i-1].Length
		} else {
			entries[i].Offset = n - 1
		}
	}
	if reader.Len() != 0 {
		return nil, errors.New("trailing map directory bytes")
	}
	return entries, nil
}

func (a *Archive) directory(ctx context.Context, offset, length uint64) ([]Entry, error) {
	if d, ok := a.directories[offset]; ok {
		return d, nil
	}
	body, err := a.Range(ctx, offset, length)
	if err != nil {
		return nil, err
	}
	body, err = Inflate(body, a.Compression, 16<<20)
	if err != nil {
		return nil, err
	}
	entries, err := DecodeDirectory(body)
	if err == nil {
		if a.directoryEntries+len(entries) > 1000000 {
			return nil, errors.New("map directory memory limit exceeded")
		}
		a.directoryEntries += len(entries)
		a.directories[offset] = entries
	}
	return entries, err
}

func (a *Archive) find(ctx context.Context, id uint64) (Entry, bool, error) {
	offset, length := a.RootOffset, a.RootLength
	for depth := 0; depth < 4; depth++ {
		entries, err := a.directory(ctx, offset, length)
		if err != nil {
			return Entry{}, false, err
		}
		pos := sort.Search(len(entries), func(i int) bool { return entries[i].ID > id }) - 1
		if pos < 0 {
			return Entry{}, false, nil
		}
		entry := entries[pos]
		if entry.Run > 0 {
			if id-entry.ID >= entry.Run {
				return Entry{}, false, nil
			}
			if !within(entry.Offset, entry.Length, a.DataLength) || entry.Length > 4<<20 {
				return Entry{}, false, errors.New("invalid map tile range")
			}
			entry.Offset += a.DataOffset
			return entry, true, nil
		}
		if !within(entry.Offset, entry.Length, a.LeafLength) {
			return Entry{}, false, errors.New("invalid map leaf directory")
		}
		offset = a.LeafOffset + entry.Offset
		length = entry.Length
	}
	return Entry{}, false, errors.New("map directory nesting limit exceeded")
}

func TileID(z, x, y uint32) uint64 {
	n := uint32(1) << z
	result := (uint64(1)<<(2*z) - 1) / 3
	for s := n / 2; s > 0; s /= 2 {
		rx, ry := uint32(0), uint32(0)
		if x&s != 0 {
			rx = 1
		}
		if y&s != 0 {
			ry = 1
		}
		result += uint64(s) * uint64(s) * uint64((3*rx)^ry)
		if ry == 0 {
			if rx == 1 {
				x = n - 1 - x
				y = n - 1 - y
			}
			x, y = y, x
		}
	}
	return result
}

// Select enumerates only the requested bounded coverage, never the planet's directory.
func (a *Archive) Select(ctx context.Context, bounds [4]float64, maxZoom uint32) ([]Tile, error) {
	for _, v := range bounds {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return nil, errors.New("invalid map bounds")
		}
	}
	if bounds[0] < -180 || bounds[2] > 180 || bounds[1] < -85.052 || bounds[3] > 85.052 {
		return nil, errors.New("invalid map bounds")
	}
	if maxZoom > 14 || maxZoom > uint32(a.MaxZoom) || bounds[0] >= bounds[2] || bounds[1] >= bounds[3] {
		return nil, errors.New("invalid map coverage or depth")
	}
	wanted := []Tile{}
	for z := uint32(0); z <= maxZoom; z++ {
		n := float64(uint32(1) << z)
		x := func(lon float64) uint32 { return uint32(math.Max(0, math.Min(n-1, math.Floor((lon+180)/360*n)))) }
		y := func(lat float64) uint32 {
			lat = math.Max(-85.051128, math.Min(85.051128, lat)) * math.Pi / 180
			return uint32(math.Max(0, math.Min(n-1, math.Floor((1-math.Asinh(math.Tan(lat))/math.Pi)/2*n))))
		}
		west, east, north, south := x(bounds[0]), x(bounds[2]), y(bounds[3]), y(bounds[1])
		// A one-tile margin avoids blank edges around the selected area at each zoom.
		if west > 0 {
			west--
		}
		if north > 0 {
			north--
		}
		if east < uint32(n)-1 {
			east++
		}
		if south < uint32(n)-1 {
			south++
		}
		if uint64(east-west+1)*uint64(south-north+1)+uint64(len(wanted)) > 200000 {
			return nil, errors.New("this map is too large; choose less detail or a smaller area")
		}
		for tx := west; tx <= east; tx++ {
			for ty := north; ty <= south; ty++ {
				wanted = append(wanted, Tile{Z: z, X: tx, Y: ty})
			}
		}
	}
	sort.Slice(wanted, func(i, j int) bool {
		return TileID(wanted[i].Z, wanted[i].X, wanted[i].Y) < TileID(wanted[j].Z, wanted[j].X, wanted[j].Y)
	})
	result := make([]Tile, 0, len(wanted))
	for i, tile := range wanted {
		if i%128 == 0 && ctx.Err() != nil {
			return nil, ctx.Err()
		}
		entry, found, err := a.find(ctx, TileID(tile.Z, tile.X, tile.Y))
		if err != nil {
			return nil, err
		}
		if found {
			tile.Offset = entry.Offset
			tile.Length = entry.Length
			result = append(result, tile)
		}
	}
	if len(result) == 0 {
		return nil, errors.New("no map data in this selection")
	}
	return result, nil
}

type Group struct {
	Offset, Length uint64
	Tiles          []Tile
}

func Groups(tiles []Tile) []Group {
	ordered := append([]Tile(nil), tiles...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].Offset < ordered[j].Offset })
	result := []Group{}
	for _, tile := range ordered {
		if len(result) > 0 {
			last := &result[len(result)-1]
			end := last.Offset + last.Length
			if tile.Offset <= end+16384 && tile.Offset+tile.Length-last.Offset <= 8<<20 {
				if tile.Offset+tile.Length > end {
					last.Length = tile.Offset + tile.Length - last.Offset
				}
				last.Tiles = append(last.Tiles, tile)
				continue
			}
		}
		result = append(result, Group{Offset: tile.Offset, Length: tile.Length, Tiles: []Tile{tile}})
	}
	return result
}
