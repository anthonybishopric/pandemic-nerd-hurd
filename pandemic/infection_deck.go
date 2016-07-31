package pandemic

import (
	"fmt"
	"strings"
)

type InfectionDeck struct {
	Drawn      Set
	Striations []Set // all Striations still present on the infection deck. the 0th is the top
}

type InfectionCard struct {
	Name string
}

func NewInfectionDeck(cities []string) *InfectionDeck {
	firstStriation := Set{}
	for _, city := range cities {
		firstStriation.Add(strings.ToLower(city))
	}
	return &InfectionDeck{
		Drawn:      Set{},
		Striations: []Set{firstStriation},
	}
}

func (d *InfectionDeck) assertStriationCount() {
	if len(d.Striations) < 1 {
		panic("Unexpectedly didn't have any Striations - was this game set up correctly?")
	}
}

func (d *InfectionDeck) Draw(cityName string) error {
	cityName = strings.ToLower(cityName)
	d.assertStriationCount()
	for d.Striations[0].Size() == 0 {
		d.Striations = d.Striations[1:]
	}
	d.assertStriationCount()
	if _, ok := d.Striations[0].Remove(cityName); !ok {
		return fmt.Errorf("Card %v is not present in the active striation - how the fuck did you draw this card?", cityName)
	}
	d.Drawn.Add(cityName)
	return nil
}

func (d *InfectionDeck) PullFromBottom(card string) error {
	d.assertStriationCount()
	bottomStriation := d.Striations[len(d.Striations)-1]
	if _, ok := bottomStriation.Remove(card); !ok {
		return fmt.Errorf("Card %v should not be present in the bottom striation", card)
	}
	d.Drawn.Add(card)
	return nil
}

// We just prepend the currently drawn pile onto the front
// of our deck Striations. Then we reset drawn.
func (d *InfectionDeck) ShuffleDrawn() {
	d.Striations = append([]Set{d.Drawn}, d.Striations...)
	d.Drawn = Set{}
}

func (d *InfectionDeck) CurrentStriationCount() int {
	return d.Striations[0].Size()
}

func (d *InfectionDeck) DrawnCount() int {
	return d.Drawn.Size()
}

func (d *InfectionDeck) ProbabilityOfDrawing(city string, infectionRate int) float64 {
	// Has the city already been drawn?
	if d.Drawn.Contains(city) {
		return 0.0
	}

	// Clone myself so we can recurse into the future. <- coolest code comment I've ever left.
	dCopy := *d

	// Probability of ANY of the infection cards being the City is equal to 1 minus the probabilty
	// that *none* of the cards is the city card
	//
	// P(C) ~= 1 - P(!C)^numCardsRemaining
	// Assuming 10 cards in the striation and infection rate = 4
	// P(C) = 1 - (9/10)*(8/9)*(7/8)*(6/7) = 1 - 6/10 = 40%
	probability := 1.0
	curStriationSize := dCopy.Striations[0].Size()
	for draw := 0; draw < infectionRate; draw++ {
		// if we've run out of cards in this striation, pop and
		// start using the next striation down.
		for curStriationSize == 0 {
			dCopy.Striations = dCopy.Striations[1:]
			dCopy.assertStriationCount()
			curStriationSize = dCopy.Striations[0].Size()
		}

		// calculate the probability of drawing the given card
		// based on it's presence in the current striation. If
		// it is not present, the chance of not drawing the
		// target card is 100%.
		if dCopy.Striations[0].Contains(city) {
			probability *= float64(curStriationSize-1) / float64(curStriationSize)
		}
		// Reduce the count of cards we know will remain in this
		// striation after drawing a card
		curStriationSize = curStriationSize - 1
	}

	return 1 - probability
}
