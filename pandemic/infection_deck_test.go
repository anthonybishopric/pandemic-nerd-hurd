package pandemic

import (
	"math"
	"testing"
)

func testInfectionDeck() *InfectionDeck {
	return NewInfectionDeck([]string{
		"SanFrancisco",
		"NewYork",
		"Montreal",
		"Miami",
		"Washington",
	})
}

func TestInfectionDeckCountStriations(t *testing.T) {
	deck := testInfectionDeck()

	if deck.CurrentStriationCount() != 5 {
		t.Fatal("Should have expected 5 cards in the current striation")
	}
	if err := deck.Draw("SanFrancisco"); err != nil {
		t.Fatalf("Did not expect error when drawing: %v", err)
	}
	if !deck.Drawn.Contains("SanFrancisco") {
		t.Fatal("Expected the drawn deck to have san francisco")
	}
	if deck.CurrentStriationCount() != 4 {
		t.Fatal("Expected 4 remaining cards in the current striation")
	}

}

func checkProbability(t *testing.T, deck *InfectionDeck, city string, infectRate int, expected float64) {
	// round to hundredths for the comparison
	actual := deck.ProbabilityOfDrawing(city, infectRate)
	actualRounded := math.Floor(actual*100) / 100.0
	expectedRounded := math.Floor(expected*100) / 100.0
	if expectedRounded != actualRounded {
		t.Fatalf("Expected probability of drawing %v to be %v but was %v", city, expectedRounded, actualRounded)
	}
}

func TestProbabilityOfStriations(t *testing.T) {
	deck := testInfectionDeck()
	checkProbability(t, deck, "Washington", 3, 3.0/5.0)
	deck.Draw("SanFrancisco")
	deck.ShuffleDrawn()
	checkProbability(t, deck, "SanFrancisco", 3, 1.0)
	checkProbability(t, deck, "Washington", 1, 0.0)
	checkProbability(t, deck, "Washington", 2, 0.25)
}
