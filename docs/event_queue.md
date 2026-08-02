# Event-Driven Architecture & Queue Systems Specification

## 1. Executive Summary & Strategy

The `duangdee` Tarot Backend utilizes a hybrid asynchronous messaging architecture:
1. **Apache Kafka (Event Streaming / Event-Driven Architecture)**: For core domain events, inter-service notifications, analytics, and decoupled business workflows.
2. **Asynq / Redis Streams (Task Queues)**: For reliable transactional tasks requiring retries, backoff mechanisms, dead-letter queues (DLQ), and scheduled delayed executions.

---

## 2. Event-Driven Architecture (Apache Kafka Events)

Kafka is used when an event has **multiple potential consumers** or needs to be recorded as an immutable domain event log.

### 2.1 Event Map & Topology

```
+------------------+         Topic: user.registered          +------------------------+
|   Auth Service   | -------------------------------------> | Payment Service        |
+------------------+                                         | (Create Wallet)        |
         |                                                   +------------------------+
         |                                                   | Notification Service   |
         |                                                   | (Send Welcome Email)   |
         |                                                   +------------------------+

+------------------+         Topic: reading.initiated        +------------------------+
| Reading Service  | -------------------------------------> | Payment Service        |
+------------------+                                         | (Deduct Coins)         |
         ^                                                   +------------------------+
         |                   Topic: credit.deducted                      |
         +---------------------------------------------------------------+

+------------------+         Topic: reading.completed        +------------------------+
| Reading Service  | -------------------------------------> | Notification Service   |
+------------------+                                         | (Send Push Result)     |
                                                             +------------------------+
                                                             | Analytics Worker       |
                                                             +------------------------+
```

### 2.2 Detailed Kafka Events Table

| Event / Topic Name | Producer Service | Consumers | Trigger Condition | Why Asynchronous / Event-Driven? |
| :--- | :--- | :--- | :--- | :--- |
| **`user.registered`** | `Auth Service` | `Payment Service`, `Notification Service` | User registers via Email or OAuth. | Decouples account creation from wallet setup & email dispatch. Prevents slow email providers from delaying registration response. |
| **`reading.initiated`**| `Reading Service` | `Payment Service` | User requests a paid card reading. | Allows asynchronous payment check/hold without blocking HTTP connection. |
| **`credit.deducted`** | `Payment Service` | `Reading Service` | Wallet balance successfully debited. | Signals Reading Engine to proceed with Fisher-Yates card draw. |
| **`reading.completed`**| `Reading Service` | `Notification Service`, `Analytics Worker` | Tarot cards drawn and reading synthesized. | Triggers push notification to user's device and updates user engagement metrics. |
| **`payment.completed`**| `Payment Service` | `Auth Service`, `Notification Service` | Gateway webhook confirms top-up payment. | Triggers email receipt dispatch and updates user loyalty tier asynchronously. |

---

## 3. Task Queue Systems (Asynq / Redis Queue)

Task Queues are used for **point-to-point background jobs** that require strict execution guarantees, exponential backoff retries, and scheduled/delayed execution. We use **Asynq** (Go Distributed Task Queue backed by Redis).

### 3.1 Asynq Queue Topology

```
+-----------------------------------------------------------------------------------+
| Go Microservice (Producer)                                                        |
|   - Enqueue Task: AsynqClient.Enqueue(task, asynq.MaxRetry(5), asynq.ProcessIn(time)) |
+-----------------------------------------------------------------------------------+
                                          |
                                          v
+-----------------------------------------------------------------------------------+
| Redis Task Queue (Asynq Engine)                                                   |
|   - Queue: 'default' | Queue: 'high_priority' | Queue: 'scheduled'            |
+-----------------------------------------------------------------------------------+
                                          |
                                          v
+-----------------------------------------------------------------------------------+
| Background Worker Pool (Asynq Server Workers)                                     |
|   - Retries on Failure (Exponential Backoff)                                     |
|   - Dead Letter Queue (DLQ) for Failed Tasks                                      |
+-----------------------------------------------------------------------------------+
```

### 3.2 Task Queue Breakdown

| Queue Task Name | Owner Service | Trigger & Processing Logic | Retry Policy & DLQ Strategy |
| :--- | :--- | :--- | :--- |
| **`task:send_email_welcome`** | `Notification Service` | Triggered by `user.registered` event. Renders HTML template and dispatches via SendGrid/SES. | Max 5 retries. Exponential backoff (1m, 5m, 15m, 1h). Moves to DLQ if provider fails continuously. |
| **`task:send_email_receipt`** | `Notification Service` | Triggered by `payment.completed` event. Generates PDF receipt attachment and sends email. | Max 5 retries. |
| **`task:send_fcm_push`** | `Notification Service` | Triggered by `reading.completed`. Calls Firebase Cloud Messaging API. | Max 3 retries. Prunes invalid FCM tokens on 404 response. |
| **`task:expire_payment_qr`** | `Payment Service` | Scheduled when PromptPay QR is created (`ProcessIn(15 * time.Minute)`). Checks if order remains `PENDING` after 15 mins and updates status to `EXPIRED`. | Scheduled delayed task. Executed exactly once at TTL expiry. |
| **`task:daily_horoscope_cron`**| `Notification Service` | Scheduled via Kubernetes Cron / Asynq Cron at 07:00 AM daily. Batches subscribed users and enqueues individual FCM push tasks. | Cron Recurring Schedule. |

---

## 4. Key Differences: When to use Kafka vs Task Queue (Asynq)?

| Dimension | Apache Kafka (Event Bus) | Task Queue (Asynq / Redis) |
| :--- | :--- | :--- |
| **Primary Purpose** | Publish-Subscribe Event Notification (1-to-Many). | Point-to-Point Execution of a Specific Job (1-to-1). |
| **Data Retention** | Persistent log stream (re-playable history). | Transient queue (removed upon successful completion). |
| **Retry / Backoff** | Manual partition offset handling / Dead Letter Topics. | Built-in Automatic Exponential Backoff & DLQ UI. |
| **Delayed Execution**| Not natively built-in (requires delayed topics). | Native Support (`ProcessIn(duration)`). |
| **Use Cases in Duangdee**| `user.registered`, `reading.completed`, `payment.completed` | Email sending, FCM push alerts, PromptPay QR TTL Expiration. |
