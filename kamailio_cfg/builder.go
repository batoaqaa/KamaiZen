package kamailio_cfg

import sitter "github.com/smacker/go-tree-sitter"

type ASTNode struct {
	Node *sitter.Node
}

func (a *ASTNode) IsError() bool {
	return a.Node.IsError()
}

func (a *ASTNode) StartPoint() sitter.Point {
	return a.Node.StartPoint()
}

func (a *ASTNode) EndPoint() sitter.Point {
	return a.Node.EndPoint()
}

func (a *ASTNode) ChildCount() int {
	return int(a.Node.ChildCount())
}

func (a *ASTNode) NamedChildCount() int {
	return int(a.Node.NamedChildCount())
}

func (a *ASTNode) Accept(v ASTVisitor) {
	v.Visit(a)
}

type KamailioASTBuilder struct {
	parser *Parser
}

func NewKamailioASTBuilder() *KamailioASTBuilder {
	return &KamailioASTBuilder{
		parser: NewParser(),
	}
}

func (k *KamailioASTBuilder) BuildAST(sourceCode []byte) *ASTNode {
	k.parser.SetLanguage()
	n := k.parser.Parse(sourceCode)
	return &ASTNode{
		Node: n,
	}
}

func (k *KamailioASTBuilder) UpdateAST(content []byte) *ASTNode {
	k.parser.UpdateTree(content)
	return &ASTNode{
		Node: k.parser.GetTree().RootNode(),
	}
}
