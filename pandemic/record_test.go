package pandemic

import (
	"testing"
)

func TestProbabilityOfDrawingAlreadyDrawnCard(t *testing.T) {

	gs := GameState{
		InfectionDeck: InfectionDeck{
			Drawn: []InfectionCard{
				InfectionCard{
					Name: "SanFrancisco",
				},
			},
		},
		InfectionRate: 3,
	}

	if prob := gs.ProbabilityOfCity("SanFrancisco"); prob != 0 {
		t.Fatalf("Should have had a 0%% chance of SanFranciso, got %v", prob)
	}

}

func TestProbabilityOfDrawingCardAtStart(t *testing.T) {

	gs := GameState{
		InfectionDeck: InfectionDeck{
			Drawn: []InfectionCard{},
		},
		InfectionRate: 2,
	}

	// TODO(wjs) how 2 floating point
	if prob := gs.ProbabilityOfCity("SomeCity"); prob != 0.04166666666666674 {
		t.Fatalf("Should have had a 0%% chance of pulling random card with 3 picks, got %v", prob)
	}

}

func TestCardProbabilities(t *testing.T) {
	deck := CityDeck{
		Total: 100,
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
