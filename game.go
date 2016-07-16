package main

import (
	"bufio"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	// TODO: load gamestate from argument
	go func() {
		signalCh := make(chan os.Signal, 2)
		signal.Notify(signalCh, syscall.SIGTERM, os.Interrupt)
		<-signalCh
		logger.Println("Exiting") // TODO: save gamestate to file
		os.Exit(0)
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		logger.Println("Got " + line) // TODO: implement REPL
	}
}
