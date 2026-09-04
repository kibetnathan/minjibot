# MinjiBot — AI Agent Context

Go Discord bot + REST API + web dashboard.
Module: `github.com/kibetnathan/minjibot`, Go 1.26, pgx/v5 + Echo v5 + discordgo.

## Build & test

```sh
go build ./...        # compile all packages
go vet ./...          # static analysis
gofmt -w .            # format (no CI linter yet)
make test             # go test ./...
make integration-test # integration_tests/db_test.go against TESTING_DB
```

## Database

- **Postgres 15** on host port **5434** (not 5432): `docker compose up -d database`
- All config from `.env` (gitignored): `DB_URL`, `TESTING_DB`, `POSTGRES_*`
- Migrations: goose format in `db/migrations/` (UTC-timestamp filenames, `-- +goose Up/Down`)
- `make goose-migrate-up` / `make goose-migrate-down` (uses exported `.env` vars)

## Codegen

`sqlc generate` — reads `db/queries/*.sql` + `db/migrations/` schema, writes `infrastructure/postgres/` (package `postgres`, pgx/v5, `Querier` interface, JSON tags).
**Never hand-edit** `infrastructure/postgres/`.

## Architecture

```
cmd/
  main.go            Unified entrypoint (bot + API)
  bot/main.go        Bot-only
  api/main.go        API-only

internal/
  bot/               Discord gateway connection, handler registration
    handlers/        Gateway event handlers (message, message-delete, interaction)
  commands/          ~100+ bot commands (prefix + slash), help pagination, tldr
    commands.go      CommandHandler struct, Handle/HandleSlash dispatch, thin wrappers
    moderation.go    Shared mod helpers (perm checks, audit logging, embed builders)
    mod_*.go         Per-feature command implementations (ban, jail, purge, etc.)
    slash_commands.go Slash command definitions (SlashCommands variable)
    helpers.go       Interaction option helpers (OptInt, OptBool, OptUser)
    help.go          Help sections, BuildHelpPageEmbed
    pagination.go    Reaction-based (prefix) + button-based (slash) pagination
    rp.go            Roleplay emotion/action message helpers
  api/               Echo HTTP server
    app.go           App struct, NewApp, registerRoutes, CORS
    auth.go          Discord OAuth flow (/api/auth/discord, callback, logout, me)
    logs.go          Dashboard log endpoints (/api/guilds, /api/logs/*)
    session.go       resolveSession helper
  config/            env config (caarlos0/env/v11)
  domain/            Domain entities (pure structs, no deps)
  ports/
    dto/             Data transfer objects for repository calls
    repository/      Repository interfaces + SQL implementations
  services/auth/     Discord OAuth2 client + session manager
  logger/            slog JSON logger setup

infrastructure/postgres/  Generated sqlc output (do not edit)
db/migrations/            Goose SQL schema
db/queries/               Named sqlc queries
tests/                    External unit tests for commands package
integration_tests/        DB integration tests
dashboard/minji-bot/      React + Vite + shadcn dashboard
```

## Key files to know

| File | What it does |
|---|---|
| `internal/commands/commands.go` | `CommandHandler` struct, `Handle()` + `HandleSlash()` dispatch, thin wrappers that delegate to `mod_*.go` functions |
| `internal/commands/moderation.go` | Shared mod utilities: `effectiveModPerms`, `requireModerator`, `resolveTargetUser`, `auditAction`, `logModAction`, embed builders, option helpers |
| `internal/commands/slash_commands.go` | `SlashCommands` variable — all ApplicationCommand definitions |
| `internal/commands/helpers.go` | `OptInt`, `OptBool`, `OptUser` — interaction option extractors |
| `internal/bot/app.go` | Bot `App` struct, `NewApp`, `RegisterHandlers`, gateway intents |
| `internal/api/app.go` | API `App` struct, Echo setup, CORS, route wiring |
| `internal/api/logs.go` | Dashboard endpoints: guild list, deleted messages, mod actions |
| `internal/ports/repository/store.go` | `SQLStore` wrapping generated `Querier` |

## Conventions

- Prefix commands use `-` (e.g. `-ban`, `-help`)
- Slash commands registered globally on `Ready` via `s.ApplicationCommandCreate`
- Each command has both a prefix handler (`foo(s, m, args)`) and slash handler (`fooSlash(s, i)`)
- Slash dispatch goes through `HandleSlash()` in `commands.go`
- Mod actions log to both `audit_logs` table and configured log channel (via `logModAction`)
- Dashboard uses session cookie auth (Discord OAuth2 flow in `services/auth/`)

## Gotchas

- Bot gateway intents include `IntentsGuildMessageReactions` (needed for `-help` reaction pagination)
- `State.MaxMessageCount = 2000` — rolling cache for deleted message content
- `DISCORD_CLIENT_ID` falls back to bot user ID for slash command registration if unset
- `APP_URL` must be the API origin, `FRONTEND_URL` the dashboard origin
- Vite dev proxy: `/api` → `:8080`
- Postgres on port 5434 (not 5432)
- `init-testdb.sh` creates `TESTING_DB` on first docker start
- `go.mod` has some `// indirect` deps — don't move to direct until something imports them

## Dependencies to not touch

- `infrastructure/postgres/` — generated only
- `go.mod` indirect deps — leave as indirect until explicit import
