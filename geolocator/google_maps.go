package geolocator

import (
	"encoding/json"
	"fmt"

	"github.com/carlmjohnson/requests"
)

func NewMapsClient() *requests.Builder {
	return requests.
		URL("https://maps.googleapis.com/maps/api/geocode/json").
		// Limit to PA
		Param("components", "administrative_area:PA|country:US")
}

type GoogleMapsResults struct {
	Results []Results `json:"results"`
	Status  string    `json:"status"`
}

func (r *GoogleMapsResults) UnmarshalJSON(data []byte) error {
	type nomethods GoogleMapsResults
	if err := json.Unmarshal(data, (*nomethods)(r)); err != nil {
		return err
	}
	if r.Status != "OK" && r.Status != "ZERO_RESULTS" {
		return fmt.Errorf("got status %q from Google Maps", r.Status)
	}
	return nil
}

type AddressComponents struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}
type Northeast struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
type Southwest struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
type Bounds struct {
	Northeast Northeast `json:"northeast"`
	Southwest Southwest `json:"southwest"`
}
type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
type Viewport struct {
	Northeast Northeast `json:"northeast"`
	Southwest Southwest `json:"southwest"`
}
type Geometry struct {
	Bounds       Bounds   `json:"bounds"`
	Location     Location `json:"location"`
	LocationType string   `json:"location_type"`
	Viewport     Viewport `json:"viewport"`
}
type Results struct {
	AddressComponents []AddressComponents `json:"address_components"`
	FormattedAddress  string              `json:"formatted_address"`
	Geometry          Geometry            `json:"geometry"`
	PlaceID           string              `json:"place_id"`
	Types             []string            `json:"types"`
}
