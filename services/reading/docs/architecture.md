# Reading Engine Service Architecture

## 1. Internal Clean Architecture & Domain Engine

```
+-----------------------------------------------------------------------+
|                    1. Delivery Layer (REST Endpoints)                 |
|       - POST /readings/initiate, POST /readings/:id/draw              |
|       - GET  /readings/history                                        |
+-----------------------------------------------------------------------+
                                    |
                                    v
+-----------------------------------------------------------------------+
|                    2. Reading Usecase & Domain Engines                |
|       - ReadingSessionCoordinatorUsecase                              |
|       - CryptoRandomizerEngine (`crypto/rand` Fisher-Yates Shuffle)   |
|       - RuleBasedInterpretationCombinatorEngine                       |
+-----------------------------------------------------------------------+
         |                                              |
         v                                              v
+------------------------------------+  +-------------------------------+
| gRPC Client Adapters               |  | Kafka Event Producer/Consumer |
| - TarotService gRPC Client         |  | - Publishes: reading.initiated|
| - PaymentService gRPC Client       |  | - Consumes: credit.deducted   |
+------------------------------------+  | - Publishes: reading.completed|
                                        +-------------------------------+
```

## 2. Component Design & Deep Dive Mechanics

1. **Cryptographic Randomizer Engine (`pkg/randomizer`)**:
   - Uses Go's CSPRNG (`crypto/rand`) instead of standard pseudo-random generators to guarantee zero-bias card shuffle.
   - Implements **Fisher-Yates Shuffle Algorithm** on array indices `[1..78]`.
   - Calculates binary orientation (0 = Upright, 1 = Reversed) using cryptographically secure random bits.
2. **Reading Session Lifecycle State Machine**:
   - Manages state progression: `INITIATED` -> `CHECKING_QUOTA` -> `CREDIT_HOLD` -> `CARD_DRAWN` -> `COMPLETED`.
3. **Rule-Based Interpretation Combinator Engine**:
   - Takes drawn card IDs, position indices, and orientations.
   - Fetches card position meanings from Tarot Service via gRPC.
   - Synthesizes position-by-position advice and overall reading summary based on the requested category (Love, Career, Finance).
4. **Daily Quota Manager**:
   - Checks Redis daily counter (`reading:daily_free_quota:<user_id>:<YYYY-MM-DD>`).
   - If free quota is available, bypasses coin deduction; otherwise, triggers Kafka `reading.initiated` to charge coins.
