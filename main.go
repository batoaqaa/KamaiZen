package main

import (
	"KamaiZen/analysis"
	"KamaiZen/docs"
	"KamaiZen/lsp"
	"KamaiZen/rpc"
	"KamaiZen/settings"
	"KamaiZen/utils"
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"
)

const (
	MethodInitialize = "initialize"
	MethodDidOpen    = "textDocument/didOpen"
	MethodDidChange  = "textDocument/didChange"
	MethodHover      = "textDocument/hover"
	MethodDefinition = "textDocument/definition"
	MethodFormatting = "textDocument/formatting"
)

func main() {
	logger := utils.GetLogger()
	logger.Println("Starting KamaiZen")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)
	state := analysis.NewState()
	writer := os.Stdout

	settings := settings.NewLSPSettings("/home/ibrahim/work/kamailio", "/path/to/root")
	docs.Initialise(settings)

	analyser_channel := make(chan analysis.State)

	var wg sync.WaitGroup
	wg.Add(1)

	go analysis.StartAnalyser(analyser_channel, writer, logger, &wg)

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, error := rpc.DecodeMessage(msg)
		if error != nil {
			logger.Println("Error decoding message: ", error)
			continue
		}
		handleMessage(writer, logger, state, method, contents, analyser_channel)
	}
	wg.Wait()
	close(analyser_channel)
	defer logger.Println("KamaiZen stopped")
}

func handleMessage(writer io.Writer, logger *log.Logger, state analysis.State, method string, contents []byte, analyser_channel chan analysis.State) {
	logger.Println("Received message with method: ", method)
	switch method {

	case MethodInitialize:
		var request lsp.InitializeRequest
		if error := json.Unmarshal(contents, &request); error != nil {
			logger.Println("Error unmarshalling initialize request: ", error)
			return
		}
		logger.Printf("Connected to %s with version %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)
		response := lsp.NewInitializeResponse(request.ID)
		lsp.WriteResponse(writer, response)
		logger.Println("Sent initialize response")

	case MethodDidOpen:
		var notification lsp.DidOpenTextDocumentNotification
		if error := json.Unmarshal(contents, &notification); error != nil {
			logger.Println("Error unmarshalling didOpen notification: ", error)
			return
		}
		logger.Println("Opened document with URI: ", notification.Params.TextDocument.URI)
		// logger.Println("Document content: ", notification.Params.TextDocument.Text)
		state.OpenDocument(notification.Params.TextDocument.URI, notification.Params.TextDocument.Text)
		analyser_channel <- state

	case MethodDidChange:
		var notification lsp.DidChangeTextDocumentNotification
		if error := json.Unmarshal(contents, &notification); error != nil {
			logger.Println("Error unmarshalling didChange notification: ", error)
			return
		}
		for _, change := range notification.Params.ContentChanges {
			state.UpdateDocument(notification.Params.TextDocument.URI, change.Text)
		}
		analyser_channel <- state

	case MethodHover:
		var request lsp.HoverRequest
		if error := json.Unmarshal(contents, &request); error != nil {
			logger.Println("Error unmarshalling hover request: ", error)
			return
		}
		logger.Println("Hover request for document with URI: ", request.Params.TextDocument.URI)
		logger.Println("Position: ", request.Params.Position)
		response := state.Hover(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		logger.Printf("Sent hover response %v", response)
		lsp.WriteResponse(writer, response)

	case MethodDefinition:
		var request lsp.DefinitionProviderRequest
		if error := json.Unmarshal(contents, &request); error != nil {
			logger.Println("Error unmarshalling definition request: ", error)
			return
		}
		logger.Println("Definition request for document with URI: ", request.Params.TextDocument.URI)
		logger.Println("Position: ", request.Params.Position)
		response := state.Definition(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		logger.Printf("Sent definition response %v", response)
		lsp.WriteResponse(writer, response)
	case MethodFormatting:
		var request lsp.DocumentFormattingRequest
		if error := json.Unmarshal(contents, &request); error != nil {
			logger.Println("Error unmarshalling formatting request: ", error)
			return
		}
		logger.Println("Formatting request for document with URI: ", request)
		response := state.Formatting(request.ID, request.Params.TextDocument.URI, request.Params.Options)
		// logger.Printf("Sent formatting response %v", response)
		lsp.WriteResponse(writer, response)
	}

}
