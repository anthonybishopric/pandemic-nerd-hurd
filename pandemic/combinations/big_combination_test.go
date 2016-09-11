package combinations

import (
	"math"
	"testing"
)

func TestBigCombinationLargeFactors(t *testing.T) {
	c := bigCombination{}
	// calculate 10001!/10000!
	for i := 1; i <= 10000; i++ {
		c.numeratorTerms = append(c.numeratorTerms, i+1)
		c.denominatorTerms = append(c.denominatorTerms, i)
	}
	if round(c.Float64()) != round(10001.0) {
		t.Fatalf("Expected %v, got %v", 10001.0, c.Float64())
	}
}

func round(in float64) float64 {
	return math.Floor(in*1000) / 1000
}

func TestNChooseK(t *testing.T) {
	combo := nChooseK(52, 2)
	res := combo.Float64()
	if res != 1326.0 {
		t.Fatalf("Expected result to be 1326.0, got %v", res)
	}
}

// Standard deck of cards, draw at least 2 hearts in 6 draws.
func TestDrawHeartsExample(t *testing.T) {
	actual := round(AtLeastNDraws(52, 6, 2, 13))
	expected := round(1886.0 / 3995.0)
	if actual != expected {
		t.Fatalf("Expected %v, got %v", expected, actual)
	}
}
