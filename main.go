package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Product struct {
	ProductID  string  `json:"product_id"`
	Name       string  `json:"name"`
	Category   string  `json:"category"`
	StockLevel int     `json:"stock_level"`
	InStock    bool    `json:"in_stock"`
	PriceUSD   float64 `json:"price_usd"`
	Restricted bool    `json:"-"`
}

type Order struct {
	OrderID           string `json:"order_id"`
	CustomerName      string `json:"customer_name"`
	ProductID         string `json:"product_id"`
	Quantity          int    `json:"quantity"`
	Status            string `json:"status"`
	EstimatedDelivery string `json:"estimated_delivery,omitempty"`
	Restricted        bool   `json:"-"`
}

type APIError struct {
	Status           int    `json:"status"`
	ErrorCode        string `json:"error_code"`
	Message          string `json:"message"`
	AgentInstruction string `json:"agent_instruction"`
}

var (
	productIDRegex = regexp.MustCompile(`^PROD-\d{3}$`)
	orderIDRegex   = regexp.MustCompile(`^ORD-\d{4}$`)

	// 10 Sample Products
	products = map[string]Product{
		"PROD-001": {ProductID: "PROD-001", Name: "Keychron K2 Mechanical Keyboard", Category: "Peripherals", StockLevel: 24, InStock: true, PriceUSD: 99.00},
		"PROD-002": {ProductID: "PROD-002", Name: "Logitech MX Master 3S Mouse", Category: "Peripherals", StockLevel: 15, InStock: true, PriceUSD: 99.99},
		"PROD-003": {ProductID: "PROD-003", Name: "Dell UltraSharp 27 4K Monitor", Category: "Displays", StockLevel: 0, InStock: false, PriceUSD: 549.00},
		"PROD-004": {ProductID: "PROD-004", Name: "CalDigit TS4 Thunderbolt 4 Dock", Category: "Accessories", StockLevel: 12, InStock: true, PriceUSD: 399.95},
		"PROD-005": {ProductID: "PROD-005", Name: "Sony WH-1000XM5 Headphones", Category: "Audio", StockLevel: 8, InStock: true, PriceUSD: 348.00},
		"PROD-006": {ProductID: "PROD-006", Name: "Elgato Facecam Pro 4K60", Category: "Video", StockLevel: 0, InStock: false, PriceUSD: 299.99},
		"PROD-007": {ProductID: "PROD-007", Name: "Twelve South Curve Laptop Stand", Category: "Accessories", StockLevel: 45, InStock: true, PriceUSD: 59.99},
		"PROD-008": {ProductID: "PROD-008", Name: "Shure MV7+ Podcast Microphone", Category: "Audio", StockLevel: 11, InStock: true, PriceUSD: 279.00},
		"PROD-009": {ProductID: "PROD-009", Name: "Samsung T9 Portable SSD 2TB", Category: "Storage", StockLevel: 33, InStock: true, PriceUSD: 199.99},
		"PROD-010": {ProductID: "PROD-010", Name: "Internal Prototype DevKit v2", Category: "Hardware", StockLevel: 3, InStock: true, PriceUSD: 1200.00, Restricted: true},
	}

	// 10 Sample Orders
	orders = map[string]Order{
		"ORD-1001": {OrderID: "ORD-1001", CustomerName: "Alice Mercer", ProductID: "PROD-001", Quantity: 1, Status: "delivered", EstimatedDelivery: "2026-11-10"},
		"ORD-1002": {OrderID: "ORD-1002", CustomerName: "Bob Vance", ProductID: "PROD-004", Quantity: 2, Status: "shipped", EstimatedDelivery: "2026-11-18"},
		"ORD-1003": {OrderID: "ORD-1003", CustomerName: "Charlie Kim", ProductID: "PROD-005", Quantity: 1, Status: "processing", EstimatedDelivery: "2026-11-21"},
		"ORD-1004": {OrderID: "ORD-1004", CustomerName: "Dana Scully", ProductID: "PROD-002", Quantity: 1, Status: "cancelled"},
		"ORD-1005": {OrderID: "ORD-1005", CustomerName: "Evan Wright", ProductID: "PROD-007", Quantity: 3, Status: "shipped", EstimatedDelivery: "2026-11-19"},
		"ORD-1006": {OrderID: "ORD-1006", CustomerName: "Fiona Gallagher", ProductID: "PROD-008", Quantity: 1, Status: "delivered", EstimatedDelivery: "2026-11-12"},
		"ORD-1007": {OrderID: "ORD-1007", CustomerName: "George Miller", ProductID: "PROD-009", Quantity: 2, Status: "processing", EstimatedDelivery: "2026-11-22"},
		"ORD-1008": {OrderID: "ORD-1008", CustomerName: "Hannah Abbott", ProductID: "PROD-001", Quantity: 1, Status: "shipped", EstimatedDelivery: "2026-11-20"},
		"ORD-1009": {OrderID: "ORD-1009", CustomerName: "Ian Malcolm", ProductID: "PROD-004", Quantity: 4, Status: "payment_verification_pending"},
		"ORD-1010": {OrderID: "ORD-1010", CustomerName: "Enterprise Defense Corp", ProductID: "PROD-010", Quantity: 10, Status: "shipped", Restricted: true},
	}
)

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func withCORSAndAuth(expectedKey string, requireAuth bool, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, X-API-Key, Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if requireAuth {
			bearer := strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer"))
			apiKeyHeader := strings.TrimSpace(r.Header.Get("X-API-Key"))
			queryKey := strings.TrimSpace(r.URL.Query().Get("api_key"))

			token := bearer
			if token == "" {
				token = apiKeyHeader
			}
			if token == "" {
				token = queryKey
			}

			if token != expectedKey {
				writeJSON(w, http.StatusUnauthorized, APIError{
					Status:           http.StatusUnauthorized,
					ErrorCode:        "AUTH_REFUSED_INVALID_TOKEN",
					Message:          "Request refused: Missing or invalid API authorization token.",
					AgentInstruction: "The API refused this request because the authorization credentials are missing or invalid. Tell the user the store API refused access due to an authorization failure.",
				})
				return
			}
		}

		next(w, r)
	}
}

// Read Method 1: Check Product Stock Level
func handleGetStock(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		id = strings.TrimSpace(r.URL.Query().Get("product_id"))
	}

	if id == "" || !productIDRegex.MatchString(id) {
		writeJSON(w, http.StatusBadRequest, APIError{
			Status:           http.StatusBadRequest,
			ErrorCode:        "INVALID_PRODUCT_ID",
			Message:          "Parameter validation failed: '" + id + "' is not a valid product_id. Expected format PROD-XXX (3 digits, e.g., PROD-001).",
			AgentInstruction: "Do not retry with this ID. Inform the user that the product ID format is invalid and must follow PROD-XXX (for example, PROD-001).",
		})
		return
	}

	prod, exists := products[id]
	if !exists {
		writeJSON(w, http.StatusNotFound, APIError{
			Status:           http.StatusNotFound,
			ErrorCode:        "PRODUCT_NOT_FOUND",
			Message:          "Unknown product record: '" + id + "' does not exist in the catalog.",
			AgentInstruction: "Inform the user that product " + id + " was not found in the store inventory.",
		})
		return
	}

	if prod.Restricted {
		writeJSON(w, http.StatusForbidden, APIError{
			Status:           http.StatusForbidden,
			ErrorCode:        "REQUEST_REFUSED_RESTRICTED_PRODUCT",
			Message:          "Request refused by store policy: Product '" + id + "' is an internal prototype and stock levels cannot be queried via the public store agent.",
			AgentInstruction: "The API refused this request due to store security policy. Explain to the user that stock data for " + id + " is restricted and cannot be accessed.",
		})
		return
	}

	writeJSON(w, http.StatusOK, prod)
}

// Read Method 2: Check Order Status
func handleGetOrderStatus(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		id = strings.TrimSpace(r.URL.Query().Get("order_id"))
	}

	if id == "" || !orderIDRegex.MatchString(id) {
		writeJSON(w, http.StatusBadRequest, APIError{
			Status:           http.StatusBadRequest,
			ErrorCode:        "INVALID_ORDER_ID",
			Message:          "Parameter validation failed: '" + id + "' is not a valid order_id. Expected format ORD-XXXX (4 digits, e.g., ORD-1001).",
			AgentInstruction: "Do not retry with this ID. Inform the user that the order ID format is invalid and must follow ORD-XXXX (for example, ORD-1001).",
		})
		return
	}

	ord, exists := orders[id]
	if !exists {
		writeJSON(w, http.StatusNotFound, APIError{
			Status:           http.StatusNotFound,
			ErrorCode:        "ORDER_NOT_FOUND",
			Message:          "Unknown order record: '" + id + "' does not exist.",
			AgentInstruction: "Inform the user that order " + id + " does not exist in our records.",
		})
		return
	}

	if ord.Restricted {
		writeJSON(w, http.StatusForbidden, APIError{
			Status:           http.StatusForbidden,
			ErrorCode:        "REQUEST_REFUSED_RESTRICTED_ORDER",
			Message:          "Request refused by authorization policy: Order '" + id + "' belongs to a restricted Enterprise Defense account and requires elevated clearance.",
			AgentInstruction: "The API refused this request because order " + id + " is restricted under enterprise security policy. Tell the user the store API refused access and advise them to contact their enterprise account manager.",
		})
		return
	}

	writeJSON(w, http.StatusOK, ord)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "agentplace-store-api"})
}

func main() {
	apiKey := strings.TrimSpace(os.Getenv("API_KEY"))
	if apiKey == "" {
		log.Fatal("FATAL: API_KEY environment variable is required")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	// Health check (unauthenticated for cloud platform probes)
	mux.HandleFunc("GET /health", handleHealth)

	// Path-parameter and Query-parameter routes for maximum Agentplace Builder compatibility
	mux.HandleFunc("GET /api/v1/products/{id}/stock", withCORSAndAuth(apiKey, true, handleGetStock))
	mux.HandleFunc("GET /api/v1/stock", withCORSAndAuth(apiKey, true, handleGetStock))

	mux.HandleFunc("GET /api/v1/orders/{id}/status", withCORSAndAuth(apiKey, true, handleGetOrderStatus))
	mux.HandleFunc("GET /api/v1/order-status", withCORSAndAuth(apiKey, true, handleGetOrderStatus))

	log.Printf("Store API listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}