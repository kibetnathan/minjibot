# Multi-stage build: compile the binary in a builder stage so container startup
# is fast (the API binds to :8080 almost immediately, satisfying Render's port
# scan). The runtime stage installs only goose, the migrations, and the binary.

### Builder stage ###
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src

# Cache dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source and build a static binary.
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/minjibot ./cmd/main.go

### Runtime stage ###
FROM golang:1.26-alpine AS runtime

RUN apk add --no-cache ca-certificates tzdata git

# Install the goose migration CLI (pinned for a stable CLI syntax).
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.27.3

WORKDIR /app

# Copy the prebuilt binary and startup script.
COPY --from=builder /out/minjibot /app/minjibot
COPY ./scripts/migrate.sh /app/scripts/migrate.sh
COPY ./db /app/db

RUN chmod +x /app/scripts/migrate.sh

ENV GOOSE_DRIVER=postgres \
    GOOSE_MIGRATION_DIR=/app/db/migrations

EXPOSE 8080

# The startup script migrates then runs the bot. Render overrides this via
# dockerCommand (render.yaml). Starting the already-built binary is fast, so the
# API binds to :8080 quickly and passes Render's health/port checks.
CMD ["sh", "/app/scripts/migrate.sh"]
