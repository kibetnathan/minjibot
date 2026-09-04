# Contributing to MinjiBot

Thanks for wanting to help! This guide covers the conventions and workflow for
contributing to the codebase. For architecture and operational details, see
`README.md` and `AGENTS.md`.

## Repository layout

```
cmd/                     Entrypoints (bot, api, unified main.go)
internal/
  bot/                   Discord gateway app + handlers
  commands/              Bot commands (prefix + slash), help, pagination
  api/                   Echo HTTP server (auth, dashboard endpoints)
  config/                Env config
  domain/                Domain entities (pure structs)
  ports/                 DTOs + repository interfaces & implementations
  services/auth/         Discord OAuth2 + sessions
  logger/                slog setup
infrastructure/postgres/ Generated sqlc output — DO NOT hand-edit
db/migrations/           Goose SQL migrations
db/queries/              Named sqlc queries
dashboard/minji-bot/     React + Vite + shadcn dashboard
tests/                   Unit tests for the commands package
integration_tests/       DB integration tests
```

## Environment setup

1. **Create `.env`** in the repo root (gitignored). See `README.md` →
   *Getting started* for the full variable list.
2. **Start Postgres** and apply migrations:

   ```sh
   make docker/up          # starts Postgres on port 5434
   make goose-migrate-up   # applies schema migrations
   ```

> Postgres listens on host port **5434** (not 5432).

## Branching & workflow

- **`main`** — production. Merged via pull request.
- **`develop`** — integration branch with the latest merged work.
- Feature work happens on short-lived branches, typically named
  `nk/<descriptor>` (e.g. `nk/logging`).

Basic flow:

```sh
git checkout develop
git pull
git checkout -b nk/my-feature
# ...make changes, commit in small atomic steps...
git push -u origin nk/my-feature
# open a pull request against develop (or main for releases)
```

Keep each commit focused on one logical change and reference the feature in the
message. Commit messages follow the conventional style:

```
feat: add deleted message logging
fix: handle empty log channel on status
docs: document the dashboard API
refactor: split slash command registration into its own file
```

## Building & testing

```sh
go build ./...                 # compile
go vet ./...                   # static analysis (no CI linter yet)
gofmt -w .                     # format before committing
make test                      # unit tests (go test ./...)
make integration-test          # DB integration tests (needs TESTING_DB)
```

Always run `go build`, `go vet`, `go test`, and `gofmt` before opening a PR.
For dashboard-only changes, also run `npm run typecheck` (and `npm run build`)
in `dashboard/minji-bot`.

## Code conventions

- **Go**: standard library + dependencies already in `go.mod`. Format with
  `gofmt`. Every package should have a package-level doc comment.
- **Generated code**: never hand-edit `infrastructure/postgres/`. It is
  produced by `sqlc generate`.
- **Commands**: each command has a prefix handler (`foo(s, m, args)`) and a
  slash handler (`fooSlash(s, i)`). Dispatch goes through `Handle()` /
  `HandleSlash()` in `commands.go`. Slash command definitions live in
  `slash_commands.go`.
- **Moderation**: reuse the shared helpers in `moderation.go`
  (`requireModerator`, `resolveTargetUser`, `logModAction`, etc.) instead of
  re-implementing permission checks or audit logging.

## Database changes

1. Add a **migration** in `db/migrations/` — goose format
   (`-- +goose Up` / `-- +goose Down`), UTC-timestamp filename, e.g.
   `20260904100000_add_column.sql`.
2. Add named **queries** in `db/queries/<table>.sql`.
3. Regenerate the Go code:

   ```sh
   sqlc generate
   ```

4. Apply the migration:

   ```sh
   make goose-migrate-up
   ```

5. If the repository interface changed, update the corresponding
   `internal/ports/repository/*.go` implementation and any callers.

## Dashboard changes

```sh
cd dashboard/minji-bot
npm install
npm run dev         # Vite on :5173, /api proxied to :8080
```

Before submitting dashboard changes:

```sh
npm run typecheck
npm run lint
npm run build
```

## Pull requests

- Target **`develop`** for feature work; **`main`** only for releases.
- Keep the diff small and reviewable. Split large features into several PRs or
  logical commits.
- Include a summary of the change and note any new env vars, migrations, or
  breaking changes.
- Make sure CI-equivalent checks pass locally: `go build`, `go vet`,
  `go test`, `gofmt`, and dashboard `typecheck`/`build` where relevant.

## Migrations & the test database

`TESTING_DB` (default `minjitest`) is a separate database on the same Postgres
instance, created by `init-testdb.sh` on first container start. Apply schema to
it with `make test-migrate-up` before running integration tests.

## Release notes

Releases merge `develop` into `main`. Keep `to-do.md` and `features.md`
accurate so they reflect the actual state of the codebase.
