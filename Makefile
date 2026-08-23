ifneq (,$(wildcard ./.env))
	include .env
	export
endif

GOOSE_DBSTRING ?= $(DB_URL)

.PHONY: all run goose-migrate-up goose-migrate-down

.PHONY: docker/up
docker/up: # Start DB in backgrounf
	$(SECRETS) docker compose up -d 

.PHONY: docker/down
docker/down: ## Stop all containers (keeps volumes)
	docker compose down

.PHONY: docker/down/v
docker/down/v: ## Stop all containers and delete volumes (wipes DB data)
	docker compose down -v

.PHONY: docker/logs
docker/logs: ## Tail DB logs
	$(SECRETS) docker compose logs -f database
