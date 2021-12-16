//go:build writegob

package geolocator

import (
	"bytes"
	_ "embed"
	"encoding/gob"
	"os"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

var (
	//go:embed embeds/house-2012.geojson
	house2012 []byte
	//go:embed embeds/senate-2012.geojson
	senate2012 []byte
	//go:embed embeds/house-2021.geojson
	house2021 []byte
	//go:embed embeds/senate-2021.geojson
	senate2021 []byte
)

var (
	House2012Map  = geojson2Map(house2012, "embeds/house-2012.gob", false)
	House2021Map  = geojson2Map(house2021, "embeds/house-2021.gob", true)
	Senate2012Map = geojson2Map(senate2012, "embeds/senate-2012.gob", false)
	Senate2021Map = geojson2Map(senate2021, "embeds/senate-2021.gob", true)
)

func geojson2Map(b []byte, name string, newstyle bool) Map {
	fc, err := geojson.UnmarshalFeatureCollection(b)
	if err != nil {
		panic(err)
	}

	ds := make(Map, len(fc.Features))
	for i, f := range fc.Features {
		propname := "District_1"
		if newstyle {
			propname = "District"
		}
		dist := f.Properties[propname].(string)

		mgon, ok := f.Geometry.(orb.MultiPolygon)
		if !ok {
			poly := f.Geometry.(orb.Polygon)
			mgon = []orb.Polygon{poly}
		}
		if len(mgon[0]) < 1 {
			panic(name + "-" + dist)
		}
		bound := mgon[0][0].Bound()
		for _, poly := range mgon {
			for _, ring := range poly {
				bound = bound.Union(ring.Bound())
			}
		}
		ds[i] = District{
			Name:         dist,
			MultiPolygon: mgon,
			Bound:        bound,
			NewStyle:     newstyle,
		}
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err = enc.Encode(ds); err != nil {
		panic(err)
	}
	if err = os.WriteFile(name, buf.Bytes(), 0644); err != nil {
		panic(err)
	}
	return ds
}
