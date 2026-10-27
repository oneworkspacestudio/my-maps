package mapdata

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHilbertTileIDs(t *testing.T) {
	for _, v := range [][4]uint32{{0, 0, 0, 0}, {1, 0, 0, 1}, {1, 0, 1, 2}, {1, 1, 1, 3}, {1, 1, 0, 4}, {2, 0, 0, 5}, {2, 1, 0, 6}, {2, 1, 1, 7}, {2, 0, 1, 8}, {2, 0, 2, 9}} {
		if got := TileID(v[0], v[1], v[2]); got != uint64(v[3]) {
			t.Fatalf("%v -> %d", v, got)
		}
	}
	seen := map[uint64]bool{}
	for x := uint32(0); x < 64; x++ {
		for y := uint32(0); y < 64; y++ {
			id := TileID(6, x, y)
			if seen[id] {
				t.Fatal("duplicate Hilbert ID")
			}
			seen[id] = true
		}
	}
}
func TestArchiveRangesDirectoryAndExtract(t *testing.T) {
	dir := []byte{}
	for _, v := range []uint64{1, 0, 1, 3, 1} {
		dir = binary.AppendUvarint(dir, v)
	}
	raw := make([]byte, 127)
	copy(raw, "PMTiles")
	raw[7] = 3
	raw[97] = 1
	raw[98] = 1
	raw[99] = 1
	raw[101] = 14
	binary.LittleEndian.PutUint64(raw[8:], 127)
	binary.LittleEndian.PutUint64(raw[16:], uint64(len(dir)))
	binary.LittleEndian.PutUint64(raw[56:], 127+uint64(len(dir)))
	binary.LittleEndian.PutUint64(raw[64:], 3)
	raw = append(raw, dir...)
	raw = append(raw, 1, 2, 3)
	tag := `"one"`
	wrong := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var start, end int
		fmt.Sscanf(r.Header.Get("Range"), "bytes=%d-%d", &start, &end)
		if start < 0 || end >= len(raw) {
			t.Error("unbounded range")
			w.WriteHeader(416)
			return
		}
		w.Header().Set("ETag", tag)
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, len(raw)))
		if wrong {
			w.Header().Set("Content-Range", "bytes 0-0/1")
		}
		w.WriteHeader(206)
		w.Write(raw[start : end+1])
	}))
	defer server.Close()
	a, err := Open(context.Background(), server.Client(), server.URL)
	if err != nil {
		t.Fatal(err)
	}
	tiles, err := a.Select(context.Background(), [4]float64{-1, -1, 1, 1}, 0)
	if err != nil || len(tiles) != 1 {
		t.Fatal(tiles, err)
	}
	groups := Groups(tiles)
	if len(groups) != 1 || groups[0].Length != 3 {
		t.Fatal(groups)
	}
	body, err := a.Range(context.Background(), tiles[0].Offset, tiles[0].Length)
	if err != nil || !bytes.Equal(body, []byte{1, 2, 3}) {
		t.Fatal(body, err)
	}
	tag = `"changed"`
	if _, err = a.Range(context.Background(), 0, 127); err == nil {
		t.Fatal("accepted source replacement")
	}
	tag = `"one"`
	wrong = true
	if _, err = a.Range(context.Background(), 0, 127); err == nil {
		t.Fatal("accepted wrong range")
	}
	if _, err = a.Select(context.Background(), [4]float64{math.NaN(), 0, 1, 1}, 6); err == nil {
		t.Fatal("NaN bounds accepted")
	}
}
func TestMalformedDataAndBoundedExpansion(t *testing.T) {
	for _, body := range [][]byte{{}, {255, 255}, {1, 0, 1, 0, 0}, {0, 1}} {
		if _, err := DecodeDirectory(body); err == nil {
			t.Fatalf("accepted invalid directory %v", body)
		}
	}
	var out bytes.Buffer
	g := gzip.NewWriter(&out)
	io.WriteString(g, strings.Repeat("a", 4096))
	g.Close()
	if _, err := Inflate(out.Bytes(), 2, 1024); err == nil {
		t.Fatal("unbounded expansion")
	}
	if _, err := Inflate([]byte("bad gzip"), 2, 1024); err == nil {
		t.Fatal("bad gzip")
	}
}
func TestGroupsDeduplicatedRunAndGap(t *testing.T) {
	g := Groups([]Tile{{Offset: 100, Length: 20}, {Offset: 100, Length: 20}, {Offset: 150, Length: 10}, {Offset: 100000, Length: 2}})
	if len(g) != 2 || g[0].Length != 60 || len(g[0].Tiles) != 3 {
		t.Fatal(g)
	}
}
