package main

import (
	"KamaiZen/analysis"
	"KamaiZen/lsp"
	"KamaiZen/rpc"
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
)

func main() {
	// FIXME: This is a temporary solution
	logger := getLogger("/home/ibrahim/work/KamaiZen/kamaizen.log")
	logger.Println("Starting KamaiZen")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)
	state := analysis.NewState()
	writer := os.Stdout

	analyser_channel := make(chan analysis.State)

	var wg sync.WaitGroup
	wg.Add(1)

	go analysis.StartAnalyser(analyser_channel, logger, &wg)

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
		writeResponse(writer, response)
		logger.Println("Sent initialize response")

	case MethodDidOpen:
		var notification lsp.DidOpenTextDocumentNotification
		if error := json.Unmarshal(contents, &notification); error != nil {
			logger.Println("Error unmarshalling didOpen notification: ", error)
			return
		}
		logger.Println("Opened document with URI: ", notification.Params.TextDocument.URI)
		logger.Println("Document content: ", notification.Params.TextDocument.Text)
		state.OpenDocument(notification.Params.TextDocument.URI, notification.Params.TextDocument.Text)
		logger.Println("Sending to analyser Document content after open: ", state.Documents[notification.Params.TextDocument.URI])
		analyser_channel <- state

	case MethodDidChange:
		var notification lsp.DidChangeTextDocumentNotification
		if error := json.Unmarshal(contents, &notification); error != nil {
			logger.Println("Error unmarshalling didChange notification: ", error)
			return
		}
		logger.Println("Changed document with URI: ", notification.Params.TextDocument.URI)
		logger.Println("Document content: ", notification.Params.ContentChanges[0].Text)
		state.ChangeDocument(notification.Params.TextDocument.URI, notification.Params.ContentChanges)
		logger.Println("Document content after change: ", state.Documents[notification.Params.TextDocument.URI])

	case MethodHover:
		var request lsp.HoverRequest
		if error := json.Unmarshal(contents, &request); error != nil {
			logger.Println("Error unmarshalling hover request: ", error)
			return
		}
		logger.Println("Hover request for document with URI: ", request.Params.TextDocument.URI)
		logger.Println("Position: ", request.Params.Position)
		response := state.Hover(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		writeResponse(writer, response)
		logger.Printf("Sent hover response %s", response)

	case MethodDefinition:
		var request lsp.DefinitionProviderRequest
		if error := json.Unmarshal(contents, &request); error != nil {
			logger.Println("Error unmarshalling definition request: ", error)
			return
		}
		logger.Println("Definition request for document with URI: ", request.Params.TextDocument.URI)
		logger.Println("Position: ", request.Params.Position)
		response := state.Definition(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		writeResponse(writer, response)
		logger.Printf("Sent definition response %s", response)

	}

}

func writeResponse(writer io.Writer, response interface{}) {
	reply := rpc.EncodeMessage(response)
	writer.Write([]byte(reply))
}

func getLogger(filename string) *log.Logger {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return log.New(file, "[KamaiZen]", log.Ldate|log.Ltime|log.Lshortfile)
}
