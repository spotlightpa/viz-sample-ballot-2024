package geolocator

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/planar"
)

type District struct {
	Name string
	orb.MultiPolygon
	orb.Bound
	NewStyle bool
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
		if pointInMultiPoly(p, d.Bound, d.MultiPolygon, d.NewStyle) {
			return &d
		}
	}
	return nil
}

func pointInMultiPoly(p orb.Point, bound orb.Bound, mgon orb.MultiPolygon, newstyle bool) bool {
	if !bound.Contains(p) {
		return false
	}
	if newstyle {
		return planar.MultiPolygonContains(mgon, p)
	}
	poly := mgon[0]
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
