package geolocator_test

import (
	"testing"

	"github.com/paulmach/orb"
	"github.com/spotlightpa/viz-redistricting-2020/geolocator"
)

func TestGetDistrict(t *testing.T) {
	for name, tc := range map[string]struct {
		Point orb.Point
		Name  string
	}{
		"blank":            {},
		"Harrisburg":       {orb.Point{-76.88375, 40.26444}, "103"},
		"includes islands": {orb.Point{-80.34728, 41.00326}, "9"},
		"exclusion island": {orb.Point{-80.38492, 40.96953}, "10"},
		"wampum":           {orb.Point{-80.33811, 40.88811}, "10"},
	} {
		t.Run(name, func(t *testing.T) {
			d := geolocator.House2012Map.District(tc.Point)
			if tc.Name != d.Name() {
				t.Fatalf("want %q; got %q", tc.Name, d.Name())
			}
		})
	}
}

func BenchmarkGetDistrict(b *testing.B) {
	cases := []struct {
		Point orb.Point
		Name  string
	}{
		{},
		{orb.Point{-76.88375, 40.26444}, "103"},
		{orb.Point{-80.34728, 41.00326}, "9"},
		{orb.Point{-80.38492, 40.96953}, "10"},
		{orb.Point{-80.33811, 40.88811}, "10"},
	}
	for i := 0; i < b.N; i++ {
		tc := cases[i%len(cases)]
		d := geolocator.House2012Map.District(tc.Point)
		if d.Name() != tc.Name {
			b.FailNow()
		}
	}
}
