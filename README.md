# MinjiBot

A Discord bot with a companion REST API and web dashboard, built in Go.

> **Status:** early development. The bot supports prefix commands (`-`) and slash commands. The API serves authentication (Discord OAuth) and dashboard data endpoints. The dashboard is a React + Vite + shadcn app.

## Tech stack

| Layer | Technology |
| --- | --- |
| Language | **Go 1.26** — module `github.com/kibetnathan/minjibot` |
| Discord | **discordgo** — gateway bot (WebSocket), prefix + slash commands |
| HTTP API | **Echo v5** — REST API on `:8080` |
| Database | **PostgreSQL 15** via **pgx/v5**, queries generated with **sqlc** |
| Migrations | **goose** — SQL schema migrations |
| Logging | **slog** — structured JSON logging |
| Dashboard | **React 19 + Vite + shadcn/ui** — `dashboard/minji-bot` |
| Auth | **Discord OAuth2** — session-cookie-based |

## Prerequisites

| Tool | Purpose |
| --- | --- |
| Go 1.26+ | Build and test |
| Docker | Local Postgres |
| Node.js | Dashboard dev server |
| goose | SQL migrations |
| sqlc | Query codegen |

## Getting started

### 1. Create `.env`

```dotenv
# Database
DB_URL=postgres://postgres:<password>@localhost:5434/minjibot?sslmode=disable

# Discord
DISCORD_TOKEN=
DISCORD_CLIENT_ID=
DISCORD_CLIENT_SECRET=

# API / Auth
APP_URL=http://localhost:8080
FRONTEND_URL=http://localhost:5173
SESSION_SECRET=<random-secret>

# Optional integrations
GOOGLE_FACTCHECK_API_KEY=
GEMINI_API_KEY=
GEMINI_MODEL=

# Postgres container (docker compose)
POSTGRES_USER=postgres
POSTGRES_PASSWORD=<password>
POSTGRES_DB=minjibot
TESTING_DB=minjitest

# Goose (migration targets)
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://postgres:<password>@localhost:5434/minjibot?sslmode=disable
GOOSE_MIGRATION_DIR=db/migrations
```

### 2. Start the database

```sh
make docker/up
```

Postgres is published on host port **5434**. On first start, `init-testdb.sh` also creates the `TESTING_DB` database.

### 3. Apply migrations

```sh
make goose-migrate-up
```

### 4. Generate database queries

```sh
sqlc generate
```

### 5. Run the bot + API

```sh
make run              # unified entrypoint (cmd/main.go)
# or separately:
make run-bot          # Discord bot only
make run-api          # HTTP API only
```

### 6. Run the dashboard

```sh
cd dashboard/minji-bot
npm install
npm run dev           # Vite dev server on :5173
```

The Vite dev server proxies `/api` requests to `:8080`.

## Make targets

| Target | Description |
| --- | --- |
| `make run` | Run bot + API together (`cmd/main.go`) |
| `make run-bot` | Start Discord bot (gateway) |
| `make run-api` | Start HTTP API on `:8080` |
| `make test` | `go test ./...` |
| `make integration-test` | Run `integration_tests/db_test.go` against `TESTING_DB` |
| `make goose-migrate-up` | Apply migrations to dev database |
| `make goose-migrate-down` | Rollback last migration |
| `make docker/up` | Start Postgres in background |
| `make docker/down` | Stop containers (keeps volumes) |
| `make docker/down/v` | Stop containers and wipe data |
| `make docker/logs` | Tail database logs |

The Makefile loads and exports `.env` automatically.

## Architecture

```
cmd/
  main.go            Unified entrypoint — starts both bot + API
  bot/main.go        Bot-only entrypoint
  api/main.go        API-only entrypoint

internal/
  bot/               Discord bot app (gateway connection, handler registration)
    handlers/        Gateway event handlers (message, delete, interaction)
  commands/          Bot commands (prefix + slash), help pagination, tldr
  api/               Echo HTTP server (auth, dashboard log endpoints)
  config/            Environment variable parsing (caarlos0/env)
  domain/            Domain entities (guild, guildsettings, auditlog, etc.)
    auditlog/
    birthday/
    deletedmessage/
    diary/
    guild/
    guildsettings/
    user/
    userpermission/
  ports/
    dto/             Data transfer objects for repository calls
    repository/      Repository interfaces + SQL implementations
  services/auth/     Discord OAuth2 flow + session management
  logger/            Structured slog logger setup

infrastructure/
  postgres/          Generated sqlc output (do not edit)

db/
  migrations/        Goose SQL schema files (UTC-timestamp filenames)
  queries/           Named sqlc query definitions

dashboard/
  minji-bot/         React + Vite + shadcn dashboard

tests/               External unit tests for the commands package
integration_tests/   Database integration tests
```

### Data flow

```
Discord Gateway → bot/handlers → commands → ports/repository → infrastructure/postgres → PostgreSQL
                                                                                          ↑
Echo HTTP API → api/auth, api/logs → ports/repository ─────────────────────────────────────┘
                                                                                          ↑
Dashboard (React) → /api/* ────────────────────────────────────────────────────────────────┘
```

### Database workflow

1. Edit `db/queries/*.sql` (named sqlc queries) or `db/migrations/` (schema)
2. Run `sqlc generate` to regenerate `infrastructure/postgres/`
3. Run `make goose-migrate-up` to apply schema changes

**Never edit** `infrastructure/postgres/` directly — it is generated-only.

## API endpoints

| Method | Path | Auth | Description |
| --- | --- | --- | --- |
| GET | `/healthz` | No | Health check |
| GET | `/api/auth/discord` | No | Start Discord OAuth flow |
| GET | `/api/auth/callback/discord` | No | OAuth callback |
| POST | `/api/auth/logout` | Session | Clear session |
| GET | `/api/auth/me` | Session | Current user info |
| GET | `/api/guilds` | Session | List guilds with counts |
| GET | `/api/logs/deleted?guild_id=` | Session | Deleted messages (paginated) |
| GET | `/api/logs/actions?guild_id=` | Session | Mod actions (paginated) |

## Dashboard routes

| Route | Description |
| --- | --- |
| `/` | Landing page |
| `/commands` | Command reference |
| `/dashboard` | Guild picker with logging stats |
| `/dashboard/guild/:id` | Deleted messages + mod actions for a guild |
| `/login` | Login page |
| `/signup` | Sign up page |

## Bot commands

The bot supports ~100+ commands across these categories:

- **General** — ping, echo, userinfo, test, bug
- **Moderation** — ban, kick, timeout, warn, purge, nuke, jail, role management, etc.
- **Utility** — translate, search, emoji management, reminders, polls
- **Fun** — birthday, diary, ship, polls, roleplay emotions/actions
- **Information** — avatar, banner, botinfo, guild stats, help, weather

See `features.md` for the full command reference.

## Deployment

### Railway (current production)

The bot + API run as a single service on Railway. The dashboard is deployed separately on Netlify.

**Environment variables:**

| Variable | Description |
| --- | --- |
| `DB_URL` | PostgreSQL connection string |
| `DISCORD_TOKEN` | Bot token |
| `DISCORD_CLIENT_ID` | Application ID (for slash command registration) |
| `DISCORD_CLIENT_SECRET` | OAuth2 client secret |
| `APP_URL` | API origin (e.g., `https://minjibot-bot-production.up.railway.app`) |
| `FRONTEND_URL` | Dashboard origin (e.g., `https://minji-bot.netlify.app`) |
| `SESSION_SECRET` | Random string for session cookie signing |

### Render (alternative)

Render provides free tier Web Services and managed PostgreSQL. See the Makefile and deployment configs for reference.

## License

Private — not for redistribution.
