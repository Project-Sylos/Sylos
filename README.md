# Sylos

Sylos is a self-hosted migration platform. One binary runs the API and serves the web UI.


## Quick start (developers)

```bash
git clone https://github.com/Project-Sylos/Sylos.git
cd Sylos

# Sibling repos required for local `replace` builds (see Sylos-Dev-Utils clone_repos.sh)
# Sylos-API, Sylos-UI, Migration-Engine, Sylos-FS, Spectra, go-path-linter
# Optional: ln -sfn ../Sylos-UI ui   # Vite convenience; make/CI do this too

make build
./bin/sylos
```

By default Sylos:

- Listens on port **8086**
- Serves the embedded web UI
- Stores data under the OS app-data directory (`org.sylos.app`)
- Requires an OS keyring for the install master key
- Opens a browser when healthy

First launch walks you through creating an admin account.

### CLI flags

| Flag | Description |
|------|-------------|
| `--no-browser` | Start the server without opening a browser |
| `--api-only` | Serve the API only (no embedded UI) |
| `--port PORT` | Override the HTTP port (default 8086) |

There is no `config.yaml`. Defaults are baked into the binary.

### Tests before compile

`make build` and `make release` run `make test-all` first: `go test ./... -count=1` in this repo and each present replace sibling (`Sylos-API`, `Migration-Engine`, `Sylos-FS`, `Spectra`, `go-path-linter`). Missing siblings are skipped. If any suite fails, the run finishes the rest, prints a short summary of which repos failed (full FAIL lines and traces stay above), and exits without compiling. Use `make test` for this repo only.

### Cloud OAuth apps (release builds)

Non-Google provider apps (Dropbox, Box, OneDrive, SharePoint) are baked at link time from committed files:

| File | Role |
|------|------|
| [`oauth.env`](./oauth.env) | Provider client IDs / secrets / optional Microsoft tenant IDs as `enc:…` values |
| [`oauth.key`](./oauth.key) | Build key used to reveal those values (`make oauth-genkey` once if missing) |

Workflow:

1. Put plaintext values in `oauth.env` (blank line skips that provider).
2. Run `make oauth-encrypt` so non-empty values become `enc:…`.
3. `make build` / `make release` reveal with `oauth.key` and pass them through `-ldflags`.

The key and ciphertext are both in the repo on purpose. That keeps secrets out of plain git diffs and casual scanners; it is not a secrecy boundary. A release binary still embeds the revealed values. Admins can override apps in Settings. Google Drive remains bring-your-own-app.

When ldflags did not bake credentials (for example `go run` from this tree), startup loads `oauth.env` + `oauth.key` from the working directory if present and fills empty builtin slots before seeding `sylos.api`.

The **Release** workflow runs `make release` the same way. If siblings are private, set Actions secret `SIBLING_CHECKOUT_TOKEN`.

## Building from source

Requires Go matching `go.mod`, Node.js/npm (to produce `internal/ui/dist`), and sibling repos for `replace` directives:

```
Sylos/                  ← this repo (unified binary)
Sylos-API/              ← API server library
Sylos-UI/               ← web UI (optional symlink at ui/ → ../Sylos-UI)
Migration-Engine/       ← migration SDK
Sylos-FS/               ← filesystem adapters
```

```bash
make build     # test-all → ui/dist → bin/sylos
make run       # build + run
make test-all  # sibling go test suites only
make release   # build + zip under dist/
make clean
```

Compile only needs a filled `internal/ui/dist` for `go:embed`. The `ui/` symlink is a local Vite convenience; a release builds assets from the Sylos-UI sibling.

## Developing individual components

**API only** (Sylos-API):

```bash
cd ../Sylos-API
go run .
```

**UI only** (Sylos-UI, with API elsewhere):

```bash
cd ui   # or ../Sylos-UI
npm install --legacy-peer-deps
npm run dev   # http://127.0.0.1:3000 → API at http://localhost:8086
```

## Under Construction Note

This repo is under active development. For early alpha testing and contributor coordination, use [GitHub Discussions](https://github.com/orgs/Project-Sylos/discussions).

## Related repositories

| Repo | Role |
|------|------|
| [Sylos-API](https://github.com/Project-Sylos/Sylos-API) | REST API and Migration Engine bridge |
| [Sylos-UI](https://github.com/Project-Sylos/Sylos-UI) | React web UI |
| [Migration-Engine](https://github.com/Project-Sylos/Migration-Engine) | Core migration SDK |
| [Sylos-FS](https://github.com/Project-Sylos/Sylos-FS) | Filesystem adapters |
| [Spectra](https://github.com/Project-Sylos/Spectra) | Synthetic filesystem simulator |
