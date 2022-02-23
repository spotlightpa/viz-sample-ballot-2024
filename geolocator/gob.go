//go:build !writegob

package geolocator

import (
	"bytes"
	_ "embed"
	"encoding/gob"
)

var (
	//go:embed embeds/congress-2018.gob
	congress2018 []byte
	//go:embed embeds/congress-2022.gob
	congress2022 []byte
	//go:embed embeds/house-2012.gob
	house2012 []byte
	//go:embed embeds/house-2022.gob
	house2022 []byte
	//go:embed embeds/senate-2012.gob
	senate2012 []byte
	//go:embed embeds/senate-2022.gob
	senate2022 []byte
)

var (
	Congress2018Map = gob2Map(congress2018)
	Congress2022Map = gob2Map(congress2022)
	House2012Map    = gob2Map(house2012)
	House2022Map    = gob2Map(house2022)
	Senate2012Map   = gob2Map(senate2012)
	Senate2022Map   = gob2Map(senate2022)
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
