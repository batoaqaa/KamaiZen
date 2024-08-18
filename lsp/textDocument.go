package lsp

type TextDocumentItem struct {
	/**
	 * The text document's URI.
	 */
	URI DocumentURI `json:"uri"`
	/**
	 * The text document's language identifier.
	 */
	LanguageId string `json:"languageId"`
	/**
	 * The version number of this document (it will strictly increase after each
	 * change, including undo/redo).
	 */
	Version int `json:"version"`
	/**
	 * The content of the opened text document.
	 */
	Text string `json:"text"`
}

type DocumentURI string

type Range struct {
	/**
	 * The range's start position.
	 */
	Start Position `json:"start"`
	/**
	 * The range's end position.
	 */
	End Position `json:"end"`
}

type Position struct {
	/**
	 * Line position in a document (zero-based).
	 */
	Line int `json:"line"`
	/**
	 * Character offset on a line in a document (zero-based).
	 */
	Character int `json:"character"`
}

type TextDocumentIdentifier struct {
	/**
	 * The text document's URI.
	 */
	URI DocumentURI `json:"uri"`
}

type VersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier
	Version int `json:"version"`
}

type TextDocuemntPositionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

type MarkupContent struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

type Location struct {
	URI   DocumentURI `json:"uri"`
	Range Range       `json:"range"`
}
