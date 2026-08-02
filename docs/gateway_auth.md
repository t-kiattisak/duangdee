# API Gateway & Authentication Routing Specification (Kong Gateway)

## 1. Overview & Gateway Responsibilities

The system utilizes **Kong API Gateway** as the single entry point for all client requests. The API Gateway is responsible for:
1. **SSL/TLS Termination**: Handling HTTPS encryption at the edge.
2. **Global Rate-Limiting**: Protecting downstream services from DDoS and abuse.
3. **Routing**: Forwarding requests to upstream microservices based on URI paths.
4. **Centralized Authentication & JWT Validation**: Offloading JWT validation from downstream microservices using Kong's JWT Plugin or OAuth2 Introspection.
5. **Token Transformation & Header Forwarding**: Injecting authenticated user context (`X-User-ID`, `X-User-Role`, `X-User-Email`) into request headers forwarded to backend services.

---

## 2. Authentication Strategy & Gateway Plugins (Kong)

Kong uses **JWKS (JSON Web Key Set)** or Shared Public Keys to validate JWT Access Tokens locally at the Gateway level without calling the Auth Service for every request (0 latency penalty).

```
Client (Web/App)
    |
    | 1. HTTP Request with Header: "Authorization: Bearer <JWT_ACCESS_TOKEN>"
    v
+-----------------------------------------------------------------------------------+
| Kong API Gateway                                                                  |
|                                                                                   |
|  [Step A]: Check Route Auth Policy (Is path Public or Protected?)                 |
|  [Step B]: Check Redis Token Revocation List (Is JWT Blacklisted/Logged Out?)     |
|  [Step C]: Validate JWT Signature & Expiration using JWKS / Public Key            |
|  [Step D]: Strip Bearer Token & Transform Headers:                                 |
|            - Add Header: `X-User-ID: usr_99887766-5544...`                        |
|            - Add Header: `X-User-Role: user`                                      |
|            - Add Header: `X-User-Email: user@example.com`                         |
+-----------------------------------------------------------------------------------+
                                    |
                                    | 2. Forward Request with X-User-* Headers
                                    v
            +-----------------------+-----------------------+
            |                                               |
            v                                               v
  +-------------------+                           +-------------------+
  |  Reading Service  |                           |  Payment Service  |
  | (Reads X-User-ID) |                           | (Reads X-User-ID) |
  +-------------------+                           +-------------------+
```

---

## 3. Route Access Matrix (Public vs Protected Routes)

| Path Pattern | HTTP Method | Auth Required? | Kong Action / Plugin | Forwarded Upstream Service |
| :--- | :--- | :--- | :--- | :--- |
| `/api/v1/auth/register` | `POST` | ❌ **Public** | Pass-through + Rate Limit | `Auth Service` |
| `/api/v1/auth/login` | `POST` | ❌ **Public** | Pass-through + Rate Limit | `Auth Service` |
| `/api/v1/auth/refresh` | `POST` | ❌ **Public** | Pass-through | `Auth Service` |
| `/api/v1/tarot/cards` | `GET` | ❌ **Public** | Pass-through + L2 Gateway Cache | `Tarot Core Service` |
| `/api/v1/tarot/spreads` | `GET` | ❌ **Public** | Pass-through + L2 Gateway Cache | `Tarot Core Service` |
| `/api/v1/users/me` | `GET`, `PUT` | ✅ **Protected**| Verify JWT -> Inject `X-User-ID` | `Auth Service` |
| `/api/v1/readings/initiate` | `POST` | ✅ **Protected**| Verify JWT -> Inject `X-User-ID` | `Reading Engine Service` |
| `/api/v1/readings/:id/draw` | `POST` | ✅ **Protected**| Verify JWT -> Inject `X-User-ID` | `Reading Engine Service` |
| `/api/v1/payments/balance` | `GET` | ✅ **Protected**| Verify JWT -> Inject `X-User-ID` | `Payment & Credit Service` |
| `/api/v1/payments/checkout` | `POST` | ✅ **Protected**| Verify JWT -> Inject `X-User-ID` | `Payment & Credit Service` |
| `/api/v1/payments/webhook/*` | `POST` | ❌ **Public** | Verify HMAC Gateway Signature | `Payment & Credit Service` |

---

## 4. Kong Header Transformation & Forwarding Workflow

### 4.1 What the Client Sends to Kong Gateway:
```http
POST /api/v1/readings/initiate HTTP/1.1
Host: api.duangdee.com
Authorization: Bearer eyJhbGciOiJSUzI1NiR5... (JWT Token)
Content-Type: application/json

{
  "spread_id": "three-card",
  "category": "love"
}
```

### 4.2 What Kong Does Internal Processing:
1. Matches route `/api/v1/readings/*` -> Identifies as **Protected**.
2. Extracts Bearer token `eyJhbGciOiJSUzI1Ni...`.
3. Verifies signature using Auth Service's Public RSA Key via JWKS endpoint (`http://auth-service:8080/.well-known/jwks.json`).
4. Checks Redis for Token Blacklist (Logout check).
5. Decodes JWT claims: `sub` -> `usr_99887766-5544-3322-1100-aabbccdd8899`, `role` -> `user`, `email` -> `user@example.com`.

### 4.3 What Kong Forwards to Downstream Service (e.g. Reading Service):
Kong strips the heavy JWT string and injects clean, pre-validated headers:

```http
POST /api/v1/readings/initiate HTTP/1.1
Host: reading-service.internal:8080
X-User-ID: usr_99887766-5544-3322-1100-aabbccdd8899
X-User-Role: user
X-User-Email: user@example.com
X-Trace-ID: trace_abc123xyz789
Content-Type: application/json

{
  "spread_id": "three-card",
  "category": "love"
}
```

---

## 5. Benefits of this Gateway Pattern
1. **Zero Auth Overhead on Microservices**: Downstream services (Reading, Payment, Notification) do **NOT** need to validate JWTs or talk to Auth Service. They simply read `c.GetHeader("X-User-ID")` in Go Fiber.
2. **High Performance**: Kong validates JWT signatures locally in memory via RSA Public Key, resulting in sub-millisecond overhead.
3. **Centralized Token Revocation**: Logout revocation is checked at Gateway level via Redis before hitting any backend microservices.
