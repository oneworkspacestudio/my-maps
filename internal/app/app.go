package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/oneworkspacestudio/my-maps/internal/directions"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type App struct {
	H                      *apphost.Host
	client                 *http.Client
	mu                     sync.Mutex
	placesMu               sync.Mutex
	installMu              sync.Mutex
	last                   time.Time
	photon, routing, tiles string
	maps                   *mapManager
	directions             *directions.Manager
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
	a := &App{H: h, client: &http.Client{Timeout: 20 * time.Second, CheckRedirect: func(next *http.Request, via []*http.Request) error {
		if len(via) > 3 || next.URL.Host != via[0].URL.Host || next.URL.Scheme != via[0].URL.Scheme || next.URL.User != nil {
			return errors.New("map service redirected outside its local origin")
		}
		return nil
	}}, photon: env("MYMAPS_GEOCODER", ""), routing: env("MYMAPS_ROUTER", ""), tiles: env("MYMAPS_TILES", "")}
	for _, endpoint := range []string{a.photon, a.routing, a.tiles} {
		if endpoint != "" && !localMapEndpoint(endpoint) {
			return nil, errors.New("map services must run at a literal loopback address; use map downloads for self-hosted data")
		}
	}
	h.Config.ImageOrigins = nil
	if e := migrate(h); e != nil {
		return nil, e
	}
	var e error
	a.maps, e = newMapManager(a)
	if e != nil {
		return nil, e
	}
	a.directions, e = directions.New(filepath.Join(h.Config.Data, "directions"), mapDiskFree)
	if e != nil {
		a.maps.close()
		return nil, e
	}
	h.Mux.HandleFunc("GET /api/config", h.Require(func(w http.ResponseWriter, r *http.Request) {
		a.maps.mu.RLock()
		defer a.maps.mu.RUnlock()
		info := a.maps.info
		attribution := env("MYMAPS_TILE_ATTRIBUTION", "© OpenStreetMap contributors · GeoNames (CC BY4.0)")
		if info != nil && info.Choice.Postcodes {
			attribution += " · Contains Royal Mail data ©2022"
		}
		localDirections := a.directions.Status()
		apphost.JSON(w, 200, map[string]any{"ready": info != nil || a.tiles != "", "local": info != nil, "installed": info, "canInstall": apphost.Current(r).User.Role == "owner", "localTiles": h.Base + "/api/local/tiles/{z}/{x}/{y}?v=" + packID(info), "nearbyAvailable": info != nil && info.Nearby || a.photon != "", "routingAvailable": localDirections.Ready || a.maps.routing != "", "directions": localDirections, "tiles": a.tiles, "customTiles": a.tiles != "", "attribution": attribution, "geocoder": a.photon})
	}))
	h.Mux.HandleFunc("GET /api/search", h.Require(a.search))
	h.Mux.HandleFunc("GET /api/reverse", h.Require(a.reverse))
	h.Mux.HandleFunc("GET /api/nearby", h.Require(a.nearby))
	h.Mux.HandleFunc("GET /api/map-resource", h.Require(a.mapResource))
	h.Mux.HandleFunc("POST /api/route", h.Require(a.route))
	h.Mux.HandleFunc("GET /api/places", h.Require(a.places))
	h.Mux.HandleFunc("POST /api/places", h.Require(a.save))
	h.Mux.HandleFunc("DELETE /api/places/{id}", h.Require(a.remove))
	h.Mux.HandleFunc("GET /api/recents", h.Require(a.recents))
	h.Mux.HandleFunc("POST /api/recents", h.Require(a.remember))
	h.Mux.HandleFunc("DELETE /api/recents", h.Require(a.clearRecents))
	h.Mux.HandleFunc("GET /api/map-install", h.Require(a.maps.status))
	h.Mux.HandleFunc("POST /api/map-install/plan", h.Owner(a.maps.makePlan))
	h.Mux.HandleFunc("POST /api/map-install/start", h.Owner(a.mapInstall))
	h.Mux.HandleFunc("POST /api/map-install/cancel", h.Owner(a.maps.cancelInstall))
	h.Mux.HandleFunc("POST /api/map-install/routing", h.Owner(a.maps.setRouting))
	h.Mux.HandleFunc("GET /api/directions", h.Require(a.directionsStatus))
	h.Mux.HandleFunc("POST /api/directions/plan", h.Owner(a.directionsPlan))
	h.Mux.HandleFunc("POST /api/directions/start", h.Owner(a.directionsInstall))
	h.Mux.HandleFunc("POST /api/directions/cancel", h.Owner(a.directionsCancel))
	h.Mux.HandleFunc("GET /api/local/tiles/{z}/{x}/{y}", h.Require(a.maps.tile))
	return a, nil
}

func localMapEndpoint(endpoint string) bool {
	u, err := url.Parse(endpoint)
	if err != nil || !validProviderEndpoint(endpoint) {
		return false
	}
	ip := net.ParseIP(u.Hostname())
	return ip != nil && ip.IsLoopback()
}

func validProviderEndpoint(endpoint string) bool {
	u, err := url.Parse(endpoint)
	if err != nil || u.Host == "" || u.User != nil || u.Fragment != "" {
		return false
	}
	ip := net.ParseIP(u.Hostname())
	return u.Scheme == "https" || (u.Scheme == "http" && ip != nil && ip.IsLoopback())
}
func (a *App) Run(ctx context.Context) error {
	defer a.maps.close()
	defer a.directions.Close()
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
	body, e := a.providerJSON(r.Context(), target, ttl)
	if e != nil {
		if r.Context().Err() == nil {
			apphost.Problem(w, 502, "The map provider is unavailable. Please try again shortly.")
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

// Share cached public provider data, never user cookies or credentials. Queue
// requests at one per second rather than rejecting ordinary fast UI actions.
func (a *App) providerJSON(ctx context.Context, target string, ttl time.Duration) ([]byte, error) {
	key := apphost.Hash(target)
	var body []byte
	if a.H != nil && a.H.DB.QueryRow("SELECT body FROM provider_cache WHERE key=? AND expires>?", key, time.Now().Unix()).Scan(&body) == nil {
		return body, nil
	}
	for {
		a.mu.Lock()
		delay := time.Until(a.last.Add(time.Second))
		if delay <= 0 {
			a.last = time.Now()
			a.mu.Unlock()
			break
		}
		a.mu.Unlock()
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, ctx.Err()
		case <-timer.C:
		}
	}
	req, e := http.NewRequestWithContext(ctx, "GET", target, nil)
	if e != nil {
		return nil, e
	}
	req.Header.Set("User-Agent", "OneWorkspace-MyMaps/0.1 (+https://github.com/oneworkspacestudio/my-maps)")
	req.Header.Set("Accept", "application/json")
	resp, e := a.client.Do(req)
	if e != nil {
		return nil, e
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("provider status %d", resp.StatusCode)
	}
	body, e = io.ReadAll(io.LimitReader(resp.Body, (4<<20)+1))
	if e != nil || len(body) > 4<<20 || !json.Valid(body) {
		return nil, errors.New("invalid provider JSON")
	}
	if a.H != nil {
		_, _ = a.H.DB.Exec("DELETE FROM provider_cache WHERE expires<?", time.Now().Unix())
		var count int
		a.H.DB.QueryRow("SELECT count(*) FROM provider_cache").Scan(&count)
		if count < 5000 {
			_, _ = a.H.DB.Exec("INSERT OR REPLACE INTO provider_cache VALUES(?,?,?)", key, body, time.Now().Add(ttl).Unix())
		}
	}
	return body, nil
}
func (a *App) search(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if len(q) < 2 || len(q) > 200 {
		apphost.Problem(w, 400, "Enter a place or address (2–200 characters)")
		return
	}
	if a.photon == "" {
		a.maps.localSearch(w, r, q)
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
	if a.photon == "" {
		a.maps.localReverse(w, r, lat, lon)
		return
	}
	a.fetch(w, r, fmt.Sprintf("%s/reverse?lat=%f&lon=%f", a.photon, lat, lon))
}
func (a *App) route(w http.ResponseWriter, r *http.Request) {
	var in directions.Request
	if !apphost.Decode(w, r, &in) {
		return
	}
	if err := in.Validate(); err != nil {
		apphost.Problem(w, 400, err.Error())
		return
	}
	if a.directions != nil && a.directions.Status().Ready {
		a.localRoute(w, r, in)
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
	router := a.routing
	if a.maps != nil {
		a.maps.mu.RLock()
		router = a.maps.routing
		a.maps.mu.RUnlock()
	}
	if router == "" {
		apphost.JSON(w, 503, map[string]string{"error": "Download local directions to drive, walk or cycle", "code": "directions_download"})
		return
	}
	if in.Mode != "car" || in.AvoidTolls || in.AvoidMotorways || in.AvoidUnpaved || in.Preference == "efficient" || in.Vehicle == "electric" || in.Vehicle == "hybrid" {
		apphost.JSON(w, 503, map[string]string{"error": "Download local directions for these travel modes and vehicle options", "code": "directions_download"})
		return
	}
	a.fetch(w, r, strings.TrimRight(router, "/")+"/route/v1/driving/"+strings.Join(pts, ";")+"?overview=full&steps=true&geometries=geojson&alternatives=true")
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
	p.List = canonicalList(p.List)
	a.placesMu.Lock()
	defer a.placesMu.Unlock()
	saved, status, err := persistPlace(r.Context(), a.H.DB, apphost.Current(r).User.ID, p)
	if err != nil {
		if status == 500 {
			a.H.DBError(w, err)
		} else {
			apphost.Problem(w, status, err.Error())
		}
		return
	}
	p = saved
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

func packID(info *MapInfo) string {
	if info == nil {
		return ""
	}
	return info.ID
}
