package pandemic

import (
	"encoding/json"
	"fmt"
)

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

func PanicLevelFromString(s string) (PanicLevel, error) {
	got, ok := map[string]PanicLevel{
		"Nothing":    Nothing,
		"Unstable":   Unstable,
		"Rioting2":   Rioting2,
		"Rioting3":   Rioting3,
		"Collapsing": Collapsing,
		"Fallen":     Fallen,
	}[s]
	if !ok {
		return PanicLevel(-1), fmt.Errorf("invalid PanicLevel %q", s)
	}
	return got, nil
}

func (pl PanicLevel) String() string {
	got, ok := map[PanicLevel]string{
		Nothing:    "Nothing",
		Unstable:   "Unstable",
		Rioting2:   "Rioting2",
		Rioting3:   "Rioting3",
		Collapsing: "Collapsing",
		Fallen:     "Fallen",
	}[pl]
	if !ok {
		return ""
	}
	return got
}

func (pl *PanicLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal(pl.String())
}

func (pl *PanicLevel) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("panic level should be a string, got %s", data)
	}

	got, err := PanicLevelFromString(s)
	if err != nil {
		return err
	}
	*pl = got
	return nil
}
