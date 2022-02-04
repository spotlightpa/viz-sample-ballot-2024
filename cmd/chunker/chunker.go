package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/paulmach/orb/geojson"
)

func main() {
	geojson := flag.String("geojson", "", "")
	format := flag.String("format", "", "")
	flag.Parse()
	if err := run(*geojson, *format); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(geojsonname, format string) error {
	data, err := os.ReadFile(geojsonname)
	if err != nil {
		return err
	}

	fc, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		return err
	}
	for _, feat := range fc.Features {
		district := feat.Properties["DISTRICT"].(string)
		newfc := geojson.NewFeatureCollection()
		newfeat := geojson.NewFeature(feat.Geometry)
		newfeat.Properties["district"] = district
		newfc.Append(newfeat)
		data, err = newfc.MarshalJSON()
		if err != nil {
			return err
		}
		fname := fmt.Sprintf(format, district)
		if err = os.WriteFile(fname, data, 0644); err != nil {
			return err
		}
	}

	return nil
}
