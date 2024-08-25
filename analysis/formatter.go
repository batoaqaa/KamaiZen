package analysis

import (
	"KamaiZen/kamailio_cfg"
	"KamaiZen/lsp"
	"KamaiZen/utils"
	"bytes"
	"context"

	"strings"

	sitter "github.com/smacker/go-tree-sitter"
)

func FormatKamailioCfg(content string) ([]lsp.TextEdit, error) {
	tabsize := 4
	logger := utils.GetLogger()
	sourceCode := []byte(content)
	parser := NewParser()
	parser.SetLanguage()

	// this is the root node of the tree
	root := parser.Parse(sourceCode)
	edits := []lsp.TextEdit{}
	block_level := 0
	var formatNode func(node *sitter.Node)
	formatNode = func(node *sitter.Node) {

		stringPadding := strings.Repeat(" ", block_level*tabsize)
		switch node.Type() {
		case kamailio_cfg.AssignmentExpressionNodeType:
			leftNode := node.ChildByFieldName("left")
			rightNode := node.ChildByFieldName("right")
			end_of_statment := node.NextSibling()
			if end_of_statment != nil {
				logger.Println("End of statement: ", end_of_statment.Type())
			}
			// Add the left-hand side, the equals sign with correct spacing, and the right-hand side
			leftNodeContent := string(sourceCode[leftNode.StartByte():leftNode.EndByte()])
			rightNodeContent := string(sourceCode[rightNode.StartByte():rightNode.EndByte()])
			leftNodeContent = strings.TrimSpace(leftNodeContent)
			leftNodeContent = stringPadding + leftNodeContent
			rightNodeContent = strings.TrimSpace(rightNodeContent)
			formattedContent := strings.Builder{}
			formattedContent.WriteString(leftNodeContent)
			formattedContent.WriteString(" = ") // Ensure exactly one space on both sides
			formattedContent.WriteString(rightNodeContent)
			logger.Println("Formatted content: ", formattedContent.String())
			edit := lsp.NewTextEdit(leftNode, rightNode, formattedContent.String())
			edits = append(edits, edit)
			logger.Printf("=====>Edit: \n%v", edit)
			// update the tree as well to reflect the changes
			// kamailio_cfg.UpdateTree(parser.GetTree(), leftNode, rightNode, formattedContent.String())
			// parser.UpdateTree([]byte(formattedContent.String()))
		case kamailio_cfg.ParenthesizedExpressionNodeType:
			if node.ChildCount() != 3 {
				logger.Println("Parenthesized expression does not have 3 children")
				break
			}
			leftBrace := node.Child(0)
			expression := node.Child(1)
			rightBrace := node.Child(2)
			formattedContent := strings.Builder{}
			formattedContent.WriteString("(")
			formattedContent.WriteString(string(sourceCode[expression.StartByte():expression.EndByte()]))
			formattedContent.WriteString(") ")
			edit := lsp.NewTextEdit(leftBrace, rightBrace, formattedContent.String())
			edits = append(edits, edit)
			// kamailio_cfg.UpdateTree(parser.GetTree(), leftBrace, rightBrace, formattedContent.String())
			// parser.UpdateTree([]byte(formattedContent.String()))
		case kamailio_cfg.TopLevelAssignmentNodeType:
			key := node.ChildByFieldName("key")
			value := node.ChildByFieldName("value")
			formattedContent := strings.Builder{}
			formattedContent.WriteString(string(sourceCode[key.StartByte():key.EndByte()]))
			formattedContent.WriteString("=")
			formattedContent.WriteString(string(sourceCode[value.StartByte():value.EndByte()]))
			edit := lsp.NewTextEdit(key, value, formattedContent.String())
			edits = append(edits, edit)
			// kamailio_cfg.UpdateTree(parser.GetTree(), key, value, formattedContent.String())
			// parser.UpdateTree([]byte(formattedContent.String()))
		case kamailio_cfg.CompoundStatementNodeType:
			leftBrace := node.Child(0)
			rightBrace := node.Child(int(node.ChildCount() - 1))
			edit := lsp.NewTextEdit(leftBrace, leftBrace, "{")
			edits = append(edits, edit)
			// kamailio_cfg.UpdateTree(parser.GetTree(), leftBrace, leftBrace, "{")
			// parser.UpdateTree([]byte("{"))
			edit = lsp.NewTextEdit(rightBrace, rightBrace, "}")
			edits = append(edits, edit)
			// kamailio_cfg.UpdateTree(parser.GetTree(), rightBrace, rightBrace, "}")
			// parser.UpdateTree([]byte("}"))
		case "if_statement":
			ifKeyword := node.Child(0) // if keyword
			formattedContent := strings.Builder{}
			content := stringPadding + "if "
			logger.Println("If keyword: ", ifKeyword.Content(sourceCode))
			if content == ifKeyword.Content(sourceCode) {
				break
			}
			formattedContent.WriteString(content)
			edit := lsp.NewTextEdit(ifKeyword, ifKeyword, formattedContent.String())
			edits = append(edits, edit)
		case "block_start":
			block_level++
			logger.Println("INCREASING Block level: ", block_level)
		case "block_end":
			block_level--
			logger.Println("DECREASING Block level: ", block_level)
		case "core_function":
			content := string(sourceCode[node.StartByte():node.EndByte()])
			content = stringPadding + content
			edit := lsp.NewTextEdit(node, node, content)
			edits = append(edits, edit)
			// kamailio_cfg.UpdateTree(parser.GetTree(), node, node, content)
			// parser.UpdateTree([]byte(content))
		case "call_expression":
			// content := string(sourceCode[node.StartByte():node.EndByte()])
			// content = stringPadding + content
			// edit := lsp.NewTextEdit(node, node, content)
			// edits = append(edits, edit)

			// default:
			// 	logger.Printf("Node type: %v -- continue\n", node.Type())
		}
	}
	logger.Println("Formatting Kamailio configuration file ", edits)
	TraverseNodeAndApply(root, formatNode)
	return edits, nil
}

func applyEditToContent(content string, edit lsp.TextEdit) string {
	logger := utils.GetLogger()
	logger.Println("Applying edits to content")
	sourceCode := []byte(content)
	// split source code 2d array of lines and characters
	lines := bytes.Split(sourceCode, []byte("\n"))
	var splitSourceCode [][]byte
	for _, line := range lines {
		var chars []byte
		for _, char := range line {
			chars = append(chars, char)
		}
		splitSourceCode = append(splitSourceCode, chars)
		// splitSourceCode = append(splitSourceCode, bytes.Split(line, nil))
	}

	logger.Println("Edit: ", edit)
	startLine := edit.Range.Start.Line
	startChar := edit.Range.Start.Character
	endLine := edit.Range.End.Line
	endChar := edit.Range.End.Character
	newText := edit.NewText
	// update source code with new text
	// Replace the specified range with the new text
	newTextBytes := []byte(newText)
	splitSourceCode[startLine] = append(splitSourceCode[startLine][:startChar], append(newTextBytes, splitSourceCode[endLine][endChar:]...)...)

	// Reconstruct the source code from the modified 2D array
	var result []byte
	for _, line := range splitSourceCode {
		result = append(result, line...)
		result = append(result, '\n')
	}
	result = bytes.TrimSuffix(result, []byte("\n")) // Remove the last newline

	return string(result)
}

// All siblings should be on the same level
func IndentKamailioCfg(content string, tabsize int) ([]lsp.TextEdit, error) {
	logger := utils.GetLogger()
	sourceCode := []byte(content)
	parser := NewParser()
	parser.SetLanguage()

	logger.Println("Indenting Kamailio configuration file")
	// this is the root node of the tree
	root := parser.Parse(sourceCode)
	edits := []lsp.TextEdit{}
	block_level := 0

	// go over the top level statements // each sibling is on the same level
	// top down call recursively IndentSiblings
	edits, err := IndentSiblings(root, sourceCode, tabsize, block_level)
	if err != nil {
		logger.Println("Error indenting siblings: ", err)
		return nil, err
	}

	return edits, nil

}

func IndentSiblings(root *sitter.Node, sourceCode []byte, tabsize int, block_level int) ([]lsp.TextEdit, error) {
	edits := []lsp.TextEdit{}
	for i := 0; i < int(root.NamedChildCount()); i++ {
		node := root.NamedChild(i)
		stringPadding := strings.Repeat(" ", block_level*tabsize)
		switch node.Type() {
		case "block_start":
			block_level++
		case "block_end":
			block_level--
		default:
			// Add the left-hand side, the equals sign with correct spacing, and the right-hand side
			nodeContent := string(sourceCode[node.StartByte():node.EndByte()])
			nodeContent = strings.TrimSpace(nodeContent)
			nodeContent = stringPadding + nodeContent
			formattedContent := strings.Builder{}
			formattedContent.WriteString(nodeContent)
			edit := lsp.NewTextEdit(node, node, formattedContent.String())
			edits = append(edits, edit)
		}
	}
	return edits, nil
}

func formatAssignmentExpression(tree *sitter.Tree, node *sitter.Node, parser *sitter.Parser, sourceCode []byte) ([]lsp.TextEdit, *sitter.Tree, error) {
	logger := utils.GetLogger()
	edits := []lsp.TextEdit{}
	leftNode := node.ChildByFieldName("left")
	rightNode := node.ChildByFieldName("right")
	// Add the left-hand side, the equals sign with correct spacing, and the right-hand side
	leftNodeContent := string(sourceCode[leftNode.StartByte():leftNode.EndByte()])
	rightNodeContent := string(sourceCode[rightNode.StartByte():rightNode.EndByte()])
	leftNodeContent = strings.TrimSpace(leftNodeContent)
	rightNodeContent = strings.TrimSpace(rightNodeContent)
	formattedContent := strings.Builder{}
	formattedContent.WriteString(leftNodeContent)
	formattedContent.WriteString(" = ") // Ensure exactly one space on both sides
	formattedContent.WriteString(rightNodeContent)
	logger.Println("Formatted content: ", formattedContent.String())
	edit := lsp.NewTextEdit(leftNode, rightNode, formattedContent.String())
	edits = append(edits, edit)
	// update the tree as well to reflect the changes
	kamailio_cfg.UpdateTree(tree, leftNode, rightNode, formattedContent.String())
	tree, err := parser.ParseCtx(context.Background(), tree, []byte(formattedContent.String()))
	if err != nil {
		logger.Println("Error parsing the formatted content: ", err)
		return nil, nil, err
	}
	return edits, tree, nil
}
