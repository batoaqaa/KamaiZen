package main

import (
	"KamaiZen/analysis"
	"KamaiZen/document_manager"
	"KamaiZen/logger"
	"KamaiZen/lsp"
	"KamaiZen/server"
	"KamaiZen/settings"
	"sync"
)

const CHAN_BUFFER = 24

func initialise() {
	// TODO: load settings from config file
	settings := settings.NewLSPSettings(
		"/home/ibrahim/work/kamailio/",
		"/path/to/root", logger.DEBUG)
	logger.SetLogLevel(settings.LogLevel())
	logger.Info("Starting KamaiZen")
	lsp.Initialise()
	document_manager.Initialise(settings)
}

func main() {
	initialise()
	defer logger.Info("KamaiZen stopped")

	state := analysis.NewState()
	analyser_channel := make(chan analysis.State, CHAN_BUFFER)

	var wg sync.WaitGroup
	wg.Add(3)
	go analysis.StartAnalyser(analyser_channel, &wg)
	go server.StartServer(&wg, state, analyser_channel)
	go lsp.Start(&wg)
	wg.Wait()
	close(analyser_channel)
}
