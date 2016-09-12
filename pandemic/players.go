package pandemic

import (
	"fmt"
)

type CharacterType int

const (
	Medic = CharacterType(iota)
	Dispatcher
	Researcher
	Scientist
	Civilian
	QuarantineSpecialist
	Colonel
	OperationsExpert
	Generalist
)

type Player struct {
	HumanName  string `json:"human_name"`
	Character  *Character
	Location   CityName
	StartCards []CardName `json:"start_cards"`
	Cards      []*CityCard
}

func (p *Player) Discard(cardName CardName) error {
	filtered := []*CityCard{}
	for _, card := range p.Cards {
		if card.Name() != cardName {
			filtered = append(filtered, card)
		}
	}
	if len(filtered) == len(p.Cards) {
		return fmt.Errorf("%v does not seem to have %v\n", p.HumanName, cardName)
	}
	p.Cards = filtered
	return nil
}

type Character struct {
	Name string        `json:"name"`
	Type CharacterType `json:"character_type"`
}
