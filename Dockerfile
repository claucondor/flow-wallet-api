# ==============================================================================
# Multi-stage, multi-target Dockerfile for Flow Wallet API
# Supports: api, api-v2, admin, docs
# Targets: development, production, lightweight
# ==============================================================================

# Base Go builder with all dependencies
FROM golang:1.23.7-alpine AS base-builder

RUN apk update && apk add --no-cache \
    ca-certificates \
    musl-dev \
    gcc \
    build-base \
    git \
    bash

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# ==============================================================================
# API Builder - Builds main wallet API
# ==============================================================================
FROM base-builder AS api-builder
RUN go build -ldflags="-w -s" -o bin/api ./cmd/api

# ==============================================================================
# API V2 Builder - Builds future API version
# ==============================================================================
FROM base-builder AS api-v2-builder
RUN go build -ldflags="-w -s" -o bin/api-v2 ./cmd/api-v2

# ==============================================================================
# Admin Builder - Builds admin tools
# ==============================================================================
FROM base-builder AS admin-builder
RUN go build -ldflags="-w -s" -o bin/admin ./cmd/admin

# ==============================================================================
# Legacy Builder - Builds original main.go (for compatibility)
# ==============================================================================
FROM base-builder AS legacy-builder
RUN go build -ldflags="-w -s" -o bin/main .

# ==============================================================================
# Production Runtime Base
# ==============================================================================
FROM alpine:latest AS runtime-base

RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    curl \
    bash

WORKDIR /app

# Create non-root user
RUN addgroup -g 1001 -S app && \
    adduser -S -D -H -u 1001 -h /app -s /sbin/nologin -G app -g app app

USER app

# ==============================================================================
# API Production Target
# ==============================================================================
FROM runtime-base AS api-production

COPY --from=api-builder --chown=app:app /build/bin/api /app/
COPY --from=base-builder --chown=app:app /build/flow/ /app/flow/

EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:3000/v1/health/ready || exit 1

CMD ["./api"]

# ==============================================================================
# API V2 Production Target  
# ==============================================================================
FROM runtime-base AS api-v2-production

COPY --from=api-v2-builder --chown=app:app /build/bin/api-v2 /app/
COPY --from=base-builder --chown=app:app /build/flow/ /app/flow/

EXPOSE 3001

HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:3001/v1/health/ready || exit 1

CMD ["./api-v2"]

# ==============================================================================
# Admin Production Target
# ==============================================================================
FROM runtime-base AS admin-production

COPY --from=admin-builder --chown=app:app /build/bin/admin /app/

CMD ["./admin"]

# ==============================================================================
# Legacy Production Target (for backward compatibility)
# ==============================================================================
FROM runtime-base AS legacy-production

COPY --from=legacy-builder --chown=app:app /build/bin/main /app/
COPY --from=base-builder --chown=app:app /build/flow/ /app/flow/

EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:3000/v1/health/ready || exit 1

CMD ["./main"]

# ==============================================================================
# Development Target - Includes source code and development tools
# ==============================================================================
FROM base-builder AS development

# Install development tools
RUN go install github.com/air-verse/air@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

EXPOSE 3000 3001

# Use air for live reloading in development
CMD ["air", "-c", ".air.toml"]

# ==============================================================================
# Lightweight Target - Single binary with minimal footprint
# ==============================================================================
FROM alpine:latest AS lightweight

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Create non-root user and data directory
RUN addgroup -g 1001 -S app && \
    adduser -S -D -H -u 1001 -h /app -s /sbin/nologin -G app -g app app && \
    mkdir -p /data && \
    chown -R app:app /data

USER app

COPY --from=api-builder --chown=app:app /build/bin/api /app/
COPY --from=base-builder --chown=app:app /build/flow/ /app/flow/

EXPOSE 3000

CMD ["./api"]

# ==============================================================================
# All-in-One Target - Contains all binaries (useful for testing)
# ==============================================================================
FROM runtime-base AS all-in-one

COPY --from=api-builder --chown=app:app /build/bin/api /app/
COPY --from=api-v2-builder --chown=app:app /build/bin/api-v2 /app/
COPY --from=admin-builder --chown=app:app /build/bin/admin /app/
COPY --from=legacy-builder --chown=app:app /build/bin/main /app/
COPY --from=base-builder --chown=app:app /build/flow/ /app/flow/

# Create startup script
COPY <<EOF /app/start.sh
#!/bin/bash
case "\$1" in
  api) ./api ;;
  api-v2) ./api-v2 ;;
  admin) ./admin "\${@:2}" ;;
  main|legacy) ./main ;;
  *) echo "Usage: \$0 {api|api-v2|admin|main}" && exit 1 ;;
esac
EOF

RUN chmod +x /app/start.sh

EXPOSE 3000 3001

CMD ["./start.sh", "api"]

# ==============================================================================
# Default target
# ==============================================================================
FROM api-production AS production