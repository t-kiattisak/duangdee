# Notification Service - Comprehensive Service Specification

---

## 1. Role & Purpose of the Service
The **Notification Service** handles all asynchronous user messaging, transactional emails (welcome emails, payment receipts), web/mobile push alerts (FCM), and scheduled daily tarot reading reminders via Kubernetes Cron Jobs.

---

## 2. Detailed Business Logic & Rules

1. **Transactional Email Dispatching**:
   - Uses HTML email templates via SendGrid / AWS SES.
   - Retries failed email dispatches up to 3 times with exponential backoff before marking as `failed` in `notification_logs`.
2. **Web Push & FCM Push Notifications**:
   - Manages Firebase Cloud Messaging (FCM) push tokens registered per user and device type (web, android, ios).
   - Validates user notification preferences in `notification_preferences` table before sending alerts.
3. **Daily Horoscope Cron Scheduler**:
   - Runs daily at 07:00 AM (ICT).
   - Scans users with `daily_horoscope = TRUE` preference and queues push notifications reminding them to claim their daily free tarot card draw.

---

## 3. Client Interaction & Request-Response Contracts (REST API)

### 3.1 `POST /api/v1/notifications/push-subscription` (Register Device Push Token)
- **Client Sends (Request)**:
  ```json
  {
    "fcm_token": "fcm_token_string_abc123xyz...",
    "device_type": "web"
  }
  ```
- **Client Receives (Response HTTP 201 Created)**:
  ```json
  {
    "status": "success",
    "message": "Push token registered successfully"
  }
  ```

### 3.2 `GET /api/v1/notifications/preferences` (Get User Preferences)
- **Client Sends**: `GET /api/v1/notifications/preferences`
- **Client Receives (Response HTTP 200 OK)**:
  ```json
  {
    "status": "success",
    "data": {
      "daily_horoscope": true,
      "reading_reminders": true,
      "promotions": false
    }
  }
  ```

---

## 4. Kafka Event Consumer Integration

### Consumes: `user.registered`
- **Trigger**: New user account creation.
- **Action**: Dispatches Welcome Email with getting started guide.

### Consumes: `reading.completed`
- **Trigger**: Tarot reading finished.
- **Action**: Dispatches Push Notification to user's registered device ("ผลการดูดวงไพ่ทาโรต์ของคุณพร้อมแล้ว!").

### Consumes: `payment.completed`
- **Trigger**: Top-up payment successful.
- **Action**: Dispatches Email Receipt with transaction ID and coin amount.
