# MinjiBot

Go Discord bot + REST API (module `github.com/kibetnathan/minjibot`, Go 1.26, pgx/v5 + Echo (v5) + discordgo). The app has grown past the initial empty scaffold: `cmd/` has real entrypoints (`cmd/bot` gateway bot, `cmd/api` HTTP server on :8080, and a unified `cmd/main.go` that runs both), `internal/*` packages hold working app code, and `tests/` contains external unit tests for the `commands` package. A few deps in `go.mod` are still `// indirect` — don't move those to direct until something imports them explicitly.

## Commands

- Build/check: `go build ./...`, `go vet ./...`. `make test` runs unit tests (`go test ./...`); `make integration-test` runs `integration_tests/db_test.go` against the `TESTING_DB`. Everything is formatted with gofmt (no CI/lint is wired up yet).
- Database: `docker compose up -d database` → Postgres published on host port **5434** (not 5432). All config comes from `.env` (gitignored): `DB_URL`, `TESTING_DB`, `POSTGRES_*`, plus `GOOSE_DRIVER` / `GOOSE_DBSTRING` / `GOOSE_MIGRATION_DIR`.
- Migrations: goose format (`-- +goose Up`/`Down`) in `db/migrations/`, UTC-timestamp filenames. Because `.env` exports the `GOOSE_*` vars, plain `goose up` / `goose down` works from the repo root.
- Makefile loads + exports `.env`; targets work without Infisical: `goose-migrate-up/down` (plain goose against `.env` config), `test-migrate-up/down` (same instance, database swapped to `TESTING_DB`), `docker/up|down|down/v|logs`. `make run` runs `go run .` (unified `cmd/main.go`, starts bot + API together); `make run-bot` / `make run-api` run each entrypoint individually.
- Codegen: `sqlc generate` uses `db/queries/*.sql` + `db/migrations/` as schema and writes into `infrastructure/postgres/` (package `postgres`, pgx/v5, emit_interface, JSON tags). It is generated-only — never hand-edit it.

## Known breakage / gotchas (verify before trusting)

- The bot talks to Discord over the **gateway (WebSocket)**; slash commands are registered globally on `Ready` and handled via `internal/bot/handlers`. There is currently **no HTTP `/interactions` endpoint** — `internal/api/app.go` only wires Echo with no routes yet. README's ngrok instructions assume an HTTP interaction endpoint that isn't implemented; treat that section as aspirational.
- Bot gateway intents include `IntentsGuildMessageReactions` (needed for the prefix `-help` reaction pagination).
- The prefix `-help` reaction pagination and slash `/help` button pagination live in `internal/commands/pagination.go`.

## Layout

- `cmd/{bot,api}` + `cmd/main.go` = entrypoints; `internal/{bot,api,commands,config,domain,logger,ports}` = app code (commands and its tests moved out of `domain`); `infrastructure/postgres` = generated-only sqlc output (don't hand-edit); `db/migrations` + `db/queries` = goose SQL + named sqlc queries; `integration_tests/` = DB integration tests; `dashboard/minji-bot` = React + Vite + shadcn dashboard placeholder; `tests/` = external unit tests for the `commands` package.
