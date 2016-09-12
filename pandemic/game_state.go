package pandemic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"sort"
	"strings"

	"github.com/anthonybishopric/pandemic-nerd-hurd/pandemic/combinations"
)

const EpidemicsPerGame = 5
const CityCardsPerTurn = 2

type GameState struct {
	Cities        *Cities        `json:"cities"`
	CityDeck      *CityDeck      `json:"city_deck"`
	DiseaseData   []DiseaseData  `json:"disease_data"`
	InfectionDeck *InfectionDeck `json:"infection_deck"`
	InfectionRate int            `json:"infection_rate"`
	Outbreaks     int            `json:"outbreaks"`
	GameName      string         `json:"game_name"`
	GameTurns     *GameTurns     `json:"game_turns"`
}

type NewGameSettings struct {
	Cities       Cities         `json:"cities"`
	Players      []*Player      `json:"players"`
	FundedEvents []*FundedEvent `json:"funded_events"`
}

func NewGame(newGameFile string, gameName string) (*GameState, error) {
	var newGameSettings NewGameSettings
	newGameData, err := ioutil.ReadFile(newGameFile)
	if err != nil {
		return nil, fmt.Errorf("Could not read new game file at %v: %v", newGameFile, err)
	}
	err = json.Unmarshal(newGameData, &newGameSettings)
	if err != nil {
		return nil, fmt.Errorf("Invalid new game JSON file at %v: %v", newGameFile, err)
	}
	cities := Cities(newGameSettings.Cities)
	players := newGameSettings.Players

	excludeFromCityDeck := Set{}
	for _, player := range players {
		if len(player.StartCards) != 2 {
			return nil, fmt.Errorf("Each player must start with 2 city cards")
		}
		for _, cityName := range player.StartCards {
			excludeFromCityDeck.Add(cityName)
		}
	}
	if len(excludeFromCityDeck) != 2*len(players) {
		return nil, fmt.Errorf("Duplicate cities detected, check the start information: %+v", excludeFromCityDeck)
	}

	cityDeck, err := cities.GenerateCityDeck(EpidemicsPerGame, newGameSettings.FundedEvents, excludeFromCityDeck)
	if err != nil {
		return nil, err
	}

	for _, player := range players {
		for _, startCard := range player.StartCards {
			card, err := cityDeck.GetCard(startCard)
			if err != nil {
				return nil, fmt.Errorf("%v is not a valid start city: %v", startCard, err)
			}
			player.Cards = append(player.Cards, card)
		}
	}

	infectionDeck := NewInfectionDeck(cities.CityNames())
	return &GameState{
		Cities:        &cities,
		DiseaseData:   []DiseaseData{Yellow, Red, Black, Blue, Faded},
		CityDeck:      &cityDeck,
		InfectionDeck: infectionDeck,
		InfectionRate: 2,
		Outbreaks:     0,
		GameName:      gameName,
		GameTurns:     InitGameTurns(players...),
	}, nil
}

func LoadGame(gameFile string) (*GameState, error) {
	var gameState GameState
	data, err := ioutil.ReadFile(gameFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &gameState)
	if err != nil {
		return nil, err
	}
	return &gameState, nil
}

func (gs GameState) ProbabilityOfCuring(player *Player, dt DiseaseType) float64 {
	// (diseaseColor choose requiredToCure)*(notDiseaseColor choose totalLessRequired)/(allCards choose totalExpectedDraws)
	remainingCards := gs.CityDeck.RemainingCardsWith(dt, gs.Cities)
	// TODO: make disease curability more programatic
	totalRequired := 5
	for _, card := range player.Cards {
		if !card.IsCity() {
			continue
		}
		city, err := gs.Cities.GetCity(card.CityName)
		if err != nil {
			panic("City card with no corresponding city: " + card.CityName)
		}
		if city.Disease == dt {
			totalRequired--
		}
	}
	if player.Character != nil && player.Character.Type == Scientist {
		totalRequired--
	}

	allRemaining := gs.CityDeck.RemainingCards()
	drawsRemaining := int(math.Floor(float64(allRemaining)/4.0)) - 2 // you don't get to use your last draw
	return combinations.AtLeastNDraws(allRemaining, drawsRemaining, totalRequired, remainingCards)
}

func (gs GameState) DrawCard(cn CardName) error {
	curTurn, err := gs.GameTurns.CurrentTurn()
	if err != nil {
		return err
	}
	if len(curTurn.DrawnCards) == CityCardsPerTurn {
		return fmt.Errorf("%v has already drawn %v cards this turn.", curTurn.Player.HumanName, CityCardsPerTurn)
	}
	card, err := gs.CityDeck.DrawCard(cn)
	if err != nil {
		return err
	}
	curTurn.DrawnCards = append(curTurn.DrawnCards, card)
	curTurn.Player.Cards = append(curTurn.Player.Cards, card)
	return nil
}

func (gs GameState) NextTurn() (*Turn, error) {
	return gs.GameTurns.NextTurn()
}

func (gs GameState) ExchangeCard(from, to *Player, name CardName) error {
	var senderNewCards []*CityCard
	var toGive *CityCard
	for _, card := range from.Cards {
		if card.Name() == name {
			toGive = card
		} else {
			senderNewCards = append(senderNewCards, card)
		}
	}
	if toGive == nil {
		return fmt.Errorf("%v does not seem to have the card %v", from.HumanName, name)
	}
	from.Cards = senderNewCards
	to.Cards = append(to.Cards, toGive)
	return nil
}

func (gs GameState) Infect(cn CityName) error {
	err := gs.InfectionDeck.Draw(cn)
	if err != nil {
		return err
	}
	city, err := gs.Cities.GetCity(cn)
	if err != nil {
		return err
	}
	if city.Quarantined {
		if !gs.quarantineSpecialistPresent(cn) {
			city.RemoveQuarantine()
		}
		return nil
	}
	// TODO: handle outbreaks
	city.Infect()
	return nil
}

func (gs GameState) Epidemic(cn CityName) error {
	err := gs.InfectionDeck.PullFromBottom(cn)
	if err != nil {
		return err
	}
	err = gs.CityDeck.DrawEpidemic()
	if err != nil {
		return err
	}
	city, _ := gs.Cities.GetCity(cn)

	if city.Quarantined {
		if !gs.quarantineSpecialistPresent(cn) {
			city.RemoveQuarantine()
		}
	} else {
		// TODO: handle outbreak
		city.Epidemic()
	}
	gs.InfectionDeck.ShuffleDrawn()
	return nil
}

func (gs GameState) quarantineSpecialistPresent(cityName CityName) bool {
	for _, player := range gs.GameTurns.PlayerOrder {
		if player.Location == cityName &&
			player.Character != nil && // TODO: actually support characters and remove null check
			player.Character.Type == QuarantineSpecialist {
			return true
		}
	}
	return false
}

func (gs GameState) Quarantine(cn CityName) error {
	city, err := gs.Cities.GetCity(cn)
	if err != nil {
		return err
	}
	if city.Quarantined {
		return fmt.Errorf("%v is already quarantined", cn)
	}
	city.Quarantine()
	return nil
}

func (gs GameState) RemoveQuarantine(cn CityName) error {
	city, err := gs.Cities.GetCity(cn)
	if err != nil {
		return err
	}
	if !city.Quarantined {
		return fmt.Errorf("%v is not quarantined ", cn)
	}
	city.RemoveQuarantine()
	return nil
}

// ProbabilityOfCity gives the aggregate probability of a city
// becoming infected. Quarantines make the probabilty of infection
// zero. This does not take into account the probability of infection
// due to neighboring city outbreaks.
func (gs GameState) ProbabilityOfCity(cn CityName) float64 {
	city, err := gs.Cities.GetCity(cn)
	if err != nil {
		return 0.0
	}
	if city.Quarantined {
		return 0.0
	}
	var cityDrawInfectRate float64
	// Check: does a city with 3 get additionally infected on drawing the city card?
	// Assume no, and no outbreak, for now.
	if DataForDisease(city.Disease).InfectOnCityDraw && city.NumInfections < 3 {
		cityDrawInfectRate = gs.CityDeck.ProbabilityOfDrawing(cn.CardName())
	}
	// P(epidemic)*P(pull from bottom or from infect drawn) + P(!epidemic)*P(infection deck draw)
	pEpi := gs.CityDeck.probabilityOfEpidemic()
	bottom := gs.InfectionDeck.BottomStriation()
	var pEpiDraw float64
	if bottom.Contains(cn) {
		pEpiDraw = 1.0 / float64(bottom.Size())
	} else if gs.InfectionDeck.Drawn.Contains(cn) {
		pEpiDraw = float64(gs.InfectionRate) / (1.0 + float64(len(gs.InfectionDeck.Drawn)))
	}

	pNoEpiDraw := gs.InfectionDeck.ProbabilityOfDrawing(cn, gs.InfectionRate)
	return cityDrawInfectRate + pEpi*pEpiDraw + (1.0-pEpi)*pNoEpiDraw
}

func (gs GameState) CanOutbreak(cn CityName) bool {
	city, err := gs.Cities.GetCity(cn)
	if err != nil {
		return false
	}
	if city.NumInfections == 0 && !DataForDisease(city.Disease).InfectOnCityDraw {
		return false
	}
	prob := gs.ProbabilityOfCity(cn)
	if prob == 0.0 {
		return false
	}
	return city.NumInfections == 3 || gs.InfectionDeck.BottomStriation().Contains(cn)
}

func (gs *GameState) GetCity(city CityName) (*City, error) {
	return gs.Cities.GetCity(city)
}

func (gs *GameState) GetDiseaseData(diseaseType DiseaseType) (*DiseaseData, error) {
	for _, data := range gs.DiseaseData {
		if data.Type == diseaseType {
			return &data, nil
		}
	}
	return nil, fmt.Errorf("No disease identified by %v", diseaseType)
}

func (gs *GameState) SortBySeverity(names []CityName) []CityName {
	b := bySeverity{names, gs}
	sort.Sort(&b)
	return b.names
}

type bySeverity struct {
	names []CityName
	gs    *GameState
}

func (b bySeverity) Len() int { return len(b.names) }

func (b bySeverity) Swap(i, j int) {
	b.names[i], b.names[j] = b.names[j], b.names[i]
}

func (b bySeverity) Less(i, j int) bool {
	nameI := b.names[i]
	nameJ := b.names[j]

	cityI, _ := b.gs.Cities.GetCity(nameI)
	cityJ, _ := b.gs.Cities.GetCity(nameJ)
	if cityI.NumInfections > cityJ.NumInfections {
		return true
	}
	if cityI.NumInfections < cityJ.NumInfections {
		return false
	}
	cityIProb := b.gs.ProbabilityOfCity(nameI)
	cityJProb := b.gs.ProbabilityOfCity(nameJ)
	if cityIProb > cityJProb {
		return true
	}
	if cityIProb < cityJProb {
		return false
	}
	return strings.Compare(string(nameI), string(nameJ)) < 0
}
