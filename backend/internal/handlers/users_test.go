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

func TestGetUserMe(t *testing.T) {
	mockRepo := repository.NewMockUserRepository()
	mockRepo.UsersByUsername["admin"] = &models.User{ID: 1, Nombre: "admin", RoleID: 1}
	handler := NewUserHandler(mockRepo)

	router := gin.Default()
	router.GET("/users/me", func(c *gin.Context) {
		c.Set("username", "admin")
		handler.GetUserMe(c)
	})

	req, _ := http.NewRequest("GET", "/users/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Esperaba 200, obtuvo %d", w.Code)
	}
}

func TestUpdateUser(t *testing.T) {
	mockRepo := repository.NewMockUserRepository()
	mockRepo.Users[1] = &models.User{ID: 1, Nombre: "admin", RoleID: 1}
	handler := NewUserHandler(mockRepo)

	tests := []struct {
		name   string
		userID string
		input  struct {
			Nombre string `json:"nombre"`
			Rol    int    `json:"rol"`
		}
		expected int
	}{
		{
			name:   "Update existing user",
			userID: "1",
			input: struct {
				Nombre string `json:"nombre"`
				Rol    int    `json:"rol"`
			}{
				Nombre: "Updated Admin",
				Rol:    1,
			},
			expected: 200,
		},
		{
			name:   "Update non-existent user",
			userID: "999",
			input: struct {
				Nombre string `json:"nombre"`
				Rol    int    `json:"rol"`
			}{
				Nombre: "Non Existent",
				Rol:    1,
			},
			expected: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()
			router.PUT("/users/:id", handler.UpdateUser)

			jsonBody, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("PUT", "/users/"+tt.userID, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expected {
				t.Errorf("Esperaba %d, obtuvo %d", tt.expected, w.Code)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	mockRepo := repository.NewMockUserRepository()
	mockRepo.Users[3] = &models.User{ID: 3, Nombre: "To Delete", RoleID: 2}
	handler := NewUserHandler(mockRepo)

	tests := []struct {
		name     string
		userID   string
		expected int
	}{
		{
			name:     "Delete existing user",
			userID:   "3",
			expected: 200,
		},
		{
			name:     "Delete non-existent user",
			userID:   "999",
			expected: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()
			router.DELETE("/users/:id", handler.DeleteUser)

			req, _ := http.NewRequest("DELETE", "/users/"+tt.userID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expected {
				t.Errorf("Esperaba %d, obtuvo %d", tt.expected, w.Code)
			}
		})
	}
}

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
