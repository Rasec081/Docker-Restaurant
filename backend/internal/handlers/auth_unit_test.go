package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"restaurant-backend/internal/repository"
)

func TestAuthRegister_WithMock(t *testing.T) {
	mockRepo := repository.NewMockUserRepository()
	mockKeycloak := &mockKeycloakService{}
	InitAuthHandler(mockRepo, mockKeycloak, nil)

	router := setupTestRouter()
	router.POST("/auth/register", Register)

	t.Run("Register ok", func(t *testing.T) {
		body := map[string]string{
			"username": "newuser",
			"email":    "newuser@test.com",
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusCreated, w.Code)
		}
	})

	t.Run("Register keycloak error", func(t *testing.T) {
		mockKeycloak.createErr = errTest
		defer func() { mockKeycloak.createErr = nil }()

		body := map[string]string{
			"username": "newuser",
			"email":    "newuser@test.com",
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("Register repo error", func(t *testing.T) {
		mockRepo.CreateError = errTest
		defer func() { mockRepo.CreateError = nil }()

		body := map[string]string{
			"username": "newuser",
			"email":    "newuser@test.com",
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("Register invalid json", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBufferString("{"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestAuthLogin_WithMockServer(t *testing.T) {
	t.Setenv("KEYCLOAK_REALM", "test")
	t.Setenv("KEYCLOAK_CLIENT_ID", "client")
	t.Setenv("KEYCLOAK_CLIENT_SECRET", "secret")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/realms/test/protocol/openid-connect/token" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.FormValue("username") == "bad" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid"))
			return
		}
		if r.FormValue("username") == "badjson" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("not-json"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"access_token":"token"}`))
	}))
	defer srv.Close()

	t.Setenv("KEYCLOAK_URL", srv.URL)

	mockRepo := repository.NewMockUserRepository()
	mockKeycloak := &mockKeycloakService{}
	InitAuthHandler(mockRepo, mockKeycloak, srv.Client())

	router := setupTestRouter()
	router.POST("/auth/login", Login)

	t.Run("Login ok", func(t *testing.T) {
		body := map[string]string{"username": "admin", "password": "admin"}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusOK, w.Code)
		}
	})

	t.Run("Login missing credentials", func(t *testing.T) {
		body := map[string]string{"username": "", "password": ""}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("Login invalid credentials", func(t *testing.T) {
		body := map[string]string{"username": "bad", "password": "x"}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Login bad json from server", func(t *testing.T) {
		body := map[string]string{"username": "badjson", "password": "x"}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusInternalServerError, w.Code)
		}
	})

	t.Run("Login invalid json", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBufferString("{"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusBadRequest, w.Code)
		}
	})
}
