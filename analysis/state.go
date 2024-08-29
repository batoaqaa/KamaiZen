package analysis

import (
	"KamaiZen/document_manager"
	"KamaiZen/kamailio_cfg"
	"KamaiZen/logger"
	"KamaiZen/lsp"
	"fmt"
)

type State struct {
	Documents map[lsp.DocumentURI]string
	Analyzer  *kamailio_cfg.Analyzer
}

var state State

func GetState() *State {
	return &state
}

func SetState(s State) {
	state = s
}

func InitializeState() {
	state = NewState()
	state.Analyzer = kamailio_cfg.NewAnalyzer()
	logger.Debug("State initialized")
}

func NewState() State {
	return State{
		Documents: make(map[lsp.DocumentURI]string),
	}
}

// register the subscribers for the events
// subscriber may include the Parser to update the ASY

func (s *State) NotifySubsrcibers(uri lsp.DocumentURI, changes []lsp.TextDocumentContentChangeEvent) {
	// notify the subscribers of the changes
	// the parser and the analyser will be listening to the changes
	// analyzer_channel <- s.ChangeDocument(uri, changes)
}

func (s *State) OpenDocument(uri lsp.DocumentURI, text string) []lsp.Diagnostic {
	s.Documents[uri] = text
	s.Analyzer.Initialize([]byte(text))
	visitor := kamailio_cfg.NewDiagnosticVisitor()
	s.Analyzer.GetAST().Accept(visitor)
	return visitor.GetDiagnostics()
}

func (s *State) ChangeDocument(uri lsp.DocumentURI, changes []lsp.TextDocumentContentChangeEvent) []lsp.Diagnostic {
	text := s.Documents[uri]
	for _, change := range changes {
		change.Apply(text)
	}
	s.Documents[uri] = text
	// NOTE: Make changes for incremental parsing
	s.Analyzer.Update([]byte(text))
	visitor := kamailio_cfg.NewDiagnosticVisitor()
	s.Analyzer.GetAST().Accept(visitor)
	return visitor.GetDiagnostics()
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
	visitor := kamailio_cfg.NewFormattingVisitor()
	s.Analyzer.GetAST().Accept(visitor)
	edits := visitor.GetEdits()
	return lsp.NewDocumentFormattingResponse(id, edits)
}
