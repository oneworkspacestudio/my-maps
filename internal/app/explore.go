package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
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
	// Photon offers bounded reverse/category search independently of the map
	// renderer. Explicit Overpass installations keep their configured priority.
	key := r.URL.Query().Get("category")
	providers := []string{"photon", "overpass"}
	if env("MYMAPS_OVERPASS", "") != "" {
		providers = []string{"overpass", "photon"}
	}
	for _, provider := range providers {
		ctx, cancel := context.WithTimeout(r.Context(), 6*time.Second)
		var elements []nearbyElement
		var err error
		if provider == "photon" {
			elements, err = a.photonNearby(ctx, lat, lon, key)
		} else {
			query := fmt.Sprintf(`[out:json][timeout:5];nwr(around:2500,%.3f,%.3f)%s;out center tags 60;`, lat, lon, filter)
			var body []byte
			body, err = a.providerJSON(ctx, env("MYMAPS_OVERPASS", "https://overpass-api.de")+"/api/interpreter?data="+url.QueryEscape(query), time.Hour)
			if err == nil {
				var result struct {
					Elements *[]nearbyElement `json:"elements"`
				}
				err = json.Unmarshal(body, &result)
				if err == nil {
					if result.Elements == nil {
						err = fmt.Errorf("invalid nearby element collection")
					} else {
						elements = *result.Elements
					}
				}
			}
		}
		cancel()
		if r.Context().Err() != nil {
			return
		}
		if err == nil {
			if elements == nil {
				elements = []nearbyElement{}
			}
			if len(elements) > 60 {
				elements = elements[:60]
			}
			apphost.JSON(w, 200, map[string]any{"elements": elements, "source": provider})
			return
		}
	}
	apphost.Problem(w, 502, "Nearby places are temporarily unavailable. Move the map or try again shortly.")
}

type nearbyElement struct {
	Type   string   `json:"type"`
	ID     int64    `json:"id"`
	Lat    *float64 `json:"lat,omitempty"`
	Lon    *float64 `json:"lon,omitempty"`
	Center *struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"center,omitempty"`
	Tags map[string]string `json:"tags"`
}

var nearbyTags = map[string][]string{
	"restaurants": {"amenity:restaurant", "amenity:cafe", "amenity:fast_food", "amenity:pub", "amenity:bar"},
	"shopping":    {"shop"}, "hotels": {"tourism:hotel", "tourism:guest_house", "tourism:hostel", "tourism:motel"},
	"attractions": {"tourism:attraction", "tourism:museum", "tourism:viewpoint", "tourism:gallery"},
	"fuel":        {"amenity:fuel"}, "parking": {"amenity:parking"},
	"transit": {"public_transport", "highway:bus_stop", "railway:station", "railway:tram_stop", "amenity:bus_station"},
	"cycling": {"amenity:bicycle_parking", "amenity:bicycle_rental", "amenity:bicycle_repair_station"},
}

func (a *App) photonNearby(ctx context.Context, lat, lon float64, category string) ([]nearbyElement, error) {
	v := url.Values{"lat": {fmt.Sprintf("%.3f", lat)}, "lon": {fmt.Sprintf("%.3f", lon)}, "radius": {"2.5"}, "limit": {"60"}, "osm_tag": nearbyTags[category]}
	body, err := a.providerJSON(ctx, a.photon+"/reverse?"+v.Encode(), time.Hour)
	if err != nil {
		return nil, err
	}
	var result struct {
		Features *[]struct {
			Geometry struct {
				Type        string    `json:"type"`
				Coordinates []float64 `json:"coordinates"`
			} `json:"geometry"`
			Properties struct {
				Name   string         `json:"name"`
				Street string         `json:"street"`
				House  string         `json:"housenumber"`
				City   string         `json:"city"`
				Key    string         `json:"osm_key"`
				Value  string         `json:"osm_value"`
				Type   string         `json:"osm_type"`
				ID     int64          `json:"osm_id"`
				Extra  map[string]any `json:"extra"`
			} `json:"properties"`
		} `json:"features"`
	}
	if err = json.Unmarshal(body, &result); err != nil || result.Features == nil {
		return nil, fmt.Errorf("invalid nearby feature collection")
	}
	out := []nearbyElement{}
	for _, f := range *result.Features {
		p, xy := f.Properties, f.Geometry.Coordinates
		if f.Geometry.Type != "Point" || len(xy) != 2 || !coordinate(xy[1], xy[0]) || geoDistance(lat, lon, xy[1], xy[0]) > 2650 {
			continue
		}
		matches := false
		for _, tag := range nearbyTags[category] {
			if tag == p.Key || tag == p.Key+":"+p.Value {
				matches = true
			}
		}
		kind := map[string]string{"N": "node", "W": "way", "R": "relation"}[p.Type]
		if !matches || kind == "" || p.ID <= 0 {
			continue
		}
		tags := map[string]string{"name": p.Name, p.Key: p.Value, "addr:housenumber": p.House, "addr:street": p.Street, "addr:city": p.City}
		for _, k := range []string{"opening_hours", "phone", "contact:phone", "website", "contact:website"} {
			if value, ok := p.Extra[k].(string); ok {
				tags[k] = value
			}
		}
		px, py := xy[0], xy[1]
		out = append(out, nearbyElement{Type: kind, ID: p.ID, Lat: &py, Lon: &px, Tags: tags})
	}
	return out, nil
}

func geoDistance(lat, lon, otherLat, otherLon float64) float64 {
	dLat, dLon := (otherLat-lat)*math.Pi/180, (otherLon-lon)*math.Pi/180
	x := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat*math.Pi/180)*math.Cos(otherLat*math.Pi/180)*math.Sin(dLon/2)*math.Sin(dLon/2)
	return 12742000 * math.Asin(math.Sqrt(math.Min(1, x)))
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
