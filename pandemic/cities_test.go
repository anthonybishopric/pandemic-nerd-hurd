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

func TestSortByInfect(t *testing.T) {
	cities := Cities{
		Cities: []*City{
			{
				Name:          "a",
				NumInfections: 2,
			},
			{
				Name:          "b",
				NumInfections: 3,
			},
			{
				Name:          "c",
				NumInfections: 1,
			},
		},
	}
	sorted := cities.SortByInfectionLevel([]CityName{"a", "b", "c"})
	if len(sorted) != 3 || sorted[0] != "b" || sorted[1] != "a" || sorted[2] != "c" {
		t.Fatalf("Incorrect order: %+v", sorted)
	}
}
