# AGENTS.md — Sylos umbrella

Guidance for Cursor agents (and humans) working in the Sylos monorepo layout. Prefer **reading the linked docs** over inventing architecture; this file is a map, not a second source of truth.

## Overall Context

Sylos is a migration tool, similar to Rclone in that it's a free open source migration software. But differs from Rclone in very many ways. 
1. For one, it allows for preview-before-action semantics throughout the whole software. 
2. It has DB backed persistence including resumability 
3. It's built with a GUI in mind but also offers headless scriptability for developers and scripters, with a TUI (currently in progress) as well.

## Layout

Local development uses sibling repos under one parent (see [README.md](./README.md) and Sylos-Dev-Utils `clone_repos.sh`):

```
Sylos/                  ← this repo: unified binary (API + embedded UI)
Sylos-API/              ← REST/API library the binary embeds
Sylos-UI/               ← React/Vite UI (git submodule at ui/ here)
Migration-Engine/       ← migration SDK (queues, DuckDB, scaling)
Sylos-FS/               ← FS adapters (local, cloud, Spectra, SFTP, …)
Spectra/                ← synthetic FS for chaos / integration tests
```

From this repo, Go `replace` directives point at `../Sylos-API`, `../Migration-Engine`, `../Sylos-FS`, etc. When changing shared behavior, edit the owning sibling — do not fork copies into Sylos.

Other nearby trees (docs / site / tooling, not always on the critical path): `Sylos.wiki/`, `pages/`, `project-sylos.github.io/`, `go-path-linter/`.

---

## Repo purposes

| Repo | Role | Start here |
|------|------|------------|
| **Sylos** | Product launcher: builds UI, embeds it, runs API on one port (~8086). Config, OAuth creds dir, first-run admin. | [README.md](./README.md), [ROADMAP.md](./ROADMAP.md), [config.yaml.example](./config.yaml.example) |
| **Sylos-API** | Chi REST layer, auth/JWT, corebridge → Migration Engine, provider/OAuth routes, migration lifecycle. | [../Sylos-API/README.md](../Sylos-API/README.md), [../Sylos-API/internal/routes/migrations/README.md](../Sylos-API/internal/routes/migrations/README.md) |
| **Sylos-UI** | Browser UI (Vite/React); talks to API (dev often `:3000` → API `:8086`). | [../Sylos-UI/README.md](../Sylos-UI/README.md) (submodule: [ui/README.md](./ui/README.md)) |
| **Migration-Engine** | One-way migration engine: BFS traversal/copy/delete, queues, DuckDB seal, autoscaler. **Not** bidirectional sync. | [../Migration-Engine/README.md](../Migration-Engine/README.md) + **docs below** |
| **Sylos-FS** | `FSAdapter` implementations and cloud OAuth session plumbing used by API + ME. | [../Sylos-FS/README.md](../Sylos-FS/README.md), [../Sylos-FS/pkg/fs/README.md](../Sylos-FS/pkg/fs/README.md), [../Sylos-FS/docs/cloud_provider_checklist.md](../Sylos-FS/docs/cloud_provider_checklist.md) |
| **Spectra** | Fake filesystem + chaos rate limits for ME/FS integration tests. | [../Spectra/README.md](../Spectra/README.md), [../Spectra/sdk/README.md](../Spectra/sdk/README.md) |

---

## Migration Engine — algorithms & control plane (read these)

Canonical design docs live in **`../Migration-Engine/docs/`**:

| Doc | Contents |
|-----|----------|
| [algorithms.md](../Migration-Engine/docs/algorithms.md) | BFS vs DFS, SRC/DST coordination, copy two-pass, delete reverse-BFS, completion rules |
| [autoscaler.md](../Migration-Engine/docs/autoscaler.md) | Observer → AIMD loop, FS throttle / memory / underfeed, operation profiles, knobs |
| [item_statuses.md](../Migration-Engine/docs/item_statuses.md) | Traversal / copy / delete status semantics (event-sourced in DuckDB) |
| [fs_error_classification.md](../Migration-Engine/docs/fs_error_classification.md) | Retry vs throttle axes, ambiguous local/FUSE errors, degradation bridge |

Package-level detail:

| Doc | Contents |
|-----|----------|
| [pkg/queue/README.md](../Migration-Engine/pkg/queue/README.md) | Pull / lease / seal; subpackages `mode` / `worker` / `observe` / `gpl` (children import root; migration blank-imports hooks) |
| [pkg/db/README.md](../Migration-Engine/pkg/db/README.md) | Schema, SealBuffer, status events, indexes; same root/child pyramid as queue |
| [pkg/scaling/README.md](../Migration-Engine/pkg/scaling/README.md) | Thin root contract + `profile` / `aimd` / `memory` / `backend` / `loop` (loop owns Autoscaler) |
| [pkg/migration/README.md](../Migration-Engine/pkg/migration/README.md) | Entry points, manager lifecycle, phase wiring |
| [pkg/tests/README.md](../Migration-Engine/pkg/tests/README.md) | How integration / chaos tests are organized |

---

## Hard invariants agents must respect

1. **One-way migration** — ME discovers and copies SRC→DST; do not invent two-way sync semantics.
2. **DuckDB is source of truth** for frontier / status; queues pull from DB and seal results back (see queue + db READMEs).
3. **Byte streaming on OpenWrite** — Sylos-FS adapters must stream SRC→DST in small in-memory chunks. **No spill-to-temp / full-file staging** before upload. Dropbox/GDrive-style session start/append/finish is fine; Graph/Box must PUT fragments during `Write`, not buffer the whole object until `Close`. See Sylos-FS cloud checklist + ME algorithms copy section.
4. **Autoscaler owns throughput knobs** — API/UI observe and pause/stop; they do not set worker counts. Profiles live in ME `pkg/scaling/profile`; the control loop is `pkg/scaling/loop`.
5. **Throttle ≠ stall** — Queue watchdog “STALL DETECTED” means no progress heartbeats with work leased. Active `RateLimitedUntil` / seal I/O wait **suppress** that dump. Fix missing beats (e.g. close-only uploads) rather than treating every stall banner as AIMD failure.
7. **Minimal diffs** — match existing style; no drive-by refactors; don’t commit unless asked.

---

## Typical change touchpoints

| Kind of change | Likely repos |
|----------------|--------------|
| New cloud provider | Sylos-FS (+ checklist) → Sylos-API providers/creds → Sylos-UI OAuth/browse → ME operation profiles |
| Traversal / copy / delete behavior | Migration-Engine `pkg/queue`, `pkg/migration`, docs/algorithms.md |
| Scaling / rate-limit behavior | Migration-Engine `pkg/scaling`, docs/autoscaler.md + Sylos-FS degradation/classify |
| Schema / seal / resume | Migration-Engine `pkg/db` |
| Product shell / config / embed UI | Sylos |
| REST contract / auth / migration jobs | Sylos-API |
| Screens / wizards | Sylos-UI |

## Rules for working on code
1. No pointless wrapper functions. If you have a struct that is private and then make a public wrapper, just make the original struct public.
2. Prefer using reusable functions over making a new function for every new thing.
3. Prefer using a reusable function instead of making 2 similar functions with one or two changed things.
4. No em-dashes anywhere, ever. Not even in code comments.

### Go: avoid structural slop (Go-Slop-Finder)

Agents writing or editing Go must follow the no-slop style in [`../../Go-Slop-Finder/docs/STYLE_GUIDE.md`](../../Go-Slop-Finder/docs/STYLE_GUIDE.md) (short form: [`CURSOR_INSTRUCTIONS.md`](../../Go-Slop-Finder/docs/CURSOR_INSTRUCTIONS.md)). Run `slopfinder .` from a repo root (not `./...`) when checking. Never add `//nolint`, slopfinder ignore comments, or similar suppressions.

- No wrapper-alias / thin-wrapper / trivial single-call forwarders; the exported API owns the logic.
- No constant-forwarding aliases: if two functions only differ by keys/tables/side constants passed into one helper, export the helper and call it with those constants at the call site. Do not keep `GetCopyFoo` / `GetDeleteFoo` one-liners as nice API names.
- No exported→unexported forwarder pairs; no assign-then-return / redundant temps.
- No getter/setter that only touches one field unless required by an interface or mutex.
- Prefer `return cond` over true/false branches; no `== true` / `== false`; no nil-check before `range`; no switch-on-bool.
- Early returns over deep nesting; no empty else; wrap errors with context when the function exists to add meaning.
- No duplicate function bodies (exact or same control-flow template); parameterize SRC/DST or copy/delete twins instead of cloning.
- Do not chase repeated composite literals when they are intentional parallel construction (e.g. Source/Destination) or one-off SQL/DDL fragments.
- Interface-required one-liners (e.g. FSAdapter `NormalizePath`) are fine; inventing an extra private alias on top is not.
- When fixing slop: no new clarifying comments; do not change print/log strings.

---

## Build / run (launcher)

```bash
# from Sylos/
cp config.yaml.example config.yaml   # once
make build && ./bin/sylos
```

UI submodule: `ui/` (or sibling Sylos-UI). API-only / no-browser flags are in [README.md](./README.md).

When iterating on ME or FS locally, rebuild the Sylos binary (or run tests in the sibling) so `replace`d modules are picked up.
