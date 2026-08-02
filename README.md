# Duangdee (ดวงดี) - Tarot Reading Backend (Go)

ระบบ Backend บริการดูดวงไพ่ทาโรต์พัฒนาด้วยภาษา **Go (Golang)** ออกแบบในรูปแบบ **Microservices Architecture**

## Documentation Index

### 📐 Global System Specs
- **[System Architecture & Kafka Flow](docs/architecture.md)**
- **[Tarot Reading Modes & Pre-Reading Intention Spec](docs/reading_modes.md)**
- **[Credit System & Reading Quota Engine Spec](docs/credit_quota.md)**
- **[Docker Containerization & Local Dev Stack Spec](docs/docker.md)**
- **[Event-Driven Architecture & Task Queue Spec (Kafka vs Asynq)](docs/event_queue.md)**
- **[Kong API Gateway & Auth Header Forwarding Spec](docs/gateway_auth.md)**
- **[Infrastructure & Observability Stack](docs/infrastructure.md)**

### 📦 Services Breakdown

| Service | Architecture (`architecture.md`) | Functional Spec (`spec.md`) |
| :--- | :--- | :--- |
| 🔑 **Auth Service** | [services/auth/docs/architecture.md](services/auth/docs/architecture.md) | [services/auth/docs/spec.md](services/auth/docs/spec.md) |
| 🃏 **Tarot Core Service** | [services/tarot/docs/architecture.md](services/tarot/docs/architecture.md) | [services/tarot/docs/spec.md](services/tarot/docs/spec.md) |
| 🔮 **Reading Service** | [services/reading/docs/architecture.md](services/reading/docs/architecture.md) | [services/reading/docs/spec.md](services/reading/docs/spec.md) |
| 💳 **Payment Service** | [services/payment/docs/architecture.md](services/payment/docs/architecture.md) | [services/payment/docs/spec.md](services/payment/docs/spec.md) |
| 🔔 **Notification Service** | [services/notification/docs/architecture.md](services/notification/docs/architecture.md) | [services/notification/docs/spec.md](services/notification/docs/spec.md) |