package lsp

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

type CompletionItem struct {
	Label         string `json:"label"`
	Detail        string `json:"detail"`
	Documentation string `json:"documentation"`
	// insert kind
	// detail
	// add kind
}

func NewCompletionResponse(id int, items []CompletionItem) CompletionResponse {
	return CompletionResponse{
		Response: Response{
			RPC: "2.0",
			ID:  id,
		},
		Result: items,
	}
}
