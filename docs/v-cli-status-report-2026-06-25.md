# v-cli — comprehensive status report

**Date:** 2026-06-25
**Repo:** `github.com/vista-cloud-dev/v-cli` · binary `v` · layer `v` · Go 1.26.3 · AGPL-3.0
**Branch:** `main` (clean) · 10 commits · tests green (`go test ./...` → ok)

---

## 0. Executive summary

`v` is the **VistA developer-tools umbrella CLI** — a single Go binary that
aggregates VistA subsystem tooling under plain-language commands (`v pkg`,
`v db`, `v config`, …), the **VistA-side mirror of `m`** (the engine-neutral M
toolchain busybox). Its reason to exist is one sentence from the platform spec:
*wrap vista-ese in modern plain language so that insider knowledge isn't required
to be productive.* KIDS/FileMan/XPAR/Broker become `v pkg` / `v db` / `v config`
/ `v rpc`, and a mechanical lint forbids any insider VistA term from ever
appearing in a command or flag name.

**Where it stands today:** v-cli is a **working, well-disciplined single-domain
proof of the platform pattern.** It mounts exactly one real domain — `v pkg`
(KIDS lifecycle, from the `v-pkg` repo) — plus the meta commands `new`, `schema`,
`version`, `install-completions`. The umbrella mechanics that matter (static
in-process composition, a generated drift-gated registry, the plain-language
conformance gate, the `v new` scaffolder, the waterline layer tag + CI gate) are
all built and green. The repo is ~140 lines of Go over four small source files;
its weight is almost entirely in the discipline, not the code.

**The single thing blocking growth:** mounting a **second** domain requires first
extracting `clikit` (the shared CLI convention kit) into a standalone module.
Today v-cli borrows `v-pkg`'s vendored `clikit`, which works for one domain but
not two (each domain would carry a divergent `clikit.Context` type). This is a
known, recorded prerequisite — flagged in `main.go:11-15` — not a surprise.

**Maturity verdict:** the *pattern* is proven and the *governance* is unusually
complete for the repo's age; the *surface* is at "first command" stage. v-cli
leads `m-cli` on registry-as-a-file and on having a scaffolder, but trails it
massively on breadth (1 domain vs m-cli's ~14 command namespaces).

---

## 1. What v-cli is

A **busybox-style umbrella**: one binary, many "applets" (domains), each a plain-
language wrapper over an insider VistA subsystem. The model is explicit in the
package doc (`main.go:1-15`) and the platform spec:

> `v` is a single umbrella CLI (exactly as `m` is for the M toolchain):
> `v <domain> <verb> [args] [--flags]`.

The naming discipline is the heart of the project. The prefix split across the
whole org is **by scope, not language** (both families are Go):

| Prefix | Means | Assumes VistA installed? |
|---|---|---|
| `m-*` | engine-neutral M toolchain | No — runs on a bare M engine |
| `v-*` | VistA-specific | Yes — needs Kernel/FileMan/KIDS |

`v pkg` is the canonical proof: KIDS is meaningless against a bare engine, so the
tool — though pure Go — is VistA-specific and lives in the `v` family (the repo
was renamed `m-kids` → `v-pkg` for exactly this reason).

### The plain-language rule (the family's reason to exist)

No domain, command, or flag name may contain insider VistA vocabulary. The
forbidden list lives in `plainlang.go:9-13`:

```
fileman, kids, xpar, mailman, taskman, duz, dic, die, diq, dik,
zis, ztload, xusec, xushsh, vistalink, rpcbroker, kernel, hl7
```

So: KIDS → `pkg`, FileMan → `db`, XPAR → `config`, Broker → `rpc`, TaskMan →
`job`, MailMan → `mail`, Device Handler → `io`. The VA product name stays in
docs only; it never reaches a command line. This is enforced **mechanically**,
not by review — see §4 (gates).

---

## 2. What the code actually does right now

The repo is four source files plus tests. Concretely:

### `main.go` — the umbrella grammar
Defines the `CLI` kong struct (`main.go:29-40`): embeds `clikit.Globals` and
mounts each domain as a named field. Today there is exactly one domain field —
`Pkg pkgcli.Commands` — plus `New`, `Schema`, `Version`, and
`InstallCompletions`. `main()` calls `clikit.Run("v", …, cli, &cli.Globals)`,
delegating all parsing/output/exit-code handling to the shared kit.

The header comment (`main.go:11-15`) records the **single-domain composition
constraint**: `pkgcli`'s `Run` methods take a `*v-pkg/clikit.Context`, so `v`
must use `v-pkg/clikit`. Correct for one domain; a second domain (with its own
vendored `clikit` type) requires extracting `clikit` into a shared module first.

### `registry.go` — the aggregated, drift-gated command surface
`buildRegistry()` (`registry.go:23-31`) builds a `Registry{schemaVersion, cli,
domains[]}` by calling each pinned domain's `Contract()` **directly in-process**
(today just `pkgcli.Contract()`). Because composition is static and in-process
(CQ1), the registry **cannot advertise a command a pinned domain no longer
provides** — it reflects the live mounted surface. Serialized to the committed,
gitignored-except-`*.json` artifact `dist/v-registry.json`.

### `new.go` — `v new <domain>` scaffolder
`NewCmd.Run` (`new.go:22-55`) scaffolds a new domain repo `v-<domain>/` with a
minimal, convention-conforming skeleton built **into the binary** (there is no
separate template repo — the former `v-tool-template` was deleted, commit
`94b96d2`). It first runs the plain-language check on the proposed domain name
(rejecting e.g. `v new fileman` with a `VISTA_ESE` usage error), then emits
`go.mod`, `README.md`, `Makefile`, and `<domain>cli/commands.go` (an empty
`Commands` struct ready to grow verbs). On success it prints the next step:
"pin it in the v umbrella's go.mod and mount `<domain>cli.Commands`."

### `plainlang.go` — the conformance lint engine
`plainLangViolations(names)` (`plainlang.go:19-34`) does a case-insensitive
substring match of every name against the `vistaEse` list. Used both by `v new`
(vet one domain name) and by the conformance test (vet the whole command tree).

### What `v` exposes today (verified live from `./dist/v --help`)

```
v pkg parse | decompose | assemble | roundtrip | canonicalize | lint
        | build | install | verify | uninstall
v new <domain>
v schema | version | install-completions
-o/--output text|json|auto · --no-color · -v/--verbose
```

`pkg parse/decompose/assemble/roundtrip/canonicalize/lint/build` are **offline**
(pure file work); `pkg install/verify/uninstall` reach a **live engine**
(ydb|iris) and do so **only** through the `m-driver-sdk` seam (`mdriver.Client`,
pinned `m-driver-sdk v0.3.0` as an indirect dep) — honoring the waterline's
transport monopoly.

> **⚠ Pin lag (finding).** The umbrella pins `v-pkg v0.1.0`, which advertises
> **10** pkg verbs. The current `v-pkg` HEAD has **13** — it adds `classify`,
> `snapshot`, and `restore` (the patch-reversibility lifecycle built since the
> tag). So `v` is one `go.mod` bump behind v-pkg's real surface. Re-pinning
> v-pkg and regenerating the registry would surface those three verbs. This is
> the static-pin discipline working as designed (adoption is a deliberate bump),
> but the lag is worth closing.

---

## 3. What it plans to do (vision & roadmap)

The governing spec is **`v-cli-platform.md`** — currently at
`m-stdlib/docs/plans/v-cli-platform.md` (DRAFT, last touched 2026-06-11). Per
its own "Home note" and the org CLAUDE.md, it is meant to **graduate into this
repo** once it exists. **That graduation has not happened** — see §7. The spec is
11 sections; the load-bearing ones:

- **§2 Naming** — the `v <plain-noun>` mapping table; VistA term in docs only.
- **§3 Composition** — own-repo/own-cadence domains, statically pinned into a
  thin `v` (CQ1). No `replace` directives in production.
- **§4 Per-domain contract** — `clikit` envelope + exit-code ladder + a
  `dist/v-contract.json` **generated from the kong defs**, drift-gated.
- **§5 Registry** — `v`'s surface is generated from the aggregate of pinned
  contracts; `v help`, completion, dispatch all read it; never hand-maintained.
- **§6 Template + `v new`** — both shipped (the standalone template repo was
  later dropped in favor of the inline skeleton).
- **§7 Conformance** — three gates: contract drift, envelope conformance,
  plain-language lint.
- **§9 First domain `v pkg`** — "the first proof of the whole `v` platform."

### Planned domains

| `v` domain | VistA subsystem | Status |
|---|---|---|
| **`v pkg`** | KIDS | **BUILT** (the only mounted domain; live on YDB + IRIS) |
| `v db` | FileMan | Planned — named next candidate; blocked on clikit extraction |
| `v config` | Parameter Tools / XPAR | Planned |
| `v rpc` | RPC Broker (#8994) | Planned |
| `v job` | TaskMan (`^%ZTLOAD`) | Planned |
| `v mail` | MailMan | Planned |
| `v io` | Device Handler (`^%ZIS`) | Planned |
| `v hl7` · `v fhir` | HL7 · FHIR (names already modern) | Planned |

All five platform design questions (CQ1–CQ5) were resolved 2026-06-11:
static-pinned composition · reuse the m-driver transport (don't grow a own) ·
thin `v` over pinned domain repos · both a generator and (originally) a template ·
wrapper-only (don't absorb m-cli internals).

### Roadmap / sequencing

1. **T0a.0 — platform foundation: DONE** (2026-06-12). The `v` umbrella mounting
   `v pkg`, `v new`, plain-language lint, contract + registry generators.
2. **M0a — `v pkg` live lifecycle: CLOSED.** install/verify/uninstall proven on
   both engines (YDB `vehu`, IRIS `foia`) over the driver; dual-engine exit gate
   green.
3. **Waterline G1 gate: BUILT** (`m arch check`), wired into org CI; all
   ecosystem repos declare a `layer` and gate clean.
4. **Next:** the **clikit extraction** to a shared module → then the **second
   domain** (`v db` / FileMan is the named candidate, or a Go-fronted
   `vista-info-hub` as `v vista`).

---

## 4. Gates, CI, and the registry discipline

v-cli embodies the org's one discipline — **`source-tag → generate → registry →
red-gate`** — on the tooling surface itself:

- **source-tag:** `repo.meta.json` declares `"layer": "v"` and
  `"exposes": {"registry": "dist/v-registry.json"}`.
- **generate:** `buildRegistry()` reflects the live mounted kong tree (via each
  domain's `Contract()`), so it can't drift from real verbs.
- **registry:** `dist/v-registry.json` is committed (`.gitignore` un-ignores
  `dist/*.json`).
- **red-gate:** three tests + Makefile targets enforce it:
  - `TestRegistry_Golden` (`registry_test.go:13`) — the §5 drift gate; regen with
    `make registry` (`UPDATE_GOLDEN=1`), check with `make check-registry`.
  - `TestRegistry_Invariants` — header sane, ≥1 domain, `pkg` present with verbs.
  - `TestPlainLanguage_CommandSurface` (`plainlang_test.go:14`) — walks the
    **entire** kong tree (umbrella + every mounted domain) and fails if any
    command/flag name contains vista-ese. This is the family's reason to exist,
    enforced mechanically.

**CI** (`.github/workflows/ci.yml`) calls two org-central reusable workflows:
`go-ci.yml` (lint/test/build/vuln) and `arch-waterline.yml` (the m/v waterline
G1–G5 gates). `repo.meta.json` lists `verification_commands`:
`go test ./...` and `./dist/m arch check .` — v-cli has **no arch command of its
own**; it borrows `m arch check` from m-cli.

**Makefile:** `make check` = vet + lint + test + build; plus `registry` /
`check-registry` for the drift gate. Builds are `-trimpath`.

---

## 5. The m-cli comparison (the pattern being mirrored)

`m-cli` is the busybox `v` replicates. Two differences matter:

1. **Composition.** m-cli is an **exec-dispatch** busybox — native `internal/`
   commands *plus* forwarding subtrees to separate sibling binaries (`irissync`,
   `kids-vc`) resolved via `M_<NAME>_BIN`/`PATH`. v-cli is a **static-link**
   busybox — it imports domains as Go modules and mounts them in-process (CQ1).
   This is why v-cli hits the clikit-extraction wall that m-cli never does.

2. **Breadth.** m-cli surfaces ~14 command namespaces (`fmt`, `lint`, `lsp`,
   `test`, `coverage`, `watch`, `vista status/exec`, `arch check`, `version`,
   `schema`, + dispatched `list/pull/push/verify`, `kids …`). v-cli surfaces
   **one** real domain. v-cli is at m-cli's "first command" stage.

Where v-cli **leads** m-cli: it commits a generated **registry file**
(`dist/v-registry.json`); m-cli only exposes a schema in-memory. And v-cli ships
a **scaffolder** (`v new`); m-cli has no built-in `m new` despite the org CLAUDE.md
referencing one.

---

## 6. The ecosystem (what could become `v` domains)

The "busybox content" for `v` — analogous to m-cli's siblings — comes from these
repos under `~/vista-cloud-dev/`:

### Directly relevant / consumed today

- **v-pkg** (layer `v`, Go, AGPL-3.0, 47 commits) — **the one mounted domain.**
  The Go port of the KIDS round-trip tool (lineage: Sam Habiel's MUMPS XPDK2VC →
  Python py-kids-vc → Go). A mature offline core (decompose/assemble/canonicalize
  byte-identical to py-kids-vc on the corpus) **plus** a feature-complete,
  dual-engine-proven live KIDS lifecycle (build/install/verify/uninstall +
  classify/snapshot/restore patch-reversibility). It vendors **`clikit`**
  (the shared CLI kit), **`pkgcli`** (the mounted verb surface), and
  **`vcontract`** (the `Manifest`/`Contract` registry type) — all three imported
  by v-cli. Open frontier: live-VistA package *extraction* (P4, design-only) and
  an IRIS re-validation of the chunked install path. *Housekeeping:* its README
  is stale (still "m-kids", Apache-2.0).

- **m-driver-sdk** (layer `m`, Go, pinned **v0.3.0**) — **the transport seam.**
  Exports `mdriver.Client` + a `Transport` interface; the *only* sanctioned path
  to a live engine. v-pkg's live verbs reach engines through it; v-cli pins it
  transitively. Not a domain — it's the waterline's single seam artifact.

### Strong future domain candidates

- **vista-info-hub** (layer `v`, Go, substantial — multi-face CLI/TUI/REST/MCP) —
  the strongest near-term second domain. Already Go, already has a command
  surface for routine/RPC/option/FileMan queries → a natural `v vista`.
- **v-web** (layer `v`, M `VWEB*`, 7 routines + tests + KIDS) — real, but
  M-routine-based; needs a Go control surface before it could mount as `v web`.

### Related M-layer "VistA Standard Library" (not a `v` CLI domain)

- **v-stdlib** (layer `v`, **pure M** `VSL*`, 73 commits, ~2 weeks old) — the
  **VistA Standard Library**, the `v`-layer counterpart to m-stdlib (`STD*`).
  13 implemented `VSL*` modules (101 public labels), each an **adapter** binding a
  portable m-stdlib seam to concrete VistA infra: `VSLCFG` (XPAR config), `VSLFS`
  (FileMan storage), `VSLIO` (Kernel TCP), `VSLLOG`, `VSLSEC` (Kernel identity),
  `VSLTASK` (TaskMan), plus the headline **RPC + HL7 → S3 traffic tap**
  (`VSLTAP`/`VSLRPCTAP`/`VSLRPCWRAP`/`VSLHL7TAP`/`VSLS3`/`VSLTAPFC`/`VSLTAPHL`).
  Dependency is strictly one-way (`VSL*` may call `STD*`, never the reverse).
  **Its relationship to v-cli:** v-stdlib ships as a **KIDS package**
  (`kids/vsl.build.json`, package `VSL` patch *4) built and installed **via the
  `v-pkg` binary** (`make kids` calls v-pkg; install = `v pkg install`, reversal
  = `v pkg uninstall`). So v-stdlib is a **consumer of `v pkg`**, not a CLI
  domain. *Housekeeping:* its README ("17 modules") and CLAUDE.md ("no modules
  yet, scaffold") are both stale — the real count is 13.

### Support / infra (not domains)

- **go-cli-template** — the scaffold + `clikit` source both families derive from;
  the reference for what `v new` emits.
- **m-parse** (m, Go) — tree-sitter-m parse substrate under m-cli.
- **m-dev-tools-mcp** (m, Go) — MCP server reflecting `m schema`; the model for an
  eventual v-MCP bridge.
- **vista-iris** (v) — a containerized VistA-on-IRIS *runtime target* that
  `v pkg install --engine iris` runs against.
- **doc-framework** — documentation-lifecycle scaffold; possible future `v docs`.
- **vpng** (v, M, 1 routine) — throwaway M1 proof-of-concept, not a domain.

---

## 7. Findings & recommended follow-ups

1. **Re-pin v-pkg.** The umbrella pins `v-pkg v0.1.0` (10 verbs) but v-pkg HEAD
   has 13 (`classify`, `snapshot`, `restore`). Tag a new v-pkg, bump the pin, run
   `make registry`. *(Closes the §2 pin lag.)*
2. **Graduate `v-cli-platform.md` into this repo.** It still lives in
   `m-stdlib/docs/plans/` and is DRAFT/stale on two points (it references the
   deleted `v-tool-template` repo and the "not-yet-existing" v-cli home). The repo
   now exists and has **no `docs/` of its own** (this report is the first doc).
   Move the spec here and refresh the two stale references.
3. **Extract `clikit` into a shared module** (e.g. `vista-cloud-dev/clikit`).
   This is the one hard prerequisite for a second domain — flagged in
   `main.go:11-15`. It is the highest-leverage next engineering step.
4. **Pick the second domain** once clikit is shared. `vista-info-hub` → `v vista`
   is the lowest-friction (already Go, already command-shaped).
5. **Add a `CLAUDE.md` to v-cli.** The repo inherits org rules but has no
   per-repo CLAUDE.md; one would pin the clikit-extraction constraint and the
   registry/plain-language gates locally.
6. **Stale-README sweep** across the family (v-pkg, v-stdlib) — low effort, high
   clarity payoff given how fast these repos moved.

---

## Appendix — file map (v-cli)

| File | Role |
|---|---|
| `main.go` | Umbrella kong grammar; mounts `v pkg` + meta commands; clikit constraint note |
| `registry.go` | `buildRegistry()` — aggregates pinned domains' contracts → `dist/v-registry.json` |
| `registry_test.go` | §5 golden drift gate + registry invariants |
| `new.go` | `v new <domain>` scaffolder (inline skeleton) |
| `plainlang.go` | `vistaEse` list + `plainLangViolations()` lint engine |
| `plainlang_test.go` | Whole-surface plain-language conformance gate |
| `repo.meta.json` | `layer: v`, exposes the registry, `verification_commands` |
| `go.mod` | Pins `v-pkg v0.1.0`, kong, kongplete; `m-driver-sdk v0.3.0` (indirect) |
| `Makefile` | `check` / `registry` / `check-registry` |
| `.github/workflows/ci.yml` | Reusable `go-ci.yml` + `arch-waterline.yml` |
| `dist/v-registry.json` | Committed generated command surface (10 pkg verbs) |
| `dist/v` | Built binary (8.1 MB) |
