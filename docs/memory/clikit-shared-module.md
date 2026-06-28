---
name: clikit-shared-module
description: clikit is now a standalone module (github.com/vista-cloud-dev/clikit); v-cli + v-pkg consume it, unblocking multi-domain
metadata:
  type: project
---

# clikit extracted to a standalone shared module (2026-06-25)

`clikit` (the shared CLI convention layer — Kong grammar + `Run`, `--output
text|json|auto` envelope, exit-code ladder, `schema`/`version`, lipgloss styling,
kongplete completions) is now its own module: **`github.com/vista-cloud-dev/clikit`
v0.1.0** (public, Apache-2.0, branch `main`). Extracted from the canonical
`go-cli-template/clikit` copy.

**Why:** the `v` umbrella composes domains in-process, and Go requires all mounted
domains to share **one** `clikit.Context` type. While each repo vendored its own
`clikit`, the umbrella could mount only one domain (it borrowed `v-pkg/clikit`).
The shared module removes that blocker — a second `v` domain now mounts as another
named CLI field in `v-cli/main.go` with no further clikit work.

**Consumers (v-family migrated; m-family not yet):**
- `v-pkg` **≥ v0.3.0** consumes it (deleted its `clikit/` dir; v0.3.0 was a breaking
  Go API bump — the `Context` type identity moved — but the command surface/13 verbs
  are unchanged so its `ContractVersion` stayed `1.0`).
- `v-cli` pins `clikit v0.1.0` + `v-pkg v0.3.0`.
- `m-cli`, `m-dev-tools-mcp`, `go-cli-template` **still vendor** their own byte-identical
  copies — they're single-binary so they don't need the shared type; migrate them when
  next touched. (Scope chosen 2026-06-25: v-family only.)

**Two gotchas worth not rediscovering:**
1. **Makefile `-ldflags -X` path** must target the shared module
   (`github.com/vista-cloud-dev/clikit.Version`), not the old `$(PKG)/clikit`. And a
   trailing `# comment` after a `LDPKG := …` value injects trailing **spaces** into
   the make variable, which corrupts the `-X importpath.name=value` flag (`link: -X
   flag requires argument of the form …`). Put the comment on its own line.
2. **A freshly created GitHub repo/tag isn't on `proxy.golang.org` immediately.** To
   pin it, `go get …@vX` via the proxy fails with "unknown revision"; use
   `go mod edit -require=github.com/vista-cloud-dev/clikit@v0.1.0` then
   `GOFLAGS=-mod=mod GOPROXY=direct go mod tidy`.

See [[../v-cli-platform.md]] §6 and the [status report](../archive/v-cli-status-report-2026-06-25.md) finding #3.
