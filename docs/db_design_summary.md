# Complete Database Design & Data Modeling Specification

## 1. Core Database Architecture Principles

The `duangdee` Tarot Backend follows these fundamental database design rules:
1. **Database-Per-Service Pattern**: Each microservice owns its isolated PostgreSQL database. Cross-database foreign keys are strictly prohibited; inter-service references use immutable UUIDs.
2. **Read-Heavy Optimization Strategy**: Static metadata (Tarot Deck & Card Meanings) uses PostgreSQL `TEXT[]` arrays and Redis L2 caching to eliminate JOIN overhead and achieve sub-millisecond query response times.
3. **ACID Financial Integrity**: Financial transactions (Coins & Balances) employ PostgreSQL `CHECK` constraints, immutable audit logs, and row-level locking (`SELECT ... FOR UPDATE`) to guarantee zero double-spending.

---

## 2. Global Database Inventory

| Database Name | Service Owner | Primary Engine | Storage Characteristics & Key Tables |
| :--- | :--- | :--- | :--- |
| **`auth_db`** | Auth & User Service | PostgreSQL + Redis | Identity, Auth Credentials, Social OAuth, Refresh Tokens, Zodiac Signs (`users`, `user_oauth_providers`). |
| **`tarot_db`** | Tarot Core Service | PostgreSQL + Redis L2 | Static Tarot Catalog, 10-Meaning Dictionary per Card, Spreads (`tarot_cards`, `card_meanings`, `tarot_spreads`). |
| **`reading_db`** | Reading Engine Service | PostgreSQL + Redis TTL | Reading Sessions, Drawn Cards, Synthesized Results, Daily Free Quota Counters (`reading_sessions`, `reading_drawn_cards`, `reading_results`). |
| **`payment_db`** | Payment Service | PostgreSQL (Acid) | User Coin Wallets, Double-Entry Audit Ledger, Payment Orders (`user_wallets`, `credit_transactions`, `payment_orders`). |
| **`notification_db`**| Notification Service | PostgreSQL + Redis Queue | FCM Push Tokens, User Notification Preferences, Delivery Logs (`push_subscriptions`, `notification_preferences`, `notification_logs`). |

---

## 3. Comprehensive Entity Schemas by Service

### 3.1 `auth_db` (Auth Service Database)

```sql
-- Main Users Table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255), -- NULL if social OAuth only
    display_name VARCHAR(100) NOT NULL,
    avatar_url TEXT,
    birth_date DATE,
    birth_time TIME,
    zodiac_sign VARCHAR(30),
    astrological_element VARCHAR(20),
    role VARCHAR(20) DEFAULT 'user', -- 'user', 'admin'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Social OAuth Accounts (Google / Line)
CREATE TABLE user_oauth_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL, -- 'google', 'line'
    provider_user_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(provider, provider_user_id)
);
```

---

### 3.2 `tarot_db` (Tarot Core Knowledge Base)

```sql
-- 78 Tarot Cards Master Catalog
CREATE TABLE tarot_cards (
    id INT PRIMARY KEY, -- 0 to 77
    name VARCHAR(100) NOT NULL, -- Card Name (e.g. 'The Fool', 'The Lovers')
    arcana_type VARCHAR(20) NOT NULL, -- 'major', 'minor'
    suit VARCHAR(20), -- 'wands', 'cups', 'swords', 'pentacles', NULL for Major
    number INT NOT NULL,
    element VARCHAR(20), -- 'fire', 'water', 'air', 'earth'
    image_url TEXT NOT NULL
);

-- Card Meanings Dictionary (10 Entries Per Card)
CREATE TABLE card_meanings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id INT NOT NULL REFERENCES tarot_cards(id),
    orientation VARCHAR(10) NOT NULL, -- 'upright', 'reversed'
    category VARCHAR(30) NOT NULL, -- 'general', 'love', 'work', 'finance', 'health'
    meaning TEXT NOT NULL, -- Full detailed interpretation
    keywords TEXT[] NOT NULL DEFAULT '{}' -- Stored as Native PostgreSQL Text Array (No JOINs required)
);

-- Spread Layout Configurations
CREATE TABLE tarot_spreads (
    id VARCHAR(50) PRIMARY KEY, -- e.g. 'single-card', 'three-card', 'celtic-cross'
    name VARCHAR(100) NOT NULL,
    card_count INT NOT NULL,
    coin_cost INT NOT NULL DEFAULT 0,
    positions_json JSONB NOT NULL -- Details of position meanings
);
```

---

### 3.3 `reading_db` (Reading Engine Database)

```sql
-- Reading Sessions Lifecycle
CREATE TABLE reading_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL, -- UUID reference to auth_db (No FK cross DB)
    spread_id VARCHAR(50) NOT NULL,
    category VARCHAR(30) DEFAULT 'general',
    question TEXT, -- Custom user intention prompt
    status VARCHAR(30) NOT NULL, -- 'initiated', 'pending_payment', 'drawn', 'completed'
    is_free_quota BOOLEAN DEFAULT FALSE,
    coins_charged INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Individual Drawn Cards per Session
CREATE TABLE reading_drawn_cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES reading_sessions(id) ON DELETE CASCADE,
    card_id INT NOT NULL,
    position_index INT NOT NULL,
    orientation VARCHAR(10) NOT NULL, -- 'upright', 'reversed'
    drawn_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Synthesized Reading Outcome & Rule Results
CREATE TABLE reading_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID UNIQUE NOT NULL REFERENCES reading_sessions(id) ON DELETE CASCADE,
    overall_summary TEXT NOT NULL,
    position_details JSONB NOT NULL, -- Matched meanings & keywords per position
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

---

### 3.4 `payment_db` (Payment & Wallet Ledger Database)

```sql
-- User Coin Wallets (Current Available Balance)
CREATE TABLE user_wallets (
    user_id UUID PRIMARY KEY, -- UUID reference to auth_db
    coin_balance INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_positive_balance CHECK (coin_balance >= 0) -- DB-level safeguard against negative balance
);

-- Immutable Double-Entry Audit Ledger
CREATE TABLE credit_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    amount INT NOT NULL, -- Positive (+100 for top-up), Negative (-10 for reading fee)
    balance_after INT NOT NULL, -- Snapshot balance after transaction
    transaction_type VARCHAR(30) NOT NULL, -- 'TOPUP', 'READING_FEE', 'DAILY_BONUS', 'REFUND'
    reference_id VARCHAR(255), -- session_id or payment_order_id
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Payment Top-up Orders (Stripe / PromptPay)
CREATE TABLE payment_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    provider VARCHAR(30) NOT NULL, -- 'stripe', 'promptpay'
    amount_thb DECIMAL(10, 2) NOT NULL,
    coins_rewarded INT NOT NULL,
    status VARCHAR(20) NOT NULL, -- 'pending', 'paid', 'failed', 'expired'
    external_charge_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

---

### 3.5 `notification_db` (Notification Database)

```sql
-- FCM / Web Push Tokens
CREATE TABLE push_subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    fcm_token TEXT NOT NULL,
    device_type VARCHAR(20) NOT NULL, -- 'web', 'ios', 'android'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, fcm_token)
);

-- User Notification Settings
CREATE TABLE notification_preferences (
    user_id UUID PRIMARY KEY,
    daily_horoscope BOOLEAN DEFAULT TRUE,
    promotions BOOLEAN DEFAULT FALSE,
    reading_reminders BOOLEAN DEFAULT TRUE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Notification Delivery Logs
CREATE TABLE notification_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    type VARCHAR(30) NOT NULL, -- 'email', 'push'
    subject TEXT,
    status VARCHAR(20) NOT NULL, -- 'sent', 'failed'
    sent_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

---

## 4. Key Performance & Reliability Techniques

1. **Indexing Strategy**:
   - `users(email)`: B-Tree unique index for fast authentication lookups.
   - `card_meanings(card_id, orientation, category)`: Composite index for instant zero-lag card dictionary matching.
   - `reading_sessions(user_id, created_at DESC)`: B-Tree index for loading user reading history efficiently.
   - `credit_transactions(user_id, created_at DESC)`: B-Tree index for fetching wallet audit logs.
2. **PostgreSQL Native Array (`TEXT[]`)**:
   - Used for `card_meanings.keywords` to eliminate 2 multi-table JOINs per reading request, delivering sub-millisecond query execution.
3. **ACID Row Locking (`SELECT FOR UPDATE`)**:
   - Used in `payment_db` during coin deductions to prevent race conditions or double-spending under high concurrency.
