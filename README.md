# MinjiBot

A Discord bot with a companion REST API and web dashboard, built in Go.

> **Status:** early development. Both the API (`cmd/api`) and bot (`cmd/bot`) entrypoints are functional. The bot supports prefix commands (`!`) and slash commands.

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
cmd/bot            Discord bot entrypoint (gateway + prefix/slash commands)
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

## Running the Bot locally

```sh
go run ./cmd/bot
```

The bot connects to Discord gateway and listens for both prefix commands (default `!`) and slash commands.

## Local development with ngrok (for slash commands & webhooks)

Slash commands require a public HTTPS endpoint for Discord to send interactions. Use ngrok to tunnel your local API:

1. Install ngrok: `brew install ngrok/ngrok/ngrok` (macOS) or download from https://ngrok.com/download
2. Authenticate: `ngrok config add-authtoken <your-token>`
3. Start the API locally: `make run` (runs on `:8080`)
4. In another terminal, start the tunnel:

   ```sh
   ngrok http 8080
   ```

5. Copy the HTTPS forwarding URL (e.g., `https://abc123.ngrok-free.app`)
6. Set it as your Discord application's **Interactions Endpoint URL** in the Developer Portal:
   - Go to Applications → Your Bot → General Information
   - Set "Interactions Endpoint URL" to `https://abc123.ngrok-free.app/interactions`
   - Save changes (Discord will send a verification request)

Now slash commands will route to your local machine.

## Deploying to Render

Render provides free tier Web Services (spin down after 15 min inactivity) and managed PostgreSQL.

### 1. Database (PostgreSQL)

- Create a new **PostgreSQL** database on Render
- Note the **External Database URL** (looks like `postgres://user:pass@host:port/db`)
- Run migrations against it:

  ```sh
  GOOSE_DBSTRING="<external-db-url>" goose -dir db/migrations up
  ```

### 2. Web Service (API + Bot)

- Create a new **Web Service** on Render
- Connect your GitHub repo
- Build command: `go build ./cmd/bot` (or `go build ./cmd/api` for API only)
- Start command: `./bot` (or `./api`)
- Environment variables:

  | Key            | Value                                    |
  | -------------- | ---------------------------------------- |
  | `DB_URL`       | Render's External Database URL           |
  | `DISCORD_TOKEN`| Your bot token from Discord Developer Portal |

- **Health Check Path**: `/health` (add a simple health endpoint to your API if needed)
- **Auto-Deploy**: Yes

> **Note:** The bot uses Discord gateway (WebSocket), not HTTP. Render Web Services support long-running WebSocket connections. If you need a separate worker for the bot, use a **Background Worker** service type instead (no public port required).

### 3. Interactions Endpoint

After deploy, your service gets a URL like `https://minjibot.onrender.com`.

Set Discord **Interactions Endpoint URL** to:
```
https://minjibot.onrender.com/interactions
```

## Uptime Monitoring

Free options to monitor your deployed bot/API:

| Service          | Free Tier                     | Setup                                    |
| ---------------- | ----------------------------- | ---------------------------------------- |
| **UptimeRobot**  | 50 monitors, 5 min interval   | Add HTTP(s) monitor → `https://your-app.onrender.com/health` |
| **Better Uptime**| 10 monitors, 3 min interval   | Similar to UptimeRobot                   |
| **Cronitor**     | 5 monitors, 1 min interval    | Add heartbeat monitor, ping from app     |
| **Healthchecks.io**| 20 checks, 1 min interval  | Add check, curl from cron or app         |

**Recommended:** UptimeRobot
1. Create account at https://uptimerobot.com
2. Add Monitor → HTTP(s)
3. URL: `https://your-app.onrender.com/health` (or `/` if no health endpoint)
4. Interval: 5 minutes
5. Alert contacts: Email, Discord webhook, Slack, etc.

For Render free tier (spins down after 15 min), uptime monitors will show "down" during cold starts. This is expected — the service spins up on first request (adds ~30-60s latency).