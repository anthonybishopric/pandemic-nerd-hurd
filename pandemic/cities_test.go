package pandemic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func TestProbabilityOfCard(t *testing.T) {
}
