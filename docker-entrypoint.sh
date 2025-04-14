#!/bin/sh

# Run migrations
echo "Running migrations..."
goose -dir db postgres "host=postgres user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=disable" up

# Start the application
echo "Starting application..."
exec "$@"
