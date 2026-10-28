package app

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLocalIndexSearchPostcodeAndNearby(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	if _, err = db.Exec(packSchema); err != nil {
		t.Fatal(err)
	}
	insert, err := db.Prepare("INSERT OR IGNORE INTO features(feature_key,name,search_key,lat,lon,category,address) VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		t.Fatal(err)
	}
	defer insert.Close()
	bounds := [4]float64{-.2, 51.4, .2, 51.7}
	ctx := context.Background()
	if err = readGeo(ctx, insert, strings.NewReader("GB\tE4 6NW\tSynthetic district\tEngland\tx\tLondon\tx\tx\tx\t51.62\t0.01\t6\nGB\tZZ1 1ZZ\tOutside\tEngland\tx\tLondon\tx\tx\tx\t0\t0\t6"), "postcodes", bounds); err != nil {
		t.Fatal(err)
	}
	for _, f := range []localFeature{{Name: "Synthetic cafe", Category: "restaurants", Lat: 51.62, Lon: .01}, {Name: "Synthetic shop", Category: "shopping", Lat: 51.6201, Lon: .01}, {Name: "Distant cafe", Category: "restaurants", Lat: 52, Lon: 0}} {
		if err = addMapFeature(ctx, insert, f.Name, "", f.Category, f.Lat, f.Lon); err != nil {
			t.Fatal(err)
		}
	}
	db.Exec("INSERT INTO features_fts(features_fts) VALUES('rebuild')")
	m := &mapManager{db: db}
	for _, q := range []string{"e46nw", "E4%206NW", "Synthetic%20caf"} {
		w := httptest.NewRecorder()
		m.localSearch(w, httptest.NewRequest("GET", "/api/search?q="+q, nil), strings.ReplaceAll(q, "%20", " "))
		if w.Code != 200 || !strings.Contains(w.Body.String(), "coordinates") {
			t.Fatal(q, w.Code, w.Body.String())
		}
	}
	w := httptest.NewRecorder()
	m.localNearby(w, httptest.NewRequest("GET", "/api/nearby", nil), 51.62, .01, "restaurants")
	if w.Code != 200 || !strings.Contains(w.Body.String(), "Synthetic cafe") || strings.Contains(w.Body.String(), "Synthetic shop") || strings.Contains(w.Body.String(), "Distant") {
		t.Fatal(w.Code, w.Body.String())
	}
	if strings.Contains(w.Body.String(), `"amenity"`) {
		t.Fatal("fabricated source tags")
	}
}
func TestCoverageAndLocalEndpointConstraints(t *testing.T) {
	for _, c := range []MapChoice{{Coverage: "world", Depth: "streets"}, {Coverage: "area", Depth: "streets", Radius: 1000}, {Coverage: "area", Depth: "streets", Radius: 10, Lat: 89}, {Coverage: "world", Depth: "overview", Postcodes: true}, {Coverage: "uk", Depth: "fake"}} {
		if _, _, err := mapBounds(c); err == nil {
			t.Fatal("unbounded choice accepted", c)
		}
	}
	for _, s := range []string{"https://photon.komoot.io", "https://router.project-osrm.org", "http://localhost:5000", "http://127.0.0.1.evil.test", "http://user:pass@127.0.0.1:5000"} {
		if localMapEndpoint(s) {
			t.Fatal("non-local endpoint", s)
		}
	}
	for _, s := range []string{"http://127.0.0.1:5000", "http://[::1]:5000"} {
		if !localMapEndpoint(s) {
			t.Fatal(s)
		}
	}
}

func TestLocalOSRMStandardDrivingContract(t *testing.T) {
	calls := 0
	a := &App{routing: "http://127.0.0.1:5000", client: &http.Client{Transport: providerTransport(func(r *http.Request) (*http.Response, error) {
		calls++
		if !strings.HasPrefix(r.URL.Path, "/route/v1/driving/") || r.URL.Host != "127.0.0.1:5000" {
			t.Fatal("incorrect local OSRM route", r.URL)
		}
		return providerResponse(`{"code":"Ok","routes":[]}`, 200), nil
	})}}
	w := httptest.NewRecorder()
	a.route(w, httptest.NewRequest("POST", "/api/route", strings.NewReader(`{"mode":"car","points":[[0,51],[0.01,51.01]]}`)))
	if w.Code != 200 || calls != 1 {
		t.Fatal(w.Code, w.Body.String())
	}
	w = httptest.NewRecorder()
	a.route(w, httptest.NewRequest("POST", "/api/route", strings.NewReader(`{"mode":"foot","points":[[0,51],[0.01,51.01]]}`)))
	if w.Code != 503 || !strings.Contains(w.Body.String(), "directions_download") || calls != 1 {
		t.Fatal("unsupported mode reached driving graph")
	}
}
