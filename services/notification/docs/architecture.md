# Notification Service Architecture

## 1. Event Consumer & Dispatcher Architecture

```
+-----------------------------------------------------------------------+
|                    1. Delivery & API Layer                            |
|       - POST /notifications/push-subscription (Register Token)        |
|       - GET/PUT /notifications/preferences (User Push Settings)      |
+-----------------------------------------------------------------------+
                                    |
                                    v
+-----------------------------------------------------------------------+
|                    2. Usecase & Dispatcher Layer                      |
|       - EmailNotificationDispatcherUsecase                            |
|       - WebPushDispatcherUsecase (FCM / WebPush)                      |
|       - DailyHoroscopeCronSchedulerUsecase                            |
+-----------------------------------------------------------------------+
         |                                              |
         v                                              v
+------------------------------------+  +-------------------------------+
| External Providers & DB Repos      |  | Kafka Event Consumer Workers  |
| - SendGrid / AWS SES Mailer        |  | - Consumer: user.registered   |
| - FCM Push SDK Adapter             |  | - Consumer: reading.completed |
| - NotificationLogRepository        |  | - Consumer: payment.completed |
+------------------------------------+  +-------------------------------+
```

## 2. Component Design & Dispatch Mechanics

1. **Transactional Email Engine**:
   - Integrates with SendGrid / AWS SES to send HTML template-based emails:
     - Welcome & Onboarding Email (triggered by `user.registered`).
     - Payment Receipt Email (triggered by `payment.completed`).
2. **Push Notification Engine (FCM / WebPush API)**:
   - Stores FCM registration tokens per user/device.
   - Sends real-time browser/mobile push notifications when a reading is completed or daily horoscope is ready.
3. **Daily Horoscope Cron Scheduler**:
   - Background worker running on Kubernetes cron schedule.
   - Queries users with enabled `daily_horoscope` preferences and dispatches daily Tarot card notification reminders.
4. **Kafka Consumer Workers**:
   - Asynchronously consumes domain events from Kafka topics without blocking primary user request flows.
