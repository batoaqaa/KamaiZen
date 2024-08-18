package kamailio_cfg_test

import (
	"testing"

	"KamaiZen/kamailio_cfg"
	tree_sitter "github.com/smacker/go-tree-sitter"
)

func TestCanLoadGrammar(t *testing.T) {
	language := tree_sitter.NewLanguage(kamailio_cfg.Language())
	if language == nil {
		t.Errorf("Error loading KamailioCfg grammar")
	}
}
