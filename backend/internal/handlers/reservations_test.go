package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreateReservation_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// ⚠️ Asegúrate que tu DB esté corriendo
	// docker compose up -d

	router := gin.Default()
	router.POST("/reservations", CreateReservation)

	body := `{
		"table_id": 1,
		"client_id": 1,
		"fecha": "2026-03-27",
		"estado": 1
	}`

	req, _ := http.NewRequest("POST", "/reservations", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Esperaba 200, obtuvo %d", w.Code)
	}
}

func TestCreateReservation_BadJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.POST("/reservations", CreateReservation)

	body := `{ MAL JSON }`

	req, _ := http.NewRequest("POST", "/reservations", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Esperaba 400, obtuvo %d", w.Code)
	}
}

func TestDeleteReservation_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.DELETE("/reservations/:id", DeleteReservation)

	req, _ := http.NewRequest("DELETE", "/reservations/1", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Esperaba 200, obtuvo %d", w.Code)
	}
}

func TestDeleteReservation_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.DELETE("/reservations/:id", DeleteReservation)

	req, _ := http.NewRequest("DELETE", "/reservations/9999", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Esperaba 404, obtuvo %d", w.Code)
	}
}
