package main

import (
	"KamaiZen/logger"
	"KamaiZen/lsp"
	"KamaiZen/server"
	"KamaiZen/state_manager"
	"sync"
)

func initialise() {
	logger.Info("Starting KamaiZen")
	state_manager.InitializeState()
	lsp.Initialise()
}

func main() {
	initialise()
	defer logger.Info("KamaiZen stopped")
	server := server.GetServerInstance()

	var wg sync.WaitGroup
	wg.Add(2)
	go server.StartServer(&wg)
	go lsp.Start(&wg)
	wg.Wait()
}
