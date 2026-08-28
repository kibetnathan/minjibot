ifneq (,$(wildcard ./.env))
	include .env
	export
endif

# Dev/prod migrations use GOOSE_DRIVER / GOOSE_DBSTRING / GOOSE_MIGRATION_DIR from .env directly.
# Test migrations hit the same Postgres instance with the database swapped to TESTING_DB
# (created by init-testdb.sh on first container init).
TEST_DB_URL ?= $(subst /$(POSTGRES_DB),/$(TESTING_DB),$(GOOSE_DBSTRING))

.PHONY: all run test goose-migrate-up goose-migrate-down test-migrate-up test-migrate-down
.PHONY: docker/up docker/down docker/down/v docker/logs
.PHONY: run-api run-bot ngrok ngrok-url

all: test run

test:
	@echo "Running tests ..."
	go test ./tests/...

run:
	@echo "Starting Go application..."
	go run .

run-api:
	@echo "Starting API server on :8080..."
	go run ./cmd/api

run-bot:
	@echo "Starting Discord bot..."
	go run ./cmd/bot

ngrok:
	@echo "Starting ngrok tunnel to API (port 8080)..."
	@echo "Copy the HTTPS URL and set as Discord Interactions Endpoint URL"
	ngrok http 8080

ngrok-url:
	@curl -s http://localhost:4040/api/tunnels | jq -r '.tunnels[0].public_url'

goose-migrate-up:
	@echo "Running migrations..."
	goose up

goose-migrate-down:
	@echo "Rolling back migrations..."
	goose down

test-migrate-up:
	@test -n "$(TESTING_DB)" || { echo "Error: TESTING_DB is not set in .env."; exit 1; }
	@echo "Running migrations on test DB..."
	GOOSE_DBSTRING='$(TEST_DB_URL)' goose up

test-migrate-down:
	@test -n "$(TESTING_DB)" || { echo "Error: TESTING_DB is not set in .env."; exit 1; }
	@echo "Rolling back migrations on test DB..."
	GOOSE_DBSTRING='$(TEST_DB_URL)' goose down

docker/up: ## Start DB in background
	docker compose up -d

docker/down: ## Stop all containers (keeps volumes)
	docker compose down

docker/down/v: ## Stop all containers and delete volumes (wipes DB data)
	docker compose down -v

docker/logs: ## Tail DB logs
	docker compose logs -f database
