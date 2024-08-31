package lsp

import "KamaiZen/settings"

// HoverRequest represents a request for hover information.
// It contains the request metadata and the parameters for the hover request.
type HoverRequest struct {
	Request
	Params HoverParams `json:"params"`
}

// HoverParams contains the parameters for the HoverRequest.
// It includes the text document position parameters.
type HoverParams struct {
	TextDocuemntPositionParams
}

// HoverResponse represents the response to a HoverRequest.
// It contains the response metadata and the hover result.
type HoverResponse struct {
	Response
	Result *Hover `json:"result"`
}

// Hover represents the hover information.
// It includes the contents of the hover and the range within the document.
type Hover struct {
	Contents MarkupContent `json:"contents"`
	Range    *Range        `json:"range,omitempty"`
}

// NewHoverResponse creates and returns a new HoverResponse.
// It initializes the response with the given ID and sets the hover contents.
//
// Parameters:
//
//	id int - The ID of the response.
//	contents string - The contents of the hover.
//
// Returns:
//
//	HoverResponse - The initialized response.
func NewHoverResponse(id int, contents string) HoverResponse {
	return HoverResponse{
		Response: Response{
			RPC: settings.RPC_VERSION,
			ID:  id,
		},
		Result: &Hover{
			Contents: MarkupContent{
				Kind:  "markdown",
				Value: contents,
			},
		},
	}
}
