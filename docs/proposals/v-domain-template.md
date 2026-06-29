---
title: "v-domain-template — a thin, drift-gated single source for `v new` domain scaffolds"
status: draft
version: v0.1.0
created: 2026-06-29
last_modified: 2026-06-29
doc_type: [PROPOSAL]
layer: v
---

# v-domain-template — a thin, drift-gated single source for `v new` domain scaffolds

> **Decision sought.** Replace `v new`'s hand-written embedded skeleton
> (`v-cli/new.go` `scaffoldDomain`) with a real, buildable, CI-gated, drift-gated
> template repo — **`v-domain-template`** — that `v new` *renders from* rather than
> *competes with*. This **re-opens CQ4** (`docs/v-cli-platform.md` §6/§10), which on
> 2026-06-25 settled on "`v new` generator only, no template repo." This proposal
> argues the reversal is warranted because the embedded skeleton has measurably
> drifted, and shows a structure that keeps the **one-source-of-truth** property the
> 2026-06-25 decision was protecting.

---

## 1. Motivating incident

On 2026-06-29 the `v-db` repo was retired to a placeholder. Comparing it against a
fresh `v new` scaffold surfaced that **`v new` produces an incomplete, silently
non-conforming skeleton**. It emits five files:

```
go.mod   README.md   docs/README.md   Makefile   <domain>cli/commands.go
```

A convention-complete `v` domain (measured against `v-pkg`, `go-cli-template`, and
the org `CLAUDE.md` rules) additionally needs:

| File | Why it's required | In `v new`? |
|---|---|---|
| `.envrc` (`source_up`) | **Engine access** — without it direnv never loads `~/data/vista-cloud-dev/auth.env`, so the repo session can't see `M_IRIS_*`/`VISTA_*` creds and can't reach the m engines. | ❌ |
| `repo.meta.json` (`layer: v`) | `m arch check` G1–G5 + the org meta-gate completeness scan key off the layer artifact. | ❌ |
| `.golangci.yml` | Org P1 drift gate made the canonical config **mandatory**; go-ci is red without it. | ❌ |
| `.go-version` (`1.26.3`) | go-ci resolves the toolchain via `go-version-file`. | ❌ |
| `CLAUDE.md` | Per-repo waterline / engine / cred rules. | ❌ |
| `.gitignore`, `LICENSE`, dependabot, CI workflows | Org baseline. | ❌ |

A new domain scaffolded today is therefore born **unable to call the engines and
red on CI** until a human re-adds ~8 files by hand — the exact "pile of bespoke
CLIs" outcome `v new` exists to prevent (§6).

**Root cause:** the skeleton lives as Go string literals in `scaffoldDomain`. A
string map is invisible to every gate the org built — `m arch check`, the
`.golangci.yml` drift gate, `docs-validate`, `go test` — so it drifted from the
conventions with nothing to catch it.

## 2. Why the 2026-06-25 "no template repo" decision doesn't settle this

CQ4 originally (2026-06-11) called for **both** a `v-tool-template` repo **and** a
`v new` generator. That standalone repo was created and then **deleted** (commit
`94b96d2`); §6 records the rationale: *"one source of truth for the scaffold and no
separate repo to keep in sync."*

That rationale was sound against the original design — **two independently
hand-edited sources** (a template repo *and* a generator) genuinely would drift
apart. But it threw out the repo and kept the *worse* of the two sources: the
generator's hand-written strings, which then drifted anyway (§1). The 2026-06-25
decision protected "one source of truth"; what it actually produced is one source
that **no gate can see**.

The fix is not "two sources" (the thing that was rightly rejected) and not "the
string map" (the thing that drifted). It is **one source that is a real repo**, with
the generator reduced to a *renderer* over it — no second authored copy to sync.

## 3. Proposal

### 3.1 Create `v-domain-template`

A new repo `github.com/vista-cloud-dev/v-domain-template`: the **single authored
source** for a `v` domain skeleton. It is a *real, buildable, CI-green* repo (it
compiles, `make check` passes, `m arch check .` passes), not a bag of strings — so
every org gate applies to it directly.

It is **thin**: it contains only what is generic to *every* `v` domain, with
substitution tokens where a value is domain-specific. It is **not** a fork of a
working product (`v-pkg` is 191 files of KIDS-specific code; using it as a template
means subtracting 180+ files on every scaffold — rejected, §5).

Layering (each layer adds only its delta; lower layers are the canonical source for
their files):

```
go-cli-template          → generic Go-CLI baseline: .golangci.yml, .go-version,
   (canonical, gated)        Makefile, .gitignore, LICENSE, dependabot, CI workflows
        │  (drift-gated: §4 G-overlay)
        ▼
v-domain-template        → adds the v-overlay: .envrc (source_up), repo.meta.json
   (this proposal)           (layer: v), CLAUDE.md, the <domain>cli/ command-package
                             shape, docs/ skeleton, conformance-suite stub
        │  (drift-gated: §4 G-render)
        ▼
v new <domain>           → renders v-domain-template with the domain substitutions;
   (thin renderer)           emits a ready-to-build, engine-capable, CI-green repo
```

### 3.2 Repo contents (the thin skeleton)

Inherited unchanged from `go-cli-template` (NOT re-authored here — see §4 G-overlay):
`.golangci.yml`, `.go-version`, `.gitignore`, `LICENSE`, `.github/dependabot.yml`,
`.github/workflows/ci.yml`, `.github/workflows/docs-validate.yml`.

Added by `v-domain-template` (the v-overlay — the part `go-cli-template` lacks
because it is engine-/layer-neutral):

- `.envrc` → exactly `source_up` (one line).
- `repo.meta.json` → `layer: "v"`, with `id`/`repo`/`role` as substitution tokens.
- `CLAUDE.md` → the per-repo waterline + engine-access + cred rules (the short form
  every engine-touching repo carries).
- `Makefile` → the standard `build`/`test`/`lint`/`check` plus `arch` (`m arch
  check .`) and `gates`.
- `<DOMAIN>cli/commands.go` → the importable `Commands` package the umbrella mounts.
- `README.md`, `docs/README.md` → the standard front-door + the org docs-layout
  index (what `v new` already emits, kept here so it is gate-checked).
- `docs/memory/`, `docs/proposals/`, `docs/design/`, `docs/archive/` → created
  on first use (the docs-README already documents this; no empty dirs committed).

### 3.3 Substitution tokens

A small, explicit token set rendered by `v new` (and left literal in the template so
the template itself builds — e.g. a `DOMAIN=template` default):

| Token | Example | Appears in |
|---|---|---|
| `<domain>` | `db` | module path, package name, command label, README |
| `<DOMAINcli>` | `dbcli` | package dir + name |
| `<module>` | `github.com/vista-cloud-dev/v-db` | `go.mod`, `repo.meta.json.repo` |
| `<role>` | "FileMan database access — the `v db` domain" | `repo.meta.json.role` |

The plain-language gate (§7) still runs on `<domain>` at scaffold time — unchanged.

### 3.4 How `v new` consumes it — airgapped-safe, no second authored copy

`v new` must keep working with **no network and no template repo on disk** (the
property that made the embedded skeleton attractive). So:

- `v-cli` **vendors a snapshot** of `v-domain-template` and `go:embed`s it; `v new`
  renders the embedded snapshot with the §3.3 substitutions.
- The snapshot is **generated, never hand-edited** — a `make sync-template` target
  refreshes it from a pinned `v-domain-template` tag, and a **drift gate**
  (§4 G-render) red-fails CI if the vendored snapshot ≠ that tag.

This is the crux of the reconciliation with 2026-06-25: there is still exactly **one
authored source** (`v-domain-template`). The embedded copy in `v-cli` is a *build
artifact* of it, mechanically gated — not a second thing a human edits. The failure
mode that killed the original `v-tool-template` (two hand-maintained sources
diverging) is gated away, not reintroduced.

## 4. Drift gates (this is what makes it trustworthy, per the org "generate → registry → red-gate" rule)

- **G-overlay** — `v-domain-template`'s inherited files are byte-identical to
  `go-cli-template`'s. (Same gate already used org-wide for `.golangci.yml`; extend
  the set.) Prevents the template from forking the baseline.
- **G-render** — `v-cli`'s vendored/embedded snapshot is byte-identical to the
  pinned `v-domain-template` tag (modulo the substitution tokens). Prevents `v new`
  from drifting from the template — the §1 failure, now impossible to land silently.
- **G-buildable** — `v-domain-template` itself stays CI-green (`make check`,
  `m arch check .`), so "the template" is provably a conforming domain at all times.

## 5. Alternatives considered

| Option | Verdict | Why |
|---|---|---|
| **A. Complete the embedded strings + add a drift gate, no repo** | Viable, lighter | Closes the §1 gap with the least machinery (no new repo). But the "source of truth" stays a string map — readable only as Go literals, not buildable, not directly gate-checkable as a repo; the drift gate would compare strings to `go-cli-template` files, which is awkward and partial. Good fallback if a new repo is judged not worth it. |
| **B. Use `v-pkg` as the template** | Rejected | Right *shape* (a domain) but it is a 191-file working product; its `.golangci.yml` is a *downstream copy* of `go-cli-template` (sourcing from a copy inverts the SSOT); it lacks `CLAUDE.md`; product churn would leak into scaffolds. A reference, not a template. |
| **C. Use `v-cli` as the template** | Rejected | Wrong shape — `v-cli` is the umbrella binary (registry/dispatch/`new`), not a domain. Templating off it ships umbrella machinery a domain must strip. |
| **D. Status quo (embedded skeleton, no gate)** | Rejected | The §1 incident is the cost: domains born engine-incapable and CI-red, caught only by chance. |

**A vs the proposal** is the real decision. The proposal costs one thin repo + two
drift gates; in return the scaffold source is a *real domain* that every existing
gate already validates, and "is the template still correct?" becomes a CI question
instead of a code-review question. Recommended unless the extra repo is unwanted, in
which case A (complete + gate the embedded set) is the minimum acceptable fix.

## 6. Rollout

1. Create `v-domain-template` from `go-cli-template` + the v-overlay; make it
   CI-green and `m arch check`-clean (G-buildable).
2. Wire **G-overlay** in `v-domain-template` CI.
3. In `v-cli`: replace `scaffoldDomain`'s string map with `go:embed` of the vendored
   snapshot + token substitution; add `make sync-template` and **G-render**. TDD
   against `new_test.go` (assert the full file set, including `.envrc`,
   `repo.meta.json`, `.golangci.yml`, `CLAUDE.md`, `.go-version`).
4. Register `v-domain-template` in `.github/ecosystem.json` (meta-gate completeness).
5. Update `v-cli-platform.md` §6/§10: amend the CQ4 note to record this reversal and
   its rationale (one gated source; `v new` is a renderer).
6. Retire this proposal to `docs/proposals/` accepted-state (or `implemented/`) per
   the Increment Protocol when it lands.

## 7. Open questions

- **OQ-1.** Vendor-and-embed (§3.4) vs. `v new` reading a known-local template path?
  Embed favors the airgapped/standalone property; a local path is simpler but
  reintroduces a "must be present" dependency. *Lean: embed.*
- **OQ-2.** Should `go-cli-template` absorb a `layer`-tagged `repo.meta.json` stub so
  the overlay shrinks further, or does the engine-/layer-neutral baseline stay
  clean? *Lean: keep `go-cli-template` neutral; the overlay owns `repo.meta.json`.*
- **OQ-3.** Does `m new` (the M-side scaffolder) have the same embedded-drift
  problem, and should the same template+gate pattern apply there? Out of scope here,
  worth a sibling check.
- **OQ-4.** Is the lighter Alternative A preferable for this solo/airgapped setup,
  given a new repo is itself a maintenance surface? This is the decision to make.
