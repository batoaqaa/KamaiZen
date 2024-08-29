package server

import (
	"KamaiZen/analysis"
	"KamaiZen/logger"
	"KamaiZen/lsp"
	"encoding/json"
)

const (
	MethodInitialize  = "initialize"
	MethodInitialized = "initialized"
	MethodDidOpen     = "textDocument/didOpen"
	MethodDidChange   = "textDocument/didChange"
	MethodHover       = "textDocument/hover"
	MethodDefinition  = "textDocument/definition"
	MethodFormatting  = "textDocument/formatting"
	MethodCompletion  = "textDocument/completion"
)

// EventManager manages event handlers for different methods.
type EventManager struct {
	handlers map[string]func(
		state analysis.State,
		contents []byte,
		analyser_channel chan analysis.State)
}

// NewEventManager creates and returns a new EventManager instance.
func NewEventManager() *EventManager {
	return &EventManager{
		handlers: make(map[string]func(state analysis.State, contents []byte, analyser_channel chan analysis.State)),
	}
}

// RegisterHandler registers a handler function for a specific method.
// method: The name of the method for which the handler is being registered.
// handler: The function to handle the event. It takes an analysis.State, the contents as a byte slice, and a channel for analysis.State.
func (em *EventManager) RegisterHandler(method string, handler func(state analysis.State, contents []byte, analyser_channel chan analysis.State)) {
	logger.Infof("Registering handler for method: %s", method)
	em.handlers[method] = handler
}

// Dispatch calls the registered handler for the given method.
// method: The name of the method for which the handler is being dispatched.
// state: The current state of the analysis.
// contents: The contents to be passed to the handler as a byte slice.
// analyser_channel: A channel for analysis.State to communicate with the handler.
func (em *EventManager) Dispatch(method string, state analysis.State, contents []byte, analyser_channel chan analysis.State) {
	if handler, found := em.handlers[method]; found {
		handler(state, contents, analyser_channel)
	} else {
		logger.Errorf("No handler found for method: %s", method)
	}
}

// handleInitialized handles the 'initialized' notification.
// state: The current state of the analysis.
// contents: The contents of the notification as a byte slice.
// analyser_channel: A channel for analysis.State to communicate with the handler.
func handleInitialized(state analysis.State, contents []byte, analyser_channel chan analysis.State) {
	var notification lsp.InitializedNotification
	if error := json.Unmarshal(contents, &notification); error != nil {
		logger.Error("Error unmarshalling initialized notfication: ", error)
		return
	}
	logger.Infof("Received initialize request with %v", notification)
}

// handleInitialize handles the 'initialize' request.
// state: The current state of the analysis.
// contents: The contents of the request as a byte slice.
// analyser_channel: A channel for analysis.State to communicate with the handler.
func handleInitialize(state analysis.State, contents []byte, analyser_channel chan analysis.State) {
	var request lsp.InitializeRequest
	if error := json.Unmarshal(contents, &request); error != nil {
		logger.Error("Error unmarshalling initialize request: ", error)
		return
	}
	logger.Infof("Connected to %s with version %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)
	response := lsp.NewInitializeResponse(request.ID)
	lsp.WriteResponse(response)
	logger.Debug("Sent initialize response")
}

// handleDidOpen handles the 'didOpen' notification.
// state: The current state of the analysis.
// contents: The contents of the notification as a byte slice.
// analyser_channel: A channel for analysis.State to communicate with the handler.
func handleDidOpen(state analysis.State, contents []byte, analyser_channel chan analysis.State) {
	var notification lsp.DidOpenTextDocumentNotification
	if error := json.Unmarshal(contents, &notification); error != nil {
		logger.Error("Error unmarshalling didOpen notification: ", error)
		return
	}
	logger.Info("Opened document with URI: ", notification.Params.TextDocument.URI)
	state.OpenDocument(notification.Params.TextDocument.URI, notification.Params.TextDocument.Text)
	analyser_channel <- state
}

// handleMessage handles incoming messages and dispatches them to the appropriate handler.
// state: The current state of the analysis.
// method: The name of the method for which the handler is being dispatched.
// contents: The contents of the message as a byte slice.
// analyser_channel: A channel for analysis.State to communicate with the handler.
// eventManager: The EventManager instance to use for dispatching the message.
func handleMessage(state analysis.State, method string, contents []byte, analyser_channel chan analysis.State, eventManager *EventManager) {
	logger.Info("Received message with method: ", method)
	eventManager.Dispatch(method, state, contents, analyser_channel)
}

// handleDidChange handles the 'didChange' notification.
// state: The current state of the analysis.
// contents: The contents of the notification as a byte slice.
// analyser_channel: A channel for analysis.State to communicate with the handler.
func handleDidChange(state analysis.State, contents []byte, analyser_channel chan analysis.State) {
	var notification lsp.DidChangeTextDocumentNotification
	if error := json.Unmarshal(contents, &notification); error != nil {
		logger.Error("Error unmarshalling didChange notification: ", error)
		return
	}
	for _, change := range notification.Params.ContentChanges {
		state.UpdateDocument(notification.Params.TextDocument.URI, change.Text)
	}
	// TODO: remove analyser_channel
	analyser_channel <- state
}

// handleHover handles the 'hover' request.
// state: The current state of the analysis.
// contents: The contents of the request as a byte slice.
// analyser_channel: A channel for analysis.State to communicate with the handler.
func handleHover(state analysis.State, contents []byte, analyser_channel chan analysis.State) {
	_ = analyser_channel
	var request lsp.HoverRequest
	if error := json.Unmarshal(contents, &request); error != nil {
		logger.Error("Error unmarshalling hover request: ", error)
		return
	}
	logger.Debug("Hover request for document with URI: ", request.Params.TextDocument.URI)
	logger.Debug("Position: ", request.Params.Position)
	response := state.Hover(request.ID, request.Params.TextDocument.URI, request.Params.Position)
	logger.Infof("Sent hover response %v", response)
	lsp.WriteResponse(response)
}

// handleDefinition handles the 'definition' request.
// state: The current state of the analysis.
// contents: The contents of the request as a byte slice.
// analyser_channel: A channel for analysis.State to communicate with the handler.
func handleDefinition(state analysis.State, contents []byte, analyser_channel chan analysis.State) {
	_ = analyser_channel
	var request lsp.DefinitionProviderRequest
	if error := json.Unmarshal(contents, &request); error != nil {
		logger.Error("Error unmarshalling definition request: ", error)
		return
	}
	logger.Debug("Definition request for document with URI: ", request.Params.TextDocument.URI)
	logger.Debug("Position: ", request.Params.Position)
	response := state.Definition(request.ID, request.Params.TextDocument.URI, request.Params.Position)
	logger.Debug("Sent definition response %v", response)
	lsp.WriteResponse(response)
}

// handleFormatting handles the 'formatting' request.
// state: The current state of the analysis.
// contents: The contents of the request as a byte slice.
// analyser_channel: A channel for analysis.State to communicate with the handler.
func handleFormatting(state analysis.State, contents []byte, analyser_channel chan analysis.State) {
	_ = analyser_channel
	var request lsp.DocumentFormattingRequest
	if error := json.Unmarshal(contents, &request); error != nil {
		logger.Error("Error unmarshalling formatting request: ", error)
		return
	}
	logger.Debug("Formatting request for document with URI: ", request.Params.TextDocument.URI)
	response := state.Formatting(request.ID, request.Params.TextDocument.URI, request.Params.Options)
	logger.Debug("Sent formatting response %v", response)
	lsp.WriteResponse(response)
}

// handleCompletion handles the 'completion' request.
// state: The current state of the analysis.
// contents: The contents of the request as a byte slice.
// analyser_channel: A channel for analysis.State to communicate with the handler.
func handleCompletion(state analysis.State, contents []byte, analyser_channel chan analysis.State) {
	_ = analyser_channel
	var request lsp.CompletionRequest
	if error := json.Unmarshal(contents, &request); error != nil {
		logger.Error("Error unmarshalling completion request: ", error)
		return
	}
	logger.Debug("Completion request for document with URI: ", request.Params.TextDocument.URI)
	response := state.TextDocumentCompletion(request.ID, request.Params.TextDocument.URI, request.Params.Position)
	logger.Debug("Sent completion response %v", response)
	lsp.WriteResponse(response)
}
