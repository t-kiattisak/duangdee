# Tarot Core Service - Comprehensive Service Specification

---

## 1. Role & Purpose of the Service
The **Tarot Core Service** is the central Tarot knowledge base. It manages static metadata for all 78 tarot cards, stores single-language upright and reversed card interpretations, manages spread layout definitions, and serves high-speed card details to clients and other services via Redis L2 cache and gRPC.

---

## 2. Detailed Business Logic & Rules

1. **Card Catalog Integrity**:
   - Maintains the complete 78-card Tarot deck (22 Major Arcana numbered 0-21, 56 Minor Arcana divided into Wands, Cups, Swords, Pentacles).
   - Serves high-resolution card asset URLs hosted on Cloudflare CDN / S3.
2. **Granular Meaning Dictionary**:
   - Every card has **10 distinct interpretation entries**:
     - 2 Orientations: **Upright (หัวตั้ง)** & **Reversed (หัวกลับ)**.
     - 5 Categories per orientation: **General (ทั่วไป)**, **Love (ความรัก)**, **Work (การงาน)**, **Finance (การเงิน)**, **Health (สุขภาพ)**.
3. **Native PostgreSQL Array Type (`TEXT[]`)**:
   - The `keywords` column is stored as PostgreSQL's native `TEXT[]` array.
   - Example DB Insert: `'{"ความไม่เข้าใจกัน", "ทางเลือกที่ยากลำบาก", "ความร้าวฉาน"}'`
   - In Go, `pgx` automatically scans `TEXT[]` into `[]string`.

---

## 3. Database Schema & Data Models

### Table: `card_meanings`
```sql
CREATE TABLE card_meanings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id INT NOT NULL REFERENCES tarot_cards(id),
    orientation VARCHAR(10) NOT NULL, -- 'upright', 'reversed'
    category VARCHAR(30) NOT NULL, -- 'general', 'love', 'work', 'finance', 'health'
    meaning TEXT NOT NULL,
    keywords TEXT[] NOT NULL DEFAULT '{}' -- Stored as Native PostgreSQL Text Array
);
```

### Example Database Record
```sql
INSERT INTO card_meanings (card_id, orientation, category, meaning, keywords)
VALUES (
    6, 
    'reversed', 
    'love', 
    'ความขัดแย้งในความสัมพันธ์ การตัดสินใจที่ผิดพลาด หรือทางขนาน...', 
    ARRAY['ความไม่เข้าใจกัน', 'ทางเลือกที่ยากลำบาก', 'ความร้าวฉาน']
);
```

---

## 4. Client Interaction & Request-Response Contracts (REST API)

### 4.1 `GET /api/v1/tarot/cards` (List Deck Catalog)
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
          "arcana_type": "major",
          "suit": null,
          "number": 0,
          "image_url": "https://cdn.duangdee.com/cards/0_the_fool.png"
        }
      ]
    }
  }
  ```

---

## 5. Internal Service-to-Service Contracts (gRPC)

### `rpc BatchGetMeanings(BatchGetMeaningsRequest) returns (BatchGetMeaningsResponse)`
- **Caller**: Reading Engine Service
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
