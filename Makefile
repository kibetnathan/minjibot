ifneq (,$(wildcard ./.env))
	include .env
	export
endif

GOOSE_DBSTRING ?= $(DB_URL)

.PHONY: all run goose-migrate-up goose-migrate-down


