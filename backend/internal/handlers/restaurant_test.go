package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"restaurant-backend/internal/models"

	"github.com/gin-gonic/gin"
)

func TestGetRestaurants(t *testing.T) {
	setupTestData(t)

	router := gin.Default()
	router.GET("/restaurants", GetRestaurants)

	req, _ := http.NewRequest("GET", "/restaurants", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Esperaba 200, obtuvo %d", w.Code)
	}

	var restaurants []models.Restaurant
	json.Unmarshal(w.Body.Bytes(), &restaurants)

	if len(restaurants) == 0 {
		t.Error("Debería haber al menos un restaurante")
	}
}

func TestCreateRestaurant(t *testing.T) {
	setupTestData(t)

	tests := []struct {
		name       string
		restaurant models.Restaurant
		expected   int
	}{
		{
			name: "Create valid restaurant",
			restaurant: models.Restaurant{
				Nombre:  "New Restaurant",
				Estado:  1,
				AdminID: 1,
			},
			expected: 200,
		},
		{
			name: "Create restaurant with missing name",
			restaurant: models.Restaurant{
				Nombre:  "",
				Estado:  1,
				AdminID: 1,
			},
			expected: 400,
		},
		{
			name: "Create restaurant with invalid admin",
			restaurant: models.Restaurant{
				Nombre:  "Invalid Admin",
				Estado:  1,
				AdminID: 999,
			},
			expected: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()
			router.POST("/restaurants", CreateRestaurant)

			jsonBody, _ := json.Marshal(tt.restaurant)
			req, _ := http.NewRequest("POST", "/restaurants", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expected {
				t.Errorf("Esperaba %d, obtuvo %d", tt.expected, w.Code)
			}
		})
	}
}
