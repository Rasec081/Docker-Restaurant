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
