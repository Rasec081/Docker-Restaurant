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

func TestCreateReservation(t *testing.T) {
	setupTestData(t)

	tests := []struct {
		name        string
		reservation models.Reservation
		expected    int
	}{
		{
			name: "Create valid reservation",
			reservation: models.Reservation{
				TableID:  1,
				ClientID: 2,
				Fecha:    "2026-03-27 19:00:00",
				Estado:   1,
			},
			expected: 200,
		},
		{
			name: "Create reservation with invalid table",
			reservation: models.Reservation{
				TableID:  999,
				ClientID: 2,
				Fecha:    "2026-03-27 19:00:00",
				Estado:   1,
			},
			expected: 500,
		},
		{
			name: "Create reservation with invalid client",
			reservation: models.Reservation{
				TableID:  1,
				ClientID: 999,
				Fecha:    "2026-03-27 19:00:00",
				Estado:   1,
			},
			expected: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()
			router.POST("/reservations", CreateReservation)

			jsonBody, _ := json.Marshal(tt.reservation)
			req, _ := http.NewRequest("POST", "/reservations", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expected {
				t.Errorf("Esperaba %d, obtuvo %d", tt.expected, w.Code)
			}
		})
	}
}

func TestDeleteReservation(t *testing.T) {
	setupTestData(t)

	// Crear reserva para eliminar
	var reservationID int
	testDB.QueryRow(`
		INSERT INTO Reservation (table_id, client_id, fecha, estado) 
		VALUES (1, 2, NOW(), 1) 
		RETURNING reservation_id
	`).Scan(&reservationID)

	tests := []struct {
		name          string
		reservationID string
		expected      int
	}{
		{
			name:          "Delete existing reservation",
			reservationID: strconv.Itoa(reservationID),
			expected:      200,
		},
		{
			name:          "Delete non-existent reservation",
			reservationID: "999",
			expected:      404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()
			router.DELETE("/reservations/:id", DeleteReservation)

			req, _ := http.NewRequest("DELETE", "/reservations/"+tt.reservationID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expected {
				t.Errorf("Esperaba %d, obtuvo %d", tt.expected, w.Code)
			}
		})
	}
}

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
