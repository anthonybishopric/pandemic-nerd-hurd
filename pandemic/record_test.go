package pandemic

import (
	"testing"
)

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
			"San Francisco",
			false,
		},
		{
			"Sydney",
			false,
		},
	}
	if prob := deck.probabilityOfEpidemic(); prob != 1.0/9.0 {
		t.Fatalf("Should have had a 1/9 probability of epidemic, got %v", prob)
	}
	deck.Drawn = append(deck.Drawn,
		CityCard{
			"Epidemic!",
			true,
		},
		CityCard{
			"Buenos Aires",
			false,
		},
	)
	if prob := deck.probabilityOfEpidemic(); prob != 0 {
		t.Fatalf("Should have had a 0%% probability of epidemic, got %v", prob)
	}
}
