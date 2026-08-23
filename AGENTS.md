# MinjiBot

Go Discord bot + API (module `github.com/kibetnathan/minjibot`, Go 1.26, pgx/v5 + Echo + discordgo). Early scaffold: `cmd/bot/main.go` and `cmd/api/main.go` are empty files and all `internal/*` packages are empty — don't assume code exists yet. All deps in `go.mod` are still marked `// indirect`.

## Commands

- Build/check: `go build ./...`, `go vet ./...`. Nothing else (lint/test/CI) is wired up.
- Database: `docker compose up -d database` → Postgres on host port **5433** (not 5432). All config comes from `.env` (gitignored): `POSTGRES_*`, `TESTING_DB`, plus `GOOSE_DRIVER` / `GOOSE_DBSTRING` / `GOOSE_MIGRATION_DIR`.
- Migrations: goose format (`-- +goose Up`/`Down`) in `db/migrations/`, UTC-timestamp filenames. Because `.env` exports the `GOOSE_*` vars, plain `goose up` / `goose down` works from the repo root.
- Makefile loads + exports `.env`; targets work without Infisical: `goose-migrate-up/down` (plain goose against `.env` config), `test-migrate-up/down` (same instance, database swapped to `TESTING_DB`), `docker/up|down|down/v|logs`. `make test` (`go test ./tests/...`) fails until `tests/` contains Go files; `make run` needs a real entrypoint in `cmd/`.
- Codegen: `sqlc generate` uses `db/queries/*.sql` + `db/migrations/` as schema and writes into `infrastructure/postgres/` (package `postgres`, pgx/v5, emit_interface, JSON tags).

## Known breakage (verify before trusting)

- Host port **5433** is also used by an unrelated local container (`mlinzi_db`, project "mlinzi"). If `docker compose up -d database` fails with `port is already allocated`, stop that container first (`docker rm -f mlinzi_db`) — don't change MinjiBot's compose file.
- `sqlc generate` has not been run since `db/queries/*.sql` were added; expect first-generation churn (JSONB params → `[]byte`, named params → Go arg structs).

## Layout

- `cmd/bot` and `cmd/api` = entrypoints; `internal/{bot,api,domain,ports}` = app code; `infrastructure/postgres` = generated-only sqlc output (don't hand-edit); `dashboard/` and `tests/` are empty placeholders.
