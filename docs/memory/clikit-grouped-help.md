---
name: clikit-grouped-help
description: v-cli repinned to clikit v0.2.0 + v-pkg v0.4.0 and tagged its umbrella commands with group:"" — `v help` shows Domains/Scaffold/Introspect, and `v pkg --help` shows v-pkg's KIDS-lifecycle groups through the umbrella. 2026-06-26.
metadata:
  type: project
---

**v-cli umbrella picked up clikit's grouped help (2026-06-26).** Repinned
`clikit v0.1.0 → v0.2.0` (styled grouped help renderer + pager; see clikit's
`cli-discovery-ux`) **and** `v-pkg v0.3.0 → v0.4.0` (the version that itself pins
clikit v0.2.0 and carries the KIDS group tags). Added `group:""` tags to the
umbrella `CLI` struct:

- **Domains**: pkg
- **Scaffold**: new
- **Introspect**: schema, version
- (install-completions → trailing "Commands" bucket)

**Why both pins move together:** v-cli mounts `pkgcli.Commands` in-process and
shares the `clikit.Context` type, so v-cli and v-pkg must agree on the clikit
version — pinning v-pkg v0.4.0 (built on clikit v0.2.0) keeps one Context type AND
surfaces v-pkg's own groups under `v pkg`. This is exactly the multi-domain
coordination the standalone clikit extraction was for.

Verified end-to-end: `v help` → umbrella groups; `v pkg --help` → v-pkg's
Inspect/Transform/Build & install/Back-out groups (breadcrumb title "v pkg").
Gates: build/vet clean, `go test ./...` green.

This completes the clikit discovery-UX Phase 1 rollout (clikit → m-cli → v-pkg →
v-cli).

**Phase 2 — `v explore` (2026-06-26).** Repinned clikit v0.2.0 → v0.3.2 and
v-pkg v0.4.0 → v0.5.0, and mounted `Explore clikit.ExploreCmd` in the umbrella
root (Introspect group) → `v explore` opens the interactive palette over the whole
`v` tree (incl. `pkg`). Deliberately NOT in `pkgcli.Commands`, so there is no
confusing `v pkg explore` (verified: 0 hits). Build/vet/test green; mounting +
non-TTY fallback smoke-tested. Phase 2 rollout complete across clikit/m-cli/
v-pkg/v-cli. Phase 3 (`browser`) remains a sketch.
