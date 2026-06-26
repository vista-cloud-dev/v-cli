// Command v is the VistA developer-tools CLI — a single umbrella that wraps each
// insider VistA subsystem in a plain-language domain (v-cli-platform.md). Each
// domain lives in its own repo and is composed in statically: `v` imports the
// domain as a Go module and mounts its command surface (CQ1 static-pinned). The
// first domain is `pkg` (the KIDS lifecycle), from v-pkg.
//
//	v pkg decompose OR_3.0_484.KID ./patches/
//	v pkg roundtrip OR_3.0_484.KID
//	v new <domain>          # scaffold a new v domain from the built-in skeleton
//
// Composition: `v` and every mounted domain share one clikit.Context type from
// the standalone github.com/vista-cloud-dev/clikit module (extracted from
// v-pkg/clikit 2026-06-25). That shared type is what lets the umbrella mount more
// than one domain — a second domain mounts here as another named CLI field with
// no further clikit work.
package main

import (
	"os"

	"github.com/willabides/kongplete"

	"github.com/vista-cloud-dev/clikit"
	"github.com/vista-cloud-dev/v-pkg/pkgcli"
)

// CLI is the umbrella grammar: one named field per domain (mounted as that
// domain's subcommand) plus the shared clikit meta commands.
type CLI struct {
	clikit.Globals

	Pkg pkgcli.Commands `cmd:"" name:"pkg" help:"VistA package (KIDS) tools: decompose / assemble / roundtrip / canonicalize / lint."`

	New NewCmd `cmd:"" help:"Scaffold a new v domain tool from a built-in skeleton."`

	Schema  clikit.SchemaCmd  `cmd:"" help:"Emit the aggregated command/flag/enum tree as JSON (agent discovery)."`
	Version clikit.VersionCmd `cmd:"" help:"Show version and build info."`

	InstallCompletions kongplete.InstallCompletions `cmd:"" help:"Install shell tab-completions."`
}

func main() {
	cli := &CLI{}
	os.Exit(clikit.Run(
		"v",
		"v — VistA developer tools (pkg / …): plain-language wrappers over insider VistA subsystems.",
		cli, &cli.Globals,
	))
}
