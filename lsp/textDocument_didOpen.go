package lsp

// DidOpenTextDocumentNotification represents a notification sent to the server
// when a text document is opened. It contains the notification metadata and the
// parameters for the open event.
type DidOpenTextDocumentNotification struct {
	Notification
	Params DidOpenTextDocumentParams `json:"params"`
}

// DidOpenTextDocumentParams contains the parameters for the DidOpenTextDocumentNotification.
// It includes the text document item that was opened.
type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}
