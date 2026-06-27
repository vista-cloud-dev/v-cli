# v-cli memory index

- [clikit grouped help](clikit-grouped-help.md) — 2026-06-26: repinned clikit v0.2.0 + v-pkg v0.4.0 and tagged umbrella commands with `group:""` (Domains/Scaffold/Introspect). `v help` shows umbrella groups; `v pkg --help` shows v-pkg's KIDS-lifecycle groups through the umbrella. Completes the clikit discovery-UX Phase 1 rollout. See [[clikit-shared-module]].
- [clikit shared module](clikit-shared-module.md) — clikit extracted to github.com/vista-cloud-dev/clikit v0.1.0; v-cli + v-pkg(≥0.3.0) consume it, unblocking multi-domain; Makefile -ldflags + fresh-repo go-get gotchas.
- [rpc domain mount](rpc-domain-mount.md) — 2026-06-26: v-cli's SECOND domain `v rpc` (v-rpc@v0.1.0) mounted alongside `v pkg`. The 3-edits-+-pin recipe (Contract() / main.go field / registry+golden) for any future domain; GOPRIVATE+`gh auth setup-git` gotcha for fetching fresh private tags; rebuild root `v-cli/v` (PATH symlink), not just `dist/v`.
