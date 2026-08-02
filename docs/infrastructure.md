# Complete Infrastructure & Observability Architecture Spec

## 1. Infrastructure Stack Overview

The infrastructure for `duangdee` (Tarot Backend System) is designed to be cloud-native, scalable, secure, and observable.

```
+---------------------------------------------------------------------------------------------------+
|                                  Edge Layer (Cloudflare / Route53)                                |
|                                       - WAF & DDoS Protection                                     |
|                                       - SSL/TLS Termination                                       |
+---------------------------------------------------------------------------------------------------+
                                                  |
                                                  v
+---------------------------------------------------------------------------------------------------+
|                                      Ingress / API Gateway Layer                                  |
|                                    - Envoy Proxy / Nginx Ingress                                  |
|                                    - Rate Limiting & Auth Validation                              |
+---------------------------------------------------------------------------------------------------+
                                                  |
                                                  v
+---------------------------------------------------------------------------------------------------+
|                                  Kubernetes Cluster (EKS / GKE)                                   |
|                                                                                                   |
|  +-------------------+  +-------------------+  +-------------------+  +-------------------+  |
|  |   Auth Service    |  |   Tarot Service   |  |  Reading Service  |  |  Payment Service  |  |
|  +-------------------+  +-------------------+  +-------------------+  +-------------------+  |
|  +---------------------------------------------------------------------------------------------+  |
|  | Notification Service | Analytics Consumer | Payment Worker    | Cron Job Worker            |  |
|  +---------------------------------------------------------------------------------------------+  |
+---------------------------------------------------------------------------------------------------+
         |                                         |                                         |
         v                                         v                                         v
+------------------+                    +--------------------+                    +------------------+
| Stateful Data    |                    | Event Messaging    |                    | Observability    |
| - PostgreSQL     |                    | - Apache Kafka     |                    | Stack            |
|   (Managed RDS)  |                    |   (Strimzi/MSK)    |                    | - ElasticSearch  |
| - Redis Cluster  |                    | - Zookeeper/KRaft  |                    | - Logstash       |
|   (ElastiCache)  |                    |                    |                    | - Kibana (EFK)   |
| - S3 / GCS       |                    |                    |                    | - Prometheus     |
|   (Media Assets) |                    |                    |                    | - Grafana        |
+------------------+                    +--------------------+                    +------------------+
```

---

## 2. Infrastructure Components Breakdown

### 2.1 Compute Layer
- **Container Orchestration**: **Kubernetes (K8s)**
  - Auto-scaling (HPA - Horizontal Pod Autoscaler) based on CPU/Memory and request throughput.
  - Node Pools: Separate node pools for stateless microservices vs background workers.
- **Containerization**: **Docker** (Multi-stage lightweight alpine/distroless Go binaries).

### 2.2 Core Data Stores (Stateful Infrastructure)
- **Primary Relational Database**: **PostgreSQL** (Managed e.g., AWS RDS PostgreSQL or GCP Cloud SQL)
  - Isolated database per service (`auth_db`, `tarot_db`, `reading_db`, `payment_db`, `notification_db`).
  - Read-Replicas for heavy read traffic (e.g. Tarot Deck inquiries).
- **Cache & In-Memory Store**: **Redis Cluster** (Managed e.g., AWS ElastiCache)
  - Redis Sentinel / Cluster mode for high availability.
  - Used for Session Store, Rate Limiting, Free Quota counters, and Token Revocation List.
- **Object Storage**: **AWS S3 / GCP Cloud Storage**
  - Storage for Tarot Card high-resolution image assets, user avatars, and system backups.

### 2.3 Event Bus Infrastructure
- **Message Broker**: **Apache Kafka** (Managed AWS MSK or Self-Hosted via Strimzi Kafka Operator on K8s)
  - **KRaft / Zookeeper** mode for cluster coordination.
  - Multi-broker deployment with replication factor of 3 for zero data-loss event streaming.
  - **Schema Registry**: Confluent Schema Registry (for Protobuf/Avro schema versioning on Kafka topics).

---

## 3. Observability & Monitoring Stack (EFK / PLG Stack)

To ensure high availability, fast incident debugging, and real-time monitoring, we deploy a full **Observability Stack**:

```
[ Application Logs (Go Apps JSON Logs) ] ---> Fluentbit / Logstash ---> ElasticSearch ---> Kibana (Logs Viewer)

[ Go Metrics (Prometheus Client) ]      ---> Prometheus Server      ---> Grafana (Dashboards & Alerts)

[ Distributed Tracing (OpenTelemetry) ]  ---> Jaeger / Tempo         ---> Grafana / Kibana (Trace View)
```

### 3.1 Logging Stack (EFK / ELK Stack + Kibana)
- **Log Generator**: All Go services emit structured JSON logs (`zap` or `zerolog`) to `stdout`.
- **Log Collector (Fluentbit / Logstash)**: Runs as a DaemonSet in Kubernetes to scrape container logs.
- **Log Engine (ElasticSearch)**: Indexes all system logs, error tracebacks, Kafka event logs, and API Gateway access logs.
- **Log UI (Kibana)**:
  - **Centralized Log Searching**: Search across all microservices using correlation IDs (`trace_id`).
  - **Error Rate Dashboards**: Real-time graphs for HTTP 5xx errors, failed payment attempts, DB timeouts.
  - **Security Auditing**: Track failed login attempts and suspicious API access.

### 3.2 Metrics & Alerting (Prometheus + Grafana)
- **Prometheus**: Metrics collection engine polling Go service runtime stats (Memory, Goroutines, GC), HTTP latency (p95, p99), gRPC response times, DB Connection Pool usage, and Kafka Consumer Lag.
- **Grafana**: Visual dashboards for metrics + Alertmanager sending alerts to Slack / Discord / PagerDuty when anomalies occur.

### 3.3 Distributed Tracing (OpenTelemetry + Jaeger)
- **Trace Context Propagation**: Every incoming HTTP request gets assigned a `X-Trace-ID`.
- This `trace_id` is propagated across gRPC calls and embedded inside Kafka event headers.
- Allows developers to trace a single user request end-to-end through API Gateway -> Auth -> Reading -> Payment -> Kafka -> Notification.

---

## 4. DevOps & CI/CD Pipeline Infrastructure

- **Infrastructure as Code (IaC)**: **Terraform** (Deploys Cloud resources: K8s clusters, RDS, Redis, S3, Kafka).
- **CI/CD Automation**: **GitHub Actions**
  - Linting (`golangci-lint`) & Unit/Integration testing.
  - Docker Image Build & push to Container Registry (ECR / Artifact Registry).
  - GitOps Deployment using **ArgoCD** / **Helm**.
