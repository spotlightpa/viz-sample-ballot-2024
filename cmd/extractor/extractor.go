package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/carlmjohnson/csv"
)

func main() {
	oldHouse := flag.String("old-house", "", "CSV `file` for old house district info")
	newHouse := flag.String("new-house", "", "CSV `file` for new house district info")
	oldSenate := flag.String("old-senate", "", "CSV `file` for old senate district info")
	newSenate := flag.String("new-senate", "", "CSV `file` for new senate district info")
	format := flag.String("format", "", "destination file path")
	flag.Parse()
	if err := run(*oldHouse, *newHouse, *oldSenate, *newSenate, *format); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(oldHouse, newHouse, oldSenate, newSenate, format string) error {
	for _, csvname := range []string{
		oldHouse, newHouse, oldSenate, newSenate,
	} {
		m, err := readCSV(csvname)
		if err != nil {
			return err
		}
		b, err := json.MarshalIndent(m, "", "  ")
		if err != nil {
			return err
		}
		name := path.Base(csvname)
		name = strings.TrimSuffix(name, path.Ext(name))
		name = fmt.Sprintf(format, name)
		if err = os.WriteFile(name, b, 0644); err != nil {
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
