# System Architecture Overview - Duangdee (Tarot Service)

## Tech Stack & Design Pattern
- **Language**: Go 1.22+
- **Pattern**: Microservices Architecture
- **API Gateway**: Kong Gateway (Centralized Auth, Rate Limiting, Header Injection)
- **Containerization**: Multi-Stage Docker Builds (Distroless Security Image) + Docker Compose
- **Database Architecture**: PostgreSQL (Database-per-service pattern) + Redis Cluster
- **Credit & Quota Engine**: Redis TTL (Daily Free Quota) + PostgreSQL Acid Double-Entry Ledger
- **Messaging & Event-Driven**: Apache Kafka (Event Streaming) + Asynq/Redis (Distributed Task Queue)
- **Communication**: HTTP REST (Gateway) + gRPC (Inter-service) + Apache Kafka (Event Bus)
- **Observability**: EFK Stack (ElasticSearch, Fluentbit, Kibana) + Prometheus + Grafana + OpenTelemetry

## Global Architecture Specifications
1. 📐 **[High-Level Architecture & Kafka Flow](architecture.md)**
2. 🗄️ **[Comprehensive Database Design & Schemas Summary](db_design_summary.md)**
3. 💰 **[Credit System & Reading Quota Engine Spec](credit_quota.md)**
4. 🔮 **[Tarot Reading Modes & Pre-Reading Intention Spec](reading_modes.md)**
5. 🐳 **[Docker Containerization & Docker-Compose Spec](docker.md)**
6. ⚡ **[Event-Driven Architecture & Task Queue Spec (Kafka vs Asynq)](event_queue.md)**
7. 🚪 **[Kong API Gateway & Auth Forwarding Specification](gateway_auth.md)**
8. 🏗️ **[Infrastructure & Kibana Observability](infrastructure.md)**

## Services Breakdown (Architecture & Functional Specifications)

| Service Name | Architecture Document | Functional Specification |
| :--- | :--- | :--- |
| 🔑 **Auth Service** | [services/auth/docs/architecture.md](../../services/auth/docs/architecture.md) | [services/auth/docs/spec.md](../../services/auth/docs/spec.md) |
| 🃏 **Tarot Core Service** | [services/tarot/docs/architecture.md](../../services/tarot/docs/architecture.md) | [services/tarot/docs/spec.md](../../services/tarot/docs/spec.md) |
| 🔮 **Reading Service** | [services/reading/docs/architecture.md](../../services/reading/docs/architecture.md) | [services/reading/docs/spec.md](../../services/reading/docs/spec.md) |
| 💳 **Payment Service** | [services/payment/docs/architecture.md](../../services/payment/docs/architecture.md) | [services/payment/docs/spec.md](../../services/payment/docs/spec.md) |
| 🔔 **Notification Service** | [services/notification/docs/architecture.md](../../services/notification/docs/architecture.md) | [services/notification/docs/spec.md](../../services/notification/docs/spec.md) |
