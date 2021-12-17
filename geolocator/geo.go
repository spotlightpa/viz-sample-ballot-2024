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
		if d.Contains(p) {
			return &d
		}
	}
	return nil
}

func (d *District) Contains(p orb.Point) bool {
	if !d.Bound.Contains(p) {
		return false
	}
	if d.NewStyle {
		return planar.MultiPolygonContains(d.MultiPolygon, p)
	}
	poly := d.MultiPolygon[0]
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
