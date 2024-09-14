package main

import (
	"KamaiZen/logger"
	"KamaiZen/lsp"
	"KamaiZen/server"
	"KamaiZen/state_manager"
	"sync"
)

const CHAN_BUFFER = 24

func initialise() {
	// TODO: load settings from config file
	file := "./config/config.json"
	settingsReader := settings.JSONSettingsReader{}
	settings := settingsReader.ReadSettings(file)
	logger.Info("Starting KamaiZen")
	state_manager.InitializeState()
	lsp.Initialise()
}

func main() {
	initialise()
	defer logger.Info("KamaiZen stopped")

	analyser_channel := make(chan state_manager.State, CHAN_BUFFER)
	server := server.GetServerInstance()

	var wg sync.WaitGroup
	wg.Add(2)
	go server.StartServer(&wg)
	go lsp.Start(&wg)
	wg.Wait()
	close(analyser_channel)
}
