# Auth & User Service - Comprehensive Service Specification

---

## 1. Role & Purpose of the Service
The **Auth & User Service** serves as the identity provider and user profile manager. It supports **3 Completely Free Authentication Methods** (Email/Password, Google OAuth2, LINE Login), manages access/refresh tokens, stores user profile/astrological info, and notifies other services via Kafka.

---

## 2. Cost Analysis of Authentication Methods

| Auth Method | Is It Free? | Cost Structure | Requirements / Developer Account |
| :--- | :--- | :--- | :--- |
| 1. **Email + Password** | 🟢 **100% Free** | 0 บาท (ทำระบบเองใน Go Backend) | ไม่ต้องมีบัญชีภายนอก |
| 2. **Google OAuth2** | 🟢 **100% Free** | 0 บาท (ฟรีไม่มีโควต้าจำกัด) | บัญชี Google Cloud Console (ฟรี) |
| 3. **LINE Login** | 🟢 **100% Free** | 0 บาท (LINE Login API ฟรี) | บัญชี LINE Developers Account (ฟรี) |
| ❌ *SMS OTP (Phone)* | 🔴 *Not Free* | *ต้องจ่ายประมาณ 0.30 - 0.50 บาท/SMS* | *SMS Gateway (Twilio/THSMS)* |

---

## 3. Client Interaction & Request-Response Contracts (REST API)

### 3.1 `POST /api/v1/auth/register` (Email Register - 100% Free)
- **Client Sends**: `{ "email": "user@example.com", "password": "...", "display_name": "..." }`
- **Response**: JWT Tokens + User Profile.

### 3.2 `POST /api/v1/auth/login` (Email Login - 100% Free)
- **Client Sends**: `{ "email": "user@example.com", "password": "..." }`
- **Response**: JWT Tokens + User Profile.

### 3.3 `POST /api/v1/auth/oauth/google` (Google 1-Click Login - 100% Free)
- **Client Sends**: `{ "id_token": "google_id_token_string_from_frontend" }`
- **Service Action**: Verifies Google ID Token via Google API. If user exists -> Login; If new -> Register & Publish `user.registered` to Kafka.
- **Response**: JWT Tokens + User Profile.

### 3.4 `POST /api/v1/auth/oauth/line` (LINE Login - 100% Free)
- **Client Sends**: `{ "access_token": "line_access_token_string_from_frontend" }`
- **Service Action**: Verifies token with LINE API. Creates/Links user.
- **Response**: JWT Tokens + User Profile.
