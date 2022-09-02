package geolocator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/carlmjohnson/resperr"
	"github.com/carlmjohnson/rootdown"
	"github.com/getsentry/sentry-go"
	"github.com/paulmach/orb"
	"github.com/rs/cors"
)

func (app *appEnv) routes() http.Handler {
	var mw rootdown.MiddlewareStack
	mw.Push(app.logRoute)
	if !app.isLambda() {
		mw.Push(cors.AllowAll().Handler)
	}
	var rr rootdown.Router
	rr.Get("/api/by-location", app.getByLocation, mw...)
	rr.Get("/api/by-address", app.getByAddress, mw...)
	rr.Get("/api/candidates-by-location", app.getCandidatesByLocation, mw...)
	rr.Get("/api/candidates-by-address", app.getCandidatesByAddress, mw...)

	return &rr
}

func (app *appEnv) logRoute(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		if r.URL.RawQuery != "" {
			q := r.URL.Query()
			q.Del("code")
			q.Del("state")
			url = url + "?" + q.Encode()
		}
		logger.Printf("[%s] %q - %s", r.Method, url, r.RemoteAddr)
		h.ServeHTTP(w, r)
	})
}

func (app *appEnv) replyJSON(statusCode int, w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	enc := json.NewEncoder(w)
	if err := enc.Encode(data); err != nil {
		app.logErr(r, fmt.Errorf("replyJSON problem: %v", err))
	}
}

func (app *appEnv) replyErr(w http.ResponseWriter, r *http.Request, err error) {
	app.logErr(r, err)
	code := resperr.StatusCode(err)
	msg := resperr.UserMessage(err)
	app.replyJSON(code, w, r, struct {
		Status       int    `json:"status"`
		ErrorMessage string `json:"error_message"`
	}{
		code,
		msg,
	})
}

func (app *appEnv) logErr(r *http.Request, err error) {
	ctx := r.Context()
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.CaptureException(err)
	} else {
		logger.Printf("sentry not in context")
	}

	logger.Printf("[%s] %q - err: %v", r.Method, r.URL.Path, err)
}

func (app *appEnv) getByLocation(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	lat, _ := strconv.ParseFloat(q.Get("lat"), 64)
	long, _ := strconv.ParseFloat(q.Get("long"), 64)

	w.Header().Set("Cache-Control", "public, max-age=3600, s-maxage=0")
	app.replyJSON(http.StatusOK, w, r, NewLocationInfo(lat, long))
}

func (app *appEnv) getByAddress(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		app.replyErr(w, r, resperr.New(http.StatusBadRequest, "no address"))
		return
	}
	var data GoogleMapsResults
	if err := app.googleMaps.
		Clone().
		Param("address", address).
		ToJSON(&data).
		Fetch(r.Context()); err != nil {
		err = resperr.WithStatusCode(err, http.StatusBadGateway)
		app.replyErr(w, r, err)
		return
	}
	if len(data.Results) < 1 {
		app.replyErr(w, r, resperr.New(
			http.StatusNotFound, "not found: %q", address))
		return
	}

	result := data.Results[0]
	long := result.Geometry.Location.Lng
	lat := result.Geometry.Location.Lat

	w.Header().Set("Cache-Control", "public, max-age=3600, s-maxage=0")
	app.replyJSON(http.StatusOK, w, r, struct {
		Address string  `json:"address"`
		Lat     float64 `json:"lat"`
		Long    float64 `json:"long"`
		LocationInfo
	}{
		result.FormattedAddress,
		lat,
		long,
		NewLocationInfo(lat, long),
	})
}

type LocationInfo struct {
	OldCongress string `json:"old_congress"`
	NewCongress string `json:"new_congress"`
	OldHouse    string `json:"old_house"`
	NewHouse    string `json:"new_house"`
	OldSenate   string `json:"old_senate"`
	NewSenate   string `json:"new_senate"`
}

func NewLocationInfo(lat, long float64) LocationInfo {
	p := orb.Point{long, lat}
	return LocationInfo{
		OldCongress: Congress2018Map.District(p).GetName(),
		NewCongress: Congress2022Map.District(p).GetName(),
		OldHouse:    House2012Map.District(p).GetName(),
		NewHouse:    House2022Map.District(p).GetName(),
		OldSenate:   Senate2012Map.District(p).GetName(),
		NewSenate:   Senate2022Map.District(p).GetName(),
	}
}

func (app *appEnv) getCandidatesByLocation(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	lat, _ := strconv.ParseFloat(q.Get("lat"), 64)
	long, _ := strconv.ParseFloat(q.Get("long"), 64)

	loc := NewLocationInfo(lat, long)
	data := NewCandiateInfo(loc)

	w.Header().Set("Cache-Control", "public, max-age=3600, s-maxage=0")
	app.replyJSON(http.StatusOK, w, r, data)
}

func (app *appEnv) getCandidatesByAddress(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		app.replyErr(w, r, resperr.New(http.StatusBadRequest, "no address"))
		return
	}
	var data GoogleMapsResults
	if err := app.googleMaps.
		Clone().
		Param("address", address).
		ToJSON(&data).
		Fetch(r.Context()); err != nil {
		err = resperr.WithStatusCode(err, http.StatusBadGateway)
		app.replyErr(w, r, err)
		return
	}
	if len(data.Results) < 1 {
		app.replyErr(w, r, resperr.New(
			http.StatusNotFound, "not found: %q", address))
		return
	}

	result := data.Results[0]
	long := result.Geometry.Location.Lng
	lat := result.Geometry.Location.Lat
	loc := NewLocationInfo(lat, long)
	can := NewCandiateInfo(loc)

	w.Header().Set("Cache-Control", "public, max-age=3600, s-maxage=0")
	app.replyJSON(http.StatusOK, w, r, struct {
		Address string  `json:"address"`
		Lat     float64 `json:"lat"`
		Long    float64 `json:"long"`
		CandidateInfo
	}{
		result.FormattedAddress,
		lat,
		long,
		can,
	})
}

type CandidateInfo struct {
	LocationInfo
	Governor    []Candidate `json:"governor"`
	USSenate    []Candidate `json:"us_senate"`
	USHouse     []Candidate `json:"us_house"`
	StateSenate []Candidate `json:"state_senate"`
	StateHouse  []Candidate `json:"state_house"`
}

func NewCandiateInfo(loc LocationInfo) CandidateInfo {
	return CandidateInfo{
		LocationInfo: loc,
		Governor:     CanGov,
		USSenate:     CanUSSenate,
		USHouse:      CanUSHouse[loc.NewCongress],
		StateSenate:  CanPASenate[loc.NewSenate],
		StateHouse:   CanPAHouse[loc.NewHouse],
	}
}
