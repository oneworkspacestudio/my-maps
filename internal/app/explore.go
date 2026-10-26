package app

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

func requestCoordinates(r *http.Request) (float64, float64, bool) {
	lat, a := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	lon, b := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)
	return lat, lon, a == nil && b == nil && coordinate(lat, lon)
}

// Deliberately bounded, user-triggered POI queries. Never accept arbitrary Overpass QL.
func (a *App) nearby(w http.ResponseWriter, r *http.Request) {
	lat, lon, ok := requestCoordinates(r)
	if !ok {
		apphost.Problem(w, 400, "Choose a location on the map first")
		return
	}
	filters := map[string]string{
		"restaurants": `["amenity"~"^(restaurant|cafe|fast_food|pub|bar)$"]`,
		"shopping":    `["shop"]["name"]`,
		"hotels":      `["tourism"~"^(hotel|guest_house|hostel|motel)$"]`,
		"attractions": `["tourism"~"^(attraction|museum|viewpoint|gallery)$"]`,
		"fuel":        `["amenity"="fuel"]`,
		"parking":     `["amenity"="parking"]`,
		"transit":     `["public_transport"~"^(station|platform|stop_position)$"]`,
		"cycling":     `["amenity"~"^(bicycle_parking|bicycle_rental|bicycle_repair_station)$"]`,
	}
	filter, ok := filters[r.URL.Query().Get("category")]
	if !ok {
		apphost.Problem(w, 400, "Choose a nearby category")
		return
	}
	// Round to ~100m so panning small distances reuses the same cached request.
	query := fmt.Sprintf(`[out:json][timeout:15];nwr(around:2500,%.3f,%.3f)%s;out center tags 60;`, lat, lon, filter)
	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()
	a.fetch(w, r.WithContext(ctx), env("MYMAPS_OVERPASS", "https://overpass-api.de")+"/api/interpreter?data="+url.QueryEscape(query))
}

func (a *App) weather(w http.ResponseWriter, r *http.Request) {
	lat, lon, ok := requestCoordinates(r)
	if !ok {
		apphost.Problem(w, 400, "Choose a location for the weather")
		return
	}
	if r.URL.Query().Get("air") == "1" {
		a.fetchTTL(w, r, fmt.Sprintf("https://air-quality-api.open-meteo.com/v1/air-quality?latitude=%.2f&longitude=%.2f&current=european_aqi,pm2_5,pm10&timezone=auto", lat, lon), 15*time.Minute)
		return
	}
	a.fetchTTL(w, r, fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.2f&longitude=%.2f&current=temperature_2m,weather_code&timezone=auto", lat, lon), 15*time.Minute)
}

var resourceHosts = map[string]bool{
	"tiles.openfreemap.org":   true,
	"tiles.mapterhorn.com":    true,
	"gibs.earthdata.nasa.gov": true,
}

func validMapResource(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && u.Scheme == "https" && resourceHosts[u.Host] && u.User == nil && u.Fragment == "" && u.RawQuery == "" && !strings.Contains(u.Path, "..")
}

func (a *App) customTileResource(raw string) bool {
	pattern := regexp.QuoteMeta(a.tiles)
	for _, key := range []string{"z", "x", "y"} {
		pattern = strings.ReplaceAll(pattern, regexp.QuoteMeta("{"+key+"}"), `[0-9]{1,10}`)
	}
	return a.tiles != "" && regexp.MustCompile("^"+pattern+"$").MatchString(raw)
}

// Fixed-host resource relay keeps the app's restrictive same-origin connect CSP.
// No credentials are forwarded and no arbitrary URLs, query strings or redirects
// to other hosts are allowed. Browser caches honour upstream freshness headers.
func (a *App) mapResource(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("url")
	if len(target) > 2000 || !(validMapResource(target) || a.customTileResource(target)) {
		apphost.Problem(w, 400, "Unsupported map resource")
		return
	}
	req, err := http.NewRequestWithContext(r.Context(), "GET", target, nil)
	if err != nil {
		apphost.Problem(w, 400, "Invalid map resource")
		return
	}
	req.Header.Set("User-Agent", "OneWorkspace-MyMaps/0.2 (+https://github.com/oneworkspacestudio/my-maps)")
	client := *a.client
	client.CheckRedirect = func(next *http.Request, via []*http.Request) error {
		if len(via) > 3 || !validMapResource(next.URL.String()) {
			return fmt.Errorf("unsupported map redirect")
		}
		return nil
	}
	resp, err := client.Do(req)
	if err != nil {
		apphost.Problem(w, 502, "Map imagery is unavailable; retry shortly")
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		apphost.Problem(w, 502, "This map resource is currently unavailable")
		return
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, (8<<20)+1))
	if err != nil || len(body) > 8<<20 {
		apphost.Problem(w, 502, "Map resource exceeds the size limit")
		return
	}
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Header().Set("Cache-Control", "private, max-age=86400")
	if cache := resp.Header.Get("Cache-Control"); cache != "" {
		w.Header().Set("Cache-Control", "private, "+cache)
	}
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Write(body)
}
