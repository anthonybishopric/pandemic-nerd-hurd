package pandemic

import (
	"fmt"
)

type GameTurns struct {
	CurTurn     int       `json:"cur_turn"`
	PlayerOrder []*Player `json:"player_order"`
	Turns       []*Turn   `json:"turns"`
}

type Turn struct {
	Player     *Player     `json:"player"`
	DrawnCards []*CityCard `json:"drawn_cards"`
}

func (t *GameTurns) AddPlayer(p *Player) error {
	// for _, existing := range t.PlayerOrder {
	// 	if existing.Character.Type == p.Character.Type {
	// 		return fmt.Errorf("%v cannot be %v because %v already been added to the game by %v", p.HumanName, p.Character.Type, p.Character.Type, existing.HumanName)
	// 	}
	// 	if existing.HumanName == p.HumanName {
	// 		return fmt.Errorf("%v has already been added to the game", p.HumanName)
	// 	}
	// }
	t.PlayerOrder = append(t.PlayerOrder, p)
	if len(t.PlayerOrder) == 1 {
		t.Turns = append(t.Turns, t.addTurn()) // create the first turn once we have a player
	}
	return nil
}

func (t *GameTurns) RemainingTurnsFor(remainingCityCards int, name string) int {
	index := -1
	for i, player := range t.PlayerOrder {
		if player.HumanName == name {
			index = i
		}
	}
	if index == -1 {
		return 0
	}

	lastPlayerIndex := (t.CurTurn + remainingCityCards/2) % len(t.PlayerOrder)
	base := remainingCityCards / (2 * len(t.PlayerOrder))
	var oddAdd int
	if remainingCityCards%2 == 1 {
		oddAdd = 1
	}
	if lastPlayerIndex == index {
		return base + oddAdd
	}
	turnDistance := index - lastPlayerIndex + 1
	if turnDistance < 0 {
		turnDistance += len(t.PlayerOrder)
	}
	if 2*turnDistance < (remainingCityCards)%(2*len(t.PlayerOrder)) {
		return base + 1
	}
	return base
}

func (t *GameTurns) CurrentTurn() (*Turn, error) {
	if len(t.PlayerOrder) < 2 {
		return nil, fmt.Errorf("Need at least two players before starting the game, currently have %v", len(t.PlayerOrder))
	}
	return t.Turns[t.CurTurn], nil
}

func (t *GameTurns) NextTurn() (*Turn, error) {
	if len(t.PlayerOrder) < 2 {
		return nil, fmt.Errorf("Need at least two players before starting the game, currently have %v", len(t.PlayerOrder))
	}
	t.CurTurn = t.CurTurn + 1
	t.Turns = append(t.Turns, t.addTurn())
	return t.CurrentTurn()
}

func (t *GameTurns) addTurn() *Turn {
	return &Turn{
		Player:     t.PlayerOrder[t.CurTurn%len(t.PlayerOrder)],
		DrawnCards: []*CityCard{},
	}
}

func (t *GameTurns) AddDrawnToCurrent(card *CityCard) error {
	turn, err := t.CurrentTurn()
	if err != nil {
		return err
	}
	if len(turn.DrawnCards) == CityCardsPerTurn {
		return fmt.Errorf("Already drew %v cards this turn", CityCardsPerTurn)
	}
	turn.DrawnCards = append(turn.DrawnCards, card)
	return nil
}

func InitGameTurns(ps ...*Player) *GameTurns {
	turns := &GameTurns{
		0,
		[]*Player{},
		[]*Turn{},
	}
	for _, p := range ps {
		turns.AddPlayer(p)
	}
	return turns
}
