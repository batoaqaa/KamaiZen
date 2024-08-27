package analysis

import (
	"KamaiZen/document_manager"
	"KamaiZen/kamailio_cfg"
	"KamaiZen/logger"
	"KamaiZen/lsp"
	"log"
	"sync"

	sitter "github.com/smacker/go-tree-sitter"
)

type Message struct {
	state State
}

type StateTree struct {
	nodes map[lsp.DocumentURI]*sitter.Node
}

// module level state tree
var stateTreeCache StateTree

func NewStateTree() StateTree {
	return StateTree{
		nodes: make(map[lsp.DocumentURI]*sitter.Node),
	}
}

func (s *StateTree) AddNode(uri lsp.DocumentURI, node *sitter.Node) {
	s.nodes[uri] = node
}

func GetDiagnosticsForDocument(uri lsp.DocumentURI, content string) []lsp.Diagnostic {
	// TODO: use the previous parser instance
	parser := kamailio_cfg.NewParser()
	parser.SetLanguage()
	node := parser.Parse([]byte(content))
	_ = node
	// diagnostics := getDiagnostics(node, parser)
	return []lsp.Diagnostic{}
}

func StartAnalyser(c <-chan State, wg *sync.WaitGroup) {
	defer wg.Done()
	// TODO: use the previous parser instance
	parser := kamailio_cfg.NewParser()
	parser.SetLanguage()
	stateTree := NewStateTree()
	logger.Info("=====Analyser started")
	for {
		select {
		case state, ok := <-c:
			if !ok {
				logger.Info("Channel closed! Exiting...")
				return
			}
			for uri, content := range state.Documents {
				// start the analysis of the text document here
				node := parser.Parse([]byte(content))
				stateTree.AddNode(uri, node)
				diagnostics := getDiagnostics(node, parser)
				if diagnostics != nil {

					// logger.Printf("Diagnostics: %v", diagnostics)
					lsp.WriteResponse(lsp.NewPublishDiagnosticNotification(uri, diagnostics))
				}
				stateTreeCache = stateTree
			}
		}
	}
}

// match the document lines with the node tree

func (s *StateTree) TraverseNode(uri lsp.DocumentURI, node *sitter.Node, logger *log.Logger, padding int) {
	// traverse the node and print the node
	// logger.Info(node)
	var i uint32
	childCount := node.ChildCount()
	for i = 0; i < childCount; i++ {
		// Print spaces for padding
		child := node.Child(int(i))
		s.TraverseNode(uri, child, logger, padding+2)
	}
}

// Diagnostic represents a syntax issue
type Diagnostic struct {
	Message string
	Line    uint32
	Column  uint32
}

func GetFunctionNameAtPosition(uri lsp.DocumentURI, position lsp.Position, source_code []byte) string {
	node := stateTreeCache.nodes[uri]
	return getFunctionName(node, position, source_code)
}

func getFunctionName(node *sitter.Node, position lsp.Position, source_code []byte) string {
	nodeAtPosition := node.NamedDescendantForPointRange(
		sitter.Point{
			Row:    uint32(position.Line),
			Column: uint32(position.Character),
		},
		sitter.Point{
			Row:    uint32(position.Line),
			Column: uint32(position.Character),
		})
	functionName := nodeAtPosition.Content(source_code)
	return functionName
}

// getDiagnostics traverses the parse tree to find issues
func getDiagnostics(rootNode *sitter.Node, parser *kamailio_cfg.Parser) []lsp.Diagnostic {
	var diagnostics []lsp.Diagnostic

	// get syntax errors
	// log.Printf("Root node: %v", rootNode)
	diagnostics = append(diagnostics, getSyntaxErrors(rootNode, parser)...)
	// log.Printf("Diagnostics: %v", diagnostics)

	// get deprecated comments
	diagnostics = append(diagnostics, getDeprecatedComments(rootNode, parser)...)
	// get unreachable code
	diagnostics = append(diagnostics, getUnreachableCode(rootNode, parser)...)
	// logger.Printf("Diagnostics: %v", diagnostics)
	return diagnostics
}

func getDeprecatedComments(node *sitter.Node, parser *kamailio_cfg.Parser) []lsp.Diagnostic {
	var diagnostics []lsp.Diagnostic

	queryStr := "(deprecated_comment) @deprecated"
	q, err := sitter.NewQuery([]byte(queryStr), parser.GetLanguage())
	if err != nil {
		return nil
	}
	cursor := sitter.NewQueryCursor()
	cursor.Exec(q, node)
	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}
		// Found a deprecated comment node
		for _, capture := range match.Captures {
			node := capture.Node
			diagnostics = append(diagnostics, addDiagnostic("Deprecated comment, use /* comment */", node, lsp.HINT))
		}
	}
	return diagnostics

}

func getSyntaxErrors(node *sitter.Node, parser *kamailio_cfg.Parser) []lsp.Diagnostic {
	var diagnostics []lsp.Diagnostic
	queryStr := "(ERROR) @error"
	q, err := sitter.NewQuery([]byte(queryStr), parser.GetLanguage())
	if err != nil {
		return nil
	}
	cursor := sitter.NewQueryCursor()
	cursor.Exec(q, node)
	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}
		// Found a syntax error node
		for _, capture := range match.Captures {
			node := capture.Node
			diagnostics = append(diagnostics, addDiagnostic("Syntax error", node, lsp.ERROR))
		}
	}
	return diagnostics

}

// hey chatgpt, how do i get the symbol tree and use it for the diagnostics?
func addDiagnostic(message string, node *sitter.Node, severity lsp.DiagnosticSeverity) lsp.Diagnostic {
	// add diagnostic to the node
	return lsp.Diagnostic{
		Range: lsp.Range{
			Start: lsp.Position{
				Line:      int(node.StartPoint().Row),
				Character: int(node.StartPoint().Column),
			},
			End: lsp.Position{
				Line:      int(node.EndPoint().Row),
				Character: int(node.EndPoint().Column),
			},
		},
		Severity: severity,
		Message:  message,
	}

}

func getUnreachableCode(node *sitter.Node, parser *kamailio_cfg.Parser) []lsp.Diagnostic {
	// check for unreachable code
	// get compound statement and core function
	var diagnostics []lsp.Diagnostic
	queryStr := "(compound_statement) @compound"
	q, err := sitter.NewQuery([]byte(queryStr), parser.GetLanguage())
	if err != nil {
		return nil
	}
	cursor := sitter.NewQueryCursor()
	cursor.Exec(q, node)
	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}
		for _, capture := range match.Captures {
			node := capture.Node
			for i := 0; i < int(node.NamedChildCount()); i++ {
				cmpStmtItem := node.NamedChild(i)
				// logger.Printf("Compound statement item: %v", cmpStmtItem.Type())
				if cmpStmtItem.Type() == kamailio_cfg.SatementNodeType && cmpStmtItem.NamedChildCount() > 0 {
					statement := cmpStmtItem.NamedChild(0)
					// logger.Printf("Block item statement type: %v", statement.Type())
					if (statement.Type() == kamailio_cfg.CoreFunctionNodeType || statement.Type() == kamailio_cfg.ReturnNodeType) && i != int(node.NamedChildCount())-1 {
						diagnostics = append(diagnostics, addDiagnostic("Unreachable code", node.NamedChild(i+1), lsp.WARNING))
					}
				}
			}
		}
	}
	return diagnostics
}

// traverse node and apply the function to the node
func TraverseNodeAndApply(node *sitter.Node, f func(*sitter.Node)) {
	f(node)
	if node.ChildCount() == 0 {
		// get siblings
		if node.NextSibling() != nil {
			TraverseNodeAndApply(node.NextSibling(), f)
		}
		return
	}
	for i := 0; i < int(node.NamedChildCount()); i++ {
		TraverseNodeAndApply(node.NamedChild(i), f)
	}

}

func GetCompletionItems(uri lsp.DocumentURI) []lsp.CompletionItem {
	var completionItems []lsp.CompletionItem
	functions := document_manager.GetAllAvailableFunctionDocs()
	for _, function := range functions {
		completionItems = append(completionItems, lsp.CompletionItem{
			Detail:        function.Name,
			Label:         function.Name + "(" + function.Parameters + ")",
			Documentation: function.Description + "\n" + function.Example,
		})
	}

	return completionItems
}
