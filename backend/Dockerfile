FROM golang:1.23.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o backend_app .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/backend_app .
COPY --from=builder /app/config ./config
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/i18n ./i18n

# Create uploads directory
RUN mkdir -p /app/uploads

EXPOSE 8081

ENV CONFIG_FILE="/app/config/config.docker.yaml"

CMD ["./backend_app"]