package pandemic

import (
	"fmt"
	"sort"
	"strings"
)

type CityName string

func (c CityName) String() string {
	return string(c)
}

type CityDeck struct {
	Drawn            []CityCard
	All              []CityCard
	probabilityModel *cityDeckProbabilityModel
}

type CityCard struct {
	City       City
	IsEpidemic bool `json:"is_epidemic"`
}

type City struct {
	Name          CityName    `json:"name"`
	Disease       DiseaseType `json:"disease"`
	PanicLevel    PanicLevel  `json:"panic_level"`
	Neighbors     []string    `json:"neighbors"`
	NumInfections int         `json:"num_infections"`
	Quarantined   bool        `json:"quarantined"`
}

type Cities struct {
	Cities []*City `json:"cities"`
}

type byInfectionRate struct {
	names  []CityName
	cities *Cities
}

func (b byInfectionRate) Len() int { return len(b.names) }

func (b byInfectionRate) Swap(i, j int) {
	b.names[i], b.names[j] = b.names[j], b.names[i]
}

func (b byInfectionRate) Less(i, j int) bool {
	nameI := b.names[i]
	nameJ := b.names[j]
	cityI, _ := b.cities.GetCity(nameI)
	cityJ, _ := b.cities.GetCity(nameJ)
	if cityI.NumInfections > cityJ.NumInfections {
		return true
	}
	if cityI.NumInfections < cityJ.NumInfections {
		return false
	}
	return strings.Compare(string(nameI), string(nameJ)) < 0
}

// do we need to model city specializations?
func (c *Cities) GenerateCityDeck(epidemicCount int) CityDeck {
	cards := []CityCard{}
	for _, city := range c.Cities {
		cards = append(cards, CityCard{*city, false})
	}
	for i := 0; i < epidemicCount; i++ {
		cards = append(cards, CityCard{City{}, true})
	}
	probModel := generateProbabilityModel(len(cards), epidemicCount)
	deck := CityDeck{
		Drawn:            []CityCard{},
		All:              cards,
		probabilityModel: &probModel,
	}
	return deck
}

func (c *Cities) SortByInfectionLevel(cities []CityName) []CityName {
	sorted := &byInfectionRate{cities, c}
	sort.Sort(sorted)
	return sorted.names
}

func (c *Cities) GetCityByPrefix(prefix string) (*City, error) {
	var ret *City
	for _, c := range c.Cities {
		c := c
		if strings.HasPrefix(strings.ToLower(string(c.Name)), strings.ToLower(prefix)) {
			if ret != nil {
				return nil, fmt.Errorf("'%v' is ambiguous", prefix)
			}
			ret = c
		}
	}
	if ret == nil {
		return nil, fmt.Errorf("%v is not a prefix for any city", prefix)
	}
	return ret, nil
}

func (c *Cities) GetCity(city CityName) (*City, error) {
	for _, c := range c.Cities {
		if c.Name == CityName(city) {
			return c, nil
		}
	}
	return nil, fmt.Errorf("No city named %v", city)
}

func (c Cities) WithDisease(disease DiseaseType) []*City {
	cities := []*City{}
	for _, city := range c.Cities {
		if city.Disease == disease {
			cities = append(cities, city)
		}
	}
	return cities
}

func (c Cities) CityNames() []CityName {
	names := []CityName{}
	for _, city := range c.Cities {
		names = append(names, city.Name)
	}
	return names
}

func (c *City) Infect() bool {
	if c.NumInfections == 3 {
		return true
	}
	c.NumInfections++
	return false
}

func (c *City) Epidemic() {
	c.NumInfections = 3
}

func (c *City) Quarantine() {
	c.Quarantined = true
}

func (c *City) RemoveQuarantine() {
	c.Quarantined = false
}

func (c *City) SetInfections(infections int) {
	c.NumInfections = infections
}

func (c CityDeck) Total() int {
	return len(c.All)
}

func (c *CityDeck) NumEpidemics() int {
	var totalEpis int
	for _, card := range c.All {
		if card.IsEpidemic {
			totalEpis++
		}
	}
	return totalEpis
}

func (c *CityDeck) cardsPerEpidemic() int {
	return c.Total() / c.NumEpidemics()
}

func (c *CityDeck) EpidemicsDrawn() int {
	count := 0
	for _, card := range c.Drawn {
		if card.IsEpidemic {
			count++
		}
	}
	return count
}

func (c *CityDeck) Draw(cn CityName) error {
	for _, card := range c.Drawn {
		if card.City.Name == cn {
			return fmt.Errorf("%v has already been drawn from the city deck", cn)
		}
	}
	c.probabilityModel.DrawCity(len(c.Drawn))
	for _, card := range c.All {
		if card.City.Name == cn {
			c.Drawn = append(c.Drawn, card)
			return nil
		}
	}
	return fmt.Errorf("No city called %v in the city deck", cn)
}

func (c *CityDeck) DrawEpidemic() error {
	totalEpis := c.NumEpidemics()
	var drawnEpis int
	for _, card := range c.Drawn {
		if card.IsEpidemic {
			drawnEpis++
		}
	}
	if drawnEpis >= totalEpis {
		return fmt.Errorf("Already drawn %v epidemics this game, there shouldn't be any more", drawnEpis)
	}
	c.probabilityModel.DrawEpidemic(len(c.Drawn))
	c.Drawn = append(c.Drawn, CityCard{City{}, true})
	return nil
}

// The function Pe(x) is the probabiltiy of drawing an epidemic at index x.
// Then, the probability of drawing an epidemic from the city deck is the sum
// of Pe(i) and Pe(i+1). However, Pe(i+1) is not independent from Pe(i), hence
// the sum of probabilities between two possible outcomes expressed below.
//
// Note that in perverse and upsetting circumstances, it is possible to have
// a probability of epidemic be greater than 1.0. This is entirely possible in
// the game of Pandemic.
func (c CityDeck) probabilityOfEpidemic() float64 {
	index := len(c.Drawn)
	analysis := c.probabilityModel.EpidemicAnalysis(index)
	return analysis.FirstCardProbability + analysis.SecondCardProbability
}

func (c CityDeck) EpidemicAnalysis() EpidemicAnalysis {
	index := len(c.Drawn)
	return c.probabilityModel.EpidemicAnalysis(index)
}

///////////////////////////////////
/// City Deck Probability Model ///
///////////////////////////////////

// The probability model of a given game of Pandemic: Legacy is composed of
// scenarios. Each scenario is capable of answering the question "what is the
// probability of an epidemic on card draw N?" The total probability of an
// epidemic draw is the weighed sum of probabilities of all scenarios.
type cityDeckProbabilityModel struct {
	scenarios      []cityDeckScenario
	epidemicsDrawn int
	lastIndex      int
}

// A deck scenario describes when the city deck has striations with card
// counts matching the cardCounts integer slice. As an example, consider a
// game scenario where the first 2 striations have 10 cards and the remaining
// 3 have 11 cards. This can occur in a game with 53 cards (48 cities, 5
// epidemics, no funded events). The underlying cardCounts slice will contain
// the values [10,10,11,11,11].
//
// While playing a real game of Pandemic, it is possible to draw epidemics in
// such a way that invalidate the possibility of a given scenario. In the
// above example, assume that we draw our first epidemic on turn 11. This
// would invalidate the [10,10,11,11,11] scenario because you guaranteed to
// draw exactly one epidemic in each striation. Thus, this scenario can be
// removed from the set of scenarios. As a result, weighted probabilities can
// be more precise with respect to actual possible scenarios.
type cityDeckScenario struct {
	cardCounts []int
}

type EpidemicAnalysis struct {
	FirstCardProbability  float64
	SecondCardProbability float64
	PossibleScenarios     int
	ScenariosWith100      int
}

// 1 extra is 5 possible scenarios 5!/1!(4!) = 5
// 2 extra is 10 possible scenarios (5!)/(2!)(3!) = 5*4/2 = 10
func generateProbabilityModel(cardCount int, epidemics int) cityDeckProbabilityModel {
	// (53-(53%5))/5 = (50/5) = 10
	minCardsPerStriation := (cardCount - (cardCount % epidemics)) / epidemics
	// 53 % 5 = 3
	striationsWithOneMore := cardCount % epidemics
	// we now have to calculate all permutations of scenarios that are possible.

	combinationSpace := 1 << uint(epidemics)

	scenarios := []cityDeckScenario{}
	for i := 0; i < combinationSpace; i++ {
		// find every binary string with exactly striationsWithOneMore 1s in it.
		// each one is a valid scenario
		binaryOneCount := 0
		binShrink := i
		for binShrink > 0 {
			if binShrink&1 == 1 {
				binaryOneCount++
			}
			if binaryOneCount > striationsWithOneMore {
				break
			}
			binShrink = binShrink >> 1
		}
		if binaryOneCount != striationsWithOneMore {
			continue
		}
		scenario := []int{}
		for striationAt := uint(0); striationAt < uint(epidemics); striationAt++ {
			// if the bit at striationAt in i is a 1, set to the higher value
			if (i>>striationAt)&1 == 1 {
				scenario = append(scenario, minCardsPerStriation+1)
			} else {
				scenario = append(scenario, minCardsPerStriation)
			}
		}
		scenarios = append(scenarios, cityDeckScenario{scenario})
	}
	return cityDeckProbabilityModel{scenarios, 0, -1}
}

func (c *cityDeckProbabilityModel) DrawCity(index int) {
	if index <= c.lastIndex {
		panic("Already drew this index!")
	}
	filtered := []cityDeckScenario{}
	for _, scenario := range c.scenarios {
		if scenario.EpidemicProbabilityAt(index, c.epidemicsDrawn) != 1.0 {
			filtered = append(filtered, scenario)
		}
	}
	c.scenarios = filtered
	c.lastIndex = index
}

func (c *cityDeckProbabilityModel) DrawEpidemic(index int) {
	if index <= c.lastIndex {
		panic("Already drew this index!")
	}
	filtered := []cityDeckScenario{}
	for _, scenario := range c.scenarios {
		if scenario.EpidemicProbabilityAt(index, c.epidemicsDrawn) != 0.0 {
			filtered = append(filtered, scenario)
		}
	}
	c.scenarios = filtered
	c.epidemicsDrawn++
	c.lastIndex = index
}

func (c *cityDeckProbabilityModel) EpidemicAnalysis(index int) EpidemicAnalysis {
	analysis := EpidemicAnalysis{}
	for _, scenario := range c.scenarios {
		scenProb := scenario.EpidemicProbabilityAt(index, c.epidemicsDrawn)
		scenProb2 := scenario.EpidemicProbabilityAt(index+1, c.epidemicsDrawn)
		if scenProb == 1.0 || scenProb2 == 1.0 {
			analysis.ScenariosWith100++
		}
	}
	analysis.PossibleScenarios = len(c.scenarios)
	analysis.FirstCardProbability = c.EpidemicProbabilityAt(index)

	noEpiOnFirst := *c
	(&noEpiOnFirst).DrawCity(index)
	epiOnFirst := *c
	(&epiOnFirst).DrawEpidemic(index)
	epiOnSecondAndFirst := epiOnFirst.EpidemicProbabilityAt(index + 1)
	epiOnSecondNotFirst := noEpiOnFirst.EpidemicProbabilityAt(index + 1)
	analysis.SecondCardProbability = analysis.FirstCardProbability*epiOnSecondAndFirst +
		(1.0-analysis.FirstCardProbability)*epiOnSecondNotFirst
	return analysis
}

func (c *cityDeckProbabilityModel) EpidemicProbabilityAt(index int) float64 {
	var aggregate float64
	denominator := float64(len(c.scenarios))
	for _, scenario := range c.scenarios {
		aggregate += scenario.EpidemicProbabilityAt(index, c.epidemicsDrawn) / denominator
	}
	return aggregate
}

func (c *cityDeckScenario) EpidemicProbabilityAt(index, epidemicsDrawn int) float64 {
	for i, striationCount := range c.cardCounts {
		if index >= striationCount {
			index = index - striationCount
		} else {
			if i < epidemicsDrawn {
				return 0.0
			}
			denominator := striationCount - index
			return 1.0 / float64(denominator)
		}
	}
	return 0.0
}
