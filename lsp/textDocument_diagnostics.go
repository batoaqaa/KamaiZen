package lsp

type PublishDiagnosticNotification struct {
	Notification
	Params PublishDiagnosticParams `json:"params"`
}

type PublishDiagnosticParams struct {
	URI         DocumentURI  `json:"uri"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

type Diagnostic struct {
	Range    Range              `json:"range"`
	Severity DiagnosticSeverity `json:"severity,omitempty"`
	// Code               interface{}     `json:"code,omitempty"`
	Source  string `json:"source,omitempty"`
	Message string `json:"message"`
	// Tags    []DiagnosticTag `json:"tags,omitempty"`
	// RelatedInformation []DiagnosticRelatedInformation `json:"relatedInformation,omitempty"`
}

func NewPublishDiagnosticNotification(uri DocumentURI, diagnostics []Diagnostic) PublishDiagnosticNotification {
	return PublishDiagnosticNotification{
		Notification: Notification{
			Method: "textDocument/publishDiagnostics",
			RPC:    "2.0",
		},
		Params: PublishDiagnosticParams{
			URI:         uri,
			Diagnostics: diagnostics,
		},
	}
}

type DiagnosticSeverity int

const (
	ERROR DiagnosticSeverity = iota + 1
	WARNING
	INFORMATION
	HINT
)
