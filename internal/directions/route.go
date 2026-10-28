package directions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

var ErrNoData = errors.New("download local directions for driving, walking and cycling")
var ErrCoverage = errors.New("these stops are outside the installed directions coverage")
var ErrNoRoute = errors.New("no route connects these stops with the selected vehicle and options")

type Request struct {
	Mode           string      `json:"mode"`
	Vehicle        string      `json:"vehicle"`
	Preference     string      `json:"preference"`
	AvoidTolls     bool        `json:"avoidTolls"`
	AvoidMotorways bool        `json:"avoidMotorways"`
	AvoidUnpaved   bool        `json:"avoidUnpaved"`
	Alternatives   bool        `json:"alternatives"`
	Points         [][]float64 `json:"points"`
}
type Maneuver struct {
	Type     string    `json:"type"`
	Modifier string    `json:"modifier,omitempty"`
	Location []float64 `json:"location"`
	Exit     int       `json:"exit,omitempty"`
}
type Step struct {
	Name        string   `json:"name"`
	Distance    float64  `json:"distance"`
	Duration    float64  `json:"duration"`
	Maneuver    Maneuver `json:"maneuver"`
	Index       int      `json:"index"`
	Instruction string   `json:"instruction"`
}
type Leg struct {
	Steps []Step `json:"steps"`
}
type Geometry struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}
type Route struct {
	Distance float64  `json:"distance"`
	Duration float64  `json:"duration"`
	Geometry Geometry `json:"geometry"`
	Legs     []Leg    `json:"legs"`
}
type Result struct {
	Code        string  `json:"code"`
	Routes      []Route `json:"routes"`
	Mode        string  `json:"mode"`
	Vehicle     string  `json:"vehicle"`
	LiveTraffic bool    `json:"liveTraffic"`
}

func (r Request) Profile() (string, error) {
	switch r.Mode {
	case "car":
		if r.Vehicle != "" && r.Vehicle != "car" && r.Vehicle != "electric" && r.Vehicle != "hybrid" {
			return "", errors.New("choose a supported car type")
		}
		if r.Preference == "efficient" {
			return "car-eco", nil
		}
		if r.Preference != "" && r.Preference != "fastest" {
			return "", errors.New("choose fastest or lower-energy routing")
		}
		return "car-fast", nil
	case "bike":
		switch r.Vehicle {
		case "", "touring":
			return "trekking", nil
		case "road":
			return "fastbike", nil
		case "gravel":
			return "gravel", nil
		case "mountain":
			return "mtb", nil
		default:
			return "", errors.New("choose a supported bicycle type")
		}
	case "foot":
		if r.Vehicle == "hiking" {
			return "hiking-mountain", nil
		}
		if r.Vehicle != "" && r.Vehicle != "walking" {
			return "", errors.New("choose walking or hiking")
		}
		return "shortest", nil
	default:
		return "", errors.New("choose Drive, Walk or Cycle")
	}
}
func (r Request) Validate() error {
	if _, err := r.Profile(); err != nil {
		return err
	}
	if len(r.Points) < 2 || len(r.Points) > 8 {
		return errors.New("routes need 2–8 stops")
	}
	for _, p := range r.Points {
		if len(p) != 2 || math.IsNaN(p[0]) || math.IsNaN(p[1]) || math.IsInf(p[0], 0) || math.IsInf(p[1], 0) || math.Abs(p[0]) > 180 || math.Abs(p[1]) > 85 {
			return errors.New("invalid route stop")
		}
	}
	return nil
}
func (m *Manager) Route(ctx context.Context, in Request) (*Result, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, 85*time.Second)
	defer cancel()
	select {
	case m.routeGate <- struct{}{}:
		defer func() { <-m.routeGate }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	m.mu.Lock()
	info := m.info
	m.mu.Unlock()
	if info == nil {
		return nil, ErrNoData
	}
	segments := map[string]bool{}
	for _, f := range info.Files {
		segments[f.Name] = true
	}
	for _, p := range in.Points {
		if !segments[segmentName(int(math.Floor(p[0]/5))*5, int(math.Floor(p[1]/5))*5)] {
			return nil, ErrCoverage
		}
	}
	engine, err := m.engineFor(ctx)
	if err != nil {
		return nil, err
	}
	profile, _ := in.Profile()
	values := url.Values{"profile": {profile}, "alternativeidx": {"0"}}
	points := make([]string, len(in.Points))
	for i, p := range in.Points {
		points[i] = fmt.Sprintf("%.6f,%.6f", p[0], p[1])
	}
	values.Set("lonlats", strings.Join(points, "|"))
	if in.Mode == "car" {
		if in.AvoidTolls {
			values.Set("profile:avoid_toll", "1")
		}
		if in.AvoidMotorways {
			values.Set("profile:avoid_motorways", "1")
		}
		if in.AvoidUnpaved {
			values.Set("profile:avoid_unpaved", "1")
		}
		recuperation := "0"
		if in.Vehicle == "electric" {
			recuperation = "0.7"
		} else if in.Vehicle == "hybrid" {
			recuperation = "0.45"
		}
		values.Set("profile:recup_efficiency", recuperation)
	}
	result := &Result{Code: "Ok", Mode: in.Mode, Vehicle: in.Vehicle}
	count := 1
	if in.Alternatives {
		count = 2
	}
	for alt := 0; alt < count; alt++ {
		values.Set("alternativeidx", strconv.Itoa(alt))
		requestID := randomID()
		values.Set("requestId", requestID)
		req, _ := http.NewRequestWithContext(ctx, "POST", engine.url+"/route", strings.NewReader(values.Encode()))
		req.Header.Set("Authorization", "Bearer "+engine.token)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		client := localClient()
		stopCancel := context.AfterFunc(ctx, func() {
			abortCtx, stop := context.WithTimeout(context.Background(), 2*time.Second)
			defer stop()
			body := url.Values{"requestId": {requestID}}.Encode()
			abort, _ := http.NewRequestWithContext(abortCtx, "POST", engine.url+"/cancel", strings.NewReader(body))
			abort.Header.Set("Authorization", "Bearer "+engine.token)
			if res, e := client.Do(abort); e == nil {
				res.Body.Close()
			}
		})
		response, callErr := client.Do(req)
		if callErr != nil {
			stopCancel()
			client.CloseIdleConnections()
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			return nil, errors.New("the local directions engine could not be reached; retry")
		}
		body, readErr := io.ReadAll(io.LimitReader(response.Body, 8<<20+1))
		response.Body.Close()
		stopCancel()
		client.CloseIdleConnections()
		if readErr != nil || len(body) > 8<<20 {
			return nil, errors.New("local route response was too large or incomplete")
		}
		if response.StatusCode != 200 {
			if alt > 0 {
				break
			}
			return nil, ErrNoRoute
		}
		route, parseErr := normalize(body)
		if parseErr != nil {
			return nil, parseErr
		}
		if alt == 0 || math.Abs(result.Routes[0].Distance-route.Distance) > 10 {
			result.Routes = append(result.Routes, *route)
		}
	}
	if len(result.Routes) == 0 {
		return nil, ErrNoRoute
	}
	return result, nil
}
func localClient() *http.Client {
	return &http.Client{Transport: &http.Transport{Proxy: nil}, CheckRedirect: func(*http.Request, []*http.Request) error { return errors.New("routing cannot leave this process") }}
}
func number(raw json.RawMessage) float64 {
	var n float64
	if json.Unmarshal(raw, &n) == nil {
		return n
	}
	var s string
	if json.Unmarshal(raw, &s) == nil {
		v, _ := strconv.ParseFloat(s, 64)
		return v
	}
	return 0
}
func normalize(body []byte) (*Route, error) {
	var data struct {
		Features []struct {
			Geometry   Geometry                   `json:"geometry"`
			Properties map[string]json.RawMessage `json:"properties"`
		} `json:"features"`
	}
	if json.Unmarshal(body, &data) != nil || len(data.Features) != 1 {
		return nil, errors.New("invalid local route response")
	}
	f := data.Features[0]
	r := &Route{Distance: number(f.Properties["track-length"]), Duration: number(f.Properties["total-time"]), Geometry: f.Geometry}
	if r.Geometry.Type != "LineString" || len(r.Geometry.Coordinates) < 2 || len(r.Geometry.Coordinates) > 100000 || r.Distance <= 0 || math.IsNaN(r.Distance) || math.IsInf(r.Distance, 0) || r.Duration <= 0 || math.IsNaN(r.Duration) || math.IsInf(r.Duration, 0) {
		return nil, ErrNoRoute
	}
	for i, p := range r.Geometry.Coordinates {
		if len(p) < 2 || math.IsNaN(p[0]) || math.IsNaN(p[1]) || math.IsInf(p[0], 0) || math.IsInf(p[1], 0) || math.Abs(p[0]) > 180 || math.Abs(p[1]) > 85 {
			return nil, errors.New("invalid local route geometry")
		}
		r.Geometry.Coordinates[i] = p[:2]
	}
	var hints [][]float64
	_ = json.Unmarshal(f.Properties["voicehints"], &hints)
	if len(hints) > 10000 {
		return nil, errors.New("too many local route instructions")
	}
	sort.SliceStable(hints, func(i, j int) bool { return len(hints[i]) > 0 && (len(hints[j]) == 0 || hints[i][0] < hints[j][0]) })
	steps := []Step{{Index: 0, Instruction: "Start", Maneuver: Maneuver{Type: "depart", Location: r.Geometry.Coordinates[0]}}}
	for _, h := range hints {
		if len(h) < 5 || h[0] < 1 || h[0] >= float64(len(r.Geometry.Coordinates)-1) || math.IsNaN(h[0]) || h[0] != math.Trunc(h[0]) {
			continue
		}
		index := int(h[0])
		if steps[len(steps)-1].Index == index {
			continue
		}
		instruction, kind, modifier := turn(int(h[1]), int(h[2]))
		steps = append(steps, Step{Index: index, Instruction: instruction, Maneuver: Maneuver{Type: kind, Modifier: modifier, Location: r.Geometry.Coordinates[index], Exit: int(h[2])}})
	}
	steps = append(steps, Step{Index: len(r.Geometry.Coordinates) - 1, Instruction: "Arrive at your destination", Maneuver: Maneuver{Type: "arrive", Location: r.Geometry.Coordinates[len(r.Geometry.Coordinates)-1]}})
	for i := 0; i < len(steps)-1; i++ {
		for j := steps[i].Index; j < steps[i+1].Index; j++ {
			steps[i].Distance += geoDistance(r.Geometry.Coordinates[j], r.Geometry.Coordinates[j+1])
		}
		steps[i].Duration = r.Duration * steps[i].Distance / r.Distance
	}
	r.Legs = []Leg{{Steps: steps}}
	return r, nil
}
func turn(command, exit int) (string, string, string) {
	switch command {
	case 2:
		return "Turn left", "turn", "left"
	case 3:
		return "Bear left", "turn", "slight left"
	case 4:
		return "Turn sharply left", "turn", "sharp left"
	case 5:
		return "Turn right", "turn", "right"
	case 6:
		return "Bear right", "turn", "slight right"
	case 7:
		return "Turn sharply right", "turn", "sharp right"
	case 8, 17:
		return "Keep left", "fork", "left"
	case 9, 18:
		return "Keep right", "fork", "right"
	case 10, 11, 15:
		return "Make a U-turn", "turn", "uturn"
	case 13, 14:
		return fmt.Sprintf("At the roundabout, take exit %d", exit), "roundabout", ""
	default:
		return "Continue", "continue", ""
	}
}
func geoDistance(a, b []float64) float64 {
	rad := math.Pi / 180
	dy, dx := (b[1]-a[1])*rad, (b[0]-a[0])*rad
	s := math.Sin(dy/2)*math.Sin(dy/2) + math.Cos(a[1]*rad)*math.Cos(b[1]*rad)*math.Sin(dx/2)*math.Sin(dx/2)
	return 12742000 * math.Asin(math.Sqrt(math.Min(1, s)))
}
