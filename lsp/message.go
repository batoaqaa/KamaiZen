package lsp

type Request struct {
	RPC    string `json:"jsonrpc"`
	ID     int    `json:"id"`
	Method string `json:"method"`

	// We will specify the params later
	// Params interface{} `json:"params"`
}

type Response struct {
	RPC string `json:"jsonrpc"`
	ID  int    `json:"id,omitempty"`
	// Result
	// Error
}

type Notification struct {
	RPC    string `json:"jsonrpc"`
	Method string `json:"method"`
}
