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

func getCityNameByPrefix(entry string, gs *pandemic.GameState) (string, error) {
	city, err := gs.Cities.GetCityByPrefix(entry)
	if err != nil {
		return "", err
	}
	return city.Name, nil
}

func (p *PandemicView) runCommand(gameState *pandemic.GameState, consoleView *gocui.View, commandView *gocui.View) error {
	commandBuffer := strings.Trim(commandView.Buffer(), "\n\t\r ")
	if commandBuffer == "" {
		return nil
	}

	defer commandView.Clear()

	commandArgs := strings.Split(commandBuffer, " ")
	cmd := commandArgs[0]

	switch cmd {
	case "infect", "i":
		if len(commandArgs) != 2 {
			fmt.Fprintln(consoleView, p.colorWarning("You must pass a city to the infect command."))
			break
		}
		city, err := getCityNameByPrefix(commandArgs[1], gameState)
		if err != nil {
			fmt.Fprintln(consoleView, p.colorWarning(err))
			break
		}
		err = gameState.InfectionDeck.Draw(city)
		if err != nil {
			fmt.Fprintln(consoleView, p.colorWarning(err))
		} else {
			cityData, err := gameState.GetCity(city)
			if err != nil {
				panic(fmt.Sprintf("%v present in infection deck but not game state cities", city))
			}
			if cityData.Infect() {
				fmt.Fprintf(consoleView, p.colorOhFuck(fmt.Sprintf("Infected and outbroke %v\n", city)))
			} else {
				fmt.Fprintf(consoleView, "Infected %v\n", city)
			}
		}
	case "epidemic", "e":
		if len(commandArgs) != 2 {
			fmt.Fprintln(consoleView, p.colorWarning("You must pass a city to the epidemic command.\n"))
			break
		}
		city, err := getCityNameByPrefix(commandArgs[1], gameState)
		if err != nil {
			fmt.Fprintln(consoleView, p.colorWarning(err))
			break
		}
		err = gameState.InfectionDeck.PullFromBottom(city)
		if err != nil {
			fmt.Fprintln(consoleView, p.colorWarning(err))
			break
		} else {
			err := gameState.CityDeck.DrawEpidemic()
			if err != nil {
				fmt.Fprintln(consoleView, p.colorWarning(err.Error()))
			}
			fmt.Fprintf(consoleView, "Epidemic in %v. Please update the infect rate (infect-rate N)\n", city)
			cityData, _ := gameState.GetCity(city)
			cityData.Epidemic()
		}
		gameState.InfectionDeck.ShuffleDrawn()
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
		cityName, err := getCityNameByPrefix(commandArgs[1], gameState)
		if err != nil {
			fmt.Fprintln(consoleView, p.colorWarning(err))
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
			fmt.Fprintln(consoleView, p.colorWarning("You must pass a city value to draw\n"))
			break
		}
		city, err := getCityNameByPrefix(commandArgs[1], gameState)
		if err != nil {
			fmt.Fprintln(consoleView, p.colorWarning(err))
			break
		}
		err = gameState.CityDeck.Draw(city)
		if err != nil {
			fmt.Fprintln(consoleView, p.colorWarning(err))
			break
		}
		fmt.Fprintf(consoleView, "Drew %v from city deck\n", city)
	default:
		fmt.Fprintf(consoleView, p.colorWarning(fmt.Sprintf("Unrecognized command %v\n", cmd)))
		return nil
	}

	filename := filepath.Join(gameState.GameName, fmt.Sprintf("game_%v_%v.json", time.Now().UnixNano(), cmd))
	err := os.MkdirAll(gameState.GameName, 0755)
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
