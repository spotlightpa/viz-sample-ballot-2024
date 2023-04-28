package build

import (
	_ "embed"
	"net/url"
	"strings"
)

//go:embed url.txt
var embedurl string

var URL = func() url.URL {
	u, err := url.Parse(strings.TrimSpace(embedurl))
	if err != nil {
		panic(err)
	}
	return *u
}()
