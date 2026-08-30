# ---- Build stage ----
FROM golang:1.26-alpine AS build
WORKDIR /src

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build the unified entrypoint (bot + API)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/minjibot ./cmd/main.go

# Build the goose migration CLI so it can run at deploy time.
RUN CGO_ENABLED=0 GOOS=linux go install github.com/pressly/goose/v3/cmd/goose@latest

# ---- Runtime stage ----
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

# Copy the app binary, goose binary, migration script and migrations.
COPY --from=build /out/minjibot /app/minjibot
COPY --from=build /go/bin/goose /usr/local/bin/goose
COPY --from=build /src/scripts/migrate.sh /app/scripts/migrate.sh
COPY --from=build /src/db/migrations /app/db/migrations

RUN chmod +x /app/minjibot /app/scripts/migrate.sh

ENV GOOSE_DRIVER=postgres
ENV GOOSE_MIGRATION_DIR=/app/db/migrations

EXPOSE 8080

ENTRYPOINT ["/app/minjibot"]
