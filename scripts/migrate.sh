#!/bin/sh
# Container startup script: runs DB migrations, then starts the prebuilt bot
# binary. Used as the container command (see render.yaml dockerCommand).
#
# Requires:
#   GOOSE_DRIVER      (defaults to postgres)
#   GOOSE_DBSTRING     (the Postgres connection string, also used by the app as DB_URL)
#   GOOSE_MIGRATION_DIR (defaults to /app/db/migrations)
#
# Set GOOSE_VERBOSE=1 for verbose output.

set -e

cd /app

: "${GOOSE_DRIVER:=postgres}"
: "${GOOSE_MIGRATION_DIR:=/app/db/migrations}"

if [ -z "$GOOSE_DBSTRING" ]; then
  echo "[startup] ERROR: GOOSE_DBSTRING is not set. Provide the Postgres connection string." >&2
  exit 1
fi

echo "[startup] Running database migrations (driver=$GOOSE_DRIVER, dir=$GOOSE_MIGRATION_DIR)..."
if [ -n "$GOOSE_VERBOSE" ]; then
  goose -dir "$GOOSE_MIGRATION_DIR" "$GOOSE_DRIVER" "$GOOSE_DBSTRING" status
fi
goose -dir "$GOOSE_MIGRATION_DIR" "$GOOSE_DRIVER" "$GOOSE_DBSTRING" up

echo "[startup] Database migrations complete. Starting bot..."
exec /app/minjibot
