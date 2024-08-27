package lsp

import "KamaiZen/settings"

type DefinitionProviderRequest struct {
	Request
	Params DefinitionProviderParams `json:"params"`
}

type DefinitionProviderParams struct {
	TextDocuemntPositionParams
}

type DefinitionProviderResponse struct {
	Response
	Result Location `json:"result"`
}

type DefinitionProvider struct {
	Contents MarkupContent `json:"contents"`
	Range    *Range        `json:"range,omitempty"`
}

func NewDefintionProviderResponse(id int, contents string) DefinitionProviderResponse {
	return DefinitionProviderResponse{
		Response: Response{
			RPC: settings.RPC_VERSION,
			ID:  id,
		},
		Result: Location{ // Location is a struct defined in textDocument_hover.go
			URI: "https://www.kamailio.org/docs/modules/devel/modules/avpops.html",
			Range: Range{
				Start: Position{
					Line:      0,
					Character: 0,
				},
				End: Position{
					Line:      0,
					Character: 0,
				},
			},
		},
	}
}
