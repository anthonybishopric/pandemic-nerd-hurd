package pandemic

import (
	"fmt"
	"math"
	"testing"
)

func getNumCards(count int, numEpis int) []CityCard {
	cards := make([]CityCard, count)
	for x := 0; x < count-numEpis; x++ {
		cards[x] = CityCard{City{Name: CityName(fmt.Sprintf("testCity%v", x))}, false}
	}
	for x := 0; x < numEpis; x++ {
		cards[x] = CityCard{City{}, true}
	}
	return cards
}

func TestCardProbabilities(t *testing.T) {
	deck := CityDeck{
		All:   getNumCards(100, EpidemicsPerGame),
		Drawn: []CityCard{},
	}
	if prob := deck.probabilityOfEpidemic(); prob != 0.1 {
		t.Fatalf("Should have had a 10%% chance of epidemic, got %v", prob)
	}
	deck.Drawn = []CityCard{
		{
			City{Name: "SanFrancisco"},
			false,
		},
		{
			City{Name: "Sydney"},
			false,
		},
	}
	if prob := deck.probabilityOfEpidemic(); prob != 1.0/9.0 {
		t.Fatalf("Should have had a 1/9 probability of epidemic, got %v", prob)
	}
	deck.Drawn = append(deck.Drawn,
		CityCard{
			City{
				Name: "Epidemic!",
			},
			true,
		},
		CityCard{
			City{Name: "BuenosAires"},
			false,
		},
	)
	if prob := deck.probabilityOfEpidemic(); prob != 0 {
		t.Fatalf("Should have had a 0%% probability of epidemic, got %v", prob)
	}
}

func getTestCityDeck() CityDeck {
	cities := []*City{
		{
			Name:    "a",
			Disease: Blue.Type,
		},
		{
			Name:    "b",
			Disease: Blue.Type,
		},
		{
			Name:    "c",
			Disease: Blue.Type,
		},
		{
			Name:    "d",
			Disease: Yellow.Type,
		},
		{
			Name:    "e",
			Disease: Yellow.Type,
		},
		{
			Name:    "f",
			Disease: Yellow.Type,
		},
		{
			Name:    "g",
			Disease: Black.Type,
		},
		{
			Name:    "h",
			Disease: Black.Type,
		},
		{
			Name:    "i",
			Disease: Red.Type,
		},
		{
			Name:    "j",
			Disease: Red.Type,
		},
	}
	// 1/3 chance of an epidemic on a turn, since
	// we cut the 10 test cards above into 2 sections (1 for each epi)
	// and 2 cards are drawn from each set of 5+1.
	citiesStr := Cities{}
	citiesStr.Cities = cities
	deck := CityDeck{}
	deck.All = citiesStr.CityCards(2)
	return deck
}

type testState struct {
	infectRate   int
	infectDrawn  []string
	infectCustom func(infect *InfectionDeck) // if not set, will be equal to the names of all cities.
	cityCustom   func(deck *CityDeck)        // called to mutate the standard test deck
}

type testExpectation struct {
	scenario            string
	state               testState
	infectProbabilities map[string]float64 // round to hundredths
}

var infectTests = []testExpectation{
	{
		// Start of game, no cards drawn, and for simplicity, no chance of epidemic
		// in order to show probability of just drawing from infection deck dominates
		// the total probability
		scenario: "Start of game with no chance of epidemic",
		state: testState{
			infectRate: 2,
			cityCustom: func(deck *CityDeck) {
				deck.DrawEpidemic() // make it impossible to draw another epidemic for now.
			},
		},
		infectProbabilities: map[string]float64{
			"a": 0.2, // 2 draws out of 10, 1/5 chance
		},
	},
	{
		scenario: "Game with 100%% chance of epidemic and $rate-1 cards in drawn",
		state: testState{
			infectRate: 2,
			cityCustom: func(deck *CityDeck) {
				deck.Draw("a")
				deck.Draw("b")
				deck.Draw("c")
				deck.Draw("d")
				// now have 2 cards left in this striation, 100% chance of epidemic
			},
			infectCustom: func(deck *InfectionDeck) {
				deck.Draw("f") // only card in drawn is f, no matter what this should be 100% infect chance
			},
		},
		infectProbabilities: map[string]float64{
			"f": 1.0,  // 100% chance of drawing f again.
			"a": 0.11, // there is a 1/9 chance of infecting any bottom striation card.
		},
	},
	{
		scenario: "Game with 50%% chance of epidemic and $rate cards in drawn",
		state: testState{
			infectRate: 2,
			cityCustom: func(deck *CityDeck) {
				deck.Draw("a")
				deck.Draw("b")
				// 4 cards left in striation, 50% chance of epidemic
			},
			infectCustom: func(deck *InfectionDeck) {
				deck.Draw("e")
				deck.Draw("f")
				// 2 drawn infection cities makes chance of re-infecting on epidemic 2/3
			},
		},
		infectProbabilities: map[string]float64{
			"c": 0.18, // 1/4 chance of infect draw, 1/8 of epi draw
			"e": 0.33,
		},
	},
}

func TestRunInfectTests(t *testing.T) {
	for _, infectTest := range infectTests {
		// SETUP
		gs := GameState{}
		cityDeck := getTestCityDeck()
		if infectTest.state.cityCustom != nil {
			infectTest.state.cityCustom(&cityDeck)
		}
		cities := []*City{}
		for _, city := range cityDeck.All {
			if !city.IsEpidemic {
				city := city
				cities = append(cities, &city.City)
			}
		}
		gs.Cities = &Cities{Cities: cities}
		gs.CityDeck = &cityDeck
		gs.InfectionRate = infectTest.state.infectRate
		infectDeck := NewInfectionDeck(gs.Cities.CityNames())
		if infectTest.state.infectCustom != nil {
			infectTest.state.infectCustom(infectDeck)
		}
		gs.InfectionDeck = infectDeck

		// TEST
		for city, expected := range infectTest.infectProbabilities {
			prob := gs.ProbabilityOfCity(CityName(city))
			actual := math.Floor(prob*100) / 100.0
			if actual != expected {
				t.Errorf("In scenario '%v', %v did not have expected probability: wanted %v, got %v\n", infectTest.scenario, city, expected, actual)
			}
		}
	}
}
