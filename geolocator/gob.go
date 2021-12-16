//go:build !writegob

package geolocator

import (
	"bytes"
	_ "embed"
	"encoding/gob"
)

var (
	//go:embed embeds/house-2012.gob
	house2012 []byte
	//go:embed embeds/house-2021.gob
	house2021 []byte
	//go:embed embeds/senate-2012.gob
	senate2012 []byte
	//go:embed embeds/senate-2021.gob
	senate2021 []byte
)

var (
	House2012Map  = gob2Map(house2012)
	House2021Map  = gob2Map(house2021)
	Senate2012Map = gob2Map(senate2012)
	Senate2021Map = gob2Map(senate2021)
)

func gob2Map(b []byte) Map {
	r := bytes.NewReader(b)
	dec := gob.NewDecoder(r)
	var ds Map
	if err := dec.Decode(&ds); err != nil {
		panic(err)
	}
	return ds
}
