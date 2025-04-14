# Load .env if it exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

.PHONY: start dev dev-build stop prod prod-build

start:
	@# Check if Docker is installed
	@if ! command -v docker >/dev/null 2>&1; then \
		echo "🚨 Docker is not installed. Please install Docker first."; \
		exit 1; \
	fi
	@echo "📦 Starting database services..."
	@docker-compose up -d postgres pgadmin

migrate:
	@goose -dir db postgres "host=gothpg user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable" up

dev-build: start
	@echo "Starting the app build !!!"
	@BUILD_TARGET=development docker-compose build app

prod-build: start
	@echo "Starting the app build !!!"
	@BUILD_TARGET=production docker-compose build app

dev: dev-build
	@echo "Starting development environment..."
	@#First stop any existing app container to avoid conflicts
	@docker-compose stop app || true
	@#Then start the app container with the correct target
	BUILD_TARGET=development docker-compose up app


prod: prod-build
	@echo "Starting development environment..."
	@#First stop any existing app container to avoid conflicts
	@docker-compose stop app || true
	@#Then start the app container with the correct target
	BUILD_TARGET=production docker-compose up app

# Stop the container
stop:
	@docker-compose down

run:
	air
