package main

import (
	"KamaiZen/document_manager"
	"KamaiZen/logger"
	"KamaiZen/lsp"
	"KamaiZen/server"
	"KamaiZen/settings"
	"KamaiZen/state_manager"
	"os"
	"sync"
)

const CHAN_BUFFER = 24

func initialise() {
	// TODO: load settings from config file
	args := os.Args[1:]
	if len(args) == 0 {
		logger.Error("No settings file provided")
		os.Exit(1)
	}
	filepath := args[0]
	settingsReader := settings.JSONSettingsReader{}
	settings := settingsReader.ReadSettings(filepath)
	logger.Info("Starting KamaiZen")
	state_manager.InitializeState()
	lsp.Initialise()
	document_manager.Initialise(settings)
}

func main() {
	initialise()
	defer logger.Info("KamaiZen stopped")

	analyser_channel := make(chan state_manager.State, CHAN_BUFFER)

	var wg sync.WaitGroup
	wg.Add(3)
	// go state_manager.StartAnalyser(analyser_channel, &wg)
	go server.StartServer(&wg, analyser_channel)
	go lsp.Start(&wg)
	wg.Wait()
	close(analyser_channel)
}
