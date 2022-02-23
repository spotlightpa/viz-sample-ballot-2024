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
		geolocator.Map
	}{
		"blank":             {Map: geolocator.House2012Map},
		"Harrisburg":        {orb.Point{-76.88375, 40.26444}, "103", geolocator.House2012Map},
		"Harrisburg CD18":   {orb.Point{-76.88375, 40.26444}, "10", geolocator.Congress2018Map},
		"Harrisburg CD22":   {orb.Point{-76.88375, 40.26444}, "10", geolocator.Congress2022Map},
		"includes islands":  {orb.Point{-80.34728, 41.00326}, "9", geolocator.House2012Map},
		"exclusion island":  {orb.Point{-80.38492, 40.96953}, "10", geolocator.House2012Map},
		"wampum":            {orb.Point{-80.33811, 40.88811}, "10", geolocator.House2012Map},
		"senate Harrisburg": {orb.Point{-76.88375, 40.26444}, "15", geolocator.Senate2012Map},
		"new h Harrisburg":  {orb.Point{-76.88375, 40.26444}, "103", geolocator.House2022Map},
		"new s Harrisburg":  {orb.Point{-76.88375, 40.26444}, "15", geolocator.Senate2022Map},
	} {
		t.Run(name, func(t *testing.T) {
			d := tc.Map.District(tc.Point)
			if tc.Name != d.GetName() {
				t.Fatalf("want %q; got %q", tc.Name, d.GetName())
			}
		})
	}
}

func BenchmarkGetDistrict(b *testing.B) {
	cases := []struct {
		Point orb.Point
		Name  string
		geolocator.Map
	}{
		{orb.Point{-76.88375, 40.26444}, "103", geolocator.House2012Map},
		{orb.Point{-80.34728, 41.00326}, "9", geolocator.House2012Map},
		{orb.Point{-80.38492, 40.96953}, "10", geolocator.House2012Map},
		{orb.Point{-80.33811, 40.88811}, "10", geolocator.House2012Map},
		{orb.Point{-76.88375, 40.26444}, "15", geolocator.Senate2022Map},
		{orb.Point{-76.88375, 40.26444}, "103", geolocator.House2022Map},
		{orb.Point{-80.38492, 40.96953}, "9", geolocator.House2022Map},
		{orb.Point{-80.33811, 40.88811}, "8", geolocator.House2022Map},
	}
	for i := 0; i < b.N; i++ {
		tc := &cases[i%len(cases)]
		d := tc.Map.District(tc.Point)
		if d.GetName() != tc.Name {
			b.Fatal(i, d.GetName(), tc.Name)
		}
	}
}
