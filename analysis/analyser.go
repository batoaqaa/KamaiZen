package analysis

import (
	"KamaiZen/kamailio_cfg"
	"KamaiZen/lsp"
	"context"
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
	n := tree.RootNode()
	return n
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

func StartAnalyser(c <-chan State, logger *log.Logger, wg *sync.WaitGroup) {
	// receive the message from the channel
	// we will receive state struct on the channel
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
				// // stateTree.TraverseNode(uri, node, logger, 0)
				// node = node.Child(0)
				// nodeCount := node.ChildCount()
				// logger.Printf("countChildren %d", nodeCount)
				// logger.Printf("%s", node.Symbol())
				// logger.Printf("%d, %d", node.StartPoint(), node.EndPoint())
				// logger.Printf("value %s", node.Child(0).IsExtra)
				// logger.Printf("sibling %s", node.NextSibling())
				diagnostics := getDiagnostics(node, parser)
				for _, diag := range diagnostics {
					logger.Printf("Diagnostic: %s at %d:%d", diag.Message, diag.Line, diag.Column)
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
		logger.Printf("%s::%s[%d:%d]%*s", uri, node.Type, node.StartByte(), node.EndByte(), padding, node)
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
func getDiagnostics(rootNode *sitter.Node, parser *Parser) []Diagnostic {
	var diagnostics []Diagnostic

	// Example: Check for undeclared variables or any syntax issues
	queryStr := "(ERROR) @error"
	q, err := sitter.NewQuery([]byte(queryStr), parser.GetLanguage())
	if err != nil {
		log.Fatalf("Failed to create query: %v", err)
	}
	cursor := sitter.NewQueryCursor()
	cursor.Exec(q, rootNode)
	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}
		// Found a syntax error node
		for _, capture := range match.Captures {
			node := capture.Node
			diagnostics = append(diagnostics, Diagnostic{
				Message: "Syntax error",
				Line:    node.StartPoint().Row,
				Column:  node.StartPoint().Column,
			})
		}
	}
	return diagnostics
}
