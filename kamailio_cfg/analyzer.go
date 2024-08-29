package kamailio_cfg

type Analyzer struct {
	builder *KamailioASTBuilder
	ast     *ASTNode
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{
		builder: NewKamailioASTBuilder(),
	}
}

func (a *Analyzer) Initialize(content []byte) {
	a.ast = a.builder.BuildAST(content)
}

func (a *Analyzer) Update(content []byte) {
	a.ast = a.builder.UpdateAST(content)
}

func (a *Analyzer) GetAST() *ASTNode {
	return a.ast
}
