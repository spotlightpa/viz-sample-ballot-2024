// Package geolocator finds districts
package geolocator

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/carlmjohnson/flagext"
	"github.com/carlmjohnson/gateway"
	"github.com/carlmjohnson/requests"
	"github.com/carlmjohnson/versioninfo"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/spotlightpa/viz-redistricting-2020/build"
)

const App = "geolocator"

var logger = log.Default()

func CLI(args []string) error {
	var app appEnv
	err := app.ParseArgs(args)
	if err != nil {
		return err
	}
	err = app.Exec()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	return err
}

func (app *appEnv) ParseArgs(args []string) error {
	fl := flag.NewFlagSet(App, flag.ContinueOnError)
	fl.Usage = func() {
		fmt.Fprintf(fl.Output(), "%s - %s\n\n", App, versioninfo.Version)
		fl.PrintDefaults()
	}
	fl.IntVar(&app.port, "port", -1, "specify a port to use http rather than AWS Lambda")
	sentryDSN := fl.String("sentry-dsn", "", "DSN `pseudo-URL` for Sentry")

	fl.Func("api-key", "Google Maps API `key`", func(s string) error {
		app.googleMaps = requests.
			URL("https://maps.googleapis.com/maps/api/geocode/json").
			Param("key", s)
		return nil
	})

	if err := fl.Parse(args); err != nil {
		return err
	}
	if err := flagext.ParseEnv(fl, App); err != nil {
		return err
	}
	if err := app.initSentry(*sentryDSN); err != nil {
		return err
	}
	logger.SetPrefix(App + " ")
	logger.SetFlags(log.LstdFlags | log.Lshortfile)
	return nil
}

type appEnv struct {
	port       int
	googleMaps *requests.Builder
}

func (app *appEnv) Exec() (err error) {
	listener := gateway.ListenAndServe
	var portStr string
	if app.isLambda() {
		portStr = build.URL.Hostname()
	} else {
		portStr = fmt.Sprintf(":%d", app.port)
		build.URL.Host += portStr
		listener = http.ListenAndServe
	}
	routes := sentryhttp.
		New(sentryhttp.Options{
			WaitForDelivery: true,
			Timeout:         5 * time.Second,
			Repanic:         !app.isLambda(),
		}).
		Handle(app.routes())

	logger.Printf("starting on %s", portStr)
	return listener(portStr, routes)
}

func (app *appEnv) initSentry(dsn string) error {
	var transport sentry.Transport
	if app.isLambda() {
		logger.Printf("setting sentry sync with timeout")
		transport = &sentry.HTTPSyncTransport{Timeout: 5 * time.Second}
	}
	if dsn == "" {
		logger.Printf("no Sentry DSN")
		return nil
	}
	return sentry.Init(sentry.ClientOptions{
		Dsn:       dsn,
		Release:   build.Rev,
		Transport: transport,
	})
}

func (app *appEnv) isLambda() bool {
	return app.port == -1
}
