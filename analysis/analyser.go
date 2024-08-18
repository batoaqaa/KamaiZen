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
	parser *sitter.Parser
}

func NewParser() *Parser {
	return &Parser{
		parser: sitter.NewParser(),
	}
}

func (p *Parser) SetLanguage() {
	p.parser.SetLanguage(sitter.NewLanguage(kamailio_cfg.Language()))
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
				logger.Println("Analyser Parsed document with URI: ", uri)

			}

		}
	}

}
