package kamailio_cfg

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"
)

// TODO: Implement NotificationHandle for the parser and visitor functions.
// The AST should be parsed and updated in the same instance of the parser.
// This is because the parser instance keeps track of the tree and the tree
// is updated in place. If a new instance of the parser is created, the tree
// will be parsed from scratch and the previous tree will be lost.
// This is why the parser instance is stored in the analyser and the diagnostics
// functions.

// Watch document states and update the parser tree accordingly.
// Document state would watch stateTree and update the document state
// Whenever the docuemnt State is updated, the parser tree should be updated and diagnostics should be generated.

// Implement Memento for parser to undo changes in case of error.

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
	p.language = sitter.NewLanguage(Language())
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
