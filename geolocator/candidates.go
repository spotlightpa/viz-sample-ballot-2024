package geolocator

import (
	_ "embed"
	"encoding/json"
)

type Candidate struct {
	Name  string `json:"name"`
	Party string `json:"party"`
}

var (
	CanAttorneyGeneral []Candidate
	CanAuditorGeneral  []Candidate
	CanPresident       []Candidate
	CanStateTreasurer  []Candidate
	CanUSSenate        []Candidate
	CanUSHouse         map[string][]Candidate
	CanPAHouse         map[string][]Candidate
	CanPASenate        map[string][]Candidate
)

//go:embed embeds/candidates-2024.json
var canData []byte

func init() {
	var candidates map[string]map[string][]Candidate
	err := json.Unmarshal(canData, &candidates)
	if err != nil {
		panic(err)
	}
	CanAttorneyGeneral = candidates["ATTORNEY GENERAL"][""]
	CanAuditorGeneral = candidates["AUDITOR GENERAL"][""]
	CanPresident = candidates["PRESIDENT OF THE UNITED STATES"][""]
	CanStateTreasurer = candidates["STATE TREASURER"][""]
	CanUSSenate = candidates["US Senate"][""]
	CanUSHouse = candidates["US House"]
	CanPAHouse = candidates["State House"]
	CanPASenate = candidates["State Senate"]
}
