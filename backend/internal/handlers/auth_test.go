package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"restaurant-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

func TestRegister(t *testing.T) {
	mockRepo := repository.NewMockUserRepository()
	mockKeycloak := &mockKeycloakService{}
	InitAuthHandler(mockRepo, mockKeycloak, nil)

	tests := []struct {
		name       string
		body       map[string]string
		expected   int
		wantErrMsg string
	}{
		{
			name: "Successful registration",
			body: map[string]string{
				"username": "newuser",
				"email":    "newuser@test.com",
				"password": "password123",
			},
			expected: 201,
		},
		{
			name: "Missing username",
			body: map[string]string{
				"email":    "test@test.com",
				"password": "password123",
			},
			expected:   400,
			wantErrMsg: "error",
		},
		{
			name: "Missing email",
			body: map[string]string{
				"username": "testuser",
				"password": "password123",
			},
			expected:   400,
			wantErrMsg: "error",
		},
		{
			name: "Missing password",
			body: map[string]string{
				"username": "testuser",
				"email":    "test@test.com",
			},
			expected:   400,
			wantErrMsg: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()
			router.POST("/auth/register", Register)

			jsonBody, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expected {
				t.Errorf("Esperaba %d, obtuvo %d", tt.expected, w.Code)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/realms/test/protocol/openid-connect/token" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"access_token":"token"}`))
	}))
	defer srv.Close()

	t.Setenv("KEYCLOAK_URL", srv.URL)
	t.Setenv("KEYCLOAK_REALM", "test")
	t.Setenv("KEYCLOAK_CLIENT_ID", "client")
	t.Setenv("KEYCLOAK_CLIENT_SECRET", "secret")

	mockRepo := repository.NewMockUserRepository()
	mockKeycloak := &mockKeycloakService{}
	InitAuthHandler(mockRepo, mockKeycloak, srv.Client())

	tests := []struct {
		name     string
		body     map[string]string
		expected int
	}{
		{
			name: "Login with correct credentials",
			body: map[string]string{
				"username": "admin",
				"password": "admin",
			},
			expected: 200,
		},
		{
			name: "Login with missing credentials",
			body: map[string]string{
				"username": "",
				"password": "",
			},
			expected: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()
			router.POST("/auth/login", Login)

			jsonBody, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expected {
				t.Errorf("Esperaba %d, obtuvo %d", tt.expected, w.Code)
			}
		})
	}
}
