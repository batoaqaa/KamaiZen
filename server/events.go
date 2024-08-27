package server

import (
	"KamaiZen/analysis"
	"KamaiZen/logger"
	"KamaiZen/lsp"
	"encoding/json"
)

const (
	MethodInitialize = "initialize"
	MethodDidOpen    = "textDocument/didOpen"
	MethodDidChange  = "textDocument/didChange"
	MethodHover      = "textDocument/hover"
	MethodDefinition = "textDocument/definition"
	MethodFormatting = "textDocument/formatting"
	MethodCompletion = "textDocument/completion"
)

type EventManager struct {
	handlers map[string]func(state analysis.State, contents []byte, analyser_channel chan analysis.State)
}

func NewEventManager() *EventManager {
	return &EventManager{
		handlers: make(map[string]func(state analysis.State, contents []byte, analyser_channel chan analysis.State)),
	}
}

func (em *EventManager) RegisterHandler(method string, handler func(state analysis.State, contents []byte, analyser_channel chan analysis.State)) {
	logger.Infof("Registering handler for method: %s", method)
	em.handlers[method] = handler
}

func (em *EventManager) Dispatch(method string, state analysis.State, contents []byte, analyser_channel chan analysis.State) {
	if handler, found := em.handlers[method]; found {
		handler(state, contents, analyser_channel)
	} else {
		logger.Errorf("No handler found for method: %s", method)
	}
}

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

func handleMessage(state analysis.State, method string, contents []byte, analyser_channel chan analysis.State, eventManager *EventManager) {
	logger.Info("Received message with method: ", method)
	eventManager.Dispatch(method, state, contents, analyser_channel)
}

func handleDidChange(state analysis.State, contents []byte, analyser_channel chan analysis.State) {
	var notification lsp.DidChangeTextDocumentNotification
	if error := json.Unmarshal(contents, &notification); error != nil {
		logger.Error("Error unmarshalling didChange notification: ", error)
		return
	}
	for _, change := range notification.Params.ContentChanges {
		state.UpdateDocument(notification.Params.TextDocument.URI, change.Text)
	}
	analyser_channel <- state
}

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
