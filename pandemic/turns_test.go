package pandemic

import (
	"testing"
)

func TestRemainingTurns(t *testing.T) {
	scenarios := []struct {
		targetPlayer   int
		curTurnIndex   int
		remainingCards int
		expectedTurns  int
	}{
		{0, 0, 10, 2},
		{3, 0, 10, 1},
		{1, 0, 10, 1},
		{1, 0, 11, 2},
		{0, 2, 6, 1},
		{0, 1, 22, 2},
		{0, 1, 23, 3},
		{1, 0, 4, 1},
		{2, 0, 4, 0},
		{1, 0, 2, 0},
		{0, 3, 10, 1},
		{0, 0, 8, 1},
		{0, 0, 9, 2},
		{0, 3, 8, 1},
		{0, 3, 9, 1},
	}
	for _, scenario := range scenarios {
		turns := InitGameTurns(
			&Player{
				HumanName: "a",
			},
			&Player{
				HumanName: "b",
			},
			&Player{
				HumanName: "c",
			},
			&Player{
				HumanName: "d",
			})
		humanName := turns.PlayerOrder[scenario.targetPlayer].HumanName
		turns.CurTurn = scenario.curTurnIndex
		res := turns.RemainingTurnsFor(scenario.remainingCards, humanName)
		if res != scenario.expectedTurns {
			t.Errorf("%+v: Expected player %v to have %v turns, instead had %v", scenario, humanName, scenario.expectedTurns, res)
		}
	}
}
