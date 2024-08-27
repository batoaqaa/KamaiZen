package lsp

import "KamaiZen/logger"

type DidChangeTextDocumentNotification struct {
	Notification
	Params DidChangeTextDocumentParams `json:"params"`
}

type DidChangeTextDocumentParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

type TextDocumentContentChangeEvent struct {
	Range *Range `json:"range,omitempty"`
	Text  string `json:"text"`
}

func (change TextDocumentContentChangeEvent) Apply(text string) string {
	if change.Range != nil {
		start := change.Range.Start
		end := change.Range.End
		logger.Debug("Change range: ", start, end)
		text = text[:start.Line] + change.Text + text[end.Line:]
	} else {
		text = change.Text
	}
	return text
}
