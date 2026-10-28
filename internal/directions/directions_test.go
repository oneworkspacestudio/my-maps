package directions

import (
	"archive/zip"
	"context"
	"encoding/json"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTravelProfilesAndCoordinates(t *testing.T) {
	for _, tt := range []struct{ mode, vehicle, preference, want string }{{"car", "car", "fastest", "car-fast"}, {"car", "electric", "efficient", "car-eco"}, {"car", "hybrid", "fastest", "car-fast"}, {"bike", "touring", "", "trekking"}, {"bike", "road", "", "fastbike"}, {"bike", "gravel", "", "gravel"}, {"bike", "mountain", "", "mtb"}, {"foot", "walking", "", "shortest"}, {"foot", "hiking", "", "hiking-mountain"}} {
		r := Request{Mode: tt.mode, Vehicle: tt.vehicle, Preference: tt.preference, Points: [][]float64{{-.67, 50.78}, {-.66, 50.79}}}
		got, e := r.Profile()
		if e != nil || got != tt.want || r.Validate() != nil {
			t.Fatalf("%+v => %s %v", tt, got, e)
		}
	}
	for _, r := range []Request{{Mode: "transit"}, {Mode: "car", Vehicle: "moped"}, {Mode: "car", Vehicle: "truck"}, {Mode: "bike", Vehicle: "electric"}, {Mode: "foot", Vehicle: "car"}, {Mode: "car", Preference: "traffic"}} {
		if _, e := r.Profile(); e == nil {
			t.Fatal("unsupported profile accepted", r)
		}
	}
	for _, point := range [][]float64{{math.NaN(), 51}, {0, math.Inf(1)}, {181, 0}, {0, 86}, {0}, {}} {
		if e := (Request{Mode: "car", Points: [][]float64{point, {0, 51}}}).Validate(); e == nil {
			t.Fatal("invalid point accepted")
		}
	}
	if e := (Request{Mode: "car", Points: make([][]float64, 9)}).Validate(); e == nil {
		t.Fatal("unbounded stops accepted")
	}
}
func TestCoveragePlanning(t *testing.T) {
	for _, tt := range []struct {
		lon, lat int
		want     string
	}{{-5, 50, "W5_N50.rd5"}, {0, 50, "E0_N50.rd5"}, {10, -5, "E10_S5.rd5"}} {
		if got := segmentName(tt.lon, tt.lat); got != tt.want {
			t.Fatal(got)
		}
	}
	bounds, e := choiceBounds(Choice{Coverage: "uk"})
	if e != nil {
		t.Fatal(e)
	}
	names, e := segmentsFor(bounds)
	if e != nil || len(names) != 12 {
		t.Fatal(names, e)
	}
	for _, c := range []Choice{{Coverage: "world"}, {Coverage: "area", Lat: math.NaN(), Radius: 10}, {Coverage: "area", Lat: 85, Radius: 10}, {Coverage: "area", Radius: 501}, {Coverage: "area", Lon: 179.999, Radius: 20}} {
		if _, e := choiceBounds(c); e == nil {
			t.Fatal("unbounded coverage", c)
		}
	}
	if _, e := segmentsFor([4]float64{-180, -80, 180, 80}); e == nil {
		t.Fatal("large segment plan")
	}
}
func TestNormalizeTurnsSortedBoundedAndFinite(t *testing.T) {
	raw := `{"features":[{"geometry":{"type":"LineString","coordinates":[[0,51,2],[0.001,51,2],[0.001,51.001,2],[0.002,51.001,2],[0.003,51.001,2]]},"properties":{"track-length":"400","total-time":"120","voicehints":[[3,5,0,0,0],[1,2,0,0,0],[1,2,0,0,0],[-1,5,0,0,0]]}}]}`
	r, e := normalize([]byte(raw))
	if e != nil {
		t.Fatal(e)
	}
	steps := r.Legs[0].Steps
	if len(steps) != 4 || steps[1].Instruction != "Turn left" || steps[2].Instruction != "Turn right" || len(r.Geometry.Coordinates[0]) != 2 {
		t.Fatal(steps)
	}
	for _, s := range steps {
		if s.Distance < 0 || s.Duration < 0 {
			t.Fatal("negative step")
		}
	}
	for _, bad := range []string{strings.Replace(raw, `"400"`, `"NaN"`, 1), strings.Replace(raw, `"120"`, `"+Inf"`, 1), strings.Replace(raw, `"LineString"`, `"Point"`, 1)} {
		if _, e := normalize([]byte(bad)); e == nil {
			t.Fatal("invalid engine output accepted")
		}
	}
	if _, e := normalize([]byte(`{"features":[]}`)); e == nil {
		t.Fatal("missing route accepted")
	}
}
func TestRuntimeArchiveCannotEscape(t *testing.T) {
	for _, name := range []string{"jdk/../../escape.txt", "jdk/C:/escape.txt", "jdk/../escape.txt", "jdk/a\\b.txt"} {
		t.Run(strings.ReplaceAll(name, "/", "-"), func(t *testing.T) {
			dir := t.TempDir()
			archive := filepath.Join(dir, "runtime.zip")
			out, e := os.Create(archive)
			if e != nil {
				t.Fatal(e)
			}
			z := zip.NewWriter(out)
			w, _ := z.Create(name)
			w.Write([]byte("bad"))
			z.Close()
			out.Close()
			if e := extractRuntime(archive, filepath.Join(dir, "extract")); e == nil {
				t.Fatal("unsafe archive accepted", name)
			}
			if _, e := os.Stat(filepath.Join(dir, "escape.txt")); e == nil {
				t.Fatal("archive escaped")
			}
		})
	}
}
func TestNoDataAndCoverageNeverCallNetwork(t *testing.T) {
	m := &Manager{routeGate: make(chan struct{}, 1)}
	in := Request{Mode: "car", Points: [][]float64{{0, 51}, {1, 51}}}
	if _, e := m.Route(context.Background(), in); e != ErrNoData {
		t.Fatal(e)
	}
	m.info = &Info{Files: []File{{Name: "W5_N50.rd5"}}}
	if _, e := m.Route(context.Background(), in); e != ErrCoverage {
		t.Fatal(e)
	}
}
func TestDownloadIntegrityAndCancellation(t *testing.T) {
	// A synthetic round tripper avoids external network or real runtime startup.
	m := &Manager{client: &http.Client{Transport: roundTrip(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, ContentLength: 3, Body: http.NoBody, Header: http.Header{}}, nil
	})}}
	_, e := m.download(context.Background(), File{Name: "x", Bytes: 3, SHA256: strings.Repeat("0", 64), URL: "https://example.invalid/x"}, filepath.Join(t.TempDir(), "x"))
	if e == nil {
		t.Fatal("truncated download accepted")
	}
	if _, e := json.Marshal(Request{Mode: "car", Points: [][]float64{{math.NaN(), 0}}}); e == nil {
		t.Fatal("non-finite request serializable")
	}
}

type roundTrip func(*http.Request) (*http.Response, error)

func (f roundTrip) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
