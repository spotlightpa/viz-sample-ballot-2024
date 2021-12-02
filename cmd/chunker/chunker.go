package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/carlmjohnson/csv"
	"github.com/paulmach/orb/geojson"
)

func main() {
	csv := flag.String("csv", "", "")
	geojson := flag.String("geojson", "", "")
	format := flag.String("format", "", "")
	flag.Parse()
	if err := run(*csv, *geojson, *format); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(csvname, geojsonname, format string) error {
	m, err := readCSV(csvname)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(geojsonname)
	if err != nil {
		return err
	}

	fc, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		return err
	}
	for _, feat := range fc.Features {
		district := feat.Properties["District_1"].(string)
		props, ok := m[district]
		if !ok {
			return fmt.Errorf("missing district %q", district)
		}
		newfc := geojson.NewFeatureCollection()
		newfeat := geojson.NewFeature(feat.Geometry)
		for k, v := range props {
			newfeat.Properties[k] = v
		}
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

func readCSV(csvname string) (map[string]map[string]string, error) {
	csvf, err := os.Open(csvname)
	if err != nil {
		return nil, err
	}
	defer csvf.Close()

	m := map[string]map[string]string{}
	csvr := csv.NewFieldReader(csvf)
	for csvr.Scan() {
		row := csvr.Fields()
		m[row["district"]] = row
	}
	if err = csvr.Err(); err != nil {
		return nil, err
	}
	return m, nil
}
