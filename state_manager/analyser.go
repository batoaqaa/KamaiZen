package state_manager

import (
	"KamaiZen/document_manager"
	"KamaiZen/kamailio_cfg"
	"KamaiZen/lsp"
	"log"

	sitter "github.com/smacker/go-tree-sitter"
)

// TODO: major refactoring required

type StateTree struct {
	nodes map[lsp.DocumentURI]*sitter.Node
}

// module level state tree
var stateTreeCache StateTree

// NewStateTree creates and returns a new instance of StateTree.
// It initializes the nodes map.
//
// Returns:
//
//	StateTree - The initialized state tree.
func NewStateTree() StateTree {
	return StateTree{
		nodes: make(map[lsp.DocumentURI]*sitter.Node),
	}
}

// AddNode adds a node to the state tree for the given document URI.
//
// Parameters:
//
//	uri lsp.DocumentURI - The URI of the document.
//	node *sitter.Node - The AST node to be added.
func (s *StateTree) AddNode(uri lsp.DocumentURI, node *sitter.Node) {
	s.nodes[uri] = node
}

// TraverseNode traverses the AST starting from the given node and logs the traversal.
// It recursively visits all child nodes and logs their information.
//
// Parameters:
//
//	uri lsp.DocumentURI - The URI of the document.
//	node *sitter.Node - The starting node for the traversal.
//	logger *log.Logger - The logger used to log the traversal.
//	padding int - The padding used for indentation.
func (s *StateTree) TraverseNode(uri lsp.DocumentURI, node *sitter.Node, logger *log.Logger, padding int) {
	// traverse the node and print the node
	var i uint32
	childCount := node.ChildCount()
	for i = 0; i < childCount; i++ {
		// Print spaces for padding
		child := node.Child(int(i))
		s.TraverseNode(uri, child, logger, padding+2)
	}
}

// GetFunctionNameAtPosition returns the name of the function at the given position in the document.
//
// Parameters:
//
//	uri lsp.DocumentURI - The URI of the document.
//	position lsp.Position - The position within the document.
//	source_code []byte - The source code of the document.
//
// Returns:
//
//	string - The name of the function at the given position.
func GetFunctionNameAtPosition(uri lsp.DocumentURI, position lsp.Position, source_code []byte) string {
	node := GetState().Analyzer.GetAST().Node
	return getFunctionName(node, position, source_code)
}

// getFunctionName returns the name of the function at the given position in the AST node.
//
// Parameters:
//
//	node *sitter.Node - The AST node to be searched.
//	position lsp.Position - The position within the node.
//	source_code []byte - The source code of the document.
//
// Returns:
//
//	string - The name of the function at the given position.
func getFunctionName(node *sitter.Node, position lsp.Position, source_code []byte) string {
	if node == nil {
		return ""
	}
	nodeAtPosition := node.NamedDescendantForPointRange(
		sitter.Point{
			Row:    uint32(position.Line),
			Column: uint32(position.Character),
		},
		sitter.Point{
			Row:    uint32(position.Line),
			Column: uint32(position.Character),
		})
	functionName := nodeAtPosition.Content(source_code)
	return functionName
}

// TraverseNodeAndApply traverses the AST starting from the given node and applies the given function to each node.
// It recursively visits all child nodes and applies the function to each node.
//
// Parameters:
//
//	node *sitter.Node - The starting node for the traversal.
//	f func(*sitter.Node) - The function to be applied to each node.
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

// getAllAvailableKeywords returns a map of all available keywords and their descriptions.
//
// Returns:
//
//	map[string]string - A map of keywords and their descriptions.
func getAllAvailableKeywords() map[string]string {
	// right now, we are hardcoding the keywords and using only SIP headers
	return kamailio_cfg.SIPHeaders
}

// GetCompletionItems returns a list of completion items for the given document URI.
//
// Parameters:
//
//	uri lsp.DocumentURI - The URI of the document.
//
// Returns:
//
//	[]lsp.CompletionItem - A list of completion items.
func GetCompletionItems(uri lsp.DocumentURI) []lsp.CompletionItem {
	var completionItems []lsp.CompletionItem
	functions := document_manager.GetAllAvailableFunctionDocs()
	for _, function := range functions {
		completionItems = append(completionItems, lsp.CompletionItem{
			Detail:        function.Name + "(" + function.Parameters + ")",
			Label:         function.Name + "(" + ")",
			Documentation: function.Description + "\n" + function.Example,
			Kind:          lsp.FUNCTION_COMPLETION,
		})
	}

	keywords := getAllAvailableKeywords()
	for header, description := range keywords {
		completionItems = append(completionItems, lsp.CompletionItem{
			Detail:        "SIP Header",
			Label:         header,
			Documentation: description,
			Kind:          lsp.VARIABLE_COMPLETION,
		})
	}

	variables := kamailio_cfg.GetGlobalVariables()
	for variable, value := range variables {
		completionItems = append(completionItems, lsp.CompletionItem{
			Detail:        "AVP",
			Label:         variable,
			Documentation: value.GetGlobalVariableDocs(),
			Kind:          lsp.VARIABLE_COMPLETION,
		})
	}

	modules := document_manager.GetAllAvailableModules()
	for _, module := range modules {
		completionItems = append(completionItems, lsp.CompletionItem{
			Detail:        "Module",
			Label:         module,
			Documentation: "Module " + module,
			Kind:          lsp.MODULE_COMPLETION,
		})
	}
	return completionItems

}
