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

# goose needs a connection string. Prefer GOOSE_DBSTRING, then fall back to
# DB_URL (used by the app itself), then DATABASE_URL (the Railway Postgres
# plugin convention).
: "${GOOSE_DBSTRING:=${DB_URL:-$DATABASE_URL}}"

if [ -z "$GOOSE_DBSTRING" ]; then
  echo "[startup] ERROR: No Postgres connection string. Set GOOSE_DBSTRING, DB_URL, or DATABASE_URL." >&2
  exit 1
fi

# goose reads DRIVER/DBSTRING/MIGRATION_DIR from the environment. Export them so
# the plain `goose up` / `goose status` commands pick them up.
export GOOSE_DRIVER GOOSE_MIGRATION_DIR GOOSE_DBSTRING

echo "[startup] Running database migrations (driver=$GOOSE_DRIVER, dir=$GOOSE_MIGRATION_DIR)..."
if [ -n "$GOOSE_VERBOSE" ]; then
  goose status
fi
goose up

echo "[startup] Database migrations complete. Starting bot..."
exec /app/minjibot
