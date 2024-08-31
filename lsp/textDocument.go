package lsp

// TextDocumentItem represents a text document in the language server.
// It includes the document's URI, language identifier, version, and text content.
type TextDocumentItem struct {
	URI        DocumentURI `json:"uri"`
	LanguageId string      `json:"languageId"`
	Version    int         `json:"version"`
	Text       string      `json:"text"`
}

// DocumentURI represents the URI of a document.
type DocumentURI string

// Range represents a range within a text document.
// It includes the start and end positions of the range.
type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// Position represents a position within a text document.
// It includes the line and character offset of the position.
type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

// TextDocumentIdentifier identifies a text document by its URI.
type TextDocumentIdentifier struct {
	URI DocumentURI `json:"uri"`
}

// VersionedTextDocumentIdentifier identifies a text document by its URI and version number.
type VersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier
	Version int `json:"version"`
}

// TextDocumentPositionParams represents parameters for a text document position request.
// It includes the text document identifier and the position within the document.
type TextDocuemntPositionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

// MarkupContent represents content with a specific markup kind.
// It includes the kind of markup and the content value.
type MarkupContent struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

// Location represents a location within a text document.
// It includes the document's URI and the range within the document.
type Location struct {
	URI   DocumentURI `json:"uri"`
	Range Range       `json:"range"`
}
