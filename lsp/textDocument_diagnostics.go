package lsp

import "KamaiZen/settings"

// PublishDiagnosticNotification represents a notification sent to the client to publish diagnostics.
// It contains the notification metadata and the parameters for the diagnostics.
type PublishDiagnosticNotification struct {
	Notification
	Params PublishDiagnosticParams `json:"params"`
}

// PublishDiagnosticParams contains the parameters for the PublishDiagnosticNotification.
// It includes the URI of the document and the list of diagnostics.
type PublishDiagnosticParams struct {
	URI         DocumentURI  `json:"uri"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

// Diagnostic represents a diagnostic message, such as a compiler error or warning.
// It includes the range, severity, source, and message of the diagnostic.
type Diagnostic struct {
	Range    Range              `json:"range"`
	Severity DiagnosticSeverity `json:"severity,omitempty"`
	// Code               interface{}     `json:"code,omitempty"`
	Source  string `json:"source,omitempty"`
	Message string `json:"message"`
	// Tags    []DiagnosticTag `json:"tags,omitempty"`
	// RelatedInformation []DiagnosticRelatedInformation `json:"relatedInformation,omitempty"`
}

// NewPublishDiagnosticNotification creates and returns a new PublishDiagnosticNotification.
// It initializes the notification with the given URI and list of diagnostics.
//
// Parameters:
//
//	uri DocumentURI - The URI of the document.
//	diagnostics []Diagnostic - The list of diagnostics.
//
// Returns:
//
//	PublishDiagnosticNotification - The initialized notification.
func NewPublishDiagnosticNotification(uri DocumentURI, diagnostics []Diagnostic) PublishDiagnosticNotification {
	return PublishDiagnosticNotification{
		Notification: Notification{
			Method: "textDocument/publishDiagnostics",
			RPC:    settings.RPC_VERSION,
		},
		Params: PublishDiagnosticParams{
			URI:         uri,
			Diagnostics: diagnostics,
		},
	}
}

// DiagnosticSeverity represents the severity of a diagnostic message.
// It is an enumeration of various severity levels.
type DiagnosticSeverity int

const (
	ERROR DiagnosticSeverity = iota + 1
	WARNING
	INFORMATION
	HINT
)
