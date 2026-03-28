package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repository"

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

func TestRestaurantHandlers_WithMock(t *testing.T) {
	mockRepo := repository.NewMockRestaurantRepository()
	InitRestaurantHandler(mockRepo)

	router := setupTestRouter()
	router.GET("/restaurants", GetRestaurants)
	router.POST("/restaurants", CreateRestaurant)

	t.Run("GetRestaurants ok", func(t *testing.T) {
		mockRepo.Restaurants[1] = &models.Restaurant{ID: 1, Nombre: "Test", Estado: 1, AdminID: 1}

		req, _ := http.NewRequest("GET", "/restaurants", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusOK, w.Code)
		}
	})

	t.Run("GetRestaurants repo error", func(t *testing.T) {
		mockRepo.GetAllError = errTest
		defer func() { mockRepo.GetAllError = nil }()

		req, _ := http.NewRequest("GET", "/restaurants", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("CreateRestaurant ok", func(t *testing.T) {
		restaurant := models.Restaurant{Nombre: "New Restaurant", Estado: 1, AdminID: 1}
		body, _ := json.Marshal(restaurant)
		req, _ := http.NewRequest("POST", "/restaurants", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusOK, w.Code)
		}
	})

	t.Run("CreateRestaurant invalid json", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/restaurants", bytes.NewBufferString("{"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("CreateRestaurant repo error", func(t *testing.T) {
		mockRepo.CreateError = errTest
		defer func() { mockRepo.CreateError = nil }()

		restaurant := models.Restaurant{Nombre: "New Restaurant", Estado: 1, AdminID: 1}
		body, _ := json.Marshal(restaurant)
		req, _ := http.NewRequest("POST", "/restaurants", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusInternalServerError, w.Code)
		}
	})
}
