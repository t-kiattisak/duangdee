# Payment & Credit Service - Comprehensive Service Specification

---

## 1. Role & Purpose of the Service
The **Payment & Credit Service** manages user monetization, coin/credit balances, double-entry audit logging, payment gateway integrations (Thai PromptPay QR & Stripe Checkout), and payment verification webhooks.

---

## 2. Detailed Business Logic & Rules

1. **ACID Double-Entry Ledger Mechanics**:
   - Maintains strict ledger integrity in `credit_transactions`.
   - Any coin increment (top-up) or decrement (reading fee) creates an immutable transaction row.
   - Credit deduction uses **PostgreSQL Row-Level Locking (`SELECT coin_balance FROM user_wallets WHERE user_id = $1 FOR UPDATE`)** to prevent race conditions or negative balances.
2. **PromptPay QR Code Generation (Thai Payments)**:
   - Integrates with Omise / GB Prime Pay API to generate dynamic PromptPay QR codes.
   - Sets a 15-minute expiration TTL on generated QR codes.
3. **Stripe Integration & Webhook Security**:
   - Creates Stripe Checkout Sessions for credit card payments.
   - Verifies Stripe Webhook signatures (`Stripe-Signature` header) using HMAC-SHA256 before processing payment fulfillment.
4. **Kafka Inter-Service Handling**:
   - `user.registered`: Provisions user wallet.
   - `reading.initiated`: Atomically deducts coins and publishes `credit.deducted`.

---

## 3. Client Interaction & Request-Response Contracts (REST API)

### 3.1 `GET /api/v1/payments/balance` (Get User Coin Balance)
- **Client Sends**: `GET /api/v1/payments/balance`
- **Client Receives (Response HTTP 200 OK)**:
  ```json
  {
    "status": "success",
    "data": {
      "coin_balance": 150,
      "currency": "COIN"
    }
  }
  ```

### 3.2 `POST /api/v1/payments/checkout` (Create Top-Up Order)
- **Client Sends (Request)**:
  ```json
  {
    "package_id": "pkg_100_coins",
    "payment_method": "promptpay" // 'promptpay' or 'stripe'
  }
  ```
- **Service Action**: Creates `payment_orders` record and requests QR code from Gateway.
- **Client Receives (Response HTTP 201 Created)**:
  ```json
  {
    "status": "success",
    "data": {
      "order_id": "ord_1122334455",
      "amount_thb": 100.00,
      "coins_to_receive": 100,
      "payment_method": "promptpay",
      "promptpay_qr_url": "https://cdn.duangdee.com/qr/qr_1122334455.png",
      "expires_at": "2026-08-02T14:40:00Z"
    }
  }
  ```

### 3.3 `POST /api/v1/payments/webhook/:provider` (Gateway Webhook Endpoint)
- **Sender**: Omise / Stripe Payment Gateway Server
- **Service Action**: Validates HMAC signature, marks order `PAID`, credits user wallet balance, and publishes `payment.completed` to Kafka.

---

## 4. Internal Service-to-Service Contracts (gRPC)

### `rpc CheckAndDeductCredit(DeductCreditRequest) returns (DeductCreditResponse)`
- **Caller**: Reading Engine Service
- **Request Payload**: `{ "user_id": "usr_99...", "coins": 10, "session_id": "sess_88..." }`
- **Response Payload**: `{ "success": true, "remaining_balance": 140 }`

---

## 5. Event Driven Integration (Kafka)

### Outbound Event: `credit.deducted`
- **Trigger**: Issued after coin deduction is committed in DB.
- **Payload**: `{ "session_id": "sess_88...", "user_id": "usr_99...", "status": "SUCCESS" }`

### Outbound Event: `payment.completed`
- **Trigger**: Issued after top-up payment webhook is verified.
- **Payload**: `{ "order_id": "ord_1122...", "user_id": "usr_99...", "coins_added": 100 }`
- **Consuming Services**: Notification Service (Sends email receipt).
