FROM golang:1.26.4-alpine AS builder
WORKDIR /app
RUN go mod init server || true
COPY ./main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go
FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates curl && \
    curl -L --output /usr/local/bin/cloudflared https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64 && \
    chmod +x /usr/local/bin/cloudflared
COPY --from=builder /app/server .
ENV CF_TUNNEL_TOKEN="eyJhIjoiYzQ4M2YzZWI4NDlhZDdlZDQyYjk4ODljYWY2MmYxMjciLCJ0IjoiYWEzNGJhYzMtOThhMC00OTk2LWIyYmUtYWVjM2MwMjZjZTRjIiwicyI6Ik1qQTRObU0xTldZdE9HRTNNUzAwTWpVNExUZ3lNakV0Wmpnd1lXUXlNak0xTnpZNCJ9"
ENV PORT=8080
ENV HOST=0.0.0.0
EXPOSE 8080
CMD ["sh", "-c", "./server & exec cloudflared tunnel --no-autoupdate run --token $CF_TUNNEL_TOKEN"]