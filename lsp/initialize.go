package lsp

import "KamaiZen/settings"

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
	TextDocumentSync           TextDocumentSyncOptions `json:"textDocumentSync"`
	HoverProvider              bool                    `json:"hoverProvider"`
	DefinitionProvider         bool                    `json:"definitionProvider"`
	DocumentFormattingProvider bool                    `json:"documentFormattingProvider"`
	CompletionProvider         map[string]any          `json:"completionProvider"`
	DocumentHighlightProvider  bool                    `json:"documentHighlightProvider"`
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
				TextDocumentSync: TextDocumentSyncOptions{
					OpenClose: true,
					Change:    TEXT_DOCUMENT_SYNC_KIND_FULL,
					// FIXME: update this to incremental
					// Change:    TEXT_DOCUMENT_SYNC_KIND_INCREMENTAL,
				},
				HoverProvider:      true,
				DefinitionProvider: true,
				// TODO: for now formatter isnt working properly
				DocumentFormattingProvider: false,
				CompletionProvider:         map[string]any{"resolveProvider": false},
				DocumentHighlightProvider:  false,
			},
			ServerInfo: ServerInfo{
				Name:    settings.MY_NAME,
				Version: settings.KAMAIZEN_VERSION,
			},
		},
	}
}
