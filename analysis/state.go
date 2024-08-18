package analysis

import (
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

func (s *State) OpenDocument(uri lsp.DocumentURI, text string) {
	s.Documents[uri] = text
	// start the analysis of the text document here
}

func (s *State) ChangeDocument(uri lsp.DocumentURI, changes []lsp.TextDocumentContentChangeEvent) {
	text := s.Documents[uri]
	for _, change := range changes {
		text = change.Apply(text)
	}
	s.Documents[uri] = text
}

func (s *State) Hover(id int, uri lsp.DocumentURI, position lsp.Position) lsp.HoverResponse {
	// TODO: Implement hover
	// lookup the text content type from the analysis and return the documentation
	// for function or variable at the position in the document
	text := s.Documents[uri]
	_ = text
	// contents := text[position.Line]
	// change content from byte to string
	return lsp.NewHoverResponse(id, fmt.Sprintf("File: %s :: Hovering at line %d, character %d", uri, position.Line, position.Character))
}

func (s *State) Definition(id int, uri lsp.DocumentURI, position lsp.Position) lsp.DefinitionProviderResponse {
	// TODO: Implement definition
	// lookup the text content type from the analysis and return the definition
	return lsp.NewDefintionProviderResponse(id, fmt.Sprintf("File: %s :: Definition at line %d, character %d", uri, position.Line, position.Character))
}
