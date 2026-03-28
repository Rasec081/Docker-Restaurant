package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

func TestGetOrder(t *testing.T) {
	setupTestData(t)

	// Crear orden de prueba
	var orderID int
	testDB.QueryRow(`
		INSERT INTO Orders (table_id, client_id, orders_type, restaurant_id) 
		VALUES (1, 2, 'Dine in', 1) 
		RETURNING orders_id
	`).Scan(&orderID)

	tests := []struct {
		name     string
		orderID  string
		expected int
	}{
		{
			name:     "Get existing order",
			orderID:  strconv.Itoa(orderID),
			expected: 200,
		},
		{
			name:     "Get non-existent order",
			orderID:  "999",
			expected: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()
			router.GET("/orders/:id", GetOrder)

			req, _ := http.NewRequest("GET", "/orders/"+tt.orderID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expected {
				t.Errorf("Esperaba %d, obtuvo %d", tt.expected, w.Code)
			}
		})
	}
}

func TestCreateOrder(t *testing.T) {
	setupTestData(t)

	tests := []struct {
		name     string
		order    models.Order
		expected int
	}{
		{
			name: "Create valid order with table",
			order: models.Order{
				TableID:      &[]int{1}[0],
				ClientID:     2,
				OrdersType:   "Dine in",
				RestaurantID: 1,
			},
			expected: 200,
		},
		{
			name: "Create valid order without table",
			order: models.Order{
				TableID:      nil,
				ClientID:     2,
				OrdersType:   "Take away",
				RestaurantID: 1,
			},
			expected: 200,
		},
		{
			name: "Create order with invalid client",
			order: models.Order{
				TableID:      &[]int{1}[0],
				ClientID:     999,
				OrdersType:   "Dine in",
				RestaurantID: 1,
			},
			expected: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()
			router.POST("/orders", CreateOrder)

			jsonBody, _ := json.Marshal(tt.order)
			req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expected {
				t.Errorf("Esperaba %d, obtuvo %d", tt.expected, w.Code)
			}
		})
	}
}

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
