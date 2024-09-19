package main

import (
	"os"

	"github.com/carlmjohnson/exitcode"
	"github.com/spotlightpa/viz-sample-ballot-2024/geolocator"
)

func main() {
	exitcode.Exit(geolocator.CLI(os.Args[1:]))
}
