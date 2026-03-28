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
