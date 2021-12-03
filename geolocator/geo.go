package geolocator

import (
	_ "embed"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
)

var (
	//go:embed embeds/house-2012.geojson
	house2012 []byte
	//go:embed embeds/senate-2001.geojson
	senate2001 []byte
	//go:embed embeds/senate-2012.geojson
	senate2012 []byte
)

type District struct {
	name string
	orb.Polygon
	orb.Bound
}

func (d *District) Name() string {
	if d == nil {
		return ""
	}
	return d.name
}

type Map []District

var (
	House2012Map  = geojson2Map(house2012)
	Senate2001Map = geojson2Map(senate2001)
	Senate2012Map = geojson2Map(senate2012)
)

func geojson2Map(b []byte) Map {
	fc, err := geojson.UnmarshalFeatureCollection(b)
	if err != nil {
		panic(err)
	}

	ds := make(Map, len(fc.Features))
	for i, f := range fc.Features {
		name := f.Properties["District_1"].(string)
		poly := f.Geometry.(orb.Polygon)
		if len(poly) < 1 {
			panic(name)
		}
		bound := poly[0].Bound()
		for _, ring := range poly[1:] {
			bound = bound.Union(ring.Bound())
		}
		ds[i] = District{
			name:    name,
			Polygon: poly,
			Bound:   bound,
		}
	}
	return ds
}

func (m Map) District(p orb.Point) *District {
	for _, d := range m {
		if pointInPoly(p, d.Bound, d.Polygon) {
			return &d
		}
	}
	return nil
}

func pointInPoly(p orb.Point, bound orb.Bound, poly orb.Polygon) bool {
	if !bound.Contains(p) {
		return false
	}
	contained := false
	for _, ring := range poly {
		if planar.RingContains(ring, p) {
			if ring.Orientation() == orb.CCW {
				return false
			}
			contained = true
		}
	}
	return contained
}
