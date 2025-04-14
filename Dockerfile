FROM golang:1.24-alpine AS builder

RUN go install github.com/air-verse/air@latest && go install github.com/pressly/goose/v3/cmd/goose@latest

RUN go install github.com/a-h/templ/cmd/templ@latest

WORKDIR /app

COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o main .

# Development stage
FROM golang:1.24-alpine AS development

RUN apk add --no-cache make git curl


COPY --from=builder /go/bin/air /usr/local/bin/air
COPY --from=builder /go/bin/goose /usr/local/bin/goose

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

EXPOSE 8080

ENV PATH="/root/go/bin:${PATH}"

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["air"]

FROM golang:1.24-alpine AS production

RUN apk add --no-cache make git curl

RUN go install github.com/a-h/templ/cmd/templ@latest

COPY --from=builder /go/bin/air /usr/local/bin/air
COPY --from=builder /go/bin/goose /usr/local/bin/goose

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

EXPOSE 8080

ENV PATH="/root/go/bin:${PATH}"

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["make" , "run"]


