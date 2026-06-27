---
name: rpc-domain-mount
description: v-cli gained its SECOND domain — `v rpc` (v-rpc@v0.1.0) — mounted alongside `v pkg` on 2026-06-26. The proof that the umbrella's static-pinned multi-domain composition works as designed.
metadata:
  type: project
---

**`v rpc` mounted into the umbrella (2026-06-26) — v-cli's second domain.**
This is the first time the umbrella composes more than one domain, proving the
static-pinned design (one clikit.Context type shared by umbrella + all domains).

**The mount is exactly three edits + a pin** (the documented recipe for any future
domain):
1. Domain side: `rpccli/contract.go` with `Contract() vcontract.Manifest` (mirror
   `v-pkg/pkgcli/contract.go`) — imports `github.com/vista-cloud-dev/v-pkg/vcontract`
   for the shared `Manifest` type. A v→v dependency (both VistA-layer), waterline-OK.
2. `v-cli/main.go`: add one field — `Rpc rpccli.Commands \`cmd:"" name:"rpc"
   group:"Domains" ...\`` — next to `Pkg`.
3. `v-cli/registry.go`: add `rpccli.Contract()` to `buildRegistry().Domains`, then
   `UPDATE_GOLDEN=1 go test -run TestRegistry_Golden` to regenerate
   `dist/v-registry.json` (the §5 drift gate).
4. `go get github.com/vista-cloud-dev/v-rpc@v0.1.0` (static pin, NO `replace`).

**Version coordination:** v-rpc had to first repin clikit v0.4.0 (umbrella + domain
must share the same clikit.Context type) and be tagged v0.1.0 before v-cli could
pin it. See [[clikit-shared-module]] / [[clikit-grouped-help]].

**GOTCHA — fetching a fresh private-repo tag:** the default
`GOPROXY=https://proxy.golang.org,direct` can't authenticate to the private org and
fails with `fatal: could not read Username for 'https://github.com'`. Fix once:
`go env -w GOPRIVATE=github.com/vista-cloud-dev` + `gh auth setup-git` (gh is logged
in, so git then uses its token for `direct` fetches).

**The PATH `v` is the root-built `v-cli/v`** (symlinked into `~/scripts/bin/`), NOT
`dist/v` — rebuild with `go -C v-cli build -o v .` after a mount, or the symlinked
`v` stays stale (`make check` only builds `dist/v`).
