package pandemic

import (
	"math"
)

const EpidemicsPerGame = 5

type GameState struct {
	CityDeck      CityDeck      `json:"city_deck"`
	InfectionDeck InfectionDeck `json:"infection_deck"`
}

type InfectionDeck struct {
	// TODO add cards drawn
	// add probability of card redraw
}

type InfectionCard struct {
	Name string `json:"name"`
}

type CityDeck struct {
	Drawn []CityCard `json:"drawn"`
	Total int        `json:"total"`
}

type CityCard struct {
	City       City `json:"city"`
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
		return 2 * 1.0 / float64(c.cardsPerEpidemic()-(len(c.Drawn)%c.cardsPerEpidemic()))
	} else {
		return 0
	}
}
