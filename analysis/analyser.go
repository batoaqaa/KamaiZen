package analysis

import (
	"KamaiZen/kamailio_cfg"
	"KamaiZen/lsp"
	"context"
	"io"
	"log"
	"sync"

	sitter "github.com/smacker/go-tree-sitter"
)

type Message struct {
	state State
}

type Parser struct {
	parser   *sitter.Parser
	language *sitter.Language
	tree     *sitter.Tree
}

func NewParser() *Parser {
	return &Parser{
		parser: sitter.NewParser(),
	}
}

func (p *Parser) SetLanguage() {
	p.language = sitter.NewLanguage(kamailio_cfg.Language())
	p.parser.SetLanguage(p.language)
}

func (p *Parser) GetLanguage() *sitter.Language {
	return p.language
}

func (p *Parser) Parse(sourceCode []byte) *sitter.Node {
	tree, _ := p.parser.ParseCtx(context.Background(), nil, sourceCode)
	p.tree = tree
	n := p.tree.RootNode()
	return n
}

func (p *Parser) GetTree() *sitter.Tree {
	return p.tree
}

func (p *Parser) UpdateTree(content []byte) {
	newTree, _ := p.parser.ParseCtx(context.Background(), p.tree, content)
	p.tree = newTree
}

type StateTree struct {
	nodes map[lsp.DocumentURI]*sitter.Node
}

func NewStateTree() StateTree {
	return StateTree{
		nodes: make(map[lsp.DocumentURI]*sitter.Node),
	}
}

func (s *StateTree) AddNode(uri lsp.DocumentURI, node *sitter.Node) {
	s.nodes[uri] = node
}

func GetDiagnosticsForDocument(uri lsp.DocumentURI, content string) []lsp.Diagnostic {
	parser := NewParser()
	parser.SetLanguage()
	node := parser.Parse([]byte(content))
	_ = node
	// diagnostics := getDiagnostics(node, parser)
	return []lsp.Diagnostic{}
}

func StartAnalyser(c <-chan State, writer io.Writer, logger *log.Logger, wg *sync.WaitGroup) {
	defer wg.Done()
	parser := NewParser()
	parser.SetLanguage()
	stateTree := NewStateTree()
	logger.Println("=====Analyser started")

	for {
		select {
		case state, ok := <-c:
			if !ok {
				logger.Println("Channel closed! Exiting...")
				return
			}
			for uri, content := range state.Documents {
				// start the analysis of the text document here
				node := parser.Parse([]byte(content))
				stateTree.AddNode(uri, node)
				logger.Printf("-----------------")
				diagnostics := getDiagnostics(logger, node, parser)
				if diagnostics != nil {

					// logger.Printf("Diagnostics: %v", diagnostics)
					lsp.WriteResponse(writer, lsp.NewPublishDiagnosticNotification(uri, diagnostics))
				}
			}
		}
	}
}

// match the document lines with the node tree

func (s *StateTree) TraverseNode(uri lsp.DocumentURI, node *sitter.Node, logger *log.Logger, padding int) {
	// traverse the node and print the node
	// logger.Println(node)
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

// getDiagnostics traverses the parse tree to find issues
func getDiagnostics(logger *log.Logger, rootNode *sitter.Node, parser *Parser) []lsp.Diagnostic {
	var diagnostics []lsp.Diagnostic

	// get syntax errors
	// log.Printf("Root node: %v", rootNode)
	diagnostics = append(diagnostics, getSyntaxErrors(rootNode, parser)...)
	// log.Printf("Diagnostics: %v", diagnostics)

	// get deprecated comments
	diagnostics = append(diagnostics, getDeprecatedComments(rootNode, parser)...)
	// get unreachable code
	diagnostics = append(diagnostics, getUnreachableCode(logger, rootNode, parser)...)
	// logger.Printf("Diagnostics: %v", diagnostics)
	return diagnostics
}

func getDeprecatedComments(node *sitter.Node, parser *Parser) []lsp.Diagnostic {
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

func getSyntaxErrors(node *sitter.Node, parser *Parser) []lsp.Diagnostic {
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

func getUnreachableCode(logger *log.Logger, node *sitter.Node, parser *Parser) []lsp.Diagnostic {
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
