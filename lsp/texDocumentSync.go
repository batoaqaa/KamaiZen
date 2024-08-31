package lsp

const (
	TEXT_DOCUMENT_SYNC_KIND_NONE = iota
	TEXT_DOCUMENT_SYNC_KIND_FULL
	TEXT_DOCUMENT_SYNC_KIND_INCREMENTAL
)

// TextDocumentSyncOptions represents the options for text document synchronization.
// It includes settings for opening/closing documents and the type of change notifications.
//
// Fields:
//
//	OpenClose bool - Indicates whether the server should be notified when a text document is opened or closed.
//	Change int - Specifies the type of change notifications (e.g., full or incremental).
type TextDocumentSyncOptions struct {
	OpenClose bool `json:"openClose,omitempty"`
	Change    int  `json:"change"`
}
