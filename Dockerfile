FROM golang:1.24-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download


COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/server ./cmd/server

FROM alpine:3.18
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server /usr/local/bin/server

COPY --from=builder /src/config/config.yaml /etc/avito-mvp/config.yaml
COPY --from=builder /src/internal/infra/db/migrations /migrations

EXPOSE 8080

ENV PORT=8080 \
    CONFIG_PATH=/etc/avito-mvp/config.yaml \
    MIGRATIONS_DIR=/migrations

ENTRYPOINT ["/usr/local/bin/server"]
