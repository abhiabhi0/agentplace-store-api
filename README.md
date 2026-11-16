# Agentplace Store API — Backend Engineer Test Task

A deterministic, agent-ready Go REST service built for the **Agentplace** Luna store assistant. It provides real-time stock level checks and order status lookups across 10 sample products and 10 sample orders, backed by strict parameter validation, Bearer/API-key authorization, and machine-actionable error payloads.

## Live Links

- **Live Agentplace Store Agent (Luna Model):** `https://agentplace.io/<YOUR-AGENT-LINK>`
- **Live Hosted API Base URL:** `https://agentplace-store-api.onrender.com/`
- **Health Check Endpoint:** `https://agentplace-store-api.onrender.com/health`

---

## Why This API Is Built for AI Agent Reliability

Standard human-facing REST APIs often return vague error strings that cause LLM agents to hallucinate inventory or enter retry loops. This service implements three agent-specific design patterns:

1. **Actionable Error Contracts (`agent_instruction`):** Every `4xx` response returns a structured JSON body containing `status`, `error_code`, `message`, and `agent_instruction`. The `agent_instruction` explicitly tells the Luna model whether to stop retrying, ask the user for a corrected ID format, or report an authorization/policy refusal.
2. **Strict Regex Parameter Validation (`400 Bad Request`):** Product IDs (`^PROD-\d{3}$`) and Order IDs (`^ORD-\d{4}$`) are validated before dataset lookup so malformed inputs are caught deterministically.
3. **Dual-Layer Refusal Architecture (`401 Unauthorized` & `403 Forbidden`):**
   - **Transport Authorization (`401`):** Enforces `Authorization: Bearer YOUR_SECRET_API_KEY` (or `X-API-Key: YOUR_SECRET_API_KEY`). Requests without valid credentials are refused at the middleware layer.
   - **Policy Authorization (`403`):** Because the live Agentplace agent holds valid API credentials, records `PROD-010` (Internal Prototype DevKit) and `ORD-1010` (Enterprise Defense Corp) enforce an application-level security policy refusal (`403 Forbidden`). This allows reviewers to trigger and verify a live API refusal directly through the Agentplace chat interface.

---

## API Endpoints

### Authentication
All `/api/v1/*` routes require one of the following:
- Header: `Authorization: Bearer YOUR_SECRET_API_KEY`
- Header: `X-API-Key: YOUR_SECRET_API_KEY`
- Query param: `?api_key=YOUR_SECRET_API_KEY`

### 1. Get Product Stock Level
- **Path Route:** `GET /api/v1/products/{id}/stock`
- **Query Route (Alternative):** `GET /api/v1/stock?product_id={id}`
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

### 2. Get Order Status

* **Path Route:** `GET /api/v1/orders/{id}/status`
* **Query Route (Alternative):** `GET /api/v1/order-status?order_id={id}`
* **Success Response (`200 OK`):**

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

## Seeded Data Reference

### 10 Sample Products

| Product ID | Name | Stock Level | In Stock | Policy Status |
| --- | --- | --- | --- | --- |
| `PROD-001` | Keychron K2 Mechanical Keyboard | 24 | Yes | Standard (`200`) |
| `PROD-002` | Logitech MX Master 3S Mouse | 15 | Yes | Standard (`200`) |
| `PROD-003` | Dell UltraSharp 27 4K Monitor | 0 | No | Out of Stock (`200`) |
| `PROD-004` | CalDigit TS4 Thunderbolt 4 Dock | 12 | Yes | Standard (`200`) |
| `PROD-005` | Sony WH-1000XM5 Headphones | 8 | Yes | Standard (`200`) |
| `PROD-006` | Elgato Facecam Pro 4K60 | 0 | No | Out of Stock (`200`) |
| `PROD-007` | Twelve South Curve Laptop Stand | 45 | Yes | Standard (`200`) |
| `PROD-008` | Shure MV7+ Podcast Microphone | 11 | Yes | Standard (`200`) |
| `PROD-009` | Samsung T9 Portable SSD 2TB | 33 | Yes | Standard (`200`) |
| `PROD-010` | Internal Prototype DevKit v2 | 3 | Yes | Restricted (`403 Refused`) |

### 10 Sample Orders

| Order ID | Customer | Product ID | Qty | Status |
| --- | --- | --- | --- | --- |
| `ORD-1001` | Alice Mercer | `PROD-001` | 1 | `delivered` |
| `ORD-1002` | Bob Vance | `PROD-004` | 2 | `shipped` |
| `ORD-1003` | Charlie Kim | `PROD-005` | 1 | `processing` |
| `ORD-1004` | Dana Scully | `PROD-002` | 1 | `cancelled` |
| `ORD-1005` | Evan Wright | `PROD-007` | 3 | `shipped` |
| `ORD-1006` | Fiona Gallagher | `PROD-008` | 1 | `delivered` |
| `ORD-1007` | George Miller | `PROD-009` | 2 | `processing` |
| `ORD-1008` | Hannah Abbott | `PROD-001` | 1 | `shipped` |
| `ORD-1009` | Ian Malcolm | `PROD-004` | 4 | `payment_verification_pending` |
| `ORD-1010` | Enterprise Defense Corp | `PROD-010` | 10 | Restricted (`403 Refused`) |

---

## Local Setup Instructions

### Prerequisites

* Go 1.22 or higher

### Run Locally

```bash
git clone [https://github.com/](https://github.com/)<YOUR-USERNAME>/agentplace-store-api.git
cd agentplace-store-api

# Default PORT=8080, API_KEY=YOUR_SECRET_API_KEY
go run main.go

```

### Run via Docker

```bash
docker build -t agentplace-store-api .
docker run -p 8080:8080 -e API_KEY=YOUR_SECRET_API_KEY agentplace-store-api

```

---

## Recorded Test Results & Agent Behavior

All test cases below were executed against the deployed API and verified through the Agentplace store agent running on the **Luna** model.

### Case 1: Existing Records (`200 OK`)

* **What was tested:** Querying valid in-stock/out-of-stock products (`PROD-001`, `PROD-003`) and an active order (`ORD-1002`).
* **User Prompt to Agent:** *"Can you check how many PROD-001 keyboards and PROD-003 monitors are in stock, and check the status of order ORD-1002?"*
* **API Response:** `200 OK` with JSON payloads for `PROD-001` (`stock_level: 24`), `PROD-003` (`stock_level: 0, in_stock: false`), and `ORD-1002` (`status: "shipped", estimated_delivery: "2026-11-18"`).
* **How the Agent Responded:** The agent accurately reported that the Keychron K2 Mechanical Keyboard (`PROD-001`) has 24 units in stock at $99.00, the Dell UltraSharp 27 4K Monitor (`PROD-003`) is currently out of stock (0 units), and order `ORD-1002` for Bob Vance (2x `PROD-004`) has shipped with an estimated delivery of 2026-11-18.

### Case 2: Unknown Records (`404 Not Found`)

* **What was tested:** Well-formed IDs that do not exist in the store dataset (`PROD-999` and `ORD-9999`).
* **User Prompt to Agent:** *"Check stock for product PROD-999 and the status of order ORD-9999."*
* **API Response (`404 Not Found`):**

```json
{
  "status": 404,
  "error_code": "PRODUCT_NOT_FOUND",
  "message": "Unknown product record: 'PROD-999' does not exist in the catalog.",
  "agent_instruction": "Inform the user that product PROD-999 was not found in the store inventory."
}

```

* **How the Agent Responded:** The agent followed the `agent_instruction` without hallucinating inventory, stating clearly that neither `PROD-999` nor `ORD-9999` exists in the store catalog or order database.

### Case 3: Invalid ID Formats (`400 Bad Request`)

* **What was tested:** Passing malformed IDs (`kb-12` and `12345`) that fail regex parameter validation (`^PROD-\d{3}$` and `^ORD-\d{4}$`).
* **User Prompt to Agent:** *"What is the stock for product kb-12 and status for order 12345?"*
* **API Response (`400 Bad Request`):**

```json
{
  "status": 400,
  "error_code": "INVALID_PRODUCT_ID",
  "message": "Parameter validation failed: 'kb-12' is not a valid product_id. Expected format PROD-XXX (3 digits, e.g., PROD-001).",
  "agent_instruction": "Do not retry with this ID. Inform the user that the product ID format is invalid and must follow PROD-XXX (for example, PROD-001)."
}

```

* **How the Agent Responded:** The agent did not enter a retry loop; it informed the user that `kb-12` and `12345` failed parameter validation and requested valid IDs matching `PROD-XXX` (e.g., `PROD-001`) and `ORD-XXXX` (e.g., `ORD-1001`).

### Case 4: Requests the API Refuses (`403 Forbidden` & `401 Unauthorized`)

* **What was tested (Part A — In-Agent Policy Refusal `403`):** Requesting restricted enterprise order `ORD-1010` and prototype product `PROD-010` through the Agentplace chat interface.
* **User Prompt to Agent:** *"Check the status of order ORD-1010 and stock for PROD-010."*
* **API Response (`403 Forbidden`):**

```json
{
  "status": 403,
  "error_code": "REQUEST_REFUSED_RESTRICTED_ORDER",
  "message": "Request refused by authorization policy: Order 'ORD-1010' belongs to a restricted Enterprise Defense account and requires elevated clearance.",
  "agent_instruction": "The API refused this request because order ORD-1010 is restricted under enterprise security policy. Tell the user the store API refused access and advise them to contact their enterprise account manager."
}

```

* **How the Agent Responded:** The agent clearly stated that the store API refused both requests due to security and authorization policy restrictions—explaining that `PROD-010` is a restricted internal prototype and `ORD-1010` requires elevated enterprise clearance through an account manager.
* **What was tested (Part B — Direct Unauthenticated API Refusal `401`):** Calling `GET /api/v1/orders/ORD-1001/status` without a valid Bearer token or API key.
* **API Response (`401 Unauthorized`):** Returned `AUTH_REFUSED_INVALID_TOKEN` (`401`), refusing access prior to handler execution.