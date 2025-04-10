# Load .env if it exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

.PHONY: db-start
db-start:
	@# Check if Docker is installed
	@if ! command -v docker >/dev/null 2>&1; then \
		echo "🚨 Docker is not installed. Please install Docker first."; \
		exit 1; \
	fi

	sudo docker-compose up -d


# Stop the container
.PHONY: db-stop
db-stop:
	@sudo docker-compose down

# Remove the container
.PHONY: db-clean
db-clean:
	sudo docker rm -f $(CONTAINER_NAME)

# Run migrations using Goose (assumes goose is installed)
.PHONY: migrate
migrate: db-start
	@goose -dir db postgres "host=$(DB_HOST) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable" up

# Run the app (assumes `go run main.go` works)
.PHONY: run
run:
	go run main.go
