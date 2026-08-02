# Reading Engine Service - Comprehensive Service Specification

---

## 1. Role & Purpose of the Service
The **Reading Engine Service** is the central domain orchestrator for tarot reading sessions. It executes the cryptographic card shuffle/draw algorithms, coordinates reading session state transitions, triggers coin deduction via Payment Service/Kafka, runs the rule-based interpretation synthesis engine, and saves historical reading logs.

---

## 2. Detailed Business Logic & Rules

1. **Cryptographic Randomizer (`crypto/rand`)**:
   - Uses Go's CSPRNG to execute the **Fisher-Yates Shuffle** on card indices `[0..77]`.
   - Determines orientation (Upright vs Reversed) using unbiased 50/50 crypto random bits.
   - Ensures drawn cards in a spread are strictly non-repeating (unique per reading).
2. **Quota & Credit Check Business Rules**:
   - Checks Redis key `reading:daily_free_quota:<user_id>:<YYYY-MM-DD>`.
   - **If counter == 0**: Increments counter to 1, marks session as `FREE_READ`, and proceeds immediately to card draw.
   - **If counter >= 1**: Requires credit payment. Marks session as `PENDING_PAYMENT` and publishes `reading.initiated` to Kafka for Payment Service processing.
3. **Rule-Based Interpretation Synthesis**:
   - Fetches position meanings from Tarot Core Service via gRPC.
   - Fetches user zodiac/astrological element from Auth Service via gRPC.
   - Synthesizes position-by-position advice and generates a cohesive reading outcome tailored to the user's category (Love, Career, Finance, Health).

---

## 3. Client Interaction & Request-Response Contracts (REST API)

### 3.1 `POST /api/v1/readings/initiate` (Start Reading Session)
- **Client Sends (Request)**:
  ```json
  {
    "spread_id": "three-card",
    "category": "love",
    "question": "ความรักกับคนนี้ในอีก 3 เดือนข้างหน้าจะเป็นอย่างไร?"
  }
  ```
- **Service Action**:
  1. Creates session in `reading_sessions` table with status `initiated`.
  2. Checks Redis daily free quota.
  3. If paid, publishes `reading.initiated` event to Kafka.
- **Client Receives (Response HTTP 201 Created)**:
  ```json
  {
    "status": "success",
    "data": {
      "session_id": "sess_88776655-4433-2211-00aa-bbccddeeff11",
      "spread_id": "three-card",
      "category": "love",
      "is_free_quota": true,
      "required_coins": 0,
      "session_status": "quota_granted"
    }
  }
  ```

### 3.2 `POST /api/v1/readings/:session_id/draw` (Draw Cards & Calculate Result)
- **Client Sends**: `POST /api/v1/readings/sess_88776655-4433-2211-00aa-bbccddeeff11/draw`
- **Service Action**:
  1. Verifies session status is eligible for draw.
  2. Runs `crypto/rand` Fisher-Yates shuffle.
  3. Draws 3 cards & orientation.
  4. Calls Tarot Service gRPC `BatchGetMeanings`.
  5. Synthesizes reading result and saves to `reading_results` table.
  6. Publishes `reading.completed` to Kafka.
- **Client Receives (Response HTTP 200 OK)**:
  ```json
  {
    "status": "success",
    "data": {
      "session_id": "sess_88776655-4433-2211-00aa-bbccddeeff11",
      "cards_drawn": [
        {
          "position_index": 1,
          "position_name": "อดีต",
          "card_id": 0,
          "card_name_th": "The Fool",
          "orientation": "upright",
          "image_url": "https://cdn.duangdee.com/cards/0_the_fool.png",
          "meaning": "การเริ่มต้นความรักครั้งใหม่ที่เต็มไปด้วยความตื่นเต้น..."
        }
      ],
      "overall_summary": "สรุปภาพรวมความรัก: ช่วงที่ผ่านมาคุณเริ่มต้นด้วยความอิสระ..."
    }
  }
  ```

### 3.3 `GET /api/v1/readings/history` (Get Past Reading History)
- **Client Sends**: `GET /api/v1/readings/history?page=1&limit=10`
- **Client Receives**: List of past user reading sessions and results.

---

## 4. Kafka Event Integration

### Outbound Event: `reading.initiated`
- **Trigger**: When non-free reading session is created.
- **Payload**: `{ "session_id": "...", "user_id": "...", "coins_required": 10 }`

### Inbound Event: `credit.deducted` (from Payment Service)
- **Action**: Updates session status to `payment_cleared`, unlocking card draw.

### Outbound Event: `reading.completed`
- **Trigger**: Published immediately after cards are drawn and result is stored.
- **Payload**: `{ "session_id": "...", "user_id": "...", "spread_id": "three-card" }`
- **Consuming Services**: Payment Service (Finalize deduction), Notification Service (Send push alert).
