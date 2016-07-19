package pandemic

import (
	"math"
)

const EpidemicsPerGame = 5
const NumInfectionCards = 48

type GameState struct {
	CityDeck      CityDeck      `json:"city_deck"`
	InfectionDeck InfectionDeck `json:"infection_deck"`
	InfectionRate int           `json:"infection_rate"`
	Outbreaks     int
}

type InfectionDeck struct {
	Drawn []InfectionCard
	// add probability of card redraw
}

type InfectionCard struct {
	Name string
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
		return 2 * 1.0 / float64(c.cardsPerEpidemic()-(len(c.Drawn)%c.cardsPerEpidemic()))
	} else {
		return 0
	}
}

func (gs GameState) ProbabilityOfCity(cn string) float64 {
	// Has the city already been drawn?
	for _, card := range gs.InfectionDeck.Drawn {
		if card.Name == cn {
			return 0.0
		}
	}

	// How many cards are left?
	cardsRemaining := NumInfectionCards - len(gs.InfectionDeck.Drawn)

	// Probability of ANY of the infection cards being the City is equal to 1 minus the probabilty
	// that *none* of the cards is the city card
	//
	// P(C) ~= P(!C)^numCardsRemaining
	// Assuming 48 cards and infection rate = 4
	// P(C) = (47/48)*(46/47)*(45/46)*(44/45)
	probability := 1.0
	for i := 0; i < gs.InfectionRate; i++ {
		probability *= float64(cardsRemaining-1) / float64(cardsRemaining)
		cardsRemaining -= 1
	}
	return 1 - probability
}
