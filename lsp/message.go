package lsp

// Request represents a JSON-RPC request message.
// It contains the JSON-RPC version, the request ID, and the method to be invoked.
type Request struct {
	RPC    string `json:"jsonrpc"`
	ID     int    `json:"id"`
	Method string `json:"method"`
}

// Response represents a JSON-RPC response message.
// It contains the JSON-RPC version and the response ID.
type Response struct {
	RPC string `json:"jsonrpc"`
	ID  int    `json:"id,omitempty"`
	// Result
	// Error
}

// Notification represents a JSON-RPC notification message.
// It contains the JSON-RPC version and the method to be invoked.
// Notifications do not have an ID because they do not expect a response.
type Notification struct {
	RPC    string `json:"jsonrpc"`
	Method string `json:"method"`
}
