include .env
export $(shell sed 's/=.*//' .env)

.PHONY: migrate-up migrate-down sqlc build

migrate-up:
	goose -dir ./internal/infrastructure/db/migrations postgres $(DATABASE_URL) up

migrate-down:
	goose -dir ./internal/infrastructure/db/migrations postgres $(DATABASE_URL) down

sqlc:
	sqlc generate -f ./internal/infrastructure/db/sqlc/sqlc.yaml

build:
	go build -o backend ./cmd
