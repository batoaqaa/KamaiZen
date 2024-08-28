package kamailio_cfg

// #cgo CFLAGS: -std=c11 -fPIC  -fcommon
// #include "include.h"
// // #cgo LDFLAGS: -Wl,--allow-multiple-definition
// // NOTE: if your language has an external scanner, add it here.
import "C"

import "unsafe"

// Get the tree-sitter Language for this grammar.
func Language() unsafe.Pointer {
	return unsafe.Pointer(C.tree_sitter_kamailio_cfg())
}
