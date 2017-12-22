package main

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/anthonybishopric/pandemic-nerd-hurd/pandemic"
	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
)

type PandemicView struct {
	logger              *logrus.Logger
	colorWhiteHighlight func(string, ...interface{}) string
	colorAllGood        func(string, ...interface{}) string
	colorWarning        func(string, ...interface{}) string
	colorHighlight      func(string, ...interface{}) string
	colorOhFuck         func(string, ...interface{}) string
	fileSaveCounter     int
}

func NewView(logger *logrus.Logger) *PandemicView {
	return &PandemicView{
		logger:              logger,
		colorWhiteHighlight: color.New(color.FgBlack).Add(color.BgWhite).SprintfFunc(),
		colorAllGood:        color.New(color.FgGreen).Add(color.BgBlack).SprintfFunc(),
		colorWarning:        color.New(color.FgYellow).Add(color.BgBlack).SprintfFunc(),
		colorHighlight:      color.New(color.FgRed).SprintfFunc(),
		colorOhFuck:         color.New(color.FgBlack).Add(color.BgRed).Add(color.BlinkSlow).SprintfFunc(),
	}
}

func (p *PandemicView) Start(game *pandemic.GameState) {
	gui, err := gocui.NewGui(gocui.OutputNormal)

	if err != nil {
		p.logger.Errorln("Could not init GUI: %v", err)
	}
	defer gui.Close()

	gui.SetManagerFunc(func(gui *gocui.Gui) error {
		width, height := gui.Size()

		p.renderCommandsView(game, gui, width)
		p.renderStriations(game, gui, 2, height/2, width)
		p.renderCityDeckAndTurns(game, gui, 0, height/2, width/2, height)
		p.renderConsoleArea(game, gui, width/2, height/2, width, height)

		p.setUpKeyBindings(game, gui, "Commands")
		gui.Cursor = true
		gui.SetCurrentView("Commands")
		return nil
	})

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		gui.Close()
		p.logger.Fatalf("Error in game main loop: %v", err)
	}
}

func (p *PandemicView) renderCommandsView(game *pandemic.GameState, gui *gocui.Gui, maxX int) {
	commandView, err := gui.SetView("Commands", 0, 0, maxX, 2)
	if err != nil && err != gocui.ErrUnknownView {
		gui.Close()
		p.logger.Fatalf("Could not render command view: %v", err)
	}
	commandView.Editable = true
	commandView.Autoscroll = false
	commandView.Title = "Commands"
}

func (p *PandemicView) renderCityDeckAndTurns(game *pandemic.GameState, gui *gocui.Gui, topX, topY, bottomX, bottomY int) {
	cityView, err := gui.SetView("Cities", topX, topY, bottomX, topY+(bottomY-topY)/2)
	p.logger.Infoln(topX, topY, bottomX, topY+(bottomY-topY)/2)
	if err != nil && err != gocui.ErrUnknownView {
		gui.Close()
		p.logger.Fatalf("Could not render city deck view: %v %v %v %v %v", err, topX, topY, bottomX, topY+(bottomY-topY)/2)
	}
	cityView.Clear()
	cityView.Title = "City Deck"
	cityView.Editable = false
	analysis := game.CityDeck.EpidemicAnalysis()
	total := analysis.FirstCardProbability + analysis.SecondCardProbability

	fmt.Fprintf(cityView, "\U0001F912 \U0001F4A5  %.2f (%v)\n", total, p.fractionalize(total))
	scenarioGuarantee := fmt.Sprintf("%v of %v Scenarios Guarantee Epidemic", analysis.ScenariosWith100, analysis.PossibleScenarios)
	if analysis.ScenariosWith100 > 0 {
		scenarioGuarantee = p.colorOhFuck(scenarioGuarantee)
	}
	fmt.Fprintln(cityView, scenarioGuarantee)

	fmt.Fprintf(cityView, "Epidemic on First City: %v\n", p.colorEpidemicPercent(analysis.FirstCardProbability))
	fmt.Fprintf(cityView, "Epidemic on Second City: %v\n", p.colorEpidemicPercent(analysis.SecondCardProbability))
	fmt.Fprintf(cityView, " -> After First City Epidemic: %v\n", p.colorEpidemicPercent(analysis.SecondCardEpiAfterFirstEpi))

	fmt.Fprintf(cityView, "Upcoming Draws Guaranteed Safe: %v\n", p.colorUpcomingSafeCount(analysis.ComingDrawsWith0))

	fmt.Fprintf(cityView, "Card counts %v  %v  ", p.iconFor(pandemic.Black.Type), game.CityDeck.RemainingCardsWith(pandemic.Black.Type, game.Cities))
	fmt.Fprintf(cityView, "%v  %v  ", p.iconFor(pandemic.Red.Type), game.CityDeck.RemainingCardsWith(pandemic.Red.Type, game.Cities))
	fmt.Fprintf(cityView, "%v  %v  ", p.iconFor(pandemic.Blue.Type), game.CityDeck.RemainingCardsWith(pandemic.Blue.Type, game.Cities))
	fmt.Fprintf(cityView, "%v  %v  ", p.iconFor(pandemic.Yellow.Type), game.CityDeck.RemainingCardsWith(pandemic.Yellow.Type, game.Cities))
	fmt.Fprintf(cityView, "%v  %v\n", p.iconFor(pandemic.Faded.Type), game.CityDeck.RemainingCardsWith(pandemic.Faded.Type, game.Cities))

	turnView, err := gui.SetView("Turns", topX, topY+(bottomY-topY)/2, bottomX, bottomY)
	if err != nil && err != gocui.ErrUnknownView {
		gui.Close()
		p.logger.Fatalf("Could not render turn view: %v", err)
	}
	turnView.Clear()
	turnView.Editable = false
	turnView.Title = "Players"

	cur, err := game.GameTurns.CurrentTurn()
	if err != nil {
		p.logger.Fatalln(err)
	}
	for _, player := range game.GameTurns.PlayerOrder {
		if cur.Player == player {
			fmt.Fprint(turnView, p.colorWhiteHighlight(player.HumanName))
		} else {
			fmt.Fprint(turnView, player.HumanName[:1])
		}
		fmt.Fprint(turnView, " ")
	}
	fmt.Fprintln(turnView)
	fmt.Fprintf(turnView, "%v has %v turns left\n", cur.Player.HumanName, game.GameTurns.RemainingTurnsFor(game.CityDeck.RemainingCards(), cur.Player.HumanName))
	if cur.Player.Character != nil && cur.Player.Character.TurnMessage != "" {
		fmt.Fprintln(turnView, p.colorAllGood(cur.Player.Character.TurnMessage))
	}

	// print all cards
	fmt.Fprint(turnView, "Cards: ")
	for _, card := range cur.Player.Cards {
		if card.IsCity() {
			city, _ := game.Cities.GetCity(card.CityName)
			fmt.Fprintf(turnView, "%v  %v ", p.iconFor(city.Disease), card.CityName[:4])
		} else if card.IsFundedEvent() {
			fmt.Fprintf(turnView, "\U0001F4B8  %v ", card.FundedEventName)
		}
	}
	fmt.Fprintln(turnView, "\nCure Likelihood: ")

	// print curability stats
	curability := byCurability{pandemic.CurableDiseases(), make(map[pandemic.DiseaseType]float64), make(map[pandemic.DiseaseType]maxCurability)}
	for _, dt := range pandemic.CurableDiseases() {
		playerProb := game.ProbabilityOfCuring(cur.Player, dt)
		curability.curability[dt] = playerProb
		curability.maxCurability[dt] = maxCurability{playerProb, cur.Player}
		for _, player := range game.GameTurns.PlayerOrder {
			if player.HumanName != cur.Player.HumanName {
				otherPlayerProb := game.ProbabilityOfCuring(player, dt)
				if otherPlayerProb > playerProb {
					curability.maxCurability[dt] = maxCurability{otherPlayerProb, player}
					playerProb = otherPlayerProb
				}
			}
		}
	}
	sort.Sort(curability)
	for _, dt := range curability.dts {
		max := curability.maxCurability[dt]
		maxStr := ""
		if max.player.HumanName != cur.Player.HumanName {
			maxStr = fmt.Sprintf("(%v %v)", max.player.HumanName, p.colorProbabilityOfCure(max.prob))
		}
		fmt.Fprintf(turnView, "%v  \U00002697  %v %v \n", p.iconFor(dt), p.colorProbabilityOfCure(curability.curability[dt]), maxStr)
	}
}

func (p *PandemicView) iconFor(dt pandemic.DiseaseType) string {
	var diseaseEmoji string
	switch dt {
	case pandemic.Yellow.Type:
		diseaseEmoji = "\U0001f49b"
	case pandemic.Blue.Type:
		diseaseEmoji = "\U0001f499"
	case pandemic.Red.Type:
		diseaseEmoji = "\u2764\ufe0f"
	case pandemic.Black.Type:
		diseaseEmoji = "\u26ab"
	case pandemic.Faded.Type:
		diseaseEmoji = "\U0001f608"
	default:
		diseaseEmoji = string(dt)
	}
	return diseaseEmoji
}

func (p *PandemicView) colorUpcomingSafeCount(safe int) string {
	if safe > 2 {
		return p.colorAllGood(fmt.Sprintf("%v", safe))
	} else if safe > 0 {
		return p.colorWarning(fmt.Sprintf("%v", safe))
	} else {
		return p.colorOhFuck(fmt.Sprintf("%v", safe))
	}
}

func (p *PandemicView) colorEpidemicPercent(total float64) string {
	var outStr string
	if total == 0.0 {
		outStr = p.colorAllGood(fmt.Sprintf("%.3f", total))
	} else if total > 0.5 {
		outStr = p.colorOhFuck(fmt.Sprintf("%.3f", total))
	} else {
		outStr = p.colorWarning(fmt.Sprintf("%.3f", total))
	}
	return outStr
}

func (p *PandemicView) fractionalize(decimalRep float64) string {
	num, dem := int(math.Floor(decimalRep*10+0.5)), 10
	for divisor := 2; float64(divisor) <= math.Min(float64(num), float64(dem)); divisor++ {
		if num%divisor == 0 && dem%divisor == 0 {
			num = num / divisor
			dem = dem / divisor
		}
	}
	return fmt.Sprintf("about %d out of %d", num, dem)
}

func (p *PandemicView) terminateIfErr(err error, msg string, gui *gocui.Gui) {
	if err != nil && err != gocui.ErrUnknownView {
		gui.Close()
		p.logger.Fatalf("%v: %v", msg, err)
	}
}

func (p *PandemicView) setUpKeyBindings(game *pandemic.GameState, gui *gocui.Gui, commandView string) {
	err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		// when we get a ctrl-C we exit the game
		gui.Close()
		p.logger.Fatalf("Buh bye") // TODO: save
		return nil
	})
	p.terminateIfErr(err, "could not establish graceful termination keybinding", gui)
	err = gui.SetKeybinding(commandView, gocui.KeyEnter, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		consoleView, err := gui.View("Console")
		if err != nil {
			gui.Close()
			p.logger.Fatalln("Console view not found, game view not set up correctly")
			return nil
		}
		return p.runCommand(game, consoleView, view)
	})
	err = gui.SetKeybinding(commandView, gocui.KeyTab, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		cleanBuffer := strings.Trim(view.Buffer(), "\n\t\r ")
		if cleanBuffer == "" {
			return nil
		}
		words := strings.Split(cleanBuffer, " ")
		prefix := words[len(words)-1]
		city, err := game.Cities.GetCityByPrefix(prefix)
		if err != nil {
			return nil
		}
		words[len(words)-1] = city.Name.String()
		x, y := view.Cursor()
		view.Clear()
		fmt.Fprint(view, strings.Join(words, " "))
		view.SetCursor(x+len(city.Name.String())-len(prefix), y)
		return nil
	})
	p.terminateIfErr(err, "could not establish keybinding for command view", gui)
}

func (p *PandemicView) renderConsoleArea(game *pandemic.GameState, gui *gocui.Gui, topX, topY, bottomX, bottomY int) {
	view, err := gui.SetView("Console", topX, topY, bottomX, bottomY)
	view.Title = "Console"
	p.terminateIfErr(err, "Could not set up console view", gui)
	view.Wrap = true
	view.Autoscroll = true
	if err == gocui.ErrUnknownView {
		fmt.Fprintf(view, "~ %v %v %v ~\n", p.colorAllGood("Pandemic Legacy"), p.colorHighlight("NeRd hUrD"), p.colorWarning("Assist-o-tron"))
		fmt.Fprintf(view, "Starting %v, %v City Cards, %v Epidemics, %v Funded Events\n", game.GameName, game.CityDeck.Total(), game.CityDeck.NumEpidemics(), game.CityDeck.NumFundedEvents())
	}
}

// Creates a series of columns, representing the current infection deck striations. Striations closer
// to the top of the infection deck are further to the right. Cities are colored based on the probability
// of being drawn.
func (p *PandemicView) renderStriations(game *pandemic.GameState, gui *gocui.Gui, topY int, bottomY int, maxX int) error {
	// We know there will never be more than 4 striations, not including drawn.
	// Divide the horizontal space by 5 and make striations that width. The 5th
	// column will be the drawn column
	strWidth := int(math.Floor(float64(maxX) / 5.0))

	for i := len(game.InfectionDeck.Striations) - 1; i >= 0; i-- {
		widthMultiplier := len(game.InfectionDeck.Striations) - i - 1
		cityNames := game.InfectionDeck.CitiesInStriation(i)
		strName := fmt.Sprintf("Infection %v", i)
		strView, err := gui.SetView(strName, strWidth*widthMultiplier, topY, (widthMultiplier+1)*strWidth, bottomY)
		if err != nil {
			return err
		}
		strView.Clear()
		strView.Title = strName
		cityNames = game.SortBySeverity(cityNames)
		for _, city := range cityNames {
			p.terminateIfErr(p.printCityWithProb(game, strView, city), "Could not render city", gui)
		}
	}
	widthMultiplier := 4
	drawnView, err := gui.SetView("Drawn", strWidth*widthMultiplier, topY, (widthMultiplier+1)*strWidth, bottomY)
	if err != nil {
		return err
	}
	drawnView.Clear()
	drawnView.Title = "Infection Drawn"
	for _, city := range game.InfectionDeck.CitiesInDrawn() {
		p.terminateIfErr(p.printCityWithProb(game, drawnView, city), "Could not render drawn card", gui)
	}
	return nil
}

func (p *PandemicView) printCityWithProb(game *pandemic.GameState, view *gocui.View, city pandemic.CityName) error {
	cityData, err := game.GetCity(city)
	if err != nil {
		return err
	}
	// diseaseData, err := game.GetDiseaseData(cityData.Disease)
	// if err != nil {
	// 	return err
	// }
	probability := game.ProbabilityOfCity(city)

	diseaseEmoji := p.iconFor(cityData.Disease)

	infectionRateEmojis := ""
	for i := 0; i < cityData.NumInfections; i++ {
		infectionRateEmojis += "•"
	}

	quarantinedEmoji := ""
	if cityData.Quarantined {
		quarantinedEmoji = "\u26d4"
	}

	text := fmt.Sprintf("%v %s  %s  %s  %.2f", city[:4], diseaseEmoji, infectionRateEmojis, quarantinedEmoji, probability)
	if probability == 0.0 {
		fmt.Fprintln(view, p.colorAllGood(text))
	} else if game.CanOutbreak(city) {
		fmt.Fprintln(view, p.colorOhFuck(text))
	} else {
		fmt.Fprintln(view, p.colorWarning(text))
	}
	return nil
}
