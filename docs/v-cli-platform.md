---
title: The `v` CLI — A Contract- and Registry-Driven Platform for VistA Developer Tools
status: adopted
version: v1.0
created: 2026-06-11
last_modified: 2026-06-25
revisions: 4
doc_type: [PLAN, SPEC]
relates_to: https://github.com/vista-cloud-dev/docs/blob/main/vsl-msl/msl-vsl-coordination-implementation-plan.md
---

# The `v` CLI — VistA Developer Tools Platform — **SPEC v1.0 (adopted)**

> **Status:** ADOPTED & BUILT. Specifies the naming scheme, command-surface
> **contract**, the **registry**, the composition model, and domain scaffolding for
> the family of VistA-specific Go developer tools fronted by a **single `v` CLI**.
> The platform foundation (T0a.0) is built: the `v` umbrella mounts its first
> domain, **`v pkg`** (the KIDS lifecycle), statically pinned; the registry,
> plain-language gate, and `v new` scaffolder are live and gate-green. Sibling to
> the [MSL⟷VSL coordination plan](https://github.com/vista-cloud-dev/docs/blob/main/vsl-msl/msl-vsl-coordination-implementation-plan.md),
> which consumed `v pkg` as its M0a.
>
> **Home note.** This is **org-tooling infrastructure**, not MSL/VSL-specific. It
> **has graduated into the `v` CLI repo** (its canonical home is now
> `v-cli/docs/v-cli-platform.md`); the former planning location
> `m-stdlib/docs/plans/v-cli-platform.md` is a redirect tombstone. The short
> governing form of the `m-*`/`v-*` scheme + the registry-driven discipline lives
> in the org-level [`CLAUDE.md`](../../CLAUDE.md) § *Naming & registry
> conventions*; **this doc is the canonical full spec** it points to (the naming
> scheme §2, the contract §4, the registry §5, scaffolding §6). For the current
> build state see the [status report](./v-cli-status-report-2026-06-25.md).
>
> **One-line summary:** A single `v` CLI wraps each insider VistA subsystem (KIDS,
> FileMan, XPAR, the RPC Broker, TaskMan, …) in a **plain-language** Go command
> (`v pkg`, `v db`, `v config`, …). Every domain lives in its **own repo with its
> own lifecycle**, is built from **one shared template**, and exposes a **versioned
> command-surface contract** that feeds a **generated, drift-gated registry** — the
> *same* `source-tag → generate → drift-gate` discipline the seam/ICR/citation
> registries use, now applied to the tooling surface.

---

## 1. Purpose & the two tool families

VistA exposes its power through insider subsystems with insider names. A developer
who wants to *package* an app must learn "KIDS"; to read a *parameter*, "XPAR"; to
touch the *database*, "FileMan." The `v` CLI's job is to **wrap vista-ese in modern
plain language** so that knowledge isn't required to be productive.

The org now has **two tool prefixes, split by *scope* (not language — both are
Go):**

| Prefix | Means | Assumes VistA? | Examples |
|---|---|---|---|
| **`m-*`** | engine-neutral M toolchain & libs | **no** — targets a bare M engine | `m-cli`, `m-stdlib`, `m-ydb`, `m-iris`, `m-driver-sdk`, `m-parse` |
| **`v-*`** | VistA-specific repos & tools | **yes** — needs Kernel/FileMan/KIDS/… | `v-stdlib` (the VistA Standard Library M package, `VSL*`) · the `v` CLI domains `v pkg`, `v db`, `v config`, … |

(Each family spans **both** M and Go — `m-stdlib` is M, `m-cli` is Go; `v-stdlib`
is M, `v-pkg` is Go. The prefix is **scope**, not language.)

This is the **same line the architecture doc draws for M code** (`STD*` portable
vs `V*` VistA-coupled), now drawn for tooling. `v` means **VistA** at both layers:

| | Engine-neutral / portable | VistA-specific |
|---|---|---|
| **M code (runs in the engine)** | `STD*` (m-stdlib) | `V*` (VSL, VWEB) |
| **Go tools (run on the host)** | `m-*` (m-cli, m-ydb, …) | **`v-*` / the `v` CLI** |

> **Consequence — `m-kids` → `v pkg`.** KIDS is a Kernel subsystem; the tool is
> useless against a bare engine. It is therefore VistA-specific and belongs in the
> `v` family. `m-kids` (pure Go, formerly `kids-vc`) is **renamed/refiled as the
> `v pkg` domain** (repo `v-pkg`). The split that matters is scope, not language —
> the `v-*` tools are *also* Go.

---

## 2. The naming scheme — `v <plain-noun>`, never the VistA product name

The whole value is a surface a developer can **guess without knowing VistA**:

| VistA subsystem (vista-ese) | `v` domain | Reads as |
|---|---|---|
| KIDS (Kernel Installation & Distribution System) | **`v pkg`** | package / install |
| FileMan (the database) | **`v db`** | database |
| Parameter Tools / XPAR | **`v config`** | configuration |
| RPC Broker / REMOTE PROCEDURE #8994 | **`v rpc`** | remote calls |
| TaskMan / `^%ZTLOAD` | **`v job`** | background jobs |
| MailMan | **`v mail`** | mail / alerts |
| Device Handler / `^%ZIS` | **`v io`** | devices / sockets |
| HL7 · FHIR | `v hl7` · `v fhir` | (already modern — keep) |

**Rule:** the domain and every verb/flag uses the **modern generic noun**, never
the VA product name — `v db` not `v fileman`; `v pkg` not `v kids`; `v config` not
`v xpar`; `install`/`uninstall` not "load distribution / back-out." The VistA term
stays in the **docs** (precision); it never appears in the **command**. §7 makes
this a mechanical gate.

> **Naming freedom.** The VA **DBA namespace registry** governs **M routine/global
> names inside VistA** (`VSL*`, `^VSL(`) — it does **not** govern host-side Go
> binary or subcommand names. So `v pkg`/`v db`/`v config` are unconstrained by VA
> governance; choose them purely for developer-friendliness.

---

## 3. One `v` CLI, many domain repos

`v` is a **single umbrella CLI** (exactly as `m` is for the M toolchain):
`v <domain> <verb> [args] [--flags]`. Each domain (`pkg`, `db`, …) is developed in
its **own repo with its own release cadence**, scaffolded by `v new` (§6) on the
shared `clikit` conventions, and **composed into `v`** via the registry (§5).

**Composition model — static-pinned (DECIDED, [CQ1](#10-resolved-decisions)), mirroring the SDK pattern.**
`v` imports each domain as a Go module and **pins its version in `go.mod`**; a
domain ships releases independently in its repo, and `v` **bumps the pin** to adopt
one — exactly the org's *"serialize the contract, parallelize the tools"* rhythm
from `m-driver-sdk` (no `replace` directives; the pin is the coordination point).
Different lifecycles are preserved at the **development** level; integration is a
deliberate pin-bump. *(Alternative: runtime plugin-dispatch — `v` discovers `v-*`
binaries on `PATH` and dispatches, reading each one's contract. Fully decoupled but
more moving parts. See [CQ1](#10-resolved-decisions).)*

---

## 4. The command-surface contract (per domain)

Built on **`clikit`** — the shared Go CLI conventions every toolchain binary already
honors: `--output text|json|auto`, a versioned JSON envelope, deterministic error
objects, the **exit-code ladder** (`0` ok · `1` runtime · `2` usage · `3`
check/drift · `4` refused), plus `schema` and `version`. The `v` contract is an
*extension* of clikit, not a reinvention.

Each domain emits a **contract manifest** `dist/v-contract.json`, **generated from
its Go command definitions** (kong), carrying:

- `domain`, tool **SemVer**, and a **`contract_version`** (bumps only on an
  incompatible command-surface change — independent of SemVer, exactly like the
  seam `contract_version` in the coordination plan §6);
- every **command**: name, summary, args (name/type/required), flags, the output
  **schema** ref, and the exit codes it can return.

A **drift gate** asserts the manifest matches the actual command tree (the same
`make check-manifest` discipline that gates `dist/`). The contract is a *file*, not
a convention — so a domain's surface cannot silently drift from what it declares.

---

## 5. The registry (the unified, generated surface)

`v`'s whole command surface is **generated from the aggregate of the pinned
domains' contract manifests** into a **registry** (`dist/v-registry.json`). `v help`,
shell completion (`kongplete`, already a dep), and dispatch all read the registry —
`v` **never hand-maintains its command list**. The registry is drift-gated against
the pinned domains' contracts, so "`v` advertises a command a domain no longer
provides" is a **red gate**, not a runtime surprise.

**Version axes (parallel to the coordination plan §6):**

| Axis | Lives on | Bumps when | Consumed by |
|---|---|---|---|
| Domain **SemVer** | the domain repo's tags | any release | `v`'s go.mod pin |
| **`contract_version`** | the domain's `dist/v-contract.json` | command surface changes incompatibly | the registry drift gate |
| **Registry pin set** | `v`'s `go.mod` + `dist/v-registry.json` | `v` adopts a new domain version | `v` users |

---

## 6. Domain scaffolding — `v new <domain>`

**`v new <domain>`** scaffolds a new domain repo with a minimal,
convention-conforming skeleton: a Go module `github.com/vista-cloud-dev/v-<domain>`,
an importable `<domain>cli` command package the umbrella mounts, and a Makefile
with the standard gates. It enforces the §2/§7 plain-language rule on the domain
name (`v new fileman` is refused). This is the tooling parallel of
`~/claude/templates/python` and `m new`, and is what makes "a Go library of VistA
utility functions" a *standardized* library rather than a pile of bespoke CLIs —
every domain is born with the same shape and quality gates.

> **Implementation note (supersedes the original CQ4 "template repo" decision).**
> The original plan called for *both* a standalone `v-tool-template` repo and a
> `v new` generator. The standalone template repo was created and then **deleted**
> (commit `94b96d2`); the skeleton is now built **into `v new`** (see
> `v-cli/new.go`), so there is one source of truth for the scaffold and no
> separate repo to keep in sync. A domain still grows from the skeleton toward the
> full contract/registry/conformance shape (§4, §5, §7), which today is provided
> concretely by the first domain, `v-pkg` (its `clikit` / `vcontract` packages).

---

## 7. Conformance + the plain-language gate

Every domain ships (from the template) a shared **conformance suite**:

1. **Contract drift** — the `dist/v-contract.json` matches the actual command tree.
2. **Envelope conformance** — output validates against the clikit schema; the
   exit-code ladder is honored.
3. **Plain-language lint** (on-brand, the family's reason to exist) — **no domain,
   command, or flag name may contain vista-ese**: `fileman`, `kids`, `xpar`,
   `mailman`, `taskman`, `duz`, `^%zis`, `dic`, `die`, … A leak = **red**. This
   mechanically enforces §2 as the family grows, so no tool ever re-exposes the
   insider terms the platform exists to hide.

---

## 8. Relationship to the rest of the stack

- **`clikit`** is the foundation — the `v` contract/registry extends it; don't
  reinvent the envelope, output modes, or exit codes.
- **The m-cli `VistaEngine` transport** — `v-*` tools that touch a live VistA should
  drive it through the **same transport abstraction** (`DockerEngine` / `SSHEngine`
  / a live VistA) the m-cli runner already owns, so every tool reaches the engine
  one uniform way. `v pkg`'s install/uninstall lifecycle drives Kernel's KIDS
  routines over this transport.
- **`v-*` Go tools vs `V*` M packages — different layers, don't conflate:**

| | `v db` (Go, host) | `VSLFS` (`V*`, M, in VistA) |
|---|---|---|
| Runs | on the developer's host | inside the VistA engine |
| Is | a developer CLI that *talks to* FileMan from outside | the seam adapter that *binds* `STDFS` to FileMan |
| Audience | a developer at a terminal | M code calling the seam |
| Lifecycle | a `v` domain release | a KIDS-installed routine |

  They are complementary — the `v-*` tools are often *used to develop and test* the
  `V*` packages — but they are not the same thing.

---

## 9. First domain: `v pkg` (the KIDS lifecycle)

Repo **`v-pkg`** (renamed from `m-kids`). It already ships the **offline** half as
pure Go (`decompose` / `assemble` / `roundtrip` / `canonicalize` / `parse` /
`lint`, byte-identical port of XPDK2VC); the platform adds the **live lifecycle**
(`build` / `install` / `verify` / `uninstall` / `status`) driving Kernel's existing
KIDS routines over the m-cli transport — **no new MUMPS package**, just Go
orchestration of `^XPDI…`. This is **M0a** of the coordination plan and the **first
proof of the whole `v` platform**: the determinism ledger there (§12.1) becomes
literal `v pkg` invocations — `v pkg build && v pkg install && … && v pkg uninstall
&& v pkg verify` — on both engines.

---

## 10. Resolved Decisions

**All resolved 2026-06-11 — none open before implementation.**

| # | Question | Decision |
|---|---|---|
| CQ1 | **Composition** — static-pinned modules, or runtime plugin-dispatch of `v-*` binaries? | **DECIDED 2026-06-11: static-pinned.** `v` imports each domain as a Go module and pins it in `go.mod`; one binary, compile-time contract safety, the registry generated at build — exactly the `m-driver-sdk` "serialize the contract, parallelize the tools" pattern the org already runs. Different lifecycles preserved at the *development* level; integration is a deliberate pin-bump. **Escape hatch:** switch to plugin-dispatch (separate `v-*` binaries on `PATH`) only if third-party plugins or release cadences too fast for a `v` rebuild ever become a goal — neither applies at current scale. |
| CQ2 | **Transport ownership** — does `v` reuse the m-cli `VistaEngine` transport, or grow its own? | **DECIDED: reuse m-cli's** `VistaEngine`/`DockerEngine`/`SSHEngine` transport — one path to the engine for the whole toolchain; `v` never reinvents connectivity. |
| CQ3 | **Repo shape** — `v-<domain>` repos pinned into a thin `v`, or a `v` monorepo with domain packages? | **DECIDED: `v-<domain>` repos** (own lifecycles) **pinned into a thin `v`** — consistent with CQ1 static-pinned composition. |
| CQ4 | **Template home** — a `v-tool-template` repo, or a `v new` generator inside the `v` repo? | **DECIDED 2026-06-11: both.** **Revised 2026-06-25: `v new` only.** The standalone `v-tool-template` repo was created then deleted (commit `94b96d2`); the skeleton is now built into `v new` (`v-cli/new.go`) — one source of truth, no separate repo to sync. See §6. |
| CQ5 | **Does `v` subsume the M-toolchain's VistA-ish bits** (e.g. `SSHEngine`→vista-meta) or stay strictly the wrapper layer? | **DECIDED: stay the wrapper** — reuse, don't absorb, m-cli internals; the transport stays in m-cli (CQ2). |

---

## 11. References

- [`msl-vsl-coordination-implementation-plan.md`](https://github.com/vista-cloud-dev/docs/blob/main/vsl-msl/msl-vsl-coordination-implementation-plan.md)
  — consumes `v pkg` as its M0a; the KIDS lifecycle + version-controlled build
  spec live in §7.1–§7.2 there.
- [`msl-vsl-architecture.md`](https://github.com/vista-cloud-dev/docs/blob/main/vsl-msl/msl-vsl-architecture.md) — the `STD*`/`V*` line this
  scheme extends to tooling.
- `clikit` (shared Go CLI conventions, in the m-cli / v-pkg toolchain) — the
  foundation the `v` contract extends. (Originates from `go-cli-template`; the
  extraction into one shared importable module — the prerequisite for mounting a
  *second* domain — is tracked in the [status report](./v-cli-status-report-2026-06-25.md).)
- [`../../CLAUDE.md`](../../CLAUDE.md) — vista-cloud-dev org rules; the
  `m-driver-sdk` *serialize-the-contract* model this platform mirrors.

---

*End SPEC v1.0. **All §10 CQs resolved (2026-06-11; CQ4 revised 2026-06-25)** —
naming (§2), contract/registry/scaffold discipline (§4–§7), composition =
static-pinned (CQ1), transport-reuse (CQ2), `v-<domain>` repos (CQ3), `v new`
generator (CQ4), wrapper-only (CQ5). The first concrete build, `v pkg` (M0a of the
coordination plan), is done: `m-kids` refiled as `v-pkg`, the umbrella mounts it
statically, and the `build`/`install`/`verify`/`uninstall` (+ classify/snapshot/
restore) verbs are live on both engines.*
