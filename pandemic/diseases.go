package pandemic

import (
	"encoding/json"
	"fmt"
)

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
