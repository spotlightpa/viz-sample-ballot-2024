package geolocator

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/planar"
)

type District struct {
	Name string
	orb.Polygon
	orb.Bound
}

func (d *District) GetName() string {
	if d == nil {
		return ""
	}
	return d.Name
}

type Map []District

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
