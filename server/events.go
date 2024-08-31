package server

import (
	"KamaiZen/logger"
	"KamaiZen/lsp"
	"KamaiZen/state_manager"
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
		contents []byte,
		analyser_channel chan state_manager.State)
}

// NewEventManager creates and returns a new EventManager instance.
func NewEventManager() *EventManager {
	return &EventManager{
		handlers: make(map[string]func(contents []byte, analyser_channel chan state_manager.State)),
	}
}

// RegisterHandler registers a handler function for a specific method.
// method: The name of the method for which the handler is being registered.
// handler: The function to handle the event. It takes an state_manager.State, the contents as a byte slice, and a channel for state_manager.State.
func (em *EventManager) RegisterHandler(method string, handler func(contents []byte, analyser_channel chan state_manager.State)) {
	logger.Infof("Registering handler for method: %s", method)
	em.handlers[method] = handler
}

// Dispatch calls the registered handler for the given method.
// method: The name of the method for which the handler is being dispatched.
// state: The current state of the state_manager.
// contents: The contents to be passed to the handler as a byte slice.
// analyser_channel: A channel for state_manager.State to communicate with the handler.
func (em *EventManager) Dispatch(method string, contents []byte, analyser_channel chan state_manager.State) {
	if handler, found := em.handlers[method]; found {
		handler(contents, analyser_channel)
	} else {
		logger.Errorf("No handler found for method: %s", method)
	}
}

// handleInitialized handles the 'initialized' notification.
// state: The current state of the state_manager.
// contents: The contents of the notification as a byte slice.
// analyser_channel: A channel for state_manager.State to communicate with the handler.
func handleInitialized(contents []byte, analyser_channel chan state_manager.State) {
	var notification lsp.InitializedNotification
	if error := json.Unmarshal(contents, &notification); error != nil {
		logger.Error("Error unmarshalling initialized notfication: ", error)
		return
	}
	logger.Infof("Received initialize request with %v", notification)
}

// handleInitialize handles the 'initialize' request.
// state: The current state of the state_manager.
// contents: The contents of the request as a byte slice.
// analyser_channel: A channel for state_manager.State to communicate with the handler.
func handleInitialize(contents []byte, analyser_channel chan state_manager.State) {
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
// state: The current state of the state_manager.
// contents: The contents of the notification as a byte slice.
// analyser_channel: A channel for state_manager.State to communicate with the handler.
func handleDidOpen(contents []byte, analyser_channel chan state_manager.State) {
	var notification lsp.DidOpenTextDocumentNotification
	if error := json.Unmarshal(contents, &notification); error != nil {
		logger.Error("Error unmarshalling didOpen notification: ", error)
		return
	}
	logger.Info("Opened document with URI: ", notification.Params.TextDocument.URI)
	dignostics := state_manager.GetState().OpenDocument(notification.Params.TextDocument.URI, notification.Params.TextDocument.Text)
	if len(dignostics) > 0 {
		lsp.WriteResponse(lsp.NewPublishDiagnosticNotification(notification.Params.TextDocument.URI, dignostics))
	}
}

// handleMessage handles incoming messages and dispatches them to the appropriate handler.
// state: The current state of the state_manager.
// method: The name of the method for which the handler is being dispatched.
// contents: The contents of the message as a byte slice.
// analyser_channel: A channel for state_manager.State to communicate with the handler.
// eventManager: The EventManager instance to use for dispatching the message.
func handleMessage(method string, contents []byte, analyser_channel chan state_manager.State, eventManager *EventManager) {
	logger.Info("Received message with method: ", method)
	eventManager.Dispatch(method, contents, analyser_channel)
}

// handleDidChange handles the 'didChange' notification.
// state: The current state of the state_manager.
// contents: The contents of the notification as a byte slice.
// analyser_channel: A channel for state_manager.State to communicate with the handler.
func handleDidChange(contents []byte, analyser_channel chan state_manager.State) {
	var notification lsp.DidChangeTextDocumentNotification
	state := state_manager.GetState()
	if error := json.Unmarshal(contents, &notification); error != nil {
		logger.Error("Error unmarshalling didChange notification: ", error)
		return
	}
	for _, change := range notification.Params.ContentChanges {
		diagnostics := state.UpdateDocument(notification.Params.TextDocument.URI, change.Text)
		if len(diagnostics) > 0 {
			lsp.WriteResponse(lsp.NewPublishDiagnosticNotification(notification.Params.TextDocument.URI, diagnostics))
		}
	}
	// TODO: remove analyser_channel
	// analyser_channel <- *state
}

// handleHover handles the 'hover' request.
// state: The current state of the state_manager.
// contents: The contents of the request as a byte slice.
// analyser_channel: A channel for state_manager.State to communicate with the handler.
func handleHover(contents []byte, analyser_channel chan state_manager.State) {
	_ = analyser_channel
	var request lsp.HoverRequest
	if error := json.Unmarshal(contents, &request); error != nil {
		logger.Error("Error unmarshalling hover request: ", error)
		return
	}
	logger.Debug("Hover request for document with URI: ", request.Params.TextDocument.URI)
	logger.Debug("Position: ", request.Params.Position)
	response := state_manager.GetState().Hover(request.ID, request.Params.TextDocument.URI, request.Params.Position)
	logger.Infof("Sent hover response %v", response)
	lsp.WriteResponse(response)
}

// handleDefinition handles the 'definition' request.
// state: The current state of the state_manager.
// contents: The contents of the request as a byte slice.
// analyser_channel: A channel for state_manager.State to communicate with the handler.
func handleDefinition(contents []byte, analyser_channel chan state_manager.State) {
	_ = analyser_channel
	var request lsp.DefinitionProviderRequest
	if error := json.Unmarshal(contents, &request); error != nil {
		logger.Error("Error unmarshalling definition request: ", error)
		return
	}
	logger.Debug("Definition request for document with URI: ", request.Params.TextDocument.URI)
	logger.Debug("Position: ", request.Params.Position)
	response := state_manager.GetState().Definition(request.ID, request.Params.TextDocument.URI, request.Params.Position)
	logger.Debug("Sent definition response %v", response)
	lsp.WriteResponse(response)
}

// handleFormatting handles the 'formatting' request.
// state: The current state of the state_manager.
// contents: The contents of the request as a byte slice.
// analyser_channel: A channel for state_manager.State to communicate with the handler.
func handleFormatting(contents []byte, analyser_channel chan state_manager.State) {
	_ = analyser_channel
	var request lsp.DocumentFormattingRequest
	if error := json.Unmarshal(contents, &request); error != nil {
		logger.Error("Error unmarshalling formatting request: ", error)
		return
	}
	logger.Debug("Formatting request for document with URI: ", request.Params.TextDocument.URI)
	response := state_manager.GetState().Formatting(request.ID, request.Params.TextDocument.URI, request.Params.Options)
	logger.Debug("Sent formatting response %v", response)
	lsp.WriteResponse(response)
}

// handleCompletion handles the 'completion' request.
// state: The current state of the state_manager.
// contents: The contents of the request as a byte slice.
// analyser_channel: A channel for state_manager.State to communicate with the handler.
func handleCompletion(contents []byte, analyser_channel chan state_manager.State) {
	_ = analyser_channel
	var request lsp.CompletionRequest
	if error := json.Unmarshal(contents, &request); error != nil {
		logger.Error("Error unmarshalling completion request: ", error)
		return
	}
	logger.Debug("Completion request for document with URI: ", request.Params.TextDocument.URI)
	response := state_manager.GetState().TextDocumentCompletion(request.ID, request.Params.TextDocument.URI, request.Params.Position)
	logger.Debug("Sent completion response %v", response)
	lsp.WriteResponse(response)
}
