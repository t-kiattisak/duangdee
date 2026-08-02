# Docker & Containerization Architecture Specification

## 1. Per-Service Isolated Dockerfile Strategy

In a production-grade Microservices Architecture, **each microservice owns its standalone `Dockerfile`** placed directly inside its service directory:

```
services/
├── auth/
│   └── Dockerfile              # Isolated Auth Service Build
├── tarot/
│   └── Dockerfile              # Isolated Tarot Core Build
├── reading/
│   └── Dockerfile              # Isolated Reading Engine Build
├── payment/
│   └── Dockerfile              # Isolated Payment Ledger Build
└── notification/
    └── Dockerfile              # Isolated Notification Worker Build
```

### Why Place Dockerfile Inside Each Service?
1. **Independent CI/CD Pipelines**: Allows GitHub Actions / GitLab CI to build and deploy only the changed service without building the entire monorepo.
2. **Clean Ownership**: Each team/service maintains its specific build arguments and OS dependencies.
3. **Ultra-Small Distroless Runtime (< 20MB)**: Uses `golang:1.22-alpine` for building and `gcr.io/distroless/static-debian12:nonroot` for runtime.

---

## 2. Standardized Microservice Dockerfile (`services/[service]/Dockerfile`)

```dockerfile
# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

# Copy module dependencies
COPY go.mod go.sum* ./
RUN go mod download

# Copy service source code
COPY . .

# Compile static Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/bin/server ./cmd/main.go

# Minimal Security Runtime Stage
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /

COPY --from=builder /app/bin/server /server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER nonroot:nonroot

EXPOSE 8080 9090

ENTRYPOINT ["/server"]
```

---

## 3. Local Multi-Container Stack (`docker-compose.yml`)

The root `docker-compose.yml` builds each service from its isolated directory context:

```yaml
services:
  auth-service:
    build:
      context: ./services/auth
      dockerfile: Dockerfile
    container_name: auth-service
    ports:
      - "8081:8080"
    environment:
      - PORT=8080
      - DB_SOURCE=postgres://postgres:secretpassword@postgres:5432/auth_db?sslmode=disable
      - REDIS_ADDR=redis:6379

  tarot-service:
    build:
      context: ./services/tarot
      dockerfile: Dockerfile
    container_name: tarot-service
    ports:
      - "8082:8080"

  reading-service:
    build:
      context: ./services/reading
      dockerfile: Dockerfile
    container_name: reading-service
    ports:
      - "8083:8080"

  payment-service:
    build:
      context: ./services/payment
      dockerfile: Dockerfile
    container_name: payment-service
    ports:
      - "8084:8080"

  notification-service:
    build:
      context: ./services/notification
      dockerfile: Dockerfile
    container_name: notification-service
    ports:
      - "8085:8080"
```
