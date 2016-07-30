package pandemic

import (
	"testing"
)

func TestProbabilityOfDrawingAlreadyDrawnCard(t *testing.T) {

	gs := GameState{
		InfectionDeck: NewInfectionDeck([]string{"SanFrancisco", "Miami", "Washington"}),
		InfectionRate: 3,
	}

	if prob := gs.ProbabilityOfCity("SanFrancisco"); prob != 1.0 {
		t.Fatalf("Should have had a 100%% chance of SanFranciso, got %v", prob)
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
