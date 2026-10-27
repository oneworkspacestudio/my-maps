package app

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type providerTransport func(*http.Request) (*http.Response, error)

func (f providerTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func providerResponse(body string, status int) *http.Response {
	return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header)}
}

func TestNearbyPhotonIndependentOfRenderer(t *testing.T) {
	t.Setenv("MYMAPS_OVERPASS", "")
	a := &App{photon: "https://photon.example", client: &http.Client{Transport: providerTransport(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/reverse" || r.URL.Query().Get("radius") != "2.5" || r.URL.Query().Get("limit") != "60" {
			t.Fatal("nearby query not bounded")
		}
		if len(r.URL.Query()["osm_tag"]) != 5 {
			t.Fatal("restaurant filters missing")
		}
		if r.Header.Get("Cookie") != "" || r.Header.Get("Authorization") != "" {
			t.Fatal("forwarded owner credentials")
		}
		return providerResponse(`{"features":[
		{"geometry":{"type":"Point","coordinates":[-0.1,51.5]},"properties":{"osm_type":"N","osm_id":1,"osm_key":"amenity","osm_value":"cafe","name":"Synthetic cafe"}},
		{"geometry":{"type":"Point","coordinates":[-0.101,51.501]},"properties":{"osm_type":"W","osm_id":2,"osm_key":"amenity","osm_value":"restaurant","name":"Synthetic restaurant"}},
		{"geometry":{"type":"Point","coordinates":[-0.1,51.5]},"properties":{"osm_type":"N","osm_id":3,"osm_key":"shop","osm_value":"supermarket","name":"Wrong category"}},
		{"geometry":{"type":"Point","coordinates":[10,20]},"properties":{"osm_type":"N","osm_id":4,"osm_key":"amenity","osm_value":"restaurant","name":"Too far"}}
		]}`, 200), nil
	})}}
	r := httptest.NewRequest("GET", "/api/nearby?lat=51.5&lon=-0.1&category=restaurants", nil)
	r.Header.Set("Cookie", "private=synthetic")
	w := httptest.NewRecorder()
	a.nearby(w, r)
	if w.Code != 200 || !strings.Contains(w.Body.String(), `"source":"photon"`) || !strings.Contains(w.Body.String(), "Synthetic cafe") || strings.Contains(w.Body.String(), "Wrong category") || strings.Contains(w.Body.String(), "Too far") {
		t.Fatalf("unexpected result: %d %s", w.Code, w.Body.String())
	}
}

func TestNearbyProviderFailoverRetainsWayCenter(t *testing.T) {
	t.Setenv("MYMAPS_OVERPASS", "")
	calls := 0
	a := &App{photon: "https://photon.example", client: &http.Client{Transport: providerTransport(func(r *http.Request) (*http.Response, error) {
		calls++
		if calls == 1 {
			return providerResponse(`{}`, 503), nil
		}
		return providerResponse(`{"elements":[{"type":"way","id":7,"center":{"lat":51.5,"lon":-0.1},"tags":{"name":"Synthetic shop","shop":"bakery"}}]}`, 200), nil
	})}}
	w := httptest.NewRecorder()
	a.nearby(w, httptest.NewRequest("GET", "/api/nearby?lat=51.5&lon=-0.1&category=shopping", nil))
	if calls != 2 || w.Code != 200 || !strings.Contains(w.Body.String(), `"source":"overpass"`) || !strings.Contains(w.Body.String(), `"center":{"lat":51.5`) || strings.Contains(w.Body.String(), `"lat":0`) {
		t.Fatalf("failover/way center lost: %d %s", w.Code, w.Body.String())
	}
}

func TestNearbyBothProvidersFailClearly(t *testing.T) {
	t.Setenv("MYMAPS_OVERPASS", "")
	a := &App{photon: "https://photon.example", client: &http.Client{Transport: providerTransport(func(r *http.Request) (*http.Response, error) {
		return providerResponse(`{"remark":"provider failure"}`, 200), nil
	})}}
	w := httptest.NewRecorder()
	a.nearby(w, httptest.NewRequest("GET", "/api/nearby?lat=51.5&lon=-0.1&category=parking", nil))
	if w.Code != 502 || !strings.Contains(w.Body.String(), "Nearby places are temporarily unavailable") {
		t.Fatal("provider failures must not appear as empty success")
	}
}

func TestProviderQueueAndCancellation(t *testing.T) {
	var starts []time.Time
	a := &App{client: &http.Client{Transport: providerTransport(func(r *http.Request) (*http.Response, error) {
		starts = append(starts, time.Now())
		return providerResponse(`{}`, 200), nil
	})}}
	for i := 0; i < 2; i++ {
		if _, e := a.providerJSON(context.Background(), fmt.Sprintf("https://provider.example/%d", i), time.Hour); e != nil {
			t.Fatal(e)
		}
	}
	if len(starts) != 2 || starts[1].Sub(starts[0]) < 900*time.Millisecond {
		t.Fatal("provider requests not spaced")
	}
	last := a.last
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	if _, e := a.providerJSON(ctx, "https://provider.example/cancel", time.Hour); e == nil {
		t.Fatal("cancelled provider wait succeeded")
	}
	if len(starts) != 2 || !a.last.Equal(last) {
		t.Fatal("cancelled request reached provider or reserved future rate slots")
	}
}
