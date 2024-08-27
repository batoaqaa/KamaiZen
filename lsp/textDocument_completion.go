package lsp

import "KamaiZen/settings"

type CompletionRequest struct {
	Request
	Params CompletionParams `json:"params"`
}

type CompletionParams struct {
	TextDocuemntPositionParams
	// optional context
}

type CompletionResponse struct {
	Response
	Result []CompletionItem `json:"result"`
}

type CompletionItemKind int

const (
	TEXT_COMPLETION CompletionItemKind = iota
	METHOD_COMPLETION
	FUNCTION_COMPLETION
	CONSTRUCTOR_COMPLETION
	FIELD_COMPLETION
	VARIABLE_COMPLETION
	CLASS_COMPLETION
	INTERFACE_COMPLETION
	MODULE_COMPLETION
	PROPERTY_COMPLETION
	UNIT_COMPLETION
	VALUE_COMPLETION
	ENUM_COMPLETION
	KEYWORD_COMPLETION
	SNIPPET_COMPLETION
	COLOR_COMPLETION
	FILE_COMPLETION
	REFERENCE_COMPLETION
)

type CompletionItem struct {
	Label         string             `json:"label"`
	Detail        string             `json:"detail"`
	Documentation string             `json:"documentation"`
	Kind          CompletionItemKind `json:"kind"`
	// insert kind
	// detail
	// add kind
}

func NewCompletionResponse(id int, items []CompletionItem) CompletionResponse {
	return CompletionResponse{
		Response: Response{
			RPC: settings.RPC_VERSION,
			ID:  id,
		},
		Result: items,
	}
}
