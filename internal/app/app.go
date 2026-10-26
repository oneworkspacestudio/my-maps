package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type App struct {
	H                      *apphost.Host
	client                 *http.Client
	mu                     sync.Mutex
	last                   time.Time
	photon, routing, tiles string
}
type Place struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Address string  `json:"address"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	List    string  `json:"list"`
	Notes   string  `json:"notes"`
}

func env(k, d string) string {
	if s := strings.TrimSpace(os.Getenv(k)); s != "" {
		return strings.TrimRight(s, "/")
	}
	return d
}
func New(h *apphost.Host) (*App, error) {
	a := &App{H: h, client: &http.Client{Timeout: 20 * time.Second}, photon: env("MYMAPS_GEOCODER", "https://photon.komoot.io"), routing: env("MYMAPS_ROUTER", "https://routing.openstreetmap.de"), tiles: env("MYMAPS_TILES", "https://tile.openstreetmap.org/{z}/{x}/{y}.png")}
	for _, endpoint := range []string{a.photon, a.routing, a.tiles, env("MYMAPS_OVERPASS", "https://overpass-api.de")} {
		u, e := url.Parse(endpoint)
		if e != nil || u.Scheme != "https" || u.Host == "" || u.User != nil {
			return nil, errors.New("map providers must be HTTPS URLs without credentials")
		}
	}
	tile, _ := url.Parse(a.tiles)
	h.Config.ImageOrigins = []string{"https://" + tile.Host}
	if e := h.Migrate(1, `CREATE TABLE places(id TEXT PRIMARY KEY,owner TEXT NOT NULL,name TEXT NOT NULL,address TEXT NOT NULL,lat REAL NOT NULL,lon REAL NOT NULL,list TEXT NOT NULL,notes TEXT NOT NULL);
 CREATE INDEX places_owner ON places(owner,list);
 CREATE TABLE provider_cache(key TEXT PRIMARY KEY,body BLOB NOT NULL,expires INTEGER NOT NULL);`); e != nil {
		return nil, e
	}
	h.Mux.HandleFunc("GET /api/config", h.Require(func(w http.ResponseWriter, r *http.Request) {
		apphost.JSON(w, 200, map[string]any{"tiles": a.tiles, "customTiles": os.Getenv("MYMAPS_TILES") != "", "attribution": env("MYMAPS_TILE_ATTRIBUTION", "© OpenStreetMap contributors"), "geocoder": a.photon, "router": a.routing})
	}))
	h.Mux.HandleFunc("GET /api/search", h.Require(a.search))
	h.Mux.HandleFunc("GET /api/reverse", h.Require(a.reverse))
	h.Mux.HandleFunc("GET /api/nearby", h.Require(a.nearby))
	h.Mux.HandleFunc("GET /api/map-resource", h.Require(a.mapResource))
	h.Mux.HandleFunc("GET /api/weather", h.Require(a.weather))
	h.Mux.HandleFunc("POST /api/route", h.Require(a.route))
	h.Mux.HandleFunc("GET /api/places", h.Require(a.places))
	h.Mux.HandleFunc("POST /api/places", h.Require(a.save))
	h.Mux.HandleFunc("DELETE /api/places/{id}", h.Require(a.remove))
	return a, nil
}
func (a *App) Run(ctx context.Context) error {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			_, _ = a.H.DB.Exec("DELETE FROM provider_cache WHERE expires<?", time.Now().Unix())
		}
	}
}
func coordinate(lat, lon float64) bool {
	return !math.IsNaN(lat) && !math.IsNaN(lon) && !math.IsInf(lat, 0) && !math.IsInf(lon, 0) && lat >= -90 && lat <= 90 && lon >= -180 && lon <= 180
}
func (a *App) fetch(w http.ResponseWriter, r *http.Request, target string) {
	a.fetchTTL(w, r, target, 24*time.Hour)
}
func (a *App) fetchTTL(w http.ResponseWriter, r *http.Request, target string, ttl time.Duration) {
	key := apphost.Hash(target)
	var body []byte
	if a.H.DB.QueryRow("SELECT body FROM provider_cache WHERE key=? AND expires>?", key, time.Now().Unix()).Scan(&body) == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
		return
	}
	a.mu.Lock()
	if time.Since(a.last) < time.Second {
		a.mu.Unlock()
		w.Header().Set("Retry-After", "1")
		apphost.Problem(w, 429, "Please wait a second before the next search or route")
		return
	}
	a.last = time.Now()
	a.mu.Unlock()
	req, e := http.NewRequestWithContext(r.Context(), "GET", target, nil)
	if e != nil {
		apphost.Problem(w, 400, "Invalid map request")
		return
	}
	req.Header.Set("User-Agent", "OneWorkspace-MyMaps/0.1 (+https://github.com/oneworkspacestudio/my-maps)")
	req.Header.Set("Accept", "application/json")
	resp, e := a.client.Do(req)
	if e != nil {
		apphost.Problem(w, 502, "The map provider is unavailable. Saved places are still available; retry shortly.")
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		apphost.Problem(w, 502, "The community map provider is busy. Please retry later.")
		return
	}
	body, e = io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if e != nil || !json.Valid(body) {
		apphost.Problem(w, 502, "The map provider returned an incomplete response")
		return
	}
	_, _ = a.H.DB.Exec("DELETE FROM provider_cache WHERE expires<?", time.Now().Unix())
	var count int
	a.H.DB.QueryRow("SELECT count(*) FROM provider_cache").Scan(&count)
	if count < 5000 {
		_, _ = a.H.DB.Exec("INSERT OR REPLACE INTO provider_cache VALUES(?,?,?)", key, body, time.Now().Add(ttl).Unix())
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
func (a *App) search(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if len(q) < 2 || len(q) > 200 {
		apphost.Problem(w, 400, "Enter a place or address (2–200 characters)")
		return
	}
	v := url.Values{"q": {q}, "limit": {"12"}}
	if lat, lon, ok := requestCoordinates(r); ok {
		v.Set("lat", fmt.Sprint(lat))
		v.Set("lon", fmt.Sprint(lon))
	}
	a.fetch(w, r, a.photon+"/api/?"+v.Encode())
}
func (a *App) reverse(w http.ResponseWriter, r *http.Request) {
	var lat, lon float64
	if _, e := fmt.Sscan(r.URL.Query().Get("lat"), &lat); e != nil {
		apphost.Problem(w, 400, "Invalid latitude")
		return
	}
	if _, e := fmt.Sscan(r.URL.Query().Get("lon"), &lon); e != nil || !coordinate(lat, lon) {
		apphost.Problem(w, 400, "Invalid coordinates")
		return
	}
	a.fetch(w, r, fmt.Sprintf("%s/reverse?lat=%f&lon=%f", a.photon, lat, lon))
}
func (a *App) route(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Mode   string      `json:"mode"`
		Points [][]float64 `json:"points"`
	}
	if !apphost.Decode(w, r, &in) {
		return
	}
	if in.Mode != "car" && in.Mode != "bike" && in.Mode != "foot" {
		apphost.Problem(w, 400, "Choose driving, cycling or walking")
		return
	}
	if len(in.Points) < 2 || len(in.Points) > 8 {
		apphost.Problem(w, 400, "Routes need 2–8 stops")
		return
	}
	pts := []string{}
	for _, p := range in.Points {
		if len(p) != 2 || !coordinate(p[1], p[0]) {
			apphost.Problem(w, 400, "Invalid route stop")
			return
		}
		pts = append(pts, fmt.Sprintf("%.6f,%.6f", p[0], p[1]))
	}
	a.fetch(w, r, a.routing+"/routed-"+in.Mode+"/route/v1/driving/"+strings.Join(pts, ";")+"?overview=full&steps=true&geometries=geojson&alternatives=true")
}
func (a *App) places(w http.ResponseWriter, r *http.Request) {
	rows, e := a.H.DB.Query("SELECT id,name,address,lat,lon,list,notes FROM places WHERE owner=? ORDER BY list,name", apphost.Current(r).User.ID)
	if e != nil {
		a.H.DBError(w, e)
		return
	}
	defer rows.Close()
	out := []Place{}
	for rows.Next() {
		var p Place
		if e = rows.Scan(&p.ID, &p.Name, &p.Address, &p.Lat, &p.Lon, &p.List, &p.Notes); e != nil {
			a.H.DBError(w, e)
			return
		}
		out = append(out, p)
	}
	apphost.JSON(w, 200, out)
}
func (a *App) save(w http.ResponseWriter, r *http.Request) {
	var p Place
	if !apphost.Decode(w, r, &p) {
		return
	}
	p.Name = strings.TrimSpace(p.Name)
	p.List = strings.TrimSpace(p.List)
	if p.List == "" {
		p.List = "Favourites"
	}
	if !coordinate(p.Lat, p.Lon) || p.Name == "" || len(p.Name) > 200 || len(p.Address) > 1000 || len(p.List) > 100 || len(p.Notes) > 10000 {
		apphost.Problem(w, 400, "Check the place name, coordinates and notes")
		return
	}
	owner := apphost.Current(r).User.ID
	if p.ID == "" {
		p.ID = apphost.ID()
		var count int
		a.H.DB.QueryRow("SELECT count(*) FROM places WHERE owner=?", owner).Scan(&count)
		if count >= 10000 {
			apphost.Problem(w, 409, "Export or remove some saved places before adding more")
			return
		}
		_, e := a.H.DB.Exec("INSERT INTO places VALUES(?,?,?,?,?,?,?,?)", p.ID, owner, p.Name, p.Address, p.Lat, p.Lon, p.List, p.Notes)
		if e != nil {
			a.H.DBError(w, e)
			return
		}
	} else {
		res, e := a.H.DB.Exec("UPDATE places SET name=?,address=?,lat=?,lon=?,list=?,notes=? WHERE id=? AND owner=?", p.Name, p.Address, p.Lat, p.Lon, p.List, p.Notes, p.ID, owner)
		if e != nil {
			a.H.DBError(w, e)
			return
		}
		n, _ := res.RowsAffected()
		if n != 1 {
			apphost.Problem(w, 404, "Saved place not found")
			return
		}
	}
	apphost.JSON(w, 200, p)
}
func (a *App) remove(w http.ResponseWriter, r *http.Request) {
	res, e := a.H.DB.Exec("DELETE FROM places WHERE id=? AND owner=?", r.PathValue("id"), apphost.Current(r).User.ID)
	if e != nil {
		a.H.DBError(w, e)
		return
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		apphost.Problem(w, 404, "Saved place not found")
		return
	}
	apphost.JSON(w, 200, map[string]bool{"ok": true})
}
