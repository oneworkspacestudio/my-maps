package app

import (
	"net/http/httptest"
	"testing"
)

func TestMapResourceAllowlist(t *testing.T) {
	for _, raw := range []string{"https://tiles.openfreemap.org/styles/liberty", "https://tiles.mapterhorn.com/1/2/3.webp", "https://gibs.earthdata.nasa.gov/wmts/example.jpeg"} {
		if !validMapResource(raw) {
			t.Errorf("rejected intended map resource %s", raw)
		}
	}
	for _, raw := range []string{"http://tiles.openfreemap.org/a", "https://tiles.openfreemap.org.evil.test/a", "https://tiles.openfreemap.org@127.0.0.1/a", "https://tiles.openfreemap.org:443/a", "https://127.0.0.1/a", "https://tiles.openfreemap.org/a?url=https://evil.test", "https://tiles.openfreemap.org/../a", "https://tiles.openfreemap.org/%2e%2e/a", "https://user:password@tiles.openfreemap.org/a"} {
		if validMapResource(raw) {
			t.Errorf("accepted unintended resource %s", raw)
		}
	}
}

func TestExploreRejectsInvalidInputBeforeFetching(t *testing.T) {
	a := &App{}
	for _, path := range []string{"/api/nearby?lat=NaN&lon=0&category=restaurants", "/api/nearby?lat=50&lon=181&category=restaurants", "/api/nearby?lat=50&lon=0&category=arbitrary-query", "/api/map-resource?url=https://127.0.0.1/private"} {
		r := httptest.NewRequest("GET", path, nil)
		w := httptest.NewRecorder()
		switch r.URL.Path {
		case "/api/nearby":
			a.nearby(w, r)
		default:
			a.mapResource(w, r)
		}
		if w.Code != 400 {
			t.Errorf("%s returned %d", path, w.Code)
		}
	}
}

func TestCustomTileTemplateIsNotAnOpenProxy(t *testing.T) {
	a := &App{tiles: "https://example.test/maps/{z}/{x}/{y}.png?key=configured"}
	if !a.customTileResource("https://example.test/maps/12/4/9.png?key=configured") {
		t.Fatal("configured tile rejected")
	}
	for _, raw := range []string{"https://example.test/private", "https://example.test/maps/12/4/9.png?key=other", "https://example.test/maps/../4/9.png?key=configured"} {
		if a.customTileResource(raw) {
			t.Errorf("template accepted %s", raw)
		}
	}
}
