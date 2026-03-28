package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"restaurant-backend/internal/models"

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
