package pandemic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"path/filepath"
	"testing"
)

func TestLoadFromJSON(t *testing.T) {

	// Do IO
	filename, _ := filepath.Abs("../data/cities.json")
	data, _ := ioutil.ReadFile(filename)
	c := Cities{}

	// Decode JSON
	dec := json.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&c); err != nil {
		fmt.Errorf("Decode city: %v", err)
	}

}

func TestSimpleGame(t *testing.T) {
	// four possible scenarios
	// [2,1,1,1], [1,2,1,1], [1,1,2,1] and [1,1,1,2]
	model := generateProbabilityModel(5, 4)

	// 1/4*0.5 + 3/4*1 == .875
	prob := model.EpidemicProbabilityAt(1)
	if math.Floor(prob*1000)/1000 != 0.875 {
		t.Fatalf("Expected 0.875 probability of epidemic, got %v", prob)
	}

	// this invalidates all of the the [1,*] scenarios. The 2nd card must now
	// be an epidemic
	model.DrawCity(0)
	prob = model.EpidemicProbabilityAt(1)
	if prob != 1.0 {
		t.Fatalf("Expected 100%% chance of epidemic, got %v", prob)
	}
}
