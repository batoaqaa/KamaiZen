package server

import (
	"KamaiZen/analysis"
	"KamaiZen/logger"
	"KamaiZen/lsp"
	"KamaiZen/rpc"
	"bufio"
	"encoding/json"
	"os"
	"sync"
)

func StartServer(wg *sync.WaitGroup, state analysis.State, analyser_channel chan analysis.State) {
	defer wg.Done()
	scanner := bufio.NewScanner(os.Stdin)
	logger.Info("Starting server")
	scanner.Split(rpc.Split)
	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, error := rpc.DecodeMessage(msg)
		if error != nil {
			logger.Error("Error decoding message: ", error)
			continue
		}
		handleMessage(state, method, contents, analyser_channel)
	}
}

func handleMessage(state analysis.State, method string, contents []byte, analyser_channel chan analysis.State) {
	logger.Info("Received message with method: ", method)
	switch method {

	case MethodInitialize:
		var request lsp.InitializeRequest
		if error := json.Unmarshal(contents, &request); error != nil {
			logger.Error("Error unmarshalling initialize request: ", error)
			return
		}
		logger.Infof("Connected to %s with version %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)
		response := lsp.NewInitializeResponse(request.ID)
		lsp.WriteResponse(response)
		logger.Debug("Sent initialize response")

	case MethodDidOpen:
		var notification lsp.DidOpenTextDocumentNotification
		if error := json.Unmarshal(contents, &notification); error != nil {
			logger.Error("Error unmarshalling didOpen notification: ", error)
			return
		}
		logger.Info("Opened document with URI: ", notification.Params.TextDocument.URI)
		state.OpenDocument(notification.Params.TextDocument.URI, notification.Params.TextDocument.Text)
		analyser_channel <- state

	case MethodDidChange:
		var notification lsp.DidChangeTextDocumentNotification
		if error := json.Unmarshal(contents, &notification); error != nil {
			logger.Error("Error unmarshalling didChange notification: ", error)
			return
		}
		for _, change := range notification.Params.ContentChanges {
			state.UpdateDocument(notification.Params.TextDocument.URI, change.Text)
		}
		analyser_channel <- state

	case MethodHover:
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

	case MethodDefinition:
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
	case MethodFormatting:
		var request lsp.DocumentFormattingRequest
		if error := json.Unmarshal(contents, &request); error != nil {
			logger.Error("Error unmarshalling formatting request: ", error)
			return
		}
		logger.Debug("Formatting request for document with URI: ", request)
		response := state.Formatting(request.ID, request.Params.TextDocument.URI, request.Params.Options)
		lsp.WriteResponse(response)
	}

}
