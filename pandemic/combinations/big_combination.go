package combinations

import (
	"sort"
)

// represents every factor in the numerator and
// denominator of a combinatorial expression.
type bigCombination struct {
	numeratorTerms   []int
	denominatorTerms []int
}

// productOf returns a single big combination that is the product
// of the numerators and denominators respectively
func productOf(a, b bigCombination) bigCombination {
	a.numeratorTerms = append(a.numeratorTerms, b.numeratorTerms...)
	a.denominatorTerms = append(a.denominatorTerms, b.denominatorTerms...)
	return a
}

func inverseOf(a bigCombination) bigCombination {
	a.numeratorTerms, a.denominatorTerms = a.denominatorTerms, a.numeratorTerms
	return a
}

// nChooseK returns a bigCombination that resolves the n•choose*k
// operation, calculated as:
//
//      n!
//  ---------
//   k!(n-k)!
//
func nChooseK(n, k int) bigCombination {
	c := bigCombination{}
	for i := n; i > n-k; i-- {
		c.numeratorTerms = append(c.numeratorTerms, i)
	}
	for i := k; i > 0; i-- {
		c.denominatorTerms = append(c.denominatorTerms, i)
	}
	return c
}

func (b bigCombination) Float64() float64 {
	// resolves the combination as a float by first
	// canceling any redundant terms and then iteratively
	// applying division to minimize the risk of overflow
	nNum, nDem := 0, 0
	sort.Ints(b.numeratorTerms)
	sort.Ints(b.denominatorTerms)
	sum := 1.0
	for nNum < len(b.numeratorTerms) || nDem < len(b.denominatorTerms) {
		for (nNum == len(b.numeratorTerms) || sum > 1.0) && nDem < len(b.denominatorTerms) {
			sum = sum / float64(b.denominatorTerms[nDem])
			nDem++
		}
		for (nDem == len(b.denominatorTerms) || sum <= 1.0) && nNum < len(b.numeratorTerms) {
			sum = sum * float64(b.numeratorTerms[nNum])
			nNum++
		}
	}
	return sum
}

// ExactlyNCardDraws returns the exact probability of drawing exactly
// n cards from a given deck of D, given some family of N cards that
// could match the given criteria.
func ExactlyNCardDraws(totalDeckSize int, numDraws int, n int, familySize int) float64 {
	if familySize > totalDeckSize {
		return 0.0
	}
	if n > numDraws {
		return 0.0
	}
	chooseNColor := nChooseK(familySize, n)
	chooseDrawLessNOther := nChooseK(totalDeckSize-familySize, numDraws-n)
	allPossibilities := nChooseK(totalDeckSize, numDraws)

	combination := productOf(chooseNColor, productOf(chooseDrawLessNOther, inverseOf(allPossibilities)))
	return combination.Float64()
}

// AtLeastNDraws calculates the probability of drawing at least
// N of the given card type from the set of cards. It uses the
// ExactlyNCardDraws function and subtracts cases from the total
// probability.
func AtLeastNDraws(totalDeckSize int, numDraws int, n int, familySize int) float64 {
	if familySize > totalDeckSize {
		return 0.0
	}
	if n > numDraws {
		return 0.0
	}
	atLeast := 1.0
	for i := 0; i < n; i++ {
		atLeast -= ExactlyNCardDraws(totalDeckSize, numDraws, i, familySize)
	}
	return atLeast
}
