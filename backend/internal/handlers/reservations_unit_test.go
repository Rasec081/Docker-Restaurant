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

func TestReservationHandlers_WithMock(t *testing.T) {
	mockRepo := repository.NewMockReservationRepository()
	InitReservationHandler(mockRepo)

	router := setupTestRouter()
	router.POST("/reservations", CreateReservation)
	router.DELETE("/reservations/:id", DeleteReservation)

	t.Run("CreateReservation ok", func(t *testing.T) {
		reservation := models.Reservation{
			TableID:  1,
			ClientID: 2,
			Fecha:    "2026-03-27 19:00:00",
			Estado:   1,
		}
		body, _ := json.Marshal(reservation)
		req, _ := http.NewRequest("POST", "/reservations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusOK, w.Code)
		}
	})

	t.Run("CreateReservation invalid json", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/reservations", bytes.NewBufferString("{"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("CreateReservation repo error", func(t *testing.T) {
		mockRepo.CreateError = errTest
		defer func() { mockRepo.CreateError = nil }()

		reservation := models.Reservation{TableID: 1, ClientID: 2, Fecha: "2026-03-27 19:00:00", Estado: 1}
		body, _ := json.Marshal(reservation)
		req, _ := http.NewRequest("POST", "/reservations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("DeleteReservation ok", func(t *testing.T) {
		reservation := &models.Reservation{ID: 1, TableID: 1, ClientID: 2, Fecha: "2026-03-27 19:00:00", Estado: 1}
		mockRepo.Reservations[1] = reservation

		req, _ := http.NewRequest("DELETE", "/reservations/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusOK, w.Code)
		}
	})

	t.Run("DeleteReservation not found", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/reservations/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("DeleteReservation invalid id", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/reservations/abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("DeleteReservation repo error", func(t *testing.T) {
		mockRepo.DeleteError = errTest
		defer func() { mockRepo.DeleteError = nil }()

		req, _ := http.NewRequest("DELETE", "/reservations/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusInternalServerError, w.Code)
		}
	})
}
