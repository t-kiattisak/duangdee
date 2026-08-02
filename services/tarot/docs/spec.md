# Tarot Core Service - Comprehensive Service Specification

---

## 1. Role & Purpose of the Service
The **Tarot Core Service** is the central Tarot knowledge base. It manages static metadata for all 78 tarot cards, stores multi-category upright and reversed card interpretations, manages spread layout definitions, and serves high-speed card details to clients and other services via Redis L2 cache and gRPC.

---

## 2. Detailed Business Logic & Rules

1. **Card Catalog Integrity**:
   - Maintains the complete 78-card Tarot deck (22 Major Arcana numbered 0-21, 56 Minor Arcana divided into Wands, Cups, Swords, Pentacles).
   - Serves high-resolution card asset URLs hosted on Cloudflare CDN / S3.
2. **Granular Meaning Dictionary**:
   - Every card has **10 distinct interpretation entries**:
     - 2 Orientations: **Upright (หัวตั้ง)** & **Reversed (หัวกลับ)**.
     - 5 Categories per orientation: **General (ทั่วไป)**, **Love (ความรัก)**, **Work (การงาน)**, **Finance (การเงิน)**, **Health (สุขภาพ)**.
3. **Spread Layout Rule Definitions**:
   - Defines position meanings for supported spreads:
     - `single-card`: Position 1 = General Advice.
     - `three-card`: Position 1 = Past, Position 2 = Present, Position 3 = Future.
     - `celtic-cross`: 10 distinct position meanings (Current State, Obstacle, Goal, Past, Recent Past, Future, Self, Environment, Hopes/Fears, Outcome).
4. **L2 In-Memory Caching (Redis)**:
   - At startup, the service loads the entire card meanings dictionary into Redis Hash maps.
   - All internal gRPC calls from Reading Service query Redis first to achieve sub-2ms response times.

---

## 3. Client Interaction & Request-Response Contracts (REST API)

### 3.1 `GET /api/v1/tarot/cards` (List Deck Catalog)
- **Client Sends**: `GET /api/v1/tarot/cards?arcana=major`
- **Client Receives (Response HTTP 200 OK)**:
  ```json
  {
    "status": "success",
    "data": {
      "cards": [
        {
          "id": 0,
          "name": "The Fool",
          "name": "เดอะ ฟูล (ผู้โง่เขลา)",
          "arcana_type": "major",
          "suit": null,
          "number": 0,
          "image_url": "https://cdn.duangdee.com/cards/0_the_fool.png"
        }
      ]
    }
  }
  ```

### 3.2 `GET /api/v1/tarot/cards/:id` (Get Single Card Detail)
- **Client Sends**: `GET /api/v1/tarot/cards/0`
- **Client Receives (Response HTTP 200 OK)**: Complete card info + meanings for all categories.

### 3.3 `GET /api/v1/tarot/spreads` (List Available Spreads)
- **Client Sends**: `GET /api/v1/tarot/spreads`
- **Client Receives (Response HTTP 200 OK)**: List of supported layouts and position descriptions.

---

## 4. Internal Service-to-Service Contracts (gRPC)

### `rpc BatchGetMeanings(BatchGetMeaningsRequest) returns (BatchGetMeaningsResponse)`
- **Caller**: Reading Engine Service
- **Purpose**: Fetch meanings for drawn cards in a specific reading session.
- **Request Payload**:
  ```json
  {
    "category": "love",
    "queries": [
      {"card_id": 0, "orientation": "upright", "position_index": 1},
      {"card_id": 6, "orientation": "reversed", "position_index": 2}
    ]
  }
  ```
- **Response Payload**:
  ```json
  {
    "results": [
      {
        "card_id": 0,
        "position_index": 1,
        "card_name": "The Fool",
        "orientation": "upright",
        "meaning": "การเริ่มต้นความรักครั้งใหม่ที่เต็มไปด้วยความตื่นเต้นและอิสระ...",
        "keywords": ["การเริ่มต้นใหม่", "ความรักอิสระ", "ความเสี่ยงที่น่าหลงใหล"]
      },
      {
        "card_id": 6,
        "position_index": 2,
        "card_name": "The Lovers",
        "orientation": "reversed",
        "meaning": "ความขัดแย้งในความสัมพันธ์ การตัดสินใจที่ผิดพลาด หรือทางขนาน...",
        "keywords": ["ความไม่เข้าใจกัน", "ทางเลือกที่ยากลำบาก", "ความร้าวฉาน"]
      }
    ]
  }
  ```
