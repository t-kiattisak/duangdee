# Payment & Credit Service Architecture

## 1. Double-Entry Ledger & Gateway Architecture

```
+-----------------------------------------------------------------------+
|                    1. Delivery Layer (REST & Webhooks)                |
|       - GET /payments/balance, POST /payments/checkout                |
|       - POST /payments/webhook/stripe, POST /payments/webhook/omise   |
+-----------------------------------------------------------------------+
                                    |
                                    v
+-----------------------------------------------------------------------+
|                    2. Usecase / Ledger Core Layer                     |
|       - DoubleEntryLedgerUsecase (Acid Credit Top-up / Usage)        |
|       - PaymentGatewayCheckoutUsecase                                 |
|       - WebhookVerificationUsecase                                    |
+-----------------------------------------------------------------------+
         |                                              |
         v                                              v
+------------------------------------+  +-------------------------------+
| Payment Gateway SDK Adapters       |  | Kafka Event Producer/Consumer |
| - Stripe SDK Adapter               |  | - Consumes: user.registered   |
| - PromptPay QR Adapter (Omise/Opn) |  | - Consumes: reading.initiated |
| - PostgresLedgerRepository         |  | - Publishes: credit.deducted  |
+------------------------------------+  | - Publishes: payment.completed|
                                        +-------------------------------+
```

## 2. Component Design & Ledger Mechanics

1. **Double-Entry Coin Ledger Engine**:
   - Every coin transaction (top-up, reading fee, refund, daily bonus) creates an immutable record in `credit_transactions`.
   - Uses PostgreSQL `SELECT ... FOR UPDATE` row-level locks to prevent race conditions or double-spending during credit deduction.
2. **PromptPay QR Code Generator Adapter**:
   - Integrates with Thai payment gateway (e.g. Omise / GB Prime Pay) to generate dynamic PromptPay QR codes with expiration TTL (15 mins).
3. **Stripe Checkout & Card Adapter**:
   - Creates Stripe PaymentIntents for credit card transactions.
   - Handles async payment fulfillment via cryptographic HMAC webhook signature validation.
4. **Kafka Event Processing Engine**:
   - `user.registered`: Automatically provisions a zero-balance `user_wallets` record.
   - `reading.initiated`: Atomically deducts coins and publishes `credit.deducted` (or failure event if insufficient balance).
