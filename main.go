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
	logger.Info("Starting KamaiZen")
	settings := settings.NewLSPSettings("/home/ibrahim/work/kamailio", "/path/to/root")
	lsp.Initialise()
	document_manager.Initialise(settings)
}

func main() {
	logger.Info("Starting KamaiZen")
	state := analysis.NewState()
	initialise()
	// make buffered channel
	analyser_channel := make(chan analysis.State, CHAN_BUFFER)

	var wg sync.WaitGroup
	wg.Add(3)
	go analysis.StartAnalyser(analyser_channel, &wg)
	go server.StartServer(&wg, state, analyser_channel)
	go lsp.Start(&wg)
	wg.Wait()
	close(analyser_channel)
	defer logger.Info("KamaiZen stopped")
}
