package main

import (
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/anthonybishopric/pandemic-nerd-hurd/pandemic"
)

func main() {
	logger := logrus.New()
	fd, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	logger.Out = fd
	wd, _ := os.Getwd()

	gameState, err := pandemic.NewGame(filepath.Join(wd, "data/cities.json")) // Parameterize
	if err != nil {
		logger.Fatal(err)
	}
	gameState.InfectionRate = 3

	gameState.InfectionDeck.Draw("Cairo")
	gameState.InfectionDeck.Draw("Cairo")
	gameState.InfectionDeck.Draw("Bogota")
	gameState.InfectionDeck.Draw("Santiago")
	gameState.InfectionDeck.Draw("SanFrancisco")
	gameState.InfectionDeck.Draw("Montreal")
	gameState.InfectionDeck.Draw("London")
	gameState.InfectionDeck.Draw("Tehran")
	gameState.InfectionDeck.Draw("Beijing")
	gameState.InfectionDeck.ShuffleDrawn()

	gameState.InfectionDeck.Draw("Bogota")
	gameState.InfectionDeck.Draw("SanFrancisco")
	gameState.InfectionDeck.Draw("Santiago")
	gameState.InfectionDeck.Draw("Montreal")
	gameState.InfectionDeck.Draw("Tehran")
	gameState.InfectionDeck.Draw("Beijing")

	gameState.InfectionDeck.ShuffleDrawn()

	gameState.InfectionDeck.Draw("Tehran")
	gameState.InfectionDeck.Draw("Beijing")
	gameState.InfectionDeck.Draw("Santiago")

	gameState.InfectionDeck.ShuffleDrawn()

	view := NewView(logger)
	view.Start(gameState)
}
