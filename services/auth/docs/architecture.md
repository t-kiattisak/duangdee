# Auth & User Service Architecture

## 1. Internal Clean Architecture & Layering

The Auth Service strictly adopts **Clean Architecture** with four distinct layers:

```
+-----------------------------------------------------------------------+
|                    1. Delivery / Interface Layer                      |
|       - HTTP Handlers (Go Fiber Framework - gofiber/fiber)           |
|       - gRPC Server Handlers (auth.proto implementation)              |
+-----------------------------------------------------------------------+
                                    |
                                    v
+-----------------------------------------------------------------------+
|                    2. Usecase / Business Logic Layer                  |
|       - UserRegistrationUsecase                                       |
|       - UserAuthenticationUsecase                                     |
|       - TokenLifecycleUsecase (Generate/Validate/Refresh/Revoke)       |
|       - ZodiacCalculatorUsecase                                       |
+-----------------------------------------------------------------------+
                                    |
                                    v
+-----------------------------------------------------------------------+
|                    3. Domain & Repository Layer                       |
|       - User Entity, Wallet Reference Entity                          |
|       - UserRepository Interface (PostgreSQL)                         |
|       - TokenRepository Interface (Redis Session & Blacklist)         |
|       - EventPublisher Interface (Kafka Producer)                     |
+-----------------------------------------------------------------------+
                                    |
                                    v
+-----------------------------------------------------------------------+
|                    4. Infrastructure / Data Layer                     |
|       - PostgresUserRepositoryImpl (pgx / GORM)                       |
|       - RedisTokenRepositoryImpl (go-redis)                           |
|       - KafkaEventPublisherImpl (segmentio/kafka-go / confluent)      |
+-----------------------------------------------------------------------+
```

## 2. Component Design & Interactions

1. **Simple 3-Option Authentication Engine**:
   - **Method 1: Email + Password**: Bcrypt hashed (Cost 12), simple & traditional.
   - **Method 2: Google OAuth2 One-Tap**: Fast 1-click login for Web & Android users.
   - **Method 3: LINE Login**: Essential for Thai market users.
2. **Token Lifecycle**:
   - Generates signed JWT Access Tokens (~15 min expiry).
   - Generates Refresh Token UUIDs stored in Redis with 30-day TTL.
3. **Zodiac & Astrological Engine**:
   - Calculates Western & Thai Zodiac signs based on user birthdate and birth time.
4. **Kafka Publisher Component**:
   - Publishes `user.registered` event upon user creation to notify Payment Service and Notification Service.
