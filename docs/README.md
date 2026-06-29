# v-cli docs

Documentation for **`v-cli`** — the `v` VistA developer-tools umbrella CLI.

## Key documents

- **[`v-cli-platform.md`](v-cli-platform.md)** — the canonical, adopted platform
  **SPEC** (naming scheme, command-surface contract, registry, composition model,
  domain scaffolding). This is the full spec the org-level `CLAUDE.md` points to.

## Folders

- **`proposals/`** — live, decision-seeking proposals.
  - [`v-domain-template.md`](proposals/v-domain-template.md) — *draft* — replace
    `v new`'s embedded skeleton with a thin, drift-gated `v-domain-template` repo
    that `v new` renders from (re-opens CQ4; fixes the scaffold-drift that left new
    domains without `.envrc`/`repo.meta.json`/`.golangci.yml`).
- **`memory/`** — auto-memory (durable gotchas/recipes only); see
  [`memory/MEMORY.md`](memory/MEMORY.md) for the index.
- **`archive/`** — retired docs kept for history (e.g. the dated
  [2026-06-25 status report](archive/v-cli-status-report-2026-06-25.md)).

Additional standard folders (`guides/`, `design/`, `modules/`) are added when there
is content for them.
