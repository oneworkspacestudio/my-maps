package app

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/oneworkspacestudio/my-maps/internal/mapdata"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

const mapLimit = uint64(2 << 30)

type MapChoice struct {
	Coverage  string  `json:"coverage"`
	Depth     string  `json:"depth"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Radius    float64 `json:"radius"`
	Postcodes bool    `json:"postcodes"`
}
type MapInfo struct {
	ID            string     `json:"id"`
	Choice        MapChoice  `json:"choice"`
	Bounds        [4]float64 `json:"bounds"`
	MaxZoom       uint32     `json:"maxZoom"`
	Tiles         int        `json:"tiles"`
	Features      int        `json:"features"`
	Nearby        bool       `json:"nearby"`
	Source        string     `json:"source"`
	SourceVersion string     `json:"sourceVersion"`
	Installed     string     `json:"installed"`
	Bytes         uint64     `json:"bytes"`
}
type geoDownload struct {
	URL, ETag, Kind string
	Bytes           uint64
}
type mapPlan struct {
	ID            string     `json:"id"`
	Choice        MapChoice  `json:"choice"`
	Bounds        [4]float64 `json:"bounds"`
	MaxZoom       uint32     `json:"maxZoom"`
	TileCount     int        `json:"tileCount"`
	DownloadBytes uint64     `json:"downloadBytes"`
	RequiredBytes uint64     `json:"requiredBytes"`
	Source        string     `json:"source"`
	Expires       time.Time  `json:"expires"`
	archive       *mapdata.Archive
	groups        []mapdata.Group
	geo           []geoDownload
}
type mapProgress struct {
	Phase   string `json:"phase"`
	Message string `json:"message"`
	Done    uint64 `json:"done"`
	Total   uint64 `json:"total"`
	Busy    bool   `json:"busy"`
}
type mapManager struct {
	app       *App
	mu        sync.RWMutex
	plan      *mapPlan
	progress  mapProgress
	info      *MapInfo
	db        *sql.DB
	client    *http.Client
	root      string
	routing   string
	ctx       context.Context
	cancel    context.CancelFunc
	jobCancel context.CancelFunc
	wg        sync.WaitGroup
}

func newMapManager(a *App) (*mapManager, error) {
	ctx, cancel := context.WithCancel(context.Background())
	m := &mapManager{app: a, routing: a.routing, ctx: ctx, cancel: cancel, root: filepath.Join(a.H.Config.Data, "map-packs"), progress: mapProgress{Phase: "choose", Message: "Choose coverage and detail to install a local map."}}
	m.client = &http.Client{Timeout: 45 * time.Second, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) > 3 || req.URL.Scheme != "https" || req.URL.Host != via[0].URL.Host || req.URL.User != nil {
			return errors.New("map download redirected outside its source")
		}
		return nil
	}}
	if err := os.MkdirAll(m.root, 0700); err != nil {
		cancel()
		return nil, err
	}
	// Only installer-owned staging files are disposable after an interrupted run.
	if entries, err := os.ReadDir(m.root); err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && regexp.MustCompile(`^[a-f0-9]{32}\.partial(-journal)?$`).MatchString(entry.Name()) {
				os.Remove(filepath.Join(m.root, entry.Name()))
			}
		}
	}
	var id string
	err := a.H.DB.QueryRow("SELECT pack,routing FROM map_install_config WHERE singleton=1").Scan(&id, &m.routing)
	if err != nil && err != sql.ErrNoRows {
		cancel()
		return nil, err
	}
	if m.routing != "" && !localMapEndpoint(m.routing) {
		cancel()
		return nil, errors.New("saved routing URL must use a literal loopback address")
	}
	if id != "" && validPackID.MatchString(id) {
		db, info, e := openPack(filepath.Join(m.root, id+".sqlite"))
		if e == nil {
			m.db = db
			m.info = info
			m.progress = mapProgress{Phase: "ready", Message: "Installed map is served from this server."}
		} else {
			m.progress = mapProgress{Phase: "error", Message: "The installed map could not be opened. Choose a new download; saved places are unaffected."}
		}
	}
	return m, nil
}

var validPackID = regexp.MustCompile(`^[a-f0-9]{32}$`)

func openPack(path string) (*sql.DB, *MapInfo, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, nil, err
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, nil, err
	}
	db.SetMaxOpenConns(1)
	var body string
	if err = db.QueryRow("SELECT body FROM pack_info WHERE singleton=1").Scan(&body); err != nil {
		db.Close()
		return nil, nil, err
	}
	var info MapInfo
	if err = json.Unmarshal([]byte(body), &info); err != nil || !validPackID.MatchString(info.ID) || info.MaxZoom > 14 {
		db.Close()
		return nil, nil, errors.New("invalid installed map metadata")
	}
	if _, err = db.Exec("PRAGMA query_only=ON"); err != nil {
		db.Close()
		return nil, nil, err
	}
	if stat, err := os.Stat(path); err == nil {
		info.Bytes = uint64(stat.Size())
	}
	return db, &info, nil
}

func (m *mapManager) close() {
	m.cancel()
	m.mu.Lock()
	if m.jobCancel != nil {
		m.jobCancel()
	}
	m.mu.Unlock()
	m.wg.Wait()
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.db != nil {
		m.db.Close()
	}
}
func (m *mapManager) status(w http.ResponseWriter, r *http.Request) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	apphost.JSON(w, 200, map[string]any{"installed": m.info, "progress": m.progress, "canInstall": apphost.Current(r).User.Role == "owner", "routing": m.routing != ""})
}

func mapBounds(c MapChoice) ([4]float64, uint32, error) {
	var bounds [4]float64
	zoom := map[string]uint32{"overview": 6, "region": 10, "streets": 14}[c.Depth]
	if zoom == 0 {
		return bounds, 0, errors.New("choose a map detail level")
	}
	switch c.Coverage {
	case "world":
		bounds = [4]float64{-180, -85.051128, 180, 85.051128}
		if c.Depth != "overview" {
			return bounds, 0, errors.New("world coverage supports overview; select a region for greater detail")
		}
	case "uk":
		bounds = [4]float64{-8.7, 49.8, 2.1, 60.9}
		if c.Depth == "streets" {
			return bounds, 0, errors.New("choose a local area for street detail")
		}
	case "europe":
		bounds = [4]float64{-25, 34, 45, 72}
		if c.Depth == "streets" {
			return bounds, 0, errors.New("choose a local area for street detail")
		}
	case "area":
		if !coordinate(c.Lat, c.Lon) || math.Abs(c.Lat) > 80 || math.IsNaN(c.Radius) || math.IsInf(c.Radius, 0) || c.Radius < 1 || c.Radius > 100 {
			return bounds, 0, errors.New("use a valid centre and a radius of1–100km")
		}
		dy := c.Radius / 111.32
		dx := dy / math.Cos(c.Lat*math.Pi/180)
		bounds = [4]float64{c.Lon - dx, c.Lat - dy, c.Lon + dx, c.Lat + dy}
		if bounds[0] < -180 || bounds[2] > 180 || bounds[1] < -85 || bounds[3] > 85 {
			return bounds, 0, errors.New("choose a smaller area away from the map edge")
		}
	default:
		return bounds, 0, errors.New("choose world, UK, Europe or a local area")
	}
	if c.Postcodes && (c.Coverage == "world" || bounds[0] > 2.1 || bounds[2] < -8.7 || bounds[1] > 60.9 || bounds[3] < 49.8) {
		return bounds, 0, errors.New("UK postcodes require UK coverage or an overlapping area")
	}
	return bounds, zoom, nil
}

func (m *mapManager) makePlan(w http.ResponseWriter, r *http.Request) {
	var choice MapChoice
	if !apphost.Decode(w, r, &choice) {
		return
	}
	bounds, zoom, err := mapBounds(choice)
	if err != nil {
		apphost.Problem(w, 400, err.Error())
		return
	}
	m.mu.RLock()
	busy := m.progress.Busy
	m.mu.RUnlock()
	if busy {
		apphost.Problem(w, 409, "A map installation is already running")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Minute)
	defer cancel()
	var archive *mapdata.Archive
	for day := 1; day <= 8; day++ {
		date := time.Now().UTC().AddDate(0, 0, -day).Format("20060102")
		archive, err = mapdata.Open(ctx, m.client, "https://build.protomaps.com/"+date+".pmtiles")
		if err == nil {
			break
		}
		if ctx.Err() != nil {
			break
		}
	}
	if err != nil {
		apphost.Problem(w, 502, "The map download source is unavailable. Retry the plan later; no existing map was changed.")
		return
	}
	tiles, err := archive.Select(ctx, bounds, zoom)
	if err != nil {
		apphost.Problem(w, 400, err.Error())
		return
	}
	plan := &mapPlan{ID: apphost.ID(), Choice: choice, Bounds: bounds, MaxZoom: zoom, TileCount: len(tiles), Expires: time.Now().Add(20 * time.Minute), Source: archive.URL, archive: archive, groups: mapdata.Groups(tiles)}
	var storage uint64
	for _, tile := range tiles {
		storage += tile.Length*12 + 128
	}
	for _, group := range plan.groups {
		plan.DownloadBytes += group.Length
	}
	for _, kind := range []string{"cities", "postcodes"} {
		if kind == "postcodes" && !choice.Postcodes {
			continue
		}
		download, e := m.geoPlan(ctx, kind)
		if e != nil {
			apphost.Problem(w, 502, "The place-name download source is unavailable. Retry the plan later.")
			return
		}
		plan.geo = append(plan.geo, download)
		plan.DownloadBytes += download.Bytes
		storage += download.Bytes * 24
	}
	if plan.DownloadBytes > mapLimit || storage > 4*mapLimit {
		apphost.Problem(w, 400, "This download is too large. Choose a smaller region or less detail.")
		return
	}
	plan.RequiredBytes = storage*2 + (512 << 20)
	free, e := mapDiskFree(m.root)
	if e != nil || free < plan.RequiredBytes {
		apphost.Problem(w, 409, "There is not enough free space for this download, its index and safe installation. Choose less detail.")
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.progress.Busy {
		apphost.Problem(w, 409, "An installation started while this plan was prepared")
		return
	}
	m.plan = plan
	apphost.JSON(w, 200, plan)
}

func (m *mapManager) install(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Plan string `json:"plan"`
	}
	if !apphost.Decode(w, r, &in) {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	plan := m.plan
	if m.progress.Busy {
		apphost.Problem(w, 409, "An installation is already running")
		return
	}
	if plan == nil || in.Plan != plan.ID || time.Now().After(plan.Expires) {
		apphost.Problem(w, 409, "Create a fresh download plan first")
		return
	}
	free, err := mapDiskFree(m.root)
	if err != nil || free < plan.RequiredBytes {
		apphost.Problem(w, 409, "Free space changed. Create a new download plan.")
		return
	}
	ctx, cancel := context.WithCancel(m.ctx)
	m.jobCancel = cancel
	m.plan = nil
	m.progress = mapProgress{Phase: "downloading", Message: "Downloading the selected open map data...", Total: plan.DownloadBytes, Busy: true}
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		err := m.runInstall(ctx, plan)
		m.mu.Lock()
		defer m.mu.Unlock()
		m.jobCancel = nil
		if err != nil {
			log.Printf("Map installation failed: %v", err)
			message := "Map installation failed. Retry with a fresh plan; the previous map and saved places were kept."
			phase := "error"
			if ctx.Err() != nil {
				message = "Map download cancelled. The previous map and saved places were kept."
				phase = "cancelled"
			}
			m.progress = mapProgress{Phase: phase, Message: message}
		} else {
			m.progress = mapProgress{Phase: "ready", Message: "Your map is ready. Browsing and place searches now use this server's data.", Done: plan.DownloadBytes, Total: plan.DownloadBytes}
		}
	}()
	apphost.JSON(w, 202, map[string]bool{"started": true})
}

func (m *mapManager) cancelInstall(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.jobCancel != nil {
		m.jobCancel()
	}
	apphost.JSON(w, 200, map[string]bool{"cancelRequested": true})
}
func (m *mapManager) update(phase, message string, done, total uint64) {
	m.mu.Lock()
	m.progress = mapProgress{Phase: phase, Message: message, Done: done, Total: total, Busy: true}
	m.mu.Unlock()
}

const packSchema = `PRAGMA journal_mode=DELETE; PRAGMA synchronous=NORMAL;
 CREATE TABLE tiles(z INTEGER NOT NULL,x INTEGER NOT NULL,y INTEGER NOT NULL,body BLOB NOT NULL,PRIMARY KEY(z,x,y));
 CREATE TABLE features(id INTEGER PRIMARY KEY,feature_key TEXT UNIQUE NOT NULL,name TEXT NOT NULL,search_key TEXT NOT NULL,lat REAL NOT NULL,lon REAL NOT NULL,category TEXT NOT NULL,address TEXT NOT NULL);
 CREATE INDEX feature_key_search ON features(search_key); CREATE INDEX feature_nearby ON features(category,lat,lon);
 CREATE VIRTUAL TABLE features_fts USING fts5(name,address,content='features',content_rowid='id');
 CREATE TABLE pack_info(singleton INTEGER PRIMARY KEY,body TEXT NOT NULL);`

func (m *mapManager) runInstall(ctx context.Context, plan *mapPlan) error {
	key := make([]byte, 16)
	if _, err := rand.Read(key); err != nil {
		return err
	}
	id := hex.EncodeToString(key)
	partial := filepath.Join(m.root, id+".partial")
	final := filepath.Join(m.root, id+".sqlite")
	db, err := sql.Open("sqlite", partial)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(1)
	defer func() { db.Close(); os.Remove(partial); os.Remove(partial + "-journal") }()
	if _, err = db.ExecContext(ctx, packSchema); err != nil {
		return err
	}
	// Bound the staging database to the storage allowance shown in the plan.
	if plan.RequiredBytes < 512<<20 {
		return errors.New("invalid installation storage allowance")
	}
	maxPages := (plan.RequiredBytes - (512 << 20)) / 4096
	if _, err = db.ExecContext(ctx, fmt.Sprintf("PRAGMA max_page_count=%d", maxPages)); err != nil {
		return err
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	tileInsert, err := tx.PrepareContext(ctx, "INSERT INTO tiles VALUES(?,?,?,?)")
	if err != nil {
		return err
	}
	defer tileInsert.Close()
	featureInsert, err := tx.PrepareContext(ctx, "INSERT OR IGNORE INTO features(feature_key,name,search_key,lat,lon,category,address) VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	defer featureInsert.Close()
	info := MapInfo{ID: id, Choice: plan.Choice, Bounds: plan.Bounds, MaxZoom: plan.MaxZoom, Tiles: plan.TileCount, Source: plan.Source, SourceVersion: plan.archive.ETag, Installed: time.Now().UTC().Format(time.RFC3339)}
	var done uint64
	for _, group := range plan.groups {
		free, e := mapDiskFree(m.root)
		if e != nil || free < 512<<20 {
			return errors.New("free-space safety reserve reached")
		}
		body, e := plan.archive.Range(ctx, group.Offset, group.Length)
		if e != nil {
			return e
		}
		for _, tile := range group.Tiles {
			compressed := body[tile.Offset-group.Offset : tile.Offset-group.Offset+tile.Length]
			data, e := mapdata.Inflate(compressed, plan.archive.TileCompression, 8<<20)
			if e != nil {
				return e
			}
			if _, e = tileInsert.ExecContext(ctx, tile.Z, tile.X, tile.Y, data); e != nil {
				return e
			}
			if tile.Z == plan.MaxZoom {
				if e = indexMapTile(ctx, featureInsert, data, tile, plan.Bounds); e != nil {
					return e
				}
			}
		}
		done += group.Length
		m.update("downloading", "Downloading and indexing the selected map...", done, plan.DownloadBytes)
	}
	for _, download := range plan.geo {
		m.update("indexing", "Adding local town names and selected postcode coverage...", done, plan.DownloadBytes)
		if err = m.importGeo(ctx, featureInsert, download, plan.Bounds); err != nil {
			return err
		}
		done += download.Bytes
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	m.update("indexing", "Finishing your local place-search index...", done, plan.DownloadBytes)
	if _, err = db.ExecContext(ctx, "INSERT INTO features_fts(features_fts) VALUES('rebuild')"); err != nil {
		return err
	}
	if err = db.QueryRowContext(ctx, "SELECT count(*) FROM features").Scan(&info.Features); err != nil {
		return err
	}
	var near int
	if err = db.QueryRowContext(ctx, "SELECT count(*) FROM features WHERE category NOT IN ('place','road','postcode')").Scan(&near); err != nil {
		return err
	}
	info.Nearby = near > 0
	body, _ := json.Marshal(info)
	if _, err = db.ExecContext(ctx, "INSERT INTO pack_info VALUES(1,?)", string(body)); err != nil {
		return err
	}
	if err = db.Close(); err != nil {
		return err
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	file, err := os.Open(partial)
	if err != nil {
		return err
	}
	hash := sha256.New()
	_, err = io.Copy(hash, file)
	file.Close()
	if err != nil {
		return err
	}
	if err = os.Rename(partial, final); err != nil {
		return err
	}
	next, nextInfo, err := openPack(final)
	if err != nil {
		os.Remove(final)
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if ctx.Err() != nil {
		next.Close()
		os.Remove(final)
		return ctx.Err()
	}
	_, err = m.app.H.DB.Exec("INSERT INTO map_install_config(singleton,pack,sha256) VALUES(1,?,?) ON CONFLICT(singleton) DO UPDATE SET pack=excluded.pack,sha256=excluded.sha256", id, hex.EncodeToString(hash.Sum(nil)))
	if err != nil {
		next.Close()
		os.Remove(final)
		return err
	}
	old := m.info
	if m.db != nil {
		m.db.Close()
	}
	m.db = next
	m.info = nextInfo
	if old != nil && validPackID.MatchString(old.ID) {
		os.Remove(filepath.Join(m.root, old.ID+".sqlite"))
	}
	return nil
}

func (m *mapManager) tile(w http.ResponseWriter, r *http.Request) {
	var z, x, y uint32
	if _, err := fmt.Sscan(r.PathValue("z"), &z); err != nil || z > 20 {
		http.NotFound(w, r)
		return
	}
	if _, err := fmt.Sscan(r.PathValue("x"), &x); err != nil || x >= 1<<z {
		http.NotFound(w, r)
		return
	}
	if _, err := fmt.Sscan(r.PathValue("y"), &y); err != nil || y >= 1<<z {
		http.NotFound(w, r)
		return
	}
	m.mu.RLock()
	if m.info == nil || r.URL.Query().Get("v") != m.info.ID {
		m.mu.RUnlock()
		http.NotFound(w, r)
		return
	}
	var body []byte
	err := sql.ErrNoRows
	if m.db != nil {
		err = m.db.QueryRowContext(r.Context(), "SELECT body FROM tiles WHERE z=? AND x=? AND y=?", z, x, y).Scan(&body)
	}
	m.mu.RUnlock()
	if err == sql.ErrNoRows {
		w.Header().Set("Content-Type", "application/x-protobuf")
		w.Header().Set("Cache-Control", "private,max-age=60")
		w.WriteHeader(200)
		return
	}
	if err != nil {
		apphost.Problem(w, 503, "The local map could not be read")
		return
	}
	w.Header().Set("Content-Type", "application/x-protobuf")
	w.Header().Set("Cache-Control", "private,max-age=86400")
	w.Write(body)
}

func sortedFeatures(values []localFeature, lat, lon float64) {
	sort.SliceStable(values, func(i, j int) bool {
		return geoDistance(lat, lon, values[i].Lat, values[i].Lon) < geoDistance(lat, lon, values[j].Lat, values[j].Lon)
	})
}

func (m *mapManager) setRouting(w http.ResponseWriter, r *http.Request) {
	var in struct {
		URL string `json:"url"`
	}
	if !apphost.Decode(w, r, &in) {
		return
	}
	if in.URL != "" && !localMapEndpoint(in.URL) {
		apphost.Problem(w, 400, "Use a literal loopback URL for your local OSRM service, such as http://127.0.0.1:5000")
		return
	}
	if in.URL != "" {
		ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
		defer cancel()
		target := strings.TrimRight(in.URL, "/") + "/nearest/v1/driving/0,51?number=1"
		req, err := http.NewRequestWithContext(ctx, "GET", target, nil)
		if err != nil {
			apphost.Problem(w, 400, "Invalid local routing URL")
			return
		}
		response, err := m.app.client.Do(req)
		if err != nil {
			apphost.Problem(w, 502, "The local OSRM service could not be reached. Start it, then reconnect.")
			return
		}
		defer response.Body.Close()
		var result struct {
			Code string `json:"code"`
		}
		err = json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&result)
		if err != nil || (result.Code != "Ok" && result.Code != "NoSegment") {
			apphost.Problem(w, 502, "The local service did not return an OSRM driving response")
			return
		}
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, err := m.app.H.DB.Exec("INSERT INTO map_install_config(singleton,pack,sha256,routing) VALUES(1,'','',?) ON CONFLICT(singleton) DO UPDATE SET routing=excluded.routing", in.URL); err != nil {
		apphost.Problem(w, 500, "Routing setting could not be saved")
		return
	}
	m.routing = in.URL
	apphost.JSON(w, 200, map[string]bool{"enabled": in.URL != ""})
}
