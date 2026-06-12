package main

import (
	"github.com/vista-cloud-dev/v-pkg/pkgcli"
	"github.com/vista-cloud-dev/v-pkg/vcontract"
)

// Registry is the `v` umbrella's aggregated command surface (v-cli-platform.md
// §5): the union of the statically-pinned domains' contract manifests. It is the
// generated, drift-gated surface that `v help`, shell completion, and dispatch
// read — `v` never hand-maintains its command list. Because composition is
// in-process (CQ1), each domain's contract is obtained by calling its
// Contract() directly, so the registry can never advertise a command a pinned
// domain no longer provides.
type Registry struct {
	SchemaVersion string               `json:"schemaVersion"`
	CLI           string               `json:"cli"`
	Domains       []vcontract.Manifest `json:"domains"`
}

// buildRegistry aggregates each pinned domain's contract. New domains are added
// here as they are pinned into the umbrella (one line per domain).
func buildRegistry() Registry {
	return Registry{
		SchemaVersion: "1.0",
		CLI:           "v",
		Domains: []vcontract.Manifest{
			pkgcli.Contract(),
		},
	}
}
