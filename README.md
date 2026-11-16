# Agentplace Store API — Backend Engineer Test Task

A deterministic, agent-ready Go REST service built for the **Agentplace** store support assistant (powered by the **GPT-6 Luna** model). It provides real-time product stock checks and order status lookups across 10 sample products and 10 sample orders, backed by strict regex parameter validation, Bearer/API-key authentication, and machine-actionable error contracts.

## Live Deployment & Agent Links

- **Live Agentplace Store Agent (Luna Model):** [https://sl0fkwsc6v7w.run.agentplace.io/](https://sl0fkwsc6v7w.run.agentplace.io/)
  *(Pre-configured with server-side credentials via Agentplace Environment settings — no API key required to test in browser)*
- **Live Hosted API Base URL:** [https://agentplace-store-api.onrender.com](https://agentplace-store-api.onrender.com)
- **Live Health Check Endpoint:** [https://agentplace-store-api.onrender.com/health](https://agentplace-store-api.onrender.com/health)
- **GitHub Repository:** [https://github.com/abhiabhi0/agentplace-store-api](https://github.com/abhiabhi0/agentplace-store-api)

---

## Architecture & Agent-Reliability Design

Standard human-facing REST APIs often return generic error messages that cause LLM agents to hallucinate inventory or enter retry loops. This service is designed specifically for reliable agent tool-calling:

1. **Machine-Actionable Error Contracts (`agent_instruction`):** Every `4xx` response returns a structured JSON payload containing `status`, `error_code`, `message`, and `agent_instruction`. The `agent_instruction` field gives the Luna model deterministic guidance on whether to stop retrying, ask the user for a corrected ID format, or report an authorization/security refusal.
2. **Strict Parameter Validation (`400 Bad Request`):** Product IDs (`^PROD-\d{3}$`) and Order IDs (`^ORD-\d{4}$`) are validated before dataset lookup so malformed inputs are caught deterministically by the backend.
3. **Fail-Fast Secret Management:** The server requires `API_KEY` to be injected via environment variables at runtime (and `STORE_API_BEARER_TOKEN` in Agentplace's `Configure -> Environment` settings), ensuring zero credentials are hardcoded in source control or prompts.
4. **Dual-Layer Refusal Architecture (`401 Unauthorized` & `403 Forbidden`):**
   - **Transport Authorization (`401`):** Middleware enforces `Authorization: Bearer <API_KEY>` (or `X-API-Key: <API_KEY>`). Unauthenticated requests are refused with `AUTH_REFUSED_INVALID_TOKEN` prior to handler execution.
   - **Policy Authorization (`403`):** Because the published Agentplace agent holds a valid server-side Bearer token, records `PROD-010` (Internal Prototype DevKit) and `ORD-1010` (Enterprise Defense Corp) enforce an application-level security policy refusal (`403 Forbidden`). This allows reviewers to trigger and verify a live API refusal directly inside the Agentplace chat UI without needing to modify environment credentials.

---

## API Specification

### Authentication
- **Via the Published Agentplace Agent:** No key is required by the end user; the agent injects `Authorization: Bearer <STORE_API_BEARER_TOKEN>` from its protected runtime environment.
- **Direct API Calls (`/api/v1/*`):** Protected by middleware requiring `Authorization: Bearer <API_KEY>` (or `X-API-Key: <API_KEY>`). Calling the live Render endpoints without a header immediately demonstrates the `401 Unauthorized` refusal.

### 1. Check Product Stock Level
- **Endpoint:** `GET /api/v1/products/{id}/stock`
- **Alternative Query Route:** `GET /api/v1/stock?product_id={id}`
- **Success Response (`200 OK`):**
```json
{
  "product_id": "PROD-001",
  "name": "Keychron K2 Mechanical Keyboard",
  "category": "Peripherals",
  "stock_level": 24,
  "in_stock": true,
  "price_usd": 99
}
```

### 2. Check Order Status
- **Endpoint:** `GET /api/v1/orders/{id}/status`
- **Alternative Query Route:** `GET /api/v1/order-status?order_id={id}`
- **Success Response (`200 OK`):**
```json
{
  "order_id": "ORD-1002",
  "customer_name": "Bob Vance",
  "product_id": "PROD-004",
  "quantity": 2,
  "status": "shipped",
  "estimated_delivery": "2026-11-18"
}
```

---

## Seeded Dataset (10 Products & 10 Orders)

### 10 Sample Products
| Product ID | Name | Category | Stock Level | In Stock | Price (USD) | Access Policy |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| `PROD-001` | Keychron K2 Mechanical Keyboard | Peripherals | 24 | `true` | $99.00 | Standard (`200 OK`) |
| `PROD-002` | Logitech MX Master 3S Mouse | Peripherals | 15 | `true` | $99.99 | Standard (`200 OK`) |
| `PROD-003` | Dell UltraSharp 27 4K Monitor | Displays | 0 | `false` | $549.00 | Out of Stock (`200 OK`) |
| `PROD-004` | CalDigit TS4 Thunderbolt 4 Dock | Accessories | 12 | `true` | $399.95 | Standard (`200 OK`) |
| `PROD-005` | Sony WH-1000XM5 Headphones | Audio | 8 | `true` | $348.00 | Standard (`200 OK`) |
| `PROD-006` | Elgato Facecam Pro 4K60 | Video | 0 | `false` | $299.99 | Out of Stock (`200 OK`) |
| `PROD-007` | Twelve South Curve Laptop Stand | Accessories | 45 | `true` | $59.99 | Standard (`200 OK`) |
| `PROD-008` | Shure MV7+ Podcast Microphone | Audio | 11 | `true` | $279.00 | Standard (`200 OK`) |
| `PROD-009` | Samsung T9 Portable SSD 2TB | Storage | 33 | `true` | $199.99 | Standard (`200 OK`) |
| `PROD-010` | Internal Prototype DevKit v2 | Hardware | 3 | `true` | $1200.00 | Restricted (`403 Forbidden`) |

### 10 Sample Orders
| Order ID | Customer Name | Product ID | Qty | Status | Est. Delivery | Access Policy |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| `ORD-1001` | Alice Mercer | `PROD-001` | 1 | `delivered` | `2026-11-10` | Standard (`200 OK`) |
| `ORD-1002` | Bob Vance | `PROD-004` | 2 | `shipped` | `2026-11-18` | Standard (`200 OK`) |
| `ORD-1003` | Charlie Kim | `PROD-005` | 1 | `processing` | `2026-11-21` | Standard (`200 OK`) |
| `ORD-1004` | Dana Scully | `PROD-002` | 1 | `cancelled` | — | Standard (`200 OK`) |
| `ORD-1005` | Evan Wright | `PROD-007` | 3 | `shipped` | `2026-11-19` | Standard (`200 OK`) |
| `ORD-1006` | Fiona Gallagher | `PROD-008` | 1 | `delivered` | `2026-11-12` | Standard (`200 OK`) |
| `ORD-1007` | George Miller | `PROD-009` | 2 | `processing` | `2026-11-22` | Standard (`200 OK`) |
| `ORD-1008` | Hannah Abbott | `PROD-001` | 1 | `shipped` | `2026-11-20` | Standard (`200 OK`) |
| `ORD-1009` | Ian Malcolm | `PROD-004` | 4 | `payment_verification_pending` | — | Standard (`200 OK`) |
| `ORD-1010` | Enterprise Defense Corp | `PROD-010` | 10 | `shipped` | — | Restricted (`403 Forbidden`) |

---

## Local Setup & Testing Instructions

### Prerequisites
- Go 1.24+ (or Docker)

### 1. Clone and Run Locally
```bash
git clone https://github.com/abhiabhi0/agentplace-store-api.git
cd agentplace-store-api

export API_KEY="local-dev-test-key"
export PORT="8080"
go run main.go
```

### 2. Build and Run with Docker
```bash
docker build -t agentplace-store-api .
docker run -p 8080:8080 -e API_KEY="local-dev-test-key" agentplace-store-api
```

### 3. Verify Endpoints Locally with cURL
```bash
export API_KEY="local-dev-test-key"

# Health check (200 OK)
curl -i http://localhost:8080/health

# 1. Existing records (200 OK)
curl -i -H "Authorization: Bearer $API_KEY" http://localhost:8080/api/v1/products/PROD-001/stock
curl -i -H "Authorization: Bearer $API_KEY" http://localhost:8080/api/v1/orders/ORD-1002/status

# 2. Unknown records (404 Not Found)
curl -i -H "Authorization: Bearer $API_KEY" http://localhost:8080/api/v1/products/PROD-999/stock
curl -i -H "Authorization: Bearer $API_KEY" http://localhost:8080/api/v1/orders/ORD-9999/status

# 3. Invalid ID formats (400 Bad Request)
curl -i -H "Authorization: Bearer $API_KEY" http://localhost:8080/api/v1/products/kb-12/stock
curl -i -H "Authorization: Bearer $API_KEY" http://localhost:8080/api/v1/orders/12345/status

# 4A. Refused request — Missing or invalid token (401 Unauthorized)
curl -i http://localhost:8080/api/v1/orders/ORD-1001/status

# 4B. Refused request — Restricted policy records (403 Forbidden)
curl -i -H "Authorization: Bearer $API_KEY" http://localhost:8080/api/v1/orders/ORD-1010/status
curl -i -H "Authorization: Bearer $API_KEY" http://localhost:8080/api/v1/products/PROD-010/stock
```

---

## Agentplace Configuration

1. Built the **Store Support Agent** on Agentplace using the free **GPT-6 Luna** model.
2. Registered two GET tools (`get_product_stock` and `get_order_status`) connected to `[https://agentplace-store-api.onrender.com/api/v1/products/](https://agentplace-store-api.onrender.com/api/v1/products/){id}/stock` and `[https://agentplace-store-api.onrender.com/api/v1/orders/](https://agentplace-store-api.onrender.com/api/v1/orders/){id}/status`.
3. Stored the runtime credential in **Configure → Environment** under `STORE_API_BEARER_TOKEN` so no secrets exist in the repository or prompt history, and the published agent works out-of-the-box for any visitor.
4. Configured the tools to pass the user's raw submitted ID string unchanged so all parameter validation and refusal checks are executed deterministically by the Go backend.

---

## Recorded Test Results & Agent Behavior

All four required edge cases were executed against the live Render API (`[https://agentplace-store-api.onrender.com](https://agentplace-store-api.onrender.com)`) through the Agentplace Store Support Agent.

### Case 1: Existing Records (`200 OK`)

#### Test 1A — In-Stock Product & Shipped Order
- **User Prompt:** `"Can you check how many PROD-001 keyboards are in stock and what the status of order ORD-1002 is?"`
- **API Status:** `200 OK` on both endpoints.
- **Recorded Agent Response:**
  - **Current lookup results**
  - **PROD-001 (Keychron K2 Mechanical Keyboard):** 24 in stock.
  - **Order ORD-1002:** shipped.

#### Test 1B — Out-of-Stock Product & Delivered Order
- **User Prompt:** `"Check the stock level for PROD-003 and the order status for ORD-1001."`
- **API Status:** `200 OK` on both endpoints.
- **Recorded Agent Response:**
  - **Current lookup results**
  - **PROD-003 (Dell UltraSharp 27 4K Monitor):** 0 in stock; not in stock.
  - **Order ORD-1001:** delivered. The response lists estimated delivery as November 10, 2026.

---

### Case 2: Unknown Records (`404 Not Found`)
- **What was tested:** Validly formatted product and order IDs (`PROD-999` and `ORD-9999`) that do not exist in the store dataset.
- **User Prompt:** `"Check the stock level for product PROD-999 and the status of order ORD-9999."`
- **API Status & Payload (`404 Not Found`):**
```json
{
  "status": 404,
  "error_code": "PRODUCT_NOT_FOUND",
  "message": "Unknown product record: 'PROD-999' does not exist in the catalog.",
  "agent_instruction": "Inform the user that product PROD-999 was not found in the store inventory."
}
```
- **Recorded Agent Response:**
  - **Lookup results**
  - **Product PROD-999:** Not found in the store inventory (`PRODUCT_NOT_FOUND`). The backend says this product does not exist in the catalog.
  - **Order ORD-9999:** Does not exist in the store's records (`ORDER_NOT_FOUND`).

---

### Case 3: Invalid ID Formats (`400 Bad Request`)
- **What was tested:** Malformed IDs (`kb-12` and `12345`) that fail regex validation (`^PROD-\d{3}$` and `^ORD-\d{4}$`).
- **User Prompt:** `"What is the stock for product kb-12 and the status for order 12345?"`
- **API Status & Payload (`400 Bad Request`):**
```json
{
  "status": 400,
  "error_code": "INVALID_PRODUCT_ID",
  "message": "Parameter validation failed: 'kb-12' is not a valid product_id. Expected format PROD-XXX (3 digits, e.g., PROD-001).",
  "agent_instruction": "Do not retry with this ID. Inform the user that the product ID format is invalid and must follow PROD-XXX (for example, PROD-001)."
}
```
- **Recorded Agent Response:**
  - **Lookup results**
  - **Product `kb-12`:** The backend rejected the ID format (`INVALID_PRODUCT_ID`). It says product IDs must follow `PROD-XXX` (for example, `PROD-001`).
  - **Order `12345`:** The backend rejected the ID format (`INVALID_ORDER_ID`). It says order IDs must follow `ORD-XXXX` (for example, `ORD-1001`).

---

### Case 4: Requests the API Refuses (`403 Forbidden` & `401 Unauthorized`)

#### Part A — In-Agent Policy Refusal (`403 Forbidden`)
- **What was tested:** Querying restricted enterprise order `ORD-1010` and restricted internal prototype `PROD-010` through the live agent.
- **User Prompt:** `"Check the status of order ORD-1010 and the stock level for PROD-010."`
- **API Status & Payload (`403 Forbidden`):**
```json
{
  "status": 403,
  "error_code": "REQUEST_REFUSED_RESTRICTED_ORDER",
  "message": "Request refused by authorization policy: Order 'ORD-1010' belongs to a restricted Enterprise Defense account and requires elevated clearance.",
  "agent_instruction": "The API refused this request because order ORD-1010 is restricted under enterprise security policy. Tell the user the store API refused access and advise them to contact their enterprise account manager."
}
```
- **Recorded Agent Response:**
  - **Lookup results**
  - **Order ORD-1010:** The store API refused access because it is restricted under enterprise security policy (`REQUEST_REFUSED_RESTRICTED_ORDER`). Please contact your enterprise account manager.
  - **Product PROD-010:** Stock data is restricted and cannot be accessed through this lookup (`REQUEST_REFUSED_RESTRICTED_PRODUCT`).

#### Part B — Direct Unauthenticated Request Refusal (`401 Unauthorized`)
- **What was tested:** Calling `GET [https://agentplace-store-api.onrender.com/api/v1/orders/ORD-1001/status](https://agentplace-store-api.onrender.com/api/v1/orders/ORD-1001/status)` directly without an `Authorization: Bearer` header.
- **API Status & Payload (`401 Unauthorized`):**
```json
{
  "status": 401,
  "error_code": "AUTH_REFUSED_INVALID_TOKEN",
  "message": "Request refused: Missing or invalid API authorization token.",
  "agent_instruction": "The API refused this request because the authorization credentials are missing or invalid. Tell the user the store API refused access due to an authorization failure."
}
```