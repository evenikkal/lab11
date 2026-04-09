package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetOrders() {
	orders = nil
	nextID = 1
}

func sampleOrder() Order {
	return Order{
		CustomerID: 42,
		Items: []OrderItem{
			{ProductID: 1, ProductName: "Laptop", Quantity: 1, UnitPrice: 1200.00},
			{ProductID: 2, ProductName: "Mouse", Quantity: 2, UnitPrice: 25.50},
		},
		ShipTo: Address{
			Street:  "Pushkina 10",
			City:    "Moscow",
			Country: "Russia",
			Zip:     "101000",
		},
	}
}

func TestHealthEndpoint(t *testing.T) {
	router := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateOrderValid(t *testing.T) {
	resetOrders()
	router := SetupRouter()
	w := httptest.NewRecorder()

	body, _ := json.Marshal(sampleOrder())
	req, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var created Order
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &created))
	assert.Equal(t, 1, created.ID)
	assert.Equal(t, "pending", created.Status)
	assert.InDelta(t, 1251.0, created.TotalAmount, 0.01)
}

func TestCreateOrderComputesTotal(t *testing.T) {
	resetOrders()
	router := SetupRouter()
	w := httptest.NewRecorder()

	order := sampleOrder()
	body, _ := json.Marshal(order)
	req, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var created Order
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &created))
	assert.InDelta(t, 1251.0, created.TotalAmount, 0.001)
}

func TestCreateOrderMissingItems(t *testing.T) {
	resetOrders()
	router := SetupRouter()
	w := httptest.NewRecorder()

	body, _ := json.Marshal(map[string]interface{}{
		"customer_id": 42,
		"ship_to":     sampleOrder().ShipTo,
	})
	req, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateOrderMissingAddress(t *testing.T) {
	resetOrders()
	router := SetupRouter()
	w := httptest.NewRecorder()

	body, _ := json.Marshal(map[string]interface{}{
		"customer_id": 42,
		"items":       sampleOrder().Items,
	})
	req, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetOrder(t *testing.T) {
	resetOrders()
	router := SetupRouter()

	body, _ := json.Marshal(sampleOrder())
	wCreate := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(wCreate, req)
	require.Equal(t, http.StatusCreated, wCreate.Code)

	wGet := httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/orders/1", nil)
	router.ServeHTTP(wGet, req)

	require.Equal(t, http.StatusOK, wGet.Code)
	var fetched Order
	require.NoError(t, json.Unmarshal(wGet.Body.Bytes(), &fetched))
	assert.Equal(t, 42, fetched.CustomerID)
}

func TestGetOrderNotFound(t *testing.T) {
	resetOrders()
	router := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/orders/999", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestListOrders(t *testing.T) {
	resetOrders()
	router := SetupRouter()

	body, _ := json.Marshal(sampleOrder())
	wCreate := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(wCreate, req)

	wList := httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/orders", nil)
	router.ServeHTTP(wList, req)

	require.Equal(t, http.StatusOK, wList.Code)
	var list []Order
	require.NoError(t, json.Unmarshal(wList.Body.Bytes(), &list))
	assert.Len(t, list, 1)
}
