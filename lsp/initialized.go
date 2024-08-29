package lsp

type InitializedNotification struct {
	Notification
	Params InitializedParams `json:"params"`
}

type InitializedParams struct{}
