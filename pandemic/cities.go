package pandemic

import (
	"encoding/json"
	"fmt"
)

/* Panic Level */
type PanicLevel int

func (p PanicLevel) CanBuildResearchStations() bool {
	return int(p) < 2
}

const (
	Nothing = PanicLevel(iota)
	Unstable
	Rioting2
	Rioting3
	Collapsing
	Fallen
)

func (pl *PanicLevel) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("panic level should be a string, got %s", data)
	}

	got, ok := map[string]PanicLevel{
		"Nothing":    Nothing,
		"Unstable":   Unstable,
		"Rioting2":   Rioting2,
		"Rioting3":   Rioting3,
		"Collapsing": Collapsing,
		"Fallen":     Fallen,
	}[s]
	if !ok {
		return fmt.Errorf("invalid PanicLevel %q", s)
	}
	*pl = got
	return nil
}

/* Disease Type */
type DiseaseType struct {
	Color         string `json:"color"`
	Incurable     bool   `json:"incurable,omitempty"`
	Untreatable   bool   `json:"untreatable,omitempty"`
	BecomingFaded bool   `json:"becoming_faded,omitempty"`
	HasZombies    bool
}

func (dt *DiseaseType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Disease Type should be a string, got %s", data)
	}

	got, ok := map[string]DiseaseType{
		"Yellow": Yellow,
		"Blue":   Blue,
		"Red":    Red,
		"Black":  Black,
		"Faded":  Faded,
	}[s]
	if !ok {
		return fmt.Errorf("invalid DiseaseType %q", s)
	}
	*dt = got
	return nil
}

var Yellow = DiseaseType{
	Color: "Yellow",
}
var Blue = DiseaseType{
	Color:         "Blue",
	Incurable:     true, // TODO: make configurable with a gamestate
	Untreatable:   true,
	BecomingFaded: true,
}
var Red = DiseaseType{
	Color: "Red",
}
var Black = DiseaseType{
	Color: "Black",
}
var Faded = DiseaseType{
	Color:       "Faded",
	HasZombies:  true,
	Incurable:   true,
	Untreatable: true,
}

/* City + Cities */
type City struct {
	Name        string
	Epidemic    bool
	FundedEvent bool
	Disease     DiseaseType
	PanicLevel  PanicLevel
	Neighbors   []string
}

type Cities struct {
	Cities []City
}

func AllCitiesWithDisease(Cities []City, disease DiseaseType) []City {
	cities := []City{}
	for _, city := range Cities {
		if city.Disease == disease {
			cities = append(cities, city)
		}
	}
	return cities
}
