package main

import (
	"testing"

	"github.com/alecthomas/kong"
)

// TestPlainLanguage_CommandSurface is the conformance plain-language gate
// (v-cli-platform.md §7): no command or flag name in the whole `v` surface —
// including every mounted domain (v pkg …) — may contain insider VistA
// vocabulary. This is the family's reason to exist, enforced mechanically so no
// tool ever re-exposes the terms the platform hides.
func TestPlainLanguage_CommandSurface(t *testing.T) {
	parser, err := kong.New(&CLI{})
	if err != nil {
		t.Fatalf("kong.New: %v", err)
	}
	var names []string
	var walk func(n *kong.Node)
	walk = func(n *kong.Node) {
		if n.Name != "" {
			names = append(names, n.Name)
		}
		for _, f := range n.Flags {
			if f.Name != "" {
				names = append(names, f.Name)
			}
		}
		for _, c := range n.Children {
			walk(c)
		}
	}
	walk(parser.Model.Node)

	if v := plainLangViolations(names); len(v) > 0 {
		t.Errorf("the v command surface contains insider VistA vocabulary: %v", v)
	}
}

func TestPlainLangViolations_Detects(t *testing.T) {
	got := plainLangViolations([]string{"db", "fileman", "config", "xpar"})
	if _, ok := got["db"]; ok {
		t.Error("db is plain-language, must not be flagged")
	}
	if got["fileman"] == nil || got["xpar"] == nil {
		t.Errorf("fileman/xpar must be flagged as vista-ese; got %v", got)
	}
}
