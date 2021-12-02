package main

import (
	"os"

	"github.com/carlmjohnson/exitcode"
	"github.com/spotlightpa/viz-redistricting-2020/geolocator"
)

func main() {
	exitcode.Exit(geolocator.CLI(os.Args[1:]))
}
