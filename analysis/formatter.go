package analysis

import (
	"KamaiZen/kamailio_cfg"
	"KamaiZen/lsp"
	"KamaiZen/utils"

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
			formattedContent.WriteString(content)
			edit := lsp.NewTextEdit(ifKeyword, ifKeyword, formattedContent.String())
			edits = append(edits, edit)
			// kamailio_cfg.UpdateTree(parser.GetTree(), ifKeyword, ifKeyword, formattedContent.String())
			// parser.UpdateTree([]byte(formattedContent.String()))
		case "block_start":
			block_level++
		case "block_end":
			block_level--
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

		default:
			logger.Printf("Node type: %v -- continue\n", node.Type())
		}
	}
	logger.Println("Formatting Kamailio configuration file ", edits)
	TraverseNodeAndApply(root, formatNode)
	return edits, nil
}
