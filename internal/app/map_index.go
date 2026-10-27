package app

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/oneworkspacestudio/my-maps/internal/mapdata"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/mvt"
	"github.com/paulmach/orb/maptile"
)

type localFeature struct {
	ID                      int64
	Name, Category, Address string
	Lat, Lon                float64
}

func searchKey(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return unicode.ToLower(r)
		}
		return -1
	}, s)
}
func inside(bounds [4]float64, lat, lon float64) bool {
	return coordinate(lat, lon) && lon >= bounds[0] && lon <= bounds[2] && lat >= bounds[1] && lat <= bounds[3]
}

func addMapFeature(ctx context.Context, insert *sql.Stmt, name, address, category string, lat, lon float64) error {
	name = strings.TrimSpace(name)
	if name == "" || len(name) > 200 || len(address) > 1000 || !coordinate(lat, lon) {
		return nil
	}
	key := fmt.Sprintf("%s|%s|%.4f|%.4f", name, category, lat, lon)
	_, err := insert.ExecContext(ctx, key, name, searchKey(name), lat, lon, category, address)
	return err
}

func indexMapTile(ctx context.Context, insert *sql.Stmt, data []byte, tile mapdata.Tile, bounds [4]float64) error {
	layers, err := mvt.Unmarshal(data)
	if err != nil {
		return errors.New("invalid vector map tile")
	}
	for _, layer := range layers {
		if layer.Name != "pois" && layer.Name != "places" && layer.Name != "roads" {
			continue
		}
		layer.ProjectToWGS84(maptile.New(tile.X, tile.Y, maptile.Zoom(tile.Z)))
		for i, feature := range layer.Features {
			if i%128 == 0 && ctx.Err() != nil {
				return ctx.Err()
			}
			name, _ := feature.Properties["name"].(string)
			if name == "" {
				name, _ = feature.Properties["name:en"].(string)
			}
			if name == "" || feature.Geometry == nil {
				continue
			}
			point := feature.Geometry.Bound().Center()
			if _, ok := feature.Geometry.(orb.Point); ok {
				point = feature.Geometry.(orb.Point)
			}
			if !inside(bounds, point.Lat(), point.Lon()) {
				continue
			}
			category := "place"
			if layer.Name == "roads" {
				category = "road"
			} else if layer.Name == "pois" {
				kind := fmt.Sprint(feature.Properties["pmap:kind"], " ", feature.Properties["kind"], " ", feature.Properties["kind_detail"])
				kind = strings.ToLower(kind)
				switch {
				case hasAny(kind, "restaurant", "cafe", "fast_food", "pub", "bar", "food_and_drink"):
					category = "restaurants"
				case hasAny(kind, "supermarket", "shop", "shopping", "convenience", "retail"):
					category = "shopping"
				case hasAny(kind, "hotel", "hostel", "motel", "lodging", "guest_house"):
					category = "hotels"
				case hasAny(kind, "museum", "attraction", "gallery", "viewpoint"):
					category = "attractions"
				case hasAny(kind, "fuel", "gas_station"):
					category = "fuel"
				case hasAny(kind, "parking"):
					category = "parking"
				case hasAny(kind, "bicycle", "cycling"):
					category = "cycling"
				case hasAny(kind, "station", "public_transport", "bus_stop"):
					category = "transit"
				}
			}
			if err := addMapFeature(ctx, insert, name, "", category, point.Lat(), point.Lon()); err != nil {
				return err
			}
		}
	}
	return nil
}
func hasAny(value string, words ...string) bool {
	for _, word := range words {
		if strings.Contains(value, word) {
			return true
		}
	}
	return false
}

func (m *mapManager) geoPlan(ctx context.Context, kind string) (geoDownload, error) {
	target := "https://download.geonames.org/export/dump/cities15000.zip"
	if kind == "postcodes" {
		target = "https://download.geonames.org/export/zip/GB_full.csv.zip"
	}
	req, err := http.NewRequestWithContext(ctx, "HEAD", target, nil)
	if err != nil {
		return geoDownload{}, err
	}
	resp, err := m.client.Do(req)
	if err != nil {
		return geoDownload{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 || resp.ContentLength <= 0 || resp.ContentLength > 40<<20 || resp.Header.Get("ETag") == "" {
		return geoDownload{}, errors.New("invalid place-name download size/version")
	}
	return geoDownload{URL: target, ETag: resp.Header.Get("ETag"), Kind: kind, Bytes: uint64(resp.ContentLength)}, nil
}

func (m *mapManager) importGeo(ctx context.Context, insert *sql.Stmt, download geoDownload, bounds [4]float64) error {
	req, err := http.NewRequestWithContext(ctx, "GET", download.URL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("If-Match", download.ETag)
	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 || resp.Header.Get("ETag") != download.ETag {
		return errors.New("place-name source changed; make a new plan")
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, int64(download.Bytes)+1))
	if err != nil || uint64(len(body)) != download.Bytes {
		return errors.New("incomplete place-name download")
	}
	archive, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return err
	}
	total := uint64(0)
	for _, file := range archive.File {
		if !strings.HasSuffix(file.Name, ".txt") && !strings.HasSuffix(file.Name, ".csv") {
			continue
		}
		if strings.Contains(strings.ToLower(file.Name), "readme") {
			continue
		}
		if file.UncompressedSize64 > (512<<20)-total {
			return errors.New("place-name archive expansion limit exceeded")
		}
		total += file.UncompressedSize64
		stream, err := file.Open()
		if err != nil {
			return err
		}
		err = readGeo(ctx, insert, io.LimitReader(stream, int64(file.UncompressedSize64)+1), download.Kind, bounds)
		stream.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func readGeo(ctx context.Context, insert *sql.Stmt, reader io.Reader, kind string, bounds [4]float64) error {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 8192), 1<<20)
	line := 0
	for scanner.Scan() {
		line++
		if line%256 == 0 && ctx.Err() != nil {
			return ctx.Err()
		}
		fields := strings.Split(scanner.Text(), "\t")
		var name, address, category, latText, lonText string
		if kind == "cities" {
			if len(fields) < 19 {
				continue
			}
			name = fields[1]
			address = fields[8]
			category = "place"
			latText = fields[4]
			lonText = fields[5]
		} else {
			if len(fields) < 12 {
				continue
			}
			name = fields[1]
			address = strings.Join([]string{fields[2], fields[3], fields[5], fields[0]}, ", ")
			category = "postcode"
			latText = fields[9]
			lonText = fields[10]
		}
		lat, e1 := strconv.ParseFloat(latText, 64)
		lon, e2 := strconv.ParseFloat(lonText, 64)
		if e1 != nil || e2 != nil || !inside(bounds, lat, lon) {
			continue
		}
		if err := addMapFeature(ctx, insert, name, address, category, lat, lon); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func (m *mapManager) localSearch(w http.ResponseWriter, r *http.Request, q string) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.db == nil {
		apphost.Problem(w, 503, "Choose a map download first")
		return
	}
	key := searchKey(q)
	if key == "" {
		apphost.Problem(w, 400, "Enter a place, road or postcode")
		return
	}
	values, err := queryLocal(r.Context(), m.db, "SELECT id,name,category,address,lat,lon FROM features WHERE search_key=? LIMIT 20", key)
	if err == nil && len(values) == 0 {
		words := strings.FieldsFunc(q, func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) })
		terms := []string{}
		for _, word := range words {
			terms = append(terms, "\""+strings.ReplaceAll(word, "\"", "\"\"")+"\"*")
		}
		if len(terms) > 8 {
			terms = terms[:8]
		}
		values, err = queryLocal(r.Context(), m.db, "SELECT f.id,f.name,f.category,f.address,f.lat,f.lon FROM features_fts JOIN features f ON f.id=features_fts.rowid WHERE features_fts MATCH ? LIMIT 60", strings.Join(terms, " AND "))
	}
	if err != nil {
		apphost.Problem(w, 503, "The local place index could not be read")
		return
	}
	if lat, lon, ok := requestCoordinates(r); ok {
		sortedFeatures(values, lat, lon)
	}
	if len(values) > 12 {
		values = values[:12]
	}
	writeLocalFeatures(w, values)
}

func queryLocal(ctx context.Context, db *sql.DB, query string, args ...any) ([]localFeature, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	values := []localFeature{}
	for rows.Next() {
		var f localFeature
		if err = rows.Scan(&f.ID, &f.Name, &f.Category, &f.Address, &f.Lat, &f.Lon); err != nil {
			return nil, err
		}
		values = append(values, f)
	}
	return values, rows.Err()
}
func writeLocalFeatures(w http.ResponseWriter, values []localFeature) {
	features := []any{}
	for _, value := range values {
		features = append(features, map[string]any{"type": "Feature", "geometry": map[string]any{"type": "Point", "coordinates": []float64{value.Lon, value.Lat}}, "properties": map[string]any{"name": value.Name, "city": value.Address, "osm_value": value.Category}})
	}
	apphost.JSON(w, 200, map[string]any{"type": "FeatureCollection", "features": features, "source": "local"})
}

func (m *mapManager) localReverse(w http.ResponseWriter, r *http.Request, lat, lon float64) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.db == nil {
		writeLocalFeatures(w, nil)
		return
	}
	values, err := queryLocal(r.Context(), m.db, "SELECT id,name,category,address,lat,lon FROM features WHERE lat BETWEEN ? AND ? AND lon BETWEEN ? AND ? LIMIT 100", lat-.01, lat+.01, lon-.02, lon+.02)
	if err != nil {
		apphost.Problem(w, 503, "The local place index could not be read")
		return
	}
	sortedFeatures(values, lat, lon)
	if len(values) > 1 {
		values = values[:1]
	}
	writeLocalFeatures(w, values)
}

func (m *mapManager) localNearby(w http.ResponseWriter, r *http.Request, lat, lon float64, category string) {
	if _, ok := nearbyTags[category]; !ok {
		apphost.Problem(w, 400, "Choose a nearby category")
		return
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.db == nil {
		apphost.Problem(w, 503, "Choose a map download first")
		return
	}
	dx := .025 / math.Max(.1, math.Cos(lat*math.Pi/180))
	values, err := queryLocal(r.Context(), m.db, "SELECT id,name,category,address,lat,lon FROM features WHERE category=? AND lat BETWEEN ? AND ? AND lon BETWEEN ? AND ? ORDER BY (lat-?)*(lat-?)+(lon-?)*(lon-?)*? LIMIT 500", category, lat-.025, lat+.025, lon-dx, lon+dx, lat, lat, lon, lon, math.Pow(math.Cos(lat*math.Pi/180), 2))
	if err != nil {
		apphost.Problem(w, 503, "The local place index could not be read")
		return
	}
	sortedFeatures(values, lat, lon)
	elements := []any{}
	for _, value := range values {
		if geoDistance(lat, lon, value.Lat, value.Lon) > 2500 {
			continue
		}
		tags := map[string]string{"name": value.Name, "addr:city": value.Address}
		tags["category"] = category
		elements = append(elements, map[string]any{"type": "local", "id": value.ID, "lat": value.Lat, "lon": value.Lon, "tags": tags})
		if len(elements) >= 60 {
			break
		}
	}
	apphost.JSON(w, 200, map[string]any{"elements": elements, "source": "local"})
}
