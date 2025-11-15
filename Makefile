ifneq (,$(wildcard .env))
    include .env
    export
endif

DATABASE_URL ?= postgres://imperial:imperial@localhost:5433/imperial?sslmode=disable
MIGRATIONS_PATH ?= ./migrations

.PHONY: run
run:
	go run cmd/api/main.go

.PHONY: migrate-up
migrate-up:
	migrate -database $(DATABASE_URL) -path $(MIGRATIONS_PATH) up

.PHONY: migrate-down
migrate-down:
	migrate -database $(DATABASE_URL) -path $(MIGRATIONS_PATH) down

.PHONY: swagger
swagger:
	~/go/bin/swag init -g cmd/api/main.go -o docs

.PHONY: build
build:
	go build -o bin/api ./cmd/api
