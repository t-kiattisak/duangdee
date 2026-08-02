# Docker & Containerization Architecture Specification

## 1. Multi-Stage Production Dockerfile Blueprint (Go Microservices)

Every Go microservice (`auth`, `tarot`, `reading`, `payment`, `notification`) follows a standardized **Multi-Stage Build Dockerfile** using `golang:1.22-alpine` for building and `gcr.io/distroless/static-debian12` for execution.

### Advantages of Distroless Base Image:
1. **Ultra-Small Image Size**: Reduces image footprint from ~900MB to **< 20MB** per microservice.
2. **Maximum Security**: Distroless contains no shell (`/bin/sh`), package managers (`apk`, `apt`), or curl utilities, preventing RCE container breakout attacks.

### Standard Service Dockerfile (`Dockerfile.service`)
```dockerfile
# Stage 1: Build binary
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

# Cache Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked binary with CGO disabled
ARG SERVICE_NAME
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/bin/server ./services/${SERVICE_NAME}/cmd/main.go

# Stage 2: Final Minimal Security Runtime
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /

# Copy compiled binary from builder
COPY --from=builder /app/bin/server /server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Run as non-root user (UID 65532) for K8s security compliance
USER nonroot:nonroot

EXPOSE 8080 9090

ENTRYPOINT ["/server"]
```

---

## 2. Comprehensive Local Development Environment (`docker-compose.yml`)

The root `docker-compose.yml` orchestrates the entire local stack including all 5 Go services, Kong Gateway, PostgreSQL databases, Redis Cluster, Apache Kafka (KRaft), ElasticSearch, Kibana, and Prometheus.

```yaml
version: '3.8'

services:
  # =========================================================================
  # Core Infrastructure Databases & Message Brokers
  # =========================================================================
  postgres:
    image: postgres:16-alpine
    container_name: duangdee-postgres
    environment:
      POSTGRES_MULTI_DB: auth_db,tarot_db,reading_db,payment_db,notification_db
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: secretpassword
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./db/init-multiple-db.sh:/docker-entrypoint-initdb.d/init-multiple-db.sh
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    container_name: duangdee-redis
    command: redis-server --requirepass redissecret
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "-a", "redissecret", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5

  kafka:
    image: bitnami/kafka:3.7
    container_name: duangdee-kafka
    ports:
      - "9092:9092"
    environment:
      - KAFKA_CFG_NODE_ID=0
      - KAFKA_CFG_PROCESS_ROLES=controller,broker
      - KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093
      - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT
      - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=0@kafka:9093
      - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
    volumes:
      - kafka_data:/bitnami/kafka

  # =========================================================================
  # API Gateway
  # =========================================================================
  kong:
    image: kong:3.6-alpine
    container_name: duangdee-kong
    environment:
      KONG_DATABASE: "off" # DB-less declarative mode
      KONG_DECLARATIVE_CONFIG: /kong/kong.yml
      KONG_PROXY_ACCESS_LOG: /dev/stdout
      KONG_ADMIN_ACCESS_LOG: /dev/stdout
      KONG_PROXY_ERROR_LOG: /dev/stderr
      KONG_ADMIN_ERROR_LOG: /dev/stderr
      KONG_ADMIN_LISTEN: 0.0.0.0:8001
    ports:
      - "8000:8000" # HTTP Proxy (Client Entry Point)
      - "8443:8443" # HTTPS Proxy
      - "8001:8001" # Kong Admin API
    volumes:
      - ./config/kong.yml:/kong/kong.yml

  # =========================================================================
  # Observability Stack (EFK / Kibana / Prometheus)
  # =========================================================================
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.12.0
    container_name: duangdee-elasticsearch
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
    ports:
      - "9200:9200"
    volumes:
      - es_data:/usr/share/elasticsearch/data

  kibana:
    image: docker.elastic.co/kibana/kibana:8.12.0
    container_name: duangdee-kibana
    environment:
      - ELASTICSEARCH_HOSTS=http://elasticsearch:9200
    ports:
      - "5601:5601"
    depends_on:
      - elasticsearch

  # =========================================================================
  # Go Microservices
  # =========================================================================
  auth-service:
    build:
      context: .
      dockerfile: Dockerfile.service
      args:
        SERVICE_NAME: auth
    container_name: auth-service
    ports:
      - "8081:8080"
    environment:
      - PORT=8080
      - DB_SOURCE=postgres://postgres:secretpassword@postgres:5432/auth_db?sslmode=disable
      - REDIS_ADDR=redis:6379
      - KAFKA_BROKERS=kafka:9092
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy

  tarot-service:
    build:
      context: .
      dockerfile: Dockerfile.service
      args:
        SERVICE_NAME: tarot
    container_name: tarot-service
    ports:
      - "8082:8080"
    environment:
      - PORT=8080
      - DB_SOURCE=postgres://postgres:secretpassword@postgres:5432/tarot_db?sslmode=disable
      - REDIS_ADDR=redis:6379
    depends_on:
      postgres:
        condition: service_healthy

  reading-service:
    build:
      context: .
      dockerfile: Dockerfile.service
      args:
        SERVICE_NAME: reading
    container_name: reading-service
    ports:
      - "8083:8080"
    environment:
      - PORT=8080
      - DB_SOURCE=postgres://postgres:secretpassword@postgres:5432/reading_db?sslmode=disable
      - REDIS_ADDR=redis:6379
      - KAFKA_BROKERS=kafka:9092
      - TAROT_GRPC_ADDR=tarot-service:9090
      - AUTH_GRPC_ADDR=auth-service:9090
    depends_on:
      postgres:
        condition: service_healthy

  payment-service:
    build:
      context: .
      dockerfile: Dockerfile.service
      args:
        SERVICE_NAME: payment
    container_name: payment-service
    ports:
      - "8084:8080"
    environment:
      - PORT=8080
      - DB_SOURCE=postgres://postgres:secretpassword@postgres:5432/payment_db?sslmode=disable
      - REDIS_ADDR=redis:6379
      - KAFKA_BROKERS=kafka:9092
    depends_on:
      postgres:
        condition: service_healthy

  notification-service:
    build:
      context: .
      dockerfile: Dockerfile.service
      args:
        SERVICE_NAME: notification
    container_name: notification-service
    ports:
      - "8085:8080"
    environment:
      - PORT=8080
      - DB_SOURCE=postgres://postgres:secretpassword@postgres:5432/notification_db?sslmode=disable
      - REDIS_ADDR=redis:6379
      - KAFKA_BROKERS=kafka:9092
    depends_on:
      postgres:
        condition: service_healthy

volumes:
  postgres_data:
  redis_data:
  kafka_data:
  es_data:
```
