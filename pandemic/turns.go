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
	Player      *Player    `json:"player"`
	DrawnCities []CityName `json:"drawn_cities"`
}

func (t *GameTurns) AddPlayer(p *Player) error {
	for _, existing := range t.PlayerOrder {
		if existing.Character.Type == p.Character.Type {
			return fmt.Errorf("The %v has already been added to the game", p.Character.Type)
		}
		if existing.HumanName == p.HumanName {
			return fmt.Errorf("%v has already been added to the game", p.HumanName)
		}
	}
	t.PlayerOrder = append(t.PlayerOrder, p)
	if len(t.PlayerOrder) == 1 {
		t.Turns = append(t.Turns, t.addTurn()) // create the first turn once we have a player
	}
	return nil
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
		Player:      t.PlayerOrder[t.CurTurn%len(t.PlayerOrder)],
		DrawnCities: []CityName{},
	}
}

func InitGameTurns() *GameTurns {
	return &GameTurns{
		0,
		[]*Player{},
		[]*Turn{},
	}
}
