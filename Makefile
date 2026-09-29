-include .env

DB_URL = postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

createdb:
	createdb -U $(POSTGRES_USER) -h $(POSTGRES_HOST) $(POSTGRES_DB)

dropdb:
	dropdb -U $(POSTGRES_USER) -h $(POSTGRES_HOST) $(POSTGRES_DB)

migrateup:
	migrate -path db/migration -database "$(DB_URL)" up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" down -all

.PHONY: createdb dropdb migrateup migratedown
