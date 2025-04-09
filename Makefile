# Load .env if it exists
ifneq (,$(wildcard helpers/.env))
    include helpers/.env
    export
endif

.PHONY: db-start
db-start:
	@# Check if Docker is installed
	@if ! command -v docker >/dev/null 2>&1; then \
		echo "🚨 Docker is not installed. Please install Docker first."; \
		exit 1; \
	fi

		echo "🚀 Starting new PostgreSQL container: $(CONTAINER_NAME)"; \
		sudo docker run --name $(CONTAINER_NAME) \
			-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
			-e POSTGRES_USER=$(DB_USER) \
			-e POSTGRES_DB=$(DB_NAME) \
			-p $(DB_PORT):5432 \
			-d postgres; \

# Stop the container
.PHONY: db-stop
db-stop:
	sudo docker stop $(CONTAINER_NAME)

# Remove the container
.PHONY: db-clean
db-clean:
	sudo docker rm -f $(CONTAINER_NAME)

# Run migrations using Goose (assumes goose is installed)
.PHONY: migrate
migrate: db-start
	goose -dir db postgres "host=localhost user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable" up

# Run the app (assumes `go run main.go` works)
.PHONY: run
run:
	go run main.go
