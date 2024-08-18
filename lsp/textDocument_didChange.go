package lsp

type DidChangeTextDocumentNotification struct {
	Notification
	Params DidChangeTextDocumentParams `json:"params"`
}

type DidChangeTextDocumentParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

type TextDocumentContentChangeEvent struct {
	// Range       *Range `json:"range,omitempty"`
	// RangeLength int    `json:"rangeLength,omitempty"`
	Text string `json:"text"`
}

func (change TextDocumentContentChangeEvent) Apply(text string) string {
	return text + change.Text
}
