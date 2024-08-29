package analysis

import (
	"KamaiZen/document_manager"
	"KamaiZen/logger"
	"KamaiZen/lsp"
	"fmt"
)

type State struct {
	// Key: URI, Value: Text content
	Documents map[lsp.DocumentURI]string
}

func NewState() State {
	return State{
		Documents: make(map[lsp.DocumentURI]string),
	}
}

// register the subscribers for the events
// subscriber may include the Parser to update the ASY

func (s *State) OpenDocument(uri lsp.DocumentURI, text string) []lsp.Diagnostic {
	s.Documents[uri] = text
	// start the analysis of the text document here
	diagnostic := GetDiagnosticsForDocument(uri, text)
	return diagnostic
}

func (s *State) ChangeDocument(uri lsp.DocumentURI, changes []lsp.TextDocumentContentChangeEvent) []lsp.Diagnostic {
	text := s.Documents[uri]
	for _, change := range changes {

		change.Apply(text)
	}
	s.Documents[uri] = text
	// notify the subscribers of the changes
	// the parser and the analyser will be listening to the changes
	// FIXME: Make changes for incremental parsing
	diagnostic := GetDiagnosticsForDocument(uri, text)
	return diagnostic
}

func (s *State) UpdateDocument(uri lsp.DocumentURI, text string) []lsp.Diagnostic {
	s.Documents[uri] = text
	diagnostic := GetDiagnosticsForDocument(uri, text)
	return diagnostic
}

func (s *State) Hover(id int, uri lsp.DocumentURI, position lsp.Position) lsp.HoverResponse {
	text := s.Documents[uri]
	functionName := GetFunctionNameAtPosition(uri, position, []byte(text))
	documentation := document_manager.FindFunctionInAllModules(functionName)
	return lsp.NewHoverResponse(id, fmt.Sprintf("%s", documentation))
}

func (s *State) Definition(id int, uri lsp.DocumentURI, position lsp.Position) lsp.DefinitionProviderResponse {
	// TODO: Implement definition
	// lookup the text content type from the analysis and return the definition
	return lsp.NewDefintionProviderResponse(id, fmt.Sprintf("File: %s :: Definition at line %d, character %d", uri, position.Line, position.Character))
}

func (s *State) TextDocumentCompletion(id int, uri lsp.DocumentURI, position lsp.Position) lsp.CompletionResponse {
	logger.Debug("Completion request for document with URI: ", uri)
	items := GetCompletionItems(uri)
	return lsp.NewCompletionResponse(id, items)
}

func (s *State) Formatting(id int, uri lsp.DocumentURI, options lsp.FormattingOptions) lsp.DocumentFormattingResponse {
	logger.Info("===?Formatting document with URI: ", uri)
	edits, error := IndentKamailioCfg(s.Documents[uri], 4)
	if error != nil {
		return lsp.NewDocumentFormattingResponse(id, []lsp.TextEdit{})
	}
	return lsp.NewDocumentFormattingResponse(id, edits)
}
