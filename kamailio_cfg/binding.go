package kamailio_cfg

// #cgo CFLAGS: -std=c11 -fPIC
// #cgo LDFLAGS: -Wl,--allow-multiple-definition
// #include "parser.c"
import "C"

import "unsafe"

// Get the tree-sitter Language for this grammar.
func Language() unsafe.Pointer {
	return unsafe.Pointer(C.tree_sitter_kamailio_cfg())
}
