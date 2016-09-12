package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/anthonybishopric/pandemic-nerd-hurd/pandemic"
	"github.com/jroimartin/gocui"
)

func getCardByPrefix(entry string, gs *pandemic.GameState) (pandemic.CardName, error) {
	card, err := gs.CityDeck.GetCardByPrefix(entry)
	if err != nil {
		return "", err
	}
	return card.Name(), nil
}

func getCityByPrefix(entry string, gs *pandemic.GameState) (pandemic.CityName, error) {
	card, err := gs.CityDeck.GetCardByPrefix(entry)
	if err != nil {
		return pandemic.CityName(""), err
	}
	if !card.IsCity() {
		return pandemic.CityName(""), fmt.Errorf("%v is not a city", card.Name())
	}
	return card.CityName, nil
}

func getPlayerByPrefix(entry string, gs *pandemic.GameState) (*pandemic.Player, error) {
	var ret *pandemic.Player
	for _, player := range gs.GameTurns.PlayerOrder {
		if strings.HasPrefix(strings.ToLower(player.HumanName), strings.ToLower(entry)) {
			if ret != nil {
				return nil, fmt.Errorf("%v is an ambiguous human name", entry)
			} else {
				ret = player
			}
		}
	}
	return ret, nil
}

func (p *PandemicView) runCommand(gameState *pandemic.GameState, consoleView *gocui.View, commandView *gocui.View) error {
	commandBuffer := strings.Trim(commandView.Buffer(), "\n\t\r ")
	if commandBuffer == "" {
		return nil
	}
	defer commandView.SetCursor(commandView.Origin())
	defer commandView.Clear()

	commandArgs := strings.Split(commandBuffer, " ")
	cmd := commandArgs[0]

	curTurn, err := gameState.GameTurns.CurrentTurn()
	if err != nil {
		return err
	}
	curPlayer := curTurn.Player

	switch cmd {
	case "infect", "i":
		if len(commandArgs) != 2 {
			fmt.Fprintln(consoleView, p.colorWarning("You must pass a city to the infect command."))
			break
		}
		city, err := getCityByPrefix(commandArgs[1], gameState)
		if err != nil {
			fmt.Fprintln(consoleView, p.colorWarning("%v", err))
			break
		}
		err = gameState.Infect(city)
		if err != nil {
			fmt.Fprintln(consoleView, p.colorWarning("%v", err))
		} else {
			fmt.Fprintf(consoleView, "Infected %v\n", city)
		}
	case "next-turn", "n":
		turn, err := gameState.NextTurn()
		if err != nil {
			fmt.Fprintln(consoleView, p.colorWarning("Could not move on to next turn: %v", err))
		} else {
			fmt.Fprintf(consoleView, "It is now %v's turn\n", turn.Player.HumanName)
		}
	case "give-card", "g":
		if len(commandArgs) != 3 {
			fmt.Fprintln(consoleView, p.colorWarning("Usage: give-card <human-prefix> <city-prefix>"))
			break
		}
		from, err := gameState.GameTurns.CurrentTurn()
		if err != nil {
			fmt.Fprintln(consoleView, p.colorWarning("%v", err))
			break
		}
		to, err := getPlayerByPrefix(commandArgs[1], gameState)
		if err != nil {
			fmt.Fprintln(consoleView, p.colorWarning("%v", err))
			break
		}
		cardName, err := getCardByPrefix(commandArgs[2], gameState)
		if err != nil {
			fmt.Fprintln(consoleView, p.colorWarning("%v", err))
			break
		}
		err = gameState.ExchangeCard(from.Player, to, cardName)
		if err != nil {
			fmt.Fprintln(consoleView, p.colorWarning("%v", err))
			break
		} else {
			fmt.Fprintf(consoleView, "%v gave %v to %v\n", from.Player.HumanName, cardName, to.HumanName)
		}
	case "epidemic", "e":
		if len(commandArgs) != 2 {
			fmt.Fprintln(consoleView, p.colorWarning("You must pass a city to the epidemic command."))
			break
		}
		city, err := getCityByPrefix(commandArgs[1], gameState)
		if err != nil {
			fmt.Fprintln(consoleView, p.colorWarning("%v", err))
			break
		}
		err = gameState.Epidemic(city)
		if err != nil {
			fmt.Fprintln(consoleView, p.colorWarning("%v", err))
			break
		} else {
			fmt.Fprintf(consoleView, "Epidemic in %v. Please update the infect rate (infect-rate N)\n", city)
		}
	case "infect-rate", "r":
		if len(commandArgs) != 2 {
			fmt.Fprintln(consoleView, p.colorWarning("You must pass an integer value to the infect rate\n"))
			break
		}
		ir, err := strconv.ParseInt(commandArgs[1], 10, 32)
		if err != nil {
			fmt.Fprintf(consoleView, p.colorWarning(fmt.Sprintf("%v is not a valid infection rate\n", commandArgs[1])))
		} else {
			fmt.Fprintf(consoleView, "infection rate now %v\n", ir)
			gameState.InfectionRate = int(ir)
		}
	case "city-infect-level", "l":
		if len(commandArgs) != 3 {
			fmt.Fprintln(consoleView, p.colorWarning("You must pass a city and infection value"))
			break
		}
		il, err := strconv.ParseInt(commandArgs[2], 10, 32)
		if err != nil {
			fmt.Fprintf(consoleView, p.colorWarning(fmt.Sprintf("%v is not a valid infection level\n", commandArgs[1])))
			break
		}
		cityName, err := getCityByPrefix(commandArgs[1], gameState)
		if err != nil {
			fmt.Fprintln(consoleView, p.colorWarning("%v", err))
			break
		}
		city, err := gameState.GetCity(cityName)
		if err != nil {
			fmt.Fprintf(consoleView, p.colorWarning(fmt.Sprintf("Could not get city %v: %v\n", cityName, err)))
			break
		}
		city.SetInfections(int(il))
		fmt.Fprintf(consoleView, "Set infection level in %v to %v\n", city.Name, city.NumInfections)
	case "city-draw", "c":
		if len(commandArgs) != 2 {
			fmt.Fprintln(consoleView, p.colorWarning("You must pass a city or funded event name to draw\n"))
			break
		}
		cardName, err := getCardByPrefix(commandArgs[1], gameState)
		if err != nil {
			fmt.Fprintln(consoleView, p.colorWarning("%v", err))
			break
		}
		err = gameState.DrawCard(cardName)
		if err != nil {
			fmt.Fprintln(consoleView, p.colorWarning("%v", err))
			break
		}
		fmt.Fprintf(consoleView, "%v drew %v from city deck\n", curPlayer.HumanName, cardName)
	case "quarantine", "q":
		if len(commandArgs) != 2 {
			fmt.Fprintln(consoleView, p.colorWarning("quarantine must be called with a city name"))
			break
		}
		cityName, err := getCityByPrefix(commandArgs[1], gameState)
		if err != nil {
			fmt.Fprintln(consoleView, p.colorWarning("%v", err))
			break
		}
		err = gameState.Quarantine(cityName)
		if err != nil {
			fmt.Fprintln(consoleView, p.colorWarning(fmt.Sprintf("Could not quarantine %v: %v", cityName, err)))
		} else {
			fmt.Fprintf(consoleView, "Quarantined %v\n", cityName)
		}
	case "discard", "d":
		if len(commandArgs) != 2 {
			fmt.Fprintln(consoleView, p.colorWarning("discard must be called with a city name"))
			break
		}
		cardName, err := getCardByPrefix(commandArgs[1], gameState)
		if err != nil {
			fmt.Fprintln(consoleView, p.colorWarning("%v", err))
			break
		}
		err = curPlayer.Discard(cardName)
		if err != nil {
			fmt.Fprintln(consoleView, p.colorWarning("%v", err))
			break
		}
		fmt.Fprintf(consoleView, "%v discarded %v\n", curPlayer.HumanName, cardName)
	case "remove-quarantine", "rq":
		if len(commandArgs) != 2 {
			fmt.Fprintf(consoleView, p.colorWarning("remove-quarantine must be called with a city name"))
		}
		cityName, err := getCityByPrefix(commandArgs[1], gameState)
		if err != nil {
			fmt.Fprintln(consoleView, p.colorWarning("%v", err))
			break
		}
		err = gameState.RemoveQuarantine(cityName)
		if err != nil {
			fmt.Fprintln(consoleView, p.colorWarning(fmt.Sprintf("Could not remove quarantine from %v: %v", cityName, err)))
		} else {
			fmt.Fprintf(consoleView, "Removed quarantine from %v\n", cityName)
		}
	default:
		fmt.Fprintf(consoleView, p.colorWarning(fmt.Sprintf("Unrecognized command %v\n", cmd)))
		return nil
	}

	filename := filepath.Join(gameState.GameName, fmt.Sprintf("game_%v_%v.json", time.Now().UnixNano(), cmd))
	err = os.MkdirAll(gameState.GameName, 0755)
	if err != nil {
		fmt.Fprintf(consoleView, p.colorOhFuck(fmt.Sprintf("Could not create a game name folder: %v", err)))
	}
	data, err := json.Marshal(gameState)
	if err != nil {
		fmt.Fprintf(consoleView, p.colorOhFuck(fmt.Sprintf("Could not marshal gamestate as JSON: %v\n", err)))
		return nil
	}
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Fprintf(consoleView, p.colorOhFuck(fmt.Sprintf("Could not save gamestate: %v\n", err)))
		return nil
	}

	return nil
}
