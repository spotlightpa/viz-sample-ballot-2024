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
	//go:embed embeds/senate-2001.geojson
	senate2001 []byte
	//go:embed embeds/senate-2012.geojson
	senate2012 []byte
)

var (
	House2012Map  = geojson2Map(house2012, "embeds/house-2012.gob")
	Senate2001Map = geojson2Map(senate2001, "embeds/senate-2001.gob")
	Senate2012Map = geojson2Map(senate2012, "embeds/senate-2012.gob")
)

func geojson2Map(b []byte, name string) Map {
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
			Name:    name,
			Polygon: poly,
			Bound:   bound,
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
