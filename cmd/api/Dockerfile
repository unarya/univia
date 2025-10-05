# =============================================================================
# SECURE MULTI-STAGE DOCKERFILE
# =============================================================================

# -----------------------------------------------------------------------------
# Build Stage
# -----------------------------------------------------------------------------
FROM golang:latest AS builder

RUN apt update && apt install -y \
    ca-certificates \
    tzdata \
    build-essential \
    libwebp-dev \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /build

# Copy dependency files first for better caching
COPY go.mod go.sum ./
COPY main.go ./
# Download dependencies
RUN go mod download && go mod verify

# Build the application (enable CGO for libwebp)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -mod=mod \
    -ldflags='-w -s' \
    -o app ./main.go

# -----------------------------------------------------------------------------
# Development Stage
# -----------------------------------------------------------------------------
FROM golang:latest AS development
RUN apt update && apt install -y default-mysql-client
# Install dependencies
RUN go install github.com/air-verse/air@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest && \
    mv $(go env GOPATH)/bin/air /usr/local/bin/
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.19.0/migrate.linux-amd64.tar.gz \
    | tar xvz && mv migrate /usr/local/bin/
WORKDIR /app

COPY .air.toml ./
COPY go.mod go.sum ./
RUN go mod tidy && go mod download

RUN air init || true
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

ENTRYPOINT ["/bin/sh", "-c", "/app/entrypoint.sh"]

# -----------------------------------------------------------------------------
# Production Stage
# -----------------------------------------------------------------------------
FROM gcr.io/distroless/base-debian12:nonroot AS production

# Copy CA certificates & timezone
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy binary và set quyền thực thi
COPY --chmod=755 --from=builder /build/app /app

# Copy libwebp và dependency liên quan
COPY --from=builder /usr/lib/libwebp.so* /usr/lib/
COPY --from=builder /lib/x86_64-linux-gnu/libjpeg.so* /lib/x86_64-linux-gnu/
COPY --from=builder /lib/x86_64-linux-gnu/libpng*.so* /lib/x86_64-linux-gnu/
COPY --from=builder /lib/x86_64-linux-gnu/libm.so* /lib/x86_64-linux-gnu/
COPY --from=builder /lib/x86_64-linux-gnu/libc.so* /lib/x86_64-linux-gnu/

USER 65532:65532
EXPOSE 2000

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ["/app", "--health-check"] || exit 1

ENTRYPOINT ["/app"]
