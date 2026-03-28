package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

func TestUserHandlers_WithMock(t *testing.T) {
	mockRepo := repository.NewMockUserRepository()
	InitUserHandler(mockRepo)

	router := setupTestRouter()
	router.GET("/users/me", func(c *gin.Context) {
		c.Set("username", "testuser")
		GetUserMe(c)
	})
	router.GET("/users/me-missing", GetUserMe)
	router.GET("/users/me-badtype", func(c *gin.Context) {
		c.Set("username", 123)
		GetUserMe(c)
	})
	router.PUT("/users/:id", UpdateUser)
	router.DELETE("/users/:id", DeleteUser)

	mockRepo.UsersByUsername["testuser"] = &models.User{ID: 1, Nombre: "testuser", RoleID: 2}

	t.Run("GetUserMe ok", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/me", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusOK, w.Code)
		}
	})

	t.Run("GetUserMe missing username", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/me-missing", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("GetUserMe bad username type", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/me-badtype", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("GetUserMe not found", func(t *testing.T) {
		delete(mockRepo.UsersByUsername, "testuser")
		defer func() { mockRepo.UsersByUsername["testuser"] = &models.User{ID: 1, Nombre: "testuser", RoleID: 2} }()

		req, _ := http.NewRequest("GET", "/users/me", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("UpdateUser ok", func(t *testing.T) {
		mockRepo.Users[1] = &models.User{ID: 1, Nombre: "old", RoleID: 1}

		input := map[string]interface{}{"nombre": "new", "rol": 2}
		body, _ := json.Marshal(input)
		req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusOK, w.Code)
		}
	})

	t.Run("UpdateUser invalid id", func(t *testing.T) {
		input := map[string]interface{}{"nombre": "new", "rol": 2}
		body, _ := json.Marshal(input)
		req, _ := http.NewRequest("PUT", "/users/abc", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("UpdateUser not found", func(t *testing.T) {
		input := map[string]interface{}{"nombre": "new", "rol": 2}
		body, _ := json.Marshal(input)
		req, _ := http.NewRequest("PUT", "/users/999", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("UpdateUser repo error", func(t *testing.T) {
		mockRepo.UpdateError = errTest
		defer func() { mockRepo.UpdateError = nil }()

		input := map[string]interface{}{"nombre": "new", "rol": 2}
		body, _ := json.Marshal(input)
		req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("DeleteUser ok", func(t *testing.T) {
		mockRepo.Users[10] = &models.User{ID: 10, Nombre: "delete", RoleID: 2}

		req, _ := http.NewRequest("DELETE", "/users/10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusOK, w.Code)
		}
	})

	t.Run("DeleteUser invalid id", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/users/abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("DeleteUser not found", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/users/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("DeleteUser foreign key", func(t *testing.T) {
		mockRepo.DeleteError = errors.New("violates foreign key")
		defer func() { mockRepo.DeleteError = nil }()

		req, _ := http.NewRequest("DELETE", "/users/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("DeleteUser repo error", func(t *testing.T) {
		mockRepo.DeleteError = errTest
		defer func() { mockRepo.DeleteError = nil }()

		req, _ := http.NewRequest("DELETE", "/users/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusInternalServerError, w.Code)
		}
	})
}
