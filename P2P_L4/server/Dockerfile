FROM golang:1.26-alpine AS builder
WORKDIR /app

# 1. Copy dependency files first
COPY go.mod go.sum ./
RUN go mod download

# 2. Copy source code
COPY . .

# 3. Clean any existing 'server' binary from host copy, then build target
RUN rm -f ./server && CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -o /app/server_bin .

FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache iproute2 ca-certificates curl && \
    ARCH=$(uname -m) && \
    if [ "$ARCH" = "x86_64" ]; then CF_ARCH="amd64"; \
    elif [ "$ARCH" = "aarch64" ]; then CF_ARCH="arm64"; \
    else CF_ARCH="amd64"; fi && \
    curl -L "https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-${CF_ARCH}" -o /usr/local/bin/cloudflared && \
    chmod +x /usr/local/bin/cloudflared

# Copy explicitly named binary
COPY --from=builder /app/server_bin /app/server
RUN chmod +x /app/server

RUN echo '#!/bin/sh' > /entrypoint.sh && \
    echo 'if [ ! -c /dev/net/tun ]; then' >> /entrypoint.sh && \
    echo '  mkdir -p /dev/net' >> /entrypoint.sh && \
    echo '  mknod /dev/net/tun c 10 200' >> /entrypoint.sh && \
    echo '  chmod 600 /dev/net/tun' >> /entrypoint.sh && \
    echo 'fi' >> /entrypoint.sh && \
    echo 'exec "$@"' >> /entrypoint.sh && \
    chmod +x /entrypoint.sh

COPY ./config.yaml .
EXPOSE 8443
ENTRYPOINT ["/entrypoint.sh"]
CMD ["./server"]