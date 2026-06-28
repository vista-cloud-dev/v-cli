package main

import (
	"strings"
	"testing"
)

// TestScaffoldDomain_DocsBornCorrect asserts a v-domain scaffolded by `v new`
// is born with the standard vista-cloud-dev docs/ layout (docs-organization
// remediation OQ-5): a docs/README.md index that documents the standard folder
// vocabulary, so a new domain never starts with an empty or ad-hoc docs/ tree.
func TestScaffoldDomain_DocsBornCorrect(t *testing.T) {
	files, err := scaffoldDomain("db")
	if err != nil {
		t.Fatalf("scaffoldDomain: %v", err)
	}

	readme, ok := files["docs/README.md"]
	if !ok {
		t.Fatalf("scaffold is missing docs/README.md; got keys %v", keysOf(files))
	}

	// The index must name the standard folder vocabulary so forkers don't invent
	// bespoke folders (the exact rot the remediation removed).
	for _, want := range []string{"guides", "design", "memory", "archive"} {
		if !strings.Contains(readme, want) {
			t.Errorf("docs/README.md does not mention standard folder %q", want)
		}
	}
}

func keysOf(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
