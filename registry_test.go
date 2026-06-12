package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestRegistry_Golden is the §5 drift gate: the committed dist/v-registry.json
// must match the aggregate of the pinned domains' contracts. Regenerate with
// `make registry` (UPDATE_GOLDEN=1).
func TestRegistry_Golden(t *testing.T) {
	got, err := json.MarshalIndent(buildRegistry(), "", "  ")
	if err != nil {
		t.Fatalf("marshal registry: %v", err)
	}
	got = append(got, '\n')

	golden := filepath.Join("dist", "v-registry.json")
	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.MkdirAll(filepath.Dir(golden), 0o755); err != nil {
			t.Fatalf("mkdir dist: %v", err)
		}
		if err := os.WriteFile(golden, got, 0o644); err != nil {
			t.Fatalf("write golden: %v", err)
		}
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatalf("read golden (run `make registry`): %v", err)
	}
	if string(got) != string(want) {
		t.Errorf("dist/v-registry.json drift — run `make registry`\n--- got ---\n%s", got)
	}
}

func TestRegistry_Invariants(t *testing.T) {
	r := buildRegistry()
	if r.SchemaVersion == "" || r.CLI != "v" {
		t.Errorf("registry header = %+v", r)
	}
	if len(r.Domains) == 0 {
		t.Fatal("registry has no domains")
	}
	found := false
	for _, d := range r.Domains {
		if d.Domain == "pkg" {
			found = true
			if len(d.Commands) == 0 {
				t.Error("pkg domain has no commands in the registry")
			}
		}
	}
	if !found {
		t.Error("registry missing the pkg domain")
	}
}
