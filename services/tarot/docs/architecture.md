# Tarot Core Service Architecture

## 1. High-Throughput Knowledge Base Architecture

```
+-----------------------------------------------------------------------+
|                    1. Delivery / Interface Layer                      |
|       - HTTP REST Handlers (Public Deck Viewer / Admin Catalog)       |
|       - gRPC Server Handlers (High-speed internal query from Reading)  |
+-----------------------------------------------------------------------+
                                    |
                                    v
+-----------------------------------------------------------------------+
|                    2. Usecase / Business Logic Layer                  |
|       - CardQueryUsecase (Get card details & images)                  |
|       - MeaningMatcherUsecase (Match position, orientation, category) |
|       - SpreadConfigUsecase (Load layout positions & rules)           |
+-----------------------------------------------------------------------+
                                    |
                                    v
+-----------------------------------------------------------------------+
|                    3. Multi-Level Caching & Repository                |
|       - Redis L2 Cache (In-memory Card Catalog & Meanings Dict)       |
|       - Postgres DB (Source of Truth for 78 cards & meanings)         |
+-----------------------------------------------------------------------+
```

## 2. Component Design & Responsibilities

1. **Card Catalog Engine**:
   - Manages static metadata of all 78 Tarot Cards (22 Major Arcana, 56 Minor Arcana).
   - Serves high-resolution card asset URLs (hosted on Cloudflare/S3).
2. **Contextual Meanings Engine**:
   - Stores granular interpretations split by **Orientation** (Upright vs Reversed) and **Life Contexts** (General, Love, Work, Finance, Health).
3. **Spread Layout Manager**:
   - Manages positions and meanings for various card spreads:
     - 1-Card Daily Draw
     - 3-Card Spread (Past / Present / Future)
     - 10-Card Celtic Cross Spread
4. **L2 Redis Caching System**:
   - Pre-loads all card meanings into Redis Hash maps at startup.
   - Ensures gRPC requests from `Reading Service` execute in **< 2ms** without querying PostgreSQL directly.
