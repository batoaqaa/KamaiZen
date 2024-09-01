package server

import (
	"KamaiZen/document_manager"
	"KamaiZen/logger"
	"KamaiZen/rpc"
	"KamaiZen/settings"
	"bufio"
	"os"
	"sync"
)

type Server struct {
	eventManager *EventManager
}

// create a single instance of the server
var serverInstance *Server

// GetServerInstance returns the single instance of the server.
func GetServerInstance() *Server {
	if serverInstance == nil {
		serverInstance = &Server{
			eventManager: NewEventManager(),
		}
	}
	return serverInstance
}

// StartServer starts the language server and listens for incoming messages from the client.
// It initializes the event manager, registers handlers for various methods, and processes incoming messages.
//
// Parameters:
//
//	wg *sync.WaitGroup - The wait group to signal when the server is done.
//	analyser_channel chan state_manager.State - The channel for communicating with the state manager.
func (s *Server) StartServer(wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(os.Stdin)
	logger.Info("Starting server")
	scanner.Split(rpc.Split)

	// Initialize EventManager and register handlers
	s.RegisterDefaultHandlers()

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, error := rpc.DecodeMessage(msg)
		if error != nil {
			logger.Error("Error decoding message: ", error)
			continue
		}
		handleMessage(method, contents, s.eventManager)
	}
}

func (s *Server) RegisterDefaultHandlers() {
	s.RegisterHandler(MethodInitialize, handleInitialize)
	s.RegisterHandler(MethodInitialized, handleInitialized)
	s.RegisterHandler(MethodDidOpen, handleDidOpen)
	s.RegisterHandler(MethodDidChange, handleDidChange)
	s.RegisterHandler(MethodDefinition, handleDefinition)
	s.RegisterHandler(MethodFormatting, handleFormatting)
	s.RegisterHandler(MethodConfigurationResponse, handleWorkspaceConfiguration)
}

func (s *Server) StopServer() {
	logger.Info("Stopping server")
}

func (s *Server) RegisterHandler(method string, handler func(contents []byte)) {
	s.eventManager.RegisterHandler(method, handler)
}

func (s *Server) addKamailioMethods(settings settings.LSPSettings) {
	logger.Info("Kamailio src detected at: ", settings.KamailioSourcePath)
	logger.Info("Adding Hover and Completion methods")
	document_manager.Initialise(settings)
	s.RegisterHandler(MethodHover, handleHover)
	s.RegisterHandler(MethodCompletion, handleCompletion)
}
