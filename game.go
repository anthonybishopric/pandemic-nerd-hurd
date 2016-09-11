package main

import (
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/anthonybishopric/pandemic-nerd-hurd/pandemic"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app         = kingpin.New("pandemic–nerd-hurd", "Start a nerd herd game")
	startCmd    = app.Command("start", "Start a new game")
	startCities = startCmd.Flag("cities-file", "The file containing initial data about Cities").Default("data/cities.json").ExistingFile()
	startMonth  = startCmd.Flag("month", "The name of the month in the game we are playing. If playing the second time in a month, add '2' after the name").Required().Enum(
		"jan",
		"feb",
		"mar",
		"apr",
		"may",
		"jun",
		"jul",
		"aug",
		"sep",
		"oct",
		"nov",
		"dec",
		"jan2",
		"feb2",
		"mar2",
		"apr2",
		"may2",
		"jun2",
		"jul2",
		"aug2",
		"sep2",
		"oct2",
		"nov2",
		"dec2",
	)
	startNumFundedEvents = startCmd.Flag("funded-events", "The number of funded events present in the city deck.").Required().Int()
	startPlayerFile      = startCmd.Flag("players-file", "Player metadata describing who will be playing as which characters and in what order.").Default("data/players.json").ExistingFile()
	loadCmd              = app.Command("load", "Load a game from an existing saved game")
	loadFile             = loadCmd.Flag("file", "The JSON file containing the game state").Required().ExistingFile()
)

func main() {
	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	logger := logrus.New()
	fd, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	logger.Out = fd
	wd, _ := os.Getwd()

	var gameState *pandemic.GameState

	switch cmd {
	case "start":
		gameState, err = pandemic.NewGame(filepath.Join(wd, *startCities), *startMonth, *startNumFundedEvents, *startPlayerFile)
		if err != nil {
			logger.Fatalln(err)
		}
	case "load":
		gameState, err = pandemic.LoadGame(filepath.Join(wd, *loadFile))
		if err != nil {
			logger.Fatalln(err)
		}
	}

	view := NewView(logger)
	view.Start(gameState)
}
