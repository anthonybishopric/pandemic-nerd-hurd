package pandemic

import (
	"fmt"
)

type City struct {
	Name          string      `json:"name"`
	Disease       DiseaseType `json:"disease"`
	PanicLevel    PanicLevel  `json:"panic_level"`
	Neighbors     []string    `json:"neighbors"`
	NumInfections int         `json:"num_infections"`
}

type Cities struct {
	Cities []*City `json:"cities"`
}

func (c *Cities) GetCity(city string) (*City, error) {
	for _, c := range c.Cities {
		if c.Name == city {
			return c, nil
		}
	}
	return nil, fmt.Errorf("No city named %v", city)
}

func (c Cities) WithDisease(disease DiseaseType) []*City {
	cities := []*City{}
	for _, city := range c.Cities {
		if city.Disease == disease {
			cities = append(cities, city)
		}
	}
	return cities
}

func (c Cities) CityNames() []string {
	names := []string{}
	for _, city := range c.Cities {
		names = append(names, city.Name)
	}
	return names
}

func (c *City) Infect() bool {
	if c.NumInfections == 3 {
		return true
	}
	c.NumInfections++
	return false
}

func (c *City) Epidemic() {
	c.NumInfections = 3
}

func (c *City) SetInfections(infections int) {
	c.NumInfections = infections
}
