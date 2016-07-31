package pandemic

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

func (c Cities) CityNames() []string {
	names := []string{}
	for _, city := range c.Cities {
		names = append(names, city.Name)
	}
	return names
}
