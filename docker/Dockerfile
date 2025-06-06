FROM golang:1.24.3-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /api

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-w -s" -o comments-system ./cmd/comments-system/main.go
RUN go build -ldflags="-w -s" -o migrator ./cmd/migrator/main.go

FROM alpine:latest

RUN apk add --no-cache \
    postgresql15-client \
    tzdata \
    bash

WORKDIR /app

COPY --from=builder /api/comments-system .
COPY --from=builder /api/migrator .
COPY --from=builder /api/configs ./configs
COPY --from=builder /api/migrations ./migrations
COPY docker/wait-for-postgres.sh /app/wait-for-postgres.sh
RUN chmod +x /app/wait-for-postgres.sh

ENV CONFIG_PATH=/app/configs/postgres.yaml

CMD ["sh", "-c", "/app/wait-for-postgres.sh postgres 5432 && ./comments-system"]