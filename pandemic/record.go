package pandemic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"strings"
)

const EpidemicsPerGame = 5
const NumInfectionCards = 48

type GameState struct {
	Cities        *Cities        `json:"cities"`
	CityDeck      *CityDeck      `json:"city_deck"`
	DiseaseData   []DiseaseData  `json:"disease_data"`
	InfectionDeck *InfectionDeck `json:"infection_deck"`
	InfectionRate int            `json:"infection_rate"`
	Outbreaks     int            `json:"outbreaks"`
}

func NewGame(citiesFile string) (*GameState, error) {
	var cities Cities
	data, err := ioutil.ReadFile(citiesFile)
	if err != nil {
		return nil, fmt.Errorf("Could not read cities file at %v: %v", citiesFile, err)
	}
	err = json.Unmarshal(data, &cities)
	if err != nil {
		return nil, fmt.Errorf("Invalid cities JSON file at %v: %v", citiesFile, err)
	}
	cityDeck := CityDeck{}
	cityDeck.Total = len(cities.Cities)

	infectionDeck := NewInfectionDeck(cities.CityNames())
	return &GameState{
		Cities:        &cities,
		DiseaseData:   []DiseaseData{Yellow, Red, Black, Blue, Faded},
		CityDeck:      &cityDeck,
		InfectionDeck: infectionDeck,
		InfectionRate: 2,
		Outbreaks:     0,
	}, nil
}

type CityDeck struct {
	Drawn []CityCard
	Total int
}

type CityCard struct {
	City       City
	IsEpidemic bool `json:"is_epidemic"`
}

func (c CityDeck) cardsPerEpidemic() int {
	return c.Total / EpidemicsPerGame
}

func (c CityDeck) EpidemicsDrawn() int {
	count := 0
	for _, card := range c.Drawn {
		if card.IsEpidemic {
			count++
		}
	}
	return count
}

// 100 city cards, 5 epidemics
// probability of drawing an epidemic on turn 0:
//   1/20 + (1/19 * 19/20)
//
//   1/18 + (1/17 * 17/18)
//
// if an epidemic is drawn, then the probability
// of an epidemic being drawn is 0 until the 10th turn.
//
// if no epidemic is drawn, the probability of drawing
// an epidemic in the 10th turn is 1/2 + (1/1 * 1/2).
//
// on the 11th turn, the probability of the 21st and 22nd cards
// being epidemics is
// 1/20 + (1/19 * 19/20)
//
func (c CityDeck) probabilityOfEpidemic() float64 {
	currentPhase := int(math.Floor(float64(len(c.Drawn))/float64(c.cardsPerEpidemic())) + 0.5)
	if currentPhase == c.EpidemicsDrawn() {
		return 2.0 / float64(c.cardsPerEpidemic()-(len(c.Drawn)%c.cardsPerEpidemic()))
	} else {
		return 0
	}
}

func (gs GameState) ProbabilityOfCity(cn string) float64 {
	return gs.InfectionDeck.ProbabilityOfDrawing(cn, gs.InfectionRate)
}

func (gs *GameState) GetCity(city string) (*City, error) {
	city = strings.ToLower(city)
	return gs.Cities.GetCity(city)
}

func (gs *GameState) GetDiseaseData(diseaseType DiseaseType) (*DiseaseData, error) {
	for _, data := range gs.DiseaseData {
		if data.Type == diseaseType {
			return &data, nil
		}
	}
	return nil, fmt.Errorf("No disease identified by %v", diseaseType)
}
