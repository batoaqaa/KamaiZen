package kamailio_cfg

import (
	"KamaiZen/logger"
)

// keeps track of variables declared int he program

// global variables
// gloabl variables are avps
var globalVariables map[string]Variable

// local variables
var localVariables map[string]Variable

// Variable struct
type Variable struct {
	Name       string
	Value      string
	Scope      string
	Identifier string
}

func InitialiseVariables() {
	globalVariables = make(map[string]Variable)
	localVariables = make(map[string]Variable)
}

func AddGlobalVariable(name string, value string, identifier string) {
	globalVariables[name] = Variable{name, value, "AVP", identifier}
}

func AddLocalVariable(name string, value string, scope string, identifier string) {
	localVariables[name] = Variable{name, value, scope, identifier}
}

func ExtractGlobalVariables(a *Analyzer, source_code []byte) {
	// extract global variables
	q, err := NewQueryExecutor(_ASSINGMENT, a.ast.Node, a.builder.parser.language)
	if err != nil {
		logger.Error("Error creating query executor: ", err)
		return
	}
	for {
		match, ok := q.NextMatch()
		if !ok {
			break
		}
		for _, capture := range match.Captures {
			node := capture.Node
			variable := node.ChildByFieldName("left")
			if variable.Type() == "pseudo_variable" {
				pc := variable.NamedChild(0)
				if pc.Type() == "pseudo_content" {
					avp := pc.NamedChild(0)
					if avp.Type() == "avp_var" {
						identifier := avp.ChildByFieldName("name").Child(0).Content(source_code)
						avp_name := "$avp(" + identifier + ")"
						value := node.ChildByFieldName("right").Content(source_code)
						AddGlobalVariable(avp_name, value, identifier)
					}
				}
			}
		}
	}
}

func (v *Variable) GetGlobalVariableDocs() string {
	return "## User defined AVP\n\t" + v.Identifier + "\n" +
		"### Value\n\t```\n" + "\t" + v.Value + "\n```\n" +
		"### Scope\n\t" + v.Scope + "\n"
}

func GetGlobalVariables() map[string]Variable {
	return globalVariables
}
