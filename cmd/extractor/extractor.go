package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

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
		district := feat.Properties["DISTRICT"].(string)
		m[district]["per_asian"] = getPct(feat, "ADJ_NH_ASN")
		m[district]["per_black"] = getPct(feat, "ADJ_NH_BLK")
		m[district]["per_hispanic"] = getPct(feat, "ADJ_HISPAN")
		m[district]["per_white"] = getPct(feat, "ADJ_NH_WHT")
		m[district]["per_mixed"] = getPct(feat, "ADJ_NH_2_R")
		other := getProp(feat, "ADJ_NH_HWN") + getProp(feat, "ADJ_NH_IND") + getProp(feat, "ADJ_NH_OTH")
		other /= getProp(feat, "ADJ_POPULA")
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
		m[row["district"]] = map[string]string{
			"district":  row["district"],
			"per_dem":   row["per_dem"],
			"per_other": row["per_other"],
			"per_rep":   row["per_rep"],
		}
	}
	if err = csvr.Err(); err != nil {
		return nil, err
	}
	return m, nil
}

func getPct(feat *geojson.Feature, name string) string {
	popstr := feat.Properties["ADJ_POPULA"].(string)
	pop, err := strconv.ParseFloat(popstr, 64)
	if err != nil {
		panic(name)
	}
	str, ok := feat.Properties[name].(string)
	if !ok {
		panic(name)
	}
	f64, err := strconv.ParseFloat(str, 64)
	if err != nil {
		panic(name)
	}
	return fmt.Sprintf("%.2f%%", f64/pop*100)
}

func getProp(feat *geojson.Feature, name string) float64 {
	s := feat.Properties[name].(string)
	f64, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(name)
	}
	return f64
}
