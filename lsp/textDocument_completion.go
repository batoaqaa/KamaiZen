package lsp

import "KamaiZen/settings"

// CompletionRequest represents a request for code completion suggestions.
// It contains the request metadata and the parameters for the completion request.
type CompletionRequest struct {
	Request
	Params CompletionParams `json:"params"`
}

// CompletionParams contains the parameters for the CompletionRequest.
// It includes the text document position parameters and an optional context.
type CompletionParams struct {
	TextDocuemntPositionParams
	// TODO: optional context
}

// CompletionResponse represents the response to a CompletionRequest.
// It contains the response metadata and the list of completion items.
type CompletionResponse struct {
	Response
	Result []CompletionItem `json:"result"`
}

// CompletionItemKind represents the kind of a completion item.
// It is an enumeration of various kinds of completion items.
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

// CompletionItem represents a single completion item in the completion response.
// It includes the label, detail, documentation, and kind of the completion item.
type CompletionItem struct {
	Label         string             `json:"label"`
	Detail        string             `json:"detail"`
	Documentation string             `json:"documentation"`
	Kind          CompletionItemKind `json:"kind"`
	// insert kind
	// detail
	// add kind
}

// NewCompletionResponse creates and returns a new CompletionResponse.
// It initializes the response with the given ID and the list of completion items.
//
// Parameters:
//
//	id int - The ID of the response.
//	items []CompletionItem - The list of completion items.
//
// Returns:
//
//	CompletionResponse - The initialized response.
func NewCompletionResponse(id int, items []CompletionItem) CompletionResponse {
	return CompletionResponse{
		Response: Response{
			RPC: settings.RPC_VERSION,
			ID:  id,
		},
		Result: items,
	}
}
