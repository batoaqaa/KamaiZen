package lsp

import "KamaiZen/settings"

type HoverRequest struct {
	Request
	Params HoverParams `json:"params"`
}

type HoverParams struct {
	TextDocuemntPositionParams
}

type HoverResponse struct {
	Response
	Result *Hover `json:"result"`
}

type Hover struct {
	Contents MarkupContent `json:"contents"`
	Range    *Range        `json:"range,omitempty"`
}

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
