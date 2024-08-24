package lsp

import sitter "github.com/smacker/go-tree-sitter"

type FormattingOptions struct {
	TabSize                int  `json:"tabSize"`
	InsertSpaces           bool `json:"insertSpaces"`
	TrimTrailingWhitespace bool `json:"trimTrailingWhitespace"`
	InsertFinalNewline     bool `json:"insertFinalNewline"`
	TrimFinalNewlines      bool `json:"trimFinalNewlines"`
}

type DocumentFormattingParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Options      FormattingOptions      `json:"options"`
}

// type FormattingNotification struct {
// 	Notification
// 	Params DocumentFormattingParams `json:"params"`
// }

type DocumentFormattingRequest struct {
	Request
	Params DocumentFormattingParams `json:"params"`
}

type DocumentFormattingResponse struct {
	Response
	Result []TextEdit `json:"result"`
}

type TextEdit struct {
	Range   Range  `json:"range"`
	NewText string `json:"newText"`
}

func NewDocumentFormattingResponse(id int, edits []TextEdit) DocumentFormattingResponse {
	return DocumentFormattingResponse{
		Response: Response{
			RPC: "2.0",
			ID:  id,
		},
		Result: edits,
	}
}

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
