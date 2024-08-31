package server

import (
	"KamaiZen/logger"
	"KamaiZen/rpc"
	"KamaiZen/state_manager"
	"bufio"
	"os"
	"sync"
)

// StartServer starts the language server and listens for incoming messages from the client.
// It initializes the event manager, registers handlers for various methods, and processes incoming messages.
//
// Parameters:
//
//	wg *sync.WaitGroup - The wait group to signal when the server is done.
//	analyser_channel chan state_manager.State - The channel for communicating with the state manager.
func StartServer(wg *sync.WaitGroup, analyser_channel chan state_manager.State) {
	// TODO: Get rid of analyser_channel
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
		handleMessage(method, contents, analyser_channel, eventManager)
	}
}
