package kamailio_cfg

import "KamaiZen/lsp"

type DiagnosticVisitor struct {
	diagnostics []lsp.Diagnostic
}

func NewDiagnosticVisitor() *DiagnosticVisitor {
	return &DiagnosticVisitor{}
}

func (d *DiagnosticVisitor) Visit(node *ASTNode) error {
	if node.Node.IsError() {
		diagnostic := lsp.Diagnostic{
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      int(node.Node.StartPoint().Row),
					Character: int(node.Node.StartPoint().Column),
				},
				End: lsp.Position{
					Line:      int(node.Node.EndPoint().Row),
					Character: int(node.Node.EndPoint().Column),
				},
			},
			Message: "Syntax error",
		}
		d.diagnostics = append(d.diagnostics, diagnostic)
	}

	// Traverse the children
	for i := 0; i < int(node.Node.ChildCount()); i++ {
		child := node.Node.Child(i)
		d.Visit(&ASTNode{Node: child})
	}
	return nil
}

func (d *DiagnosticVisitor) GetDiagnostics() []lsp.Diagnostic {
	return d.diagnostics
}
