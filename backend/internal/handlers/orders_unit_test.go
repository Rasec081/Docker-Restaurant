package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repository"
)

func TestOrderHandlers_WithMock(t *testing.T) {
	mockRepo := repository.NewMockOrderRepository()
	InitOrderHandler(mockRepo)

	router := setupTestRouter()
	router.GET("/orders/:id", GetOrder)
	router.POST("/orders", CreateOrder)

	tableID := 1
	mockRepo.Orders[1] = &models.Order{
		ID:           1,
		TableID:      &tableID,
		ClientID:     2,
		OrdersType:   "Dine in",
		RestaurantID: 1,
	}

	t.Run("GetOrder ok", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/orders/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusOK, w.Code)
		}
	})

	t.Run("GetOrder not found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/orders/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("GetOrder invalid id", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/orders/abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("CreateOrder ok", func(t *testing.T) {
		order := models.Order{
			TableID:      &tableID,
			ClientID:     2,
			OrdersType:   "Dine in",
			RestaurantID: 1,
		}
		body, _ := json.Marshal(order)
		req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusOK, w.Code)
		}
	})

	t.Run("CreateOrder invalid json", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/orders", bytes.NewBufferString("{"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("CreateOrder repo error", func(t *testing.T) {
		mockRepo.CreateError = errTest
		defer func() { mockRepo.CreateError = nil }()

		order := models.Order{ClientID: 2, OrdersType: "Dine in", RestaurantID: 1}
		body, _ := json.Marshal(order)
		req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusInternalServerError, w.Code)
		}
	})
}
