package lsp

type InitializeRequest struct {
	Request
	Params InitializeRequestParams `json:"params"`
}

type InitializeRequestParams struct {
	ClientInfo ClientInfo `json:"clientInfo"`
	// TODO: Add capabilities
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResponse struct {
	Response
	Result InitializeResult `json:"result"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}

type ServerCapabilities struct {
	// TODO: incremental instead of full
	TextDocumentSync           int  `json:"textDocumentSync"`
	HoverProvider              bool `json:"hoverProvider"`
	DefinitionProvider         bool `json:"definitionProvider"`
	DocumentFormattingProvider bool `json:"documentFormattingProvider"`
	// TODO: Add more capabilities
	// CodeActionProvider bool `json:"codeActionProvider"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func NewInitializeResponse(id int) InitializeResponse {
	return InitializeResponse{
		Response: Response{
			RPC: "2.0",
			ID:  id,
		},
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				// TODO: incremental instead of full
				TextDocumentSync:   1,
				HoverProvider:      true,
				DefinitionProvider: true,
				// TODO: for now formatter isnt working properly
				DocumentFormattingProvider: false,
			},
			ServerInfo: ServerInfo{
				Name:    "KamaiZen",
				Version: "0.0.1",
			},
		},
	}
}
