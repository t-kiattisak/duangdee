# Duangdee - Tarot Reading Backend (Go)

High-performance, scalable Tarot Card Reading Backend application developed in **Go (Golang)** using a **Microservices Architecture**.

## Architecture & Technology Highlights
- **Core Language**: Go 1.22+
- **Web Framework**: Go Fiber (`gofiber/fiber`) for high-throughput, low-allocation REST APIs.
- **API Gateway**: Kong Gateway for centralized authentication, rate limiting, and header transformation.
- **Databases**: PostgreSQL (Database-per-service pattern) & Redis Cluster (Caching, Session Store, Daily Free Quotas).
- **Messaging & Async Execution**: Apache Kafka (Event-Driven Architecture) + Asynq/Redis Queue (Background Task Workers).
- **Containerization & Observability**: Multi-stage Distroless Docker images (< 20MB), Kubernetes-ready, EFK Stack (ElasticSearch, Fluentbit, Kibana) for centralized logging.

---

## Documentation Index

### 📐 Global System Specifications
- **[System Architecture & Kafka Flow](docs/architecture.md)**
- **[Comprehensive Database Design & Schemas Summary](docs/db_design_summary.md)**
- **[Tarot Reading Modes & Pre-Reading Intention Spec](docs/reading_modes.md)**
- **[Credit System & Reading Quota Engine Spec](docs/credit_quota.md)**
- **[Docker Containerization & Local Dev Stack Spec](docs/docker.md)**
- **[Event-Driven Architecture & Task Queue Spec (Kafka vs Asynq)](docs/event_queue.md)**
- **[Kong API Gateway & Auth Header Forwarding Spec](docs/gateway_auth.md)**
- **[Infrastructure & Observability Stack](docs/infrastructure.md)**

### 📦 Microservices Architectural & Functional Specifications

| Service | Architecture Document (`architecture.md`) | Functional Specification (`spec.md`) |
| :--- | :--- | :--- |
| 🔑 **Auth Service** | [services/auth/docs/architecture.md](services/auth/docs/architecture.md) | [services/auth/docs/spec.md](services/auth/docs/spec.md) |
| 🃏 **Tarot Core Service** | [services/tarot/docs/architecture.md](services/tarot/docs/architecture.md) | [services/tarot/docs/spec.md](services/tarot/docs/spec.md) |
| 🔮 **Reading Service** | [services/reading/docs/architecture.md](services/reading/docs/architecture.md) | [services/reading/docs/spec.md](services/reading/docs/spec.md) |
| 💳 **Payment Service** | [services/payment/docs/architecture.md](services/payment/docs/architecture.md) | [services/payment/docs/spec.md](services/payment/docs/spec.md) |
| 🔔 **Notification Service** | [services/notification/docs/architecture.md](services/notification/docs/architecture.md) | [services/notification/docs/spec.md](services/notification/docs/spec.md) |