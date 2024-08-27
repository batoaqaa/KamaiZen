package lsp

const (
	TEXT_DOCUMENT_SYNC_KIND_NONE = iota
	TEXT_DOCUMENT_SYNC_KIND_FULL
	TEXT_DOCUMENT_SYNC_KIND_INCREMENTAL
)

type TextDocumentSyncOptions struct {
	OpenClose bool `json:"openClose,omitempty"`
	Change    int  `json:"change"`
}
