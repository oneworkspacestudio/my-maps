package app

import (
	"strings"
	"testing"
)

func TestRecentKindsDoNotStoreAmbiguousOrInvalidEntries(t *testing.T) {
	lat, lon, invalid := 51.5, -0.1, 91.0
	for _, item := range []Recent{{Query: "Test query"}, {Name: "Test place", Address: "Example Lane", Lat: &lat, Lon: &lon}} {
		if !validRecent(item) {
			t.Fatal("rejected valid recent kind")
		}
	}
	for _, item := range []Recent{{}, {Query: "x"}, {Name: "No location"}, {Query: "Query", Name: "Ambiguous"}, {Query: "Query", Address: "Unneeded private field"}, {Name: "Partial", Lat: &lat}, {Name: "Invalid", Lat: &invalid, Lon: &lon}, {Query: strings.Repeat("x", 201)}} {
		if validRecent(item) {
			t.Fatalf("accepted ambiguous or invalid recent kind: %+v", item)
		}
	}
}
