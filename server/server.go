package server

import (
	"KamaiZen/analysis"
	"KamaiZen/logger"
	"KamaiZen/rpc"
	"bufio"
	"os"
	"sync"
)

func StartServer(wg *sync.WaitGroup, state analysis.State, analyser_channel chan analysis.State) {
	defer wg.Done()
	scanner := bufio.NewScanner(os.Stdin)
	logger.Info("Starting server")
	scanner.Split(rpc.Split)

	// Initialize EventManager and register handlers
	eventManager := NewEventManager()
	eventManager.RegisterHandler(MethodInitialize, handleInitialize)
	eventManager.RegisterHandler(MethodInitialized, handleInitialized)
	eventManager.RegisterHandler(MethodDidOpen, handleDidOpen)
	eventManager.RegisterHandler(MethodDidChange, handleDidChange)
	eventManager.RegisterHandler(MethodHover, handleHover)
	eventManager.RegisterHandler(MethodDefinition, handleDefinition)
	eventManager.RegisterHandler(MethodFormatting, handleFormatting)
	eventManager.RegisterHandler(MethodCompletion, handleCompletion)

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, error := rpc.DecodeMessage(msg)
		if error != nil {
			logger.Error("Error decoding message: ", error)
			continue
		}
		handleMessage(state, method, contents, analyser_channel, eventManager)
	}
}
