package lsp

import (
	"KamaiZen/settings"

	sitter "github.com/smacker/go-tree-sitter"
)

// FormattingOptions represents the options for document formatting.
// It includes settings for tab size, whether to insert spaces, and whether to trim trailing whitespace and final newlines.
type FormattingOptions struct {
	TabSize                int  `json:"tabSize"`
	InsertSpaces           bool `json:"insertSpaces"`
	TrimTrailingWhitespace bool `json:"trimTrailingWhitespace"`
	InsertFinalNewline     bool `json:"insertFinalNewline"`
	TrimFinalNewlines      bool `json:"trimFinalNewlines"`
}

// DocumentFormattingParams contains the parameters for a document formatting request.
// It includes the text document identifier and the formatting options.
type DocumentFormattingParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Options      FormattingOptions      `json:"options"`
}

// DocumentFormattingRequest represents a request for document formatting.
// It contains the request metadata and the parameters for the formatting request.
type DocumentFormattingRequest struct {
	Request
	Params DocumentFormattingParams `json:"params"`
}

// DocumentFormattingResponse represents the response to a document formatting request.
// It contains the response metadata and the list of text edits to be applied.
type DocumentFormattingResponse struct {
	Response
	Result []TextEdit `json:"result"`
}

// TextEdit represents a single text edit to be applied to a document.
// It includes the range of the edit and the new text to be inserted.
type TextEdit struct {
	Range   Range  `json:"range"`
	NewText string `json:"newText"`
}

// NewDocumentFormattingResponse creates and returns a new DocumentFormattingResponse.
// It initializes the response with the given ID and the list of text edits.
//
// Parameters:
//
//	id int - The ID of the response.
//	edits []TextEdit - The list of text edits.
//
// Returns:
//
//	DocumentFormattingResponse - The initialized response.
func NewDocumentFormattingResponse(id int, edits []TextEdit) DocumentFormattingResponse {
	return DocumentFormattingResponse{
		Response: Response{
			RPC: settings.RPC_VERSION,
			ID:  id,
		},
		Result: edits,
	}
}

// NewTextEdit creates and returns a new TextEdit.
// It initializes the text edit with the given start and end nodes and the new text to be inserted.
//
// Parameters:
//
//	start_node *sitter.Node - The starting node of the edit.
//	end_node *sitter.Node - The ending node of the edit.
//	new_text string - The new text to be inserted.
//
// Returns:
//
//	TextEdit - The initialized text edit.
func NewTextEdit(start_node *sitter.Node, end_node *sitter.Node, new_text string) TextEdit {
	return TextEdit{
		Range: Range{
			Start: Position{
				Line:      int(start_node.StartPoint().Row),
				Character: int(start_node.StartPoint().Column),
			},
			End: Position{
				Line:      int(end_node.EndPoint().Row),
				Character: int(end_node.EndPoint().Column),
			},
		},
		NewText: new_text,
	}
}
