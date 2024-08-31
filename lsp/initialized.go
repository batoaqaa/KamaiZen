package lsp

// InitializedNotification represents a notification sent to the server to indicate that the client is ready.
// It contains the notification metadata and the parameters for the initialization.
type InitializedNotification struct {
	Notification
	Params InitializedParams `json:"params"`
}

// InitializedParams contains the parameters for the InitializedNotification.
// Currently, it does not include any fields.
type InitializedParams struct{}
