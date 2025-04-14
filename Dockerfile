FROM golang:1.24-alpine AS builder

RUN go install github.com/air-verse/air@latest && go install github.com/pressly/goose/v3/cmd/goose@latest

RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN templ generate

RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o main .

# Development stage
FROM golang:1.24-alpine AS development

RUN apk add --no-cache make git curl

COPY --from=builder /go/bin/air /usr/local/bin/air
COPY --from=builder /go/bin/goose /usr/local/bin/goose

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN go mod download

EXPOSE 8080

ENV PATH="/root/go/bin:${PATH}"

CMD ["sh", "-c", "goose -dir db postgres \"host=${DB_HOST} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=disable\" up && air"]

#production stage
FROM golang:1.24-alpine AS production

RUN apk add --no-cache make git curl

RUN go install github.com/a-h/templ/cmd/templ@latest

COPY --from=builder /go/bin/air /usr/local/bin/air
COPY --from=builder /go/bin/goose /usr/local/bin/goose

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN go mod download

EXPOSE 8080

ENV PATH="/root/go/bin:${PATH}"

CMD ["sh", "-c", "goose -dir db postgres \"host=${DB_HOST} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=require\" up && make run"]


