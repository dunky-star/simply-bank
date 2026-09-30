-include .env

DB_URL = postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

.DEFAULT_GOAL := help

help: ## Show available commands
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z_-]+:.*## / {printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

createdb: ## Create the simply-bank database
	createdb -U $(POSTGRES_USER) -h $(POSTGRES_HOST) $(POSTGRES_DB)

dropdb: ## Drop the simply-bank database
	dropdb -U $(POSTGRES_USER) -h $(POSTGRES_HOST) $(POSTGRES_DB)

migrateup: ## Apply all pending migrations
	migrate -path db/migration -database "$(DB_URL)" up

migratedown: ## Roll back all migrations
	migrate -path db/migration -database "$(DB_URL)" down -all

sqlc: ## Generate SQLC code
	sqlc generate

.PHONY: help createdb dropdb migrateup migratedown sqlc
