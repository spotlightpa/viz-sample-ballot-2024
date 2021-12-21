package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/carlmjohnson/csv"
	"github.com/paulmach/orb/geojson"
)

func main() {
	districtCSV := flag.String("district-csv", "", "CSV `file` for district info")
	geojson := flag.String("geojson", "", "")
	dst := flag.String("dst", "", "destination file path")
	flag.Parse()
	if err := run(*districtCSV, *geojson, *dst); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(districtCSV, geojsonname, dst string) error {
	m, err := readCSV(districtCSV)
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
		district := feat.Properties["District"].(string)
		m[district]["per_asian"] = getProp(feat, "[% Adj_NH_Asn]")
		m[district]["per_black"] = getProp(feat, "[% Adj_NH_Blk]")
		m[district]["per_hispanic"] = getProp(feat, "[% Adj_Hispanic Origin]")
		m[district]["per_white"] = getProp(feat, "[% Adj_NH_Wht]")
		m[district]["per_mixed"] = getProp(feat, "[% Adj_NH_2+ Races]")
		other := feat.Properties["[% Adj_NH_Hwn]"].(float64) +
			feat.Properties["[% Adj_NH_Ind]"].(float64) +
			feat.Properties["[% Adj_NH_Oth]"].(float64)
		m[district]["per_misc"] = fmt.Sprintf("%.2f%%", other*100)
	}

	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	if err = os.WriteFile(dst, b, 0644); err != nil {
		return err
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

func getProp(feat *geojson.Feature, name string) string {
	f64, ok := feat.Properties[name].(float64)
	if !ok {
		panic(name)
	}
	return fmt.Sprintf("%.2f%%", f64*100)
}
