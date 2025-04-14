ifneq (,$(wildcard .env))
    include .env
    export
endif

.PHONY: start dev dev-build stop prod prod-build

start:
	@echo "📦 Starting database services..."
	@docker-compose up -d postgres pgadmin

migrate:
	@goose -dir db postgres "host=gothpg user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable" up

dev-build: start
	@echo "Starting the app build !!!"
	@docker-compose build dev

prod-build: start
	@echo "Starting the app build !!!"
	@docker-compose build prod

dev: dev-build
	@echo "Starting development environment..."
	@docker-compose up dev


prod: prod-build
	@echo "Starting development environment..."
	docker-compose up prod

stop:
	@docker-compose down

run:
	air
