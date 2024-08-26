package main

import (
	"KamaiZen/analysis"
	"KamaiZen/document_manager"
	"KamaiZen/logger"
	"KamaiZen/server"
	"KamaiZen/settings"
	"sync"
)

func initialise() {
	logger.Info("Starting KamaiZen")
	settings := settings.NewLSPSettings("/home/ibrahim/work/kamailio", "/path/to/root")
	document_manager.Initialise(settings)
}

func main() {
	logger.Info("Starting KamaiZen")
	state := analysis.NewState()
	initialise()
	analyser_channel := make(chan analysis.State)
	var wg sync.WaitGroup
	wg.Add(1)
	go server.StartServer(&wg, state, analyser_channel)
	wg.Add(1)
	go analysis.StartAnalyser(analyser_channel, &wg)
	wg.Wait()
	close(analyser_channel)
	defer logger.Info("KamaiZen stopped")
}
