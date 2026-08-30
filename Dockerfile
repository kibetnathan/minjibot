# Runtime image that keeps the Go toolchain and source so the startup script
# (scripts/migrate.sh) can build the binary, run migrations, and launch the bot.

FROM golang:1.26-alpine AS runtime

RUN apk add --no-cache ca-certificates tzdata git

# Install the goose migration CLI.
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /app

# Cache dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source.
COPY . .

RUN chmod +x /app/scripts/migrate.sh

ENV GOOSE_DRIVER=postgres \
    GOOSE_MIGRATION_DIR=/app/db/migrations \
    GOCACHE=/tmp/.gocache

EXPOSE 8080

# The startup script is the default command; Render overrides it via
# dockerCommand (render.yaml). It builds, migrates, then runs the bot.
CMD ["sh", "/app/scripts/migrate.sh"]
