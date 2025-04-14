# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git curl

# Install Air for hot reloading
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Generate templ files
RUN templ generate

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Development stage
FROM golang:1.24-alpine AS development

# Install development dependencies
RUN apk add --no-cache git curl

# Install Air for hot reloading
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s

# Install templ and ensure it's in PATH
RUN go install github.com/a-h/templ/cmd/templ@latest && \
    ln -s /root/go/bin/templ /usr/local/bin/templ

# Set working directory
WORKDIR /app

# Copy the source code
COPY . .

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Generate templ files
RUN templ generate

# Expose port
EXPOSE 8080

# Override Air configuration to use port 8080
ENV AIR_PORT=8080
ENV PATH="/root/go/bin:${PATH}"

# Run the application with Air
CMD ["air", "-c", ".air.toml"]

# Production stage
FROM alpine:3.19 AS production

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata curl

# Create non-root user
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /app

# Copy the binary and necessary files
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/css ./css

# Create necessary directories
RUN mkdir -p /app/uploads && \
    chown -R appuser:appuser /app

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Use non-root user
USER appuser

# Expose port
EXPOSE 8080

# Set environment variables
ENV GIN_MODE=release
ENV TZ=UTC

# Run the application
CMD ["./main"]
