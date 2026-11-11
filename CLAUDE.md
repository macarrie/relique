# CLAUDE.md

Relique is an rsync-based backup tool written in Go. A single binary (`relique`) provides a Cobra CLI and an HTTPS server (Gin) that exposes a JSON API and serves an embedded React web UI.

## Commands

| Task | Command |
|---|---|
| Full build (webui + binary) | `make build` → `output/relique` |
| Go binary only | `make bin` (runs `go mod tidy`, injects version via `-ldflags -X .../api.ReliqueVersion`) |
| Web UI build + copy into server embed dir | `make webui` (`npm ci && npm run build`, copies `webui/dist` → `internal/server/dist`) |
| Unit tests (local) | `go test -cover ./...` and `go vet ./...` |
| Unit tests (as CI does) | `make test` (Docker build of `test/Dockerfile_tests`) |
| Single test | `go test ./internal/repo -run TestLoadFromFile` |
| Web UI dev server | `cd webui && npm run dev` (Vite proxies `/api` → `https://localhost:8433`) |
| Web UI lint | `cd webui && npm run lint` |
| Docker image | `make docker` |
| Reset local state | `make reset` (wipes `~/.config/relique` db/storage/catalog) |
| Release | `make release TAG=x.y.z` (writes `.current_version`, commits, tags `vX.Y.Z`) |

**Gotcha:** `internal/server/routes.go` uses `//go:embed dist/*`. Go builds, `go vet` and `go test ./...` fail if `internal/server/dist` doesn't exist, so run `make webui` once (or create a placeholder file in that directory) before working on Go code. `internal/server/dist` is gitignored.

## Stack

- **Backend:** Go (go.mod 1.22, Docker builds with golang:1.24)
  - CLI: `spf13/cobra`; config: `spf13/viper` + `pelletier/go-toml` (TOML files)
  - HTTP: `gin-gonic/gin`, TLS only (cert/key paths from config, default port 8433)
  - DB: SQLite (`mattn/go-sqlite3`, CGO) with `golang-migrate` (embedded SQL migrations) and `Masterminds/squirrel` query builder
  - Logging: `log/slog` with `lmittmann/tint` handler
  - Helpers: `samber/lo`, `hashicorp/go-multierror`, `InVisionApp/tabular` (CLI tables), `charmbracelet/glamour`, `promptui`, `go-git` (module install)
  - Runtime deps: `rsync` and `ssh` on the host (backups run rsync over SSH to clients)
- **Frontend (`webui/`):** React 19 + TypeScript 5.7, Vite 6, Tailwind CSS 4 + daisyUI 5, react-router-dom 7, axios, @tanstack/react-table, react-toastify, remixicon, ESLint 9 (typescript-eslint, react-hooks)
- **CI:** GitHub Actions. `tests.yml` runs `make test` on push/PR to `master`; `docker.yaml` builds and pushes multi-arch images (`macarrie/relique`) on `v*.*.*` tags.

## Directory structure

```
cmd/relique/            main.go entry point → cli.Execute()
  cli/                  one file per Cobra command group (client, repo, module, job, image, backup, restore, config, server, version); root.go has global --debug / --config flags
api/                    Public Go API layer. Business operations used by BOTH the CLI and the HTTP server (ClientList, JobGet, Backup, Restore…). Holds ReliqueVersion.
internal/
  api_helpers/          Pagination (PaginationParams, PaginatedResponse[T]) and search/filter structs shared by api, CLI and server
  config/               Global Configuration (viper/TOML), load/init, default paths, WebUI HTTP config
  client/ module/ repo/ Domain types loaded from TOML files; repo uses a Repository interface with per-type impls (local.go, generic.go)
  module/install.go     Module install from URL/path (git or tarball)
  job/ image/           Domain types with SQLite persistence (db.go) and logic (job.go/image.go)
  job_status/ job_type/ backup_type/   Small enum-like types (iota consts + String()/FromString())
  rsync_task/ (+lib/)   Builds and runs rsync commands for full/diff backup and restore; stats parsing
  db/                   SQLite connection pool + embedded migrations (migrations/NNNNNN_name.{up,down}.sql)
  server/               Gin routes (routes.go), handlers per resource (clients.go, jobs.go…), TLS server, embedded web UI (dist/)
  logger/ consts/ utils/ Logging init, constants (DEFAULT_PAGE_SIZE), misc helpers
webui/src/
  pages/                One component per route (clients.tsx, client_details.tsx…)
  components/           Reusable UI pieces (cards, lists, tables, badges…)
  layout/               App shell (main, sidebar)
  types/                TS types mirroring Go JSON (snake_case fields) + const.tsx
  utils/api.tsx         Static API class wrapping axios, base URL /api/v1/
  routes.tsx            Route table (path, name, elt)
test/                   TOML fixtures for unit tests (clients/, modules/, repo/) + Dockerfile_tests
build/package/          Production Dockerfile (Rocky Linux 9 minimal runtime, runs `relique server start`)
```

Config search paths: `$HOME/.config/relique/`, `/usr/local/etc/relique/`, `/etc/relique/` (or `--config`). Env prefix `RELIQUE_`. Clients, repositories and modules are each stored as individual TOML files in sub-folders.

## Architecture / layering

`cmd/relique/cli` and `internal/server` are thin front-ends → both call functions in `api/` → which use `internal/*` domain packages. Put new business logic in `api/` (or the relevant internal package), never directly in a CLI command or a Gin handler.

Adding an HTTP endpoint: add the route in `internal/server/routes.go` under the `/api/v1` group, a `webAPI<Verb><Resource>` handler in the matching `internal/server/<resource>.go`, then the matching method in `webui/src/utils/api.tsx`.

## Go conventions

- Errors: wrap with context using `fmt.Errorf("cannot <do thing>: %w", err)` — lowercase, starts with "cannot". Return errors up; log + exit only at the CLI/handler edge.
- Logging: structured `slog` with attributes, e.g. `slog.With(slog.Any("error", err), slog.String("name", name)).Error("Cannot find client in config")`. Domain types expose `GetLog()` returning a pre-scoped logger. Log messages are capitalized.
- CLI commands are declared inside `init()` functions of each `cli/*.go` file; flag values are package-level vars prefixed by command (`clientListPageSize`, `clientCreateSSHUser`). Groups validate config in `PersistentPreRun`. Output tables use `tabular`.
- Handlers: read params with `c.Param`/`getPagination(c)`, respond with `c.JSON(http.StatusOK, …)`, and on error log then `c.AbortWithStatus(...)`. List endpoints return `api_helpers.PaginatedResponse[T]`.
- DB access: squirrel builders (`sq.Insert/Select…`) → `ToSql()` → `tx.Exec`/`db.Handler()`; writes use transactions with deferred rollback on error. Schema changes go in a new numbered migration pair in `internal/db/migrations/`.
- Enum types: `iota` constants in a dedicated package with a wrapper struct (`JobStatus{Status uint8}`), `String()` and `FromString()`.
- Struct tags: `json` + `toml` (and `mapstructure` for viper) in snake_case.
- Package-level globals for defaults use UPPER_SNAKE_CASE (`CLIENTS_DEFAULT_FOLDER`, `MODULES_INSTALL_PATH`).
- `samber/lo` is used for filter/map/slice operations.
- Tests: standard `testing`, table-driven (`tests := []struct{name, args, want, wantErr}` + `t.Run`), `reflect.DeepEqual` comparisons, fixtures referenced via relative paths like `../../test/repo/example.toml`.

## Frontend conventions

- Function components with hooks (`useState`, `useEffect`, `useCallback`); one default-exported component per file; snake_case file names (`client_details.tsx`), PascalCase component names.
- All files use `.tsx`, including types and utils. Types are `type X = {...}; export default X;` and use the API's snake_case field names.
- Avoid using the "any" type in typescript
- API calls go through `API.<resource>.<method>()` and use `.then/.catch`; errors are surfaced with `Utils.notify(Const.CRITICAL, ...)` (react-toastify).
- Styling: follow `webui/DESIGN.md`. Use theme tokens (defined in `src/index.css`) and the shared primitives (`PageHeader`, `Section`, `PropertyList`, `Status`, `Stat`…), never raw Tailwind palette colors, cards or name badges. Icons use Remix Icon line classes (`ri-*-line`).
- 4-space indentation, double quotes.

## Git conventions

- Conventional-commit style subjects with scope: `feat(cli): …`, `fix(docker): …`, `feat(webui): …`.
- Main branch is `master`; feature work happens on branches and merges through PRs. Do not commit directly in the `master` branch.
- Release commits are `Release vX.Y.Z` made with `make release`.
