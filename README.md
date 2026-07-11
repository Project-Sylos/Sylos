# Sylos

Sylos is a self-hosted migration platform. One binary that runs the API and serves the web UI, similar to Jellyfin or Navidrome.

## Quick start

```bash
# Clone with the UI submodule
git clone --recurse-submodules https://codeberg.org/Sylos/Sylos.git
cd Sylos

# Or, if already cloned:
git submodule update --init --recursive

# Configure (copy and edit paths/services as needed)
cp config.yaml.example config.yaml

# Build and run
make build
./bin/sylos
```

By default, `sylos` starts the API on port **8086**, serves the embedded web UI on the same port, and opens your browser once the server is healthy.

On first launch, the UI walks you through **creating an admin account** (Jellyfin/Navidrome-style). After that, all API routes require login. Admins can manage additional users under Settings → Users.

### CLI flags

| Flag | Description |
|------|-------------|
| `--no-browser` | Start the server and serve the UI, but do not open a browser |
| `--api-only` | Serve the API only (no embedded UI; for TUI or external clients) |
| `--config PATH` | Config file path (default: `config.yaml`) |
| `--port PORT` | Override the HTTP port from config |

Examples:

```bash
./bin/sylos --no-browser
./bin/sylos --api-only
./bin/sylos --config /etc/sylos/config.yaml --port 9090
```

### Configuration

Copy [`config.yaml.example`](config.yaml.example) to `config.yaml` and adjust:

- `http.port` listening port (default: 8086)
- `runtime.data_dir` persistent data directory
- `services.local` allowlisted local filesystem roots
- `jwt.access_token_ttl` session length (default: 24h)
- `auth.bcrypt_cost` password hashing cost (default: 12)
- `runtime.oauth_creds_dir` cloud OAuth JSON directory (default: `./creds`)

Environment variables use the `SYLOS_` prefix (e.g. `SYLOS_HTTP_PORT=9090`, `SYLOS_CONFIG_PATH=/path/to/config.yaml`).

Cloud OAuth credentials live in **`./creds/`** next to `config.yaml` (see [`creds/README.txt`](creds/README.txt)). The API reads them at runtime for server-side OAuth token exchange. Secrets never reach the browser. Restart the server after changing credential files. The server also falls back to `ui/creds/` if present.

## Building from source

Requires Go 1.25+, Node.js/npm (for the UI submodule), and sibling repos for local development:

```
Sylos/                  ← this repo (unified binary)
Sylos-API/              ← API server library
Sylos-UI/               ← web UI (git submodule at ui/)
Migration-Engine/       ← migration SDK
Sylos-FS/               ← filesystem adapters
```

The `go.mod` in this repo uses `replace` directives pointing at `../Sylos-API`, `../Migration-Engine`, etc. Match the layout from [Sylos-Dev-Utils](https://codeberg.org/Sylos/Sylos-Dev-Utils) `clone_repos.sh` for local dev.

```bash
make build    # builds ui/dist, embeds it, produces bin/sylos
make run      # build + run with default flags
make clean    # remove bin/ and copied UI dist
```

## Developing individual components

When working on a single repo, you can still run pieces separately:

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

This repo is under active development. For early alpha testing and contributor coordination, join the official Discord at https://chat.sylos.io/.

## Contributing

See pinned messages in the #bounties channel on Discord, or the [Active ongoing issues and roadmap](https://codeberg.org/Sylos/Discussions-and-Issues/issues/5) page for dev and creative work needed.

## Related repositories

| Repo | Role |
|------|------|
| [Sylos-API](https://codeberg.org/Sylos/Sylos-API) | REST API and Migration Engine bridge |
| [Sylos-UI](https://codeberg.org/Sylos/Sylos-UI) | React web UI (submodule) |
| [Migration-Engine](https://codeberg.org/Sylos/Migration-Engine) | Core migration SDK |
| [Sylos-FS](https://codeberg.org/Sylos/Sylos-FS) | Filesystem adapters |
| [Spectra](https://codeberg.org/Sylos/Spectra) | Synthetic filesystem simulator |
