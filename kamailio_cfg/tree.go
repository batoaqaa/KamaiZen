package kamailio_cfg

import (
	sitter "github.com/smacker/go-tree-sitter"
)

const (
	ReturnNodeType                  = "return_statement"
	CoreFunctionNodeType            = "core_function_statement"
	SatementNodeType                = "statement"
	ParenthesizedExpressionNodeType = "parenthesized_expression"
	AssignmentExpressionNodeType    = "assignment_expression"
	TopLevelAssignmentNodeType      = "top_level_assignment_expression"
	CompoundStatementNodeType       = "compound_statement"
	CallExpressionNodeType          = "call_expression"
	IdentifierNodeType              = "identifier"
	ExpressionNodeType              = "expression"
	BlockEndNodeType                = "block_end"
)

func UpdateTree(tree *sitter.Tree, leftNode, rightNode *sitter.Node, formattedContent string) {
	tree.Edit(sitter.EditInput{
		StartIndex:  leftNode.StartByte(),
		OldEndIndex: rightNode.EndByte(),
		NewEndIndex: uint32(len(formattedContent)),
		StartPoint:  leftNode.StartPoint(),
		OldEndPoint: rightNode.EndPoint(),
		NewEndPoint: sitter.Point{
			Row:    leftNode.StartPoint().Row,
			Column: uint32(len(formattedContent)),
		},
	})
}
