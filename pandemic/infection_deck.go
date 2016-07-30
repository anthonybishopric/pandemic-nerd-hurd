package pandemic

import (
	"fmt"
)

type InfectionDeck struct {
	Drawn      Set
	striations []Set // all striations still present on the infection deck. the 0th is the top
}

type InfectionCard struct {
	Name string
}

func NewInfectionDeck(cities []string) *InfectionDeck {
	firstStriation := Set{}
	for _, city := range cities {
		firstStriation.Add(city)
	}
	return &InfectionDeck{
		Drawn:      Set{},
		striations: []Set{firstStriation},
	}
}

func (d *InfectionDeck) assertStriationCount() {
	if len(d.striations) < 1 {
		panic("Unexpectedly didn't have any striations - was this game set up correctly?")
	}
}

func (d *InfectionDeck) Draw(card string) error {
	d.assertStriationCount()
	for d.striations[0].Size() == 0 {
		d.striations = d.striations[1:]
	}
	d.assertStriationCount()
	if _, ok := d.striations[0].Remove(card); !ok {
		return fmt.Errorf("Card %v is not present in the active striation - how the fuck did you draw this card?", card)
	}
	d.Drawn.Add(card)
	return nil
}

func (d *InfectionDeck) PullFromBottom(card string) error {
	d.assertStriationCount()
	bottomStriation := d.striations[len(d.striations)-1]
	if _, ok := bottomStriation.Remove(card); !ok {
		return fmt.Errorf("Card %v should not be present in the bottom striation", card)
	}
	d.Drawn.Add(card)
	return nil
}

// We just prepend the currently drawn pile onto the front
// of our deck striations. Then we reset drawn.
func (d *InfectionDeck) ShuffleDrawn() {
	d.striations = append([]Set{d.Drawn}, d.striations...)
	d.Drawn = Set{}
}

func (d *InfectionDeck) CurrentStriationCount() int {
	return d.striations[0].Size()
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
	curStriationSize := dCopy.striations[0].Size()
	for draw := 0; draw < infectionRate; draw++ {
		// if we've run out of cards in this striation, pop and
		// start using the next striation down.
		for curStriationSize == 0 {
			dCopy.striations = dCopy.striations[1:]
			dCopy.assertStriationCount()
			curStriationSize = dCopy.striations[0].Size()
		}

		// calculate the probability of drawing the given card
		// based on it's presence in the current striation. If
		// it is not present, the chance of not drawing the
		// target card is 100%.
		if dCopy.striations[0].Contains(city) {
			probability *= float64(curStriationSize-1) / float64(curStriationSize)
		}
		// Reduce the count of cards we know will remain in this
		// striation after drawing a card
		curStriationSize = curStriationSize - 1
	}

	return 1 - probability
}
