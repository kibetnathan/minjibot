# MinjiBot

A Discord bot with a companion REST API and web dashboard, built in Go.

> **Status:** early development. The API entrypoint (`cmd/api`) is functional; the bot entrypoint (`cmd/bot`) is a stub.

## Tech stack

- **Go 1.26** — module `github.com/kibetnathan/minjibot`
- **discordgo** — Discord gateway bot
- **Echo v5** — HTTP API (listens on `:8080`)
- **PostgreSQL 15** via **pgx/v5**, queries generated with **sqlc**
- **goose** — SQL migrations
- **slog** — structured JSON logging
- **React + Vite + shadcn/ui** — dashboard (`dashboard/minji-bot`)
- **react-router-dom** — client-side routing
- **lucide-react + react-icons** — icons

## Prerequisites

| Tool    | Used for                    |
| ------- | --------------------------- |
| Go 1.26+ | build / test               |
| Docker  | local Postgres              |
| goose   | migrations                  |
| sqlc    | query codegen               |
| Node.js | dashboard                   |

## Getting started

1. Create a `.env` in the repo root (gitignored):

   ```dotenv
   # App / API
   DB_URL=postgres://postgres:<password>@localhost:5434/minjibot?sslmode=disable
   DISCORD_TOKEN=

   # Postgres container (used by docker compose)
   POSTGRES_USER=postgres
   POSTGRES_PASSWORD=<password>
   POSTGRES_DB=minjibot
   TESTING_DB=minjitest

   # Goose (used by migration targets)
   GOOSE_DRIVER=postgres
   GOOSE_DBSTRING=postgres://postgres:<password>@localhost:5434/minjibot?sslmode=disable
   GOOSE_MIGRATION_DIR=db/migrations
   ```

2. Start the database:

   ```sh
   make docker/up          # docker compose up -d
   ```

   Postgres is published on host port **5434**. On first start, `init-testdb.sh` also creates the `TESTING_DB` database.

3. Apply migrations:

   ```sh
   make goose-migrate-up   # dev DB (POSTGRES_DB)
   ```

4. Generate database queries:

   ```sh
   sqlc generate
   ```

5. Run the API:

   ```sh
   make run                # go run . -> http://localhost:8080
   ```

## Make targets

| Target                | Action                                              |
| --------------------- | --------------------------------------------------- |
| `make all`            | test + run                                          |
| `make test`           | `go test ./tests/...`                               |
| `make run`            | start the API                                       |
| `make goose-migrate-up/down`     | migrate the dev database                 |
| `make test-migrate-up/down`      | migrate `TESTING_DB` on the same instance |
| `make docker/up`      | start Postgres in background                        |
| `make docker/down`    | stop containers (keeps volumes)                     |
| `make docker/down/v`  | stop containers **and wipe data**                   |
| `make docker/logs`    | tail database logs                                  |

The Makefile loads and exports `.env`, so no Infisical or extra tooling is required.

## Database & codegen workflow

- Schema lives in `db/migrations/` as goose SQL files (`-- +goose Up` / `-- +goose Down`, UTC-timestamp filenames).
- Queries live in `db/queries/*.sql` (one file per table).
- `sqlc generate` compiles migrations + queries into `infrastructure/postgres/` (package `postgres`, pgx/v5, `Querier` interface, JSON tags). That directory is generated-only — never hand-edit it.
- Repositories in `internal/ports/repository` wrap the generated `Querier` behind per-entity interfaces and map rows to domain structs:

  ```sh
  sqlc generate           # after editing db/queries or db/migrations
  ```

## Project layout

```
cmd/api            API entrypoint (:8080)
cmd/bot            Discord bot entrypoint (stub)
internal/api       Echo app wiring
internal/config    env config (caarlos0/env/v11)
internal/logger    slog setup
internal/domain    domain entities
internal/ports     dto + repository interfaces/implementations
infrastructure/postgres   generated sqlc output (do not edit)
db/migrations      goose SQL migrations
db/queries         named sqlc queries
dashboard/minji-bot       React + Vite + shadcn dashboard (with landing page)
tests              integration tests (WIP)
```

## Dashboard

```sh
cd dashboard/minji-bot
npm install
npm run dev        # Vite dev server
```

Other scripts: `build`, `preview`, `lint`, `format`, `typecheck`.

Routes: `/` (landing), `/dashboard`, `/commands`, `/docs`, `/auth`.