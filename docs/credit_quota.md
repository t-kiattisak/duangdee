# Credit & Reading Quota Engine Architecture Specification

## 1. Executive Summary & Core Responsibilities

The Credit and Quota system is designed using a **Cooperative Dual-Service Pattern**:
1. **Reading Engine Service (`services/reading`)**: Manages **Daily Free Quotas** and calculates the **Coin Cost** required per reading layout.
2. **Payment & Credit Service (`services/payment`)**: Acts as the **Wallet & Ledger Source of Truth**, managing user coin balances and executing atomic double-entry deductions.

---

## 2. Quota & Credit Ledger Flow Mechanics

```
User Requests Reading (POST /readings/initiate)
                        |
                        v
+-----------------------------------------------------------------------------------+
| 1. Reading Engine Service                                                         |
|    - Check Redis: `reading:daily_free_quota:<user_id>:<YYYY-MM-DD>`                |
+-----------------------------------------------------------------------------------+
                        |
        +---------------+---------------+
        |                               |
 (Quota == 0: Free Read)        (Quota >= 1: Paid Read)
        |                               |
        v                               v
[ Grant Free Read ]           [ Calculate Required Coins ]
  - Set Redis Quota = 1         (e.g., 3-Card = 10 Coins, Celtic Cross = 30 Coins)
  - Proceed to Card Draw                |
                                        v
                              Publish `reading.initiated` to Kafka
                                        |
                                        v
+-----------------------------------------------------------------------------------+
| 2. Payment & Credit Service                                                       |
|    - Execute PostgreSQL Transaction with Row Locking:                             |
|      `SELECT coin_balance FROM user_wallets WHERE user_id = $1 FOR UPDATE`        |
+-----------------------------------------------------------------------------------+
                                        |
                        +---------------+---------------+
                        |                               |
              (Balance >= Required Coins)       (Balance < Required Coins)
                        |                               |
                        v                               v
              [ Atomic Deduction ]             [ Reject Transaction ]
              - Write Ledger Entry               - Return Error: "INSUFFICIENT_CREDIT"
              - Publish `credit.deducted`
                        |
                        v
              [ Proceed to Card Draw ]
```

---

## 3. Detailed Data Models & Storage Locations

### 3.1 Daily Free Quota (Stored in Redis - Owned by Reading Service)
- **Redis Key Pattern**: `reading:daily_free_quota:<user_id>:<YYYY-MM-DD>`
- **Value**: Integer counter (e.g. `0` = Not used, `1` = Daily free quota consumed)
- **TTL**: Auto-expires at midnight (23:59:59 ICT)
- **Business Rule**: Every registered user gets **1 Free Reading per day** automatically.

### 3.2 Coin Balance & Audit Ledger (Stored in PostgreSQL - Owned by Payment Service)

#### Table: `user_wallets` (Current Coin Balance)
```sql
CREATE TABLE user_wallets (
    user_id UUID PRIMARY KEY,
    coin_balance INT NOT NULL DEFAULT 0, -- Current available coin balance
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_positive_balance CHECK (coin_balance >= 0) -- DB-level safeguard against negative balance
);
```

#### Table: `credit_transactions` (Immutable Ledger Audit Trail)
```sql
CREATE TABLE credit_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    amount INT NOT NULL, -- Positive (+100 for topup), Negative (-10 for reading)
    balance_after INT NOT NULL, -- Running balance snapshot
    transaction_type VARCHAR(30) NOT NULL, -- 'TOPUP', 'READING_FEE', 'DAILY_BONUS', 'REFUND'
    reference_id VARCHAR(255), -- session_id or payment_order_id
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

---

## 4. Spread Cost Table & Calculation Formula

How many readings a user can perform depends on their **Coin Balance** divided by the **Spread Coin Cost**:

| Spread Layout | Cards Drawn | Free Quota Eligible? | Coin Cost | User Capacity Example (Balance = 100 Coins) |
| :--- | :--- | :--- | :--- | :--- |
| **1-Card Daily Draw** | 1 Card | ✅ Yes (Daily Free) | **5 Coins** | 1 Free + 20 Paid Readings |
| **3-Card Past/Present/Future** | 3 Cards | ✅ Yes (Daily Free) | **10 Coins** | 1 Free + 10 Paid Readings |
| **Love Triangle Spread** | 5 Cards | ❌ No (Paid Only) | **20 Coins** | 5 Paid Readings |
| **10-Card Celtic Cross** | 10 Cards | ❌ No (Paid Only) | **30 Coins** | 3 Paid Readings |

---

## 5. API Response Contract Examples

### 5.1 Checking Balance & Capacity (`GET /api/v1/payments/balance`)
- **Client Receives (HTTP 200 OK)**:
```json
{
  "status": "success",
  "data": {
    "coin_balance": 100,
    "daily_free_quota_remaining": 1,
    "estimated_readings_available": {
      "one_card_draw": 21,
      "three_card_spread": 11,
      "celtic_cross_spread": 3
    }
  }
}
```

### 5.2 Initiating Reading when Balance is Low (`POST /api/v1/readings/initiate`)
- **Client Sends**: Celtic Cross Spread (Requires 30 Coins), User Balance = 10 Coins.
- **Client Receives (HTTP 402 Payment Required)**:
```json
{
  "status": "error",
  "error_code": "INSUFFICIENT_CREDITS",
  "message": "เหรียญไม่เพียงพอสำหรับการเปิดไพ่ชุดนี้ (ต้องการ 30 เหรียญ, คุณมี 10 เหรียญ)",
  "data": {
    "required_coins": 30,
    "current_balance": 10,
    "shortfall_coins": 20,
    "topup_url": "/api/v1/payments/checkout"
  }
}
```
