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
