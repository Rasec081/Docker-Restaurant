package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type keycloakTestConfig struct {
	tokenStatus    int
	tokenBody      map[string]interface{}
	createStatus   int
	searchUsers    []map[string]interface{}
	passwordStatus int
	roles          []map[string]interface{}
	assignStatus   int
}

func newKeycloakServer(t *testing.T, cfg keycloakTestConfig) *httptest.Server {
	if cfg.tokenStatus == 0 {
		cfg.tokenStatus = http.StatusOK
	}
	if cfg.createStatus == 0 {
		cfg.createStatus = http.StatusCreated
	}
	if cfg.passwordStatus == 0 {
		cfg.passwordStatus = http.StatusNoContent
	}
	if cfg.assignStatus == 0 {
		cfg.assignStatus = http.StatusNoContent
	}
	if cfg.tokenBody == nil {
		cfg.tokenBody = map[string]interface{}{"access_token": "admin-token"}
	}
	if cfg.searchUsers == nil {
		cfg.searchUsers = []map[string]interface{}{{"id": "user-1"}}
	}
	if cfg.roles == nil {
		cfg.roles = []map[string]interface{}{{"id": "role-1", "name": "client"}}
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/realms/master/protocol/openid-connect/token":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(cfg.tokenStatus)
			json.NewEncoder(w).Encode(cfg.tokenBody)
		case "/admin/realms/restaurant-realm/users":
			if r.Method == http.MethodPost {
				w.WriteHeader(cfg.createStatus)
				return
			}
			if r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(cfg.searchUsers)
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
		case "/admin/realms/restaurant-realm/users/user-1/reset-password":
			w.WriteHeader(cfg.passwordStatus)
		case "/admin/realms/restaurant-realm/roles":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(cfg.roles)
		case "/admin/realms/restaurant-realm/users/user-1/role-mappings/realm":
			w.WriteHeader(cfg.assignStatus)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestGetAdminToken_Success(t *testing.T) {
	srv := newKeycloakServer(t, keycloakTestConfig{})
	defer srv.Close()

	t.Setenv("KEYCLOAK_ADMIN_URL", srv.URL)

	token, err := GetAdminToken()
	if err != nil {
		t.Fatalf("no esperaba error, obtuvo %v", err)
	}
	if token == "" {
		t.Fatal("token vacio")
	}
}

func TestGetAdminToken_MissingToken(t *testing.T) {
	srv := newKeycloakServer(t, keycloakTestConfig{tokenBody: map[string]interface{}{}})
	defer srv.Close()

	t.Setenv("KEYCLOAK_ADMIN_URL", srv.URL)

	_, err := GetAdminToken()
	if err == nil {
		t.Fatal("esperaba error por token faltante")
	}
}

func TestGetAdminToken_InvalidURL(t *testing.T) {
	t.Setenv("KEYCLOAK_ADMIN_URL", "://bad-url")
	_, err := GetAdminToken()
	if err == nil {
		t.Fatal("esperaba error por URL invalida")
	}
}

func TestCreateUserInKeycloak_Success(t *testing.T) {
	srv := newKeycloakServer(t, keycloakTestConfig{})
	defer srv.Close()

	t.Setenv("KEYCLOAK_ADMIN_URL", srv.URL)

	_, err := CreateUserInKeycloak("user", "user@test.com", "pass", "client")
	if err != nil {
		t.Fatalf("no esperaba error, obtuvo %v", err)
	}
}

func TestCreateUserInKeycloak_CreateError(t *testing.T) {
	srv := newKeycloakServer(t, keycloakTestConfig{createStatus: http.StatusBadRequest})
	defer srv.Close()

	t.Setenv("KEYCLOAK_ADMIN_URL", srv.URL)

	_, err := CreateUserInKeycloak("user", "user@test.com", "pass", "client")
	if err == nil {
		t.Fatal("esperaba error por create")
	}
}

func TestCreateUserInKeycloak_UserNotFound(t *testing.T) {
	srv := newKeycloakServer(t, keycloakTestConfig{searchUsers: []map[string]interface{}{}})
	defer srv.Close()

	t.Setenv("KEYCLOAK_ADMIN_URL", srv.URL)

	_, err := CreateUserInKeycloak("user", "user@test.com", "pass", "client")
	if err == nil {
		t.Fatal("esperaba error por usuario no encontrado")
	}
}

func TestCreateUserInKeycloak_PasswordError(t *testing.T) {
	srv := newKeycloakServer(t, keycloakTestConfig{passwordStatus: http.StatusBadRequest})
	defer srv.Close()

	t.Setenv("KEYCLOAK_ADMIN_URL", srv.URL)

	_, err := CreateUserInKeycloak("user", "user@test.com", "pass", "client")
	if err == nil {
		t.Fatal("esperaba error seteando password")
	}
}

func TestCreateUserInKeycloak_RoleNotFound(t *testing.T) {
	srv := newKeycloakServer(t, keycloakTestConfig{roles: []map[string]interface{}{{"id": "role-1", "name": "admin"}}})
	defer srv.Close()

	t.Setenv("KEYCLOAK_ADMIN_URL", srv.URL)

	_, err := CreateUserInKeycloak("user", "user@test.com", "pass", "client")
	if err == nil {
		t.Fatal("esperaba error por rol no encontrado")
	}
}

func TestCreateUserInKeycloak_AssignError(t *testing.T) {
	srv := newKeycloakServer(t, keycloakTestConfig{assignStatus: http.StatusBadRequest})
	defer srv.Close()

	t.Setenv("KEYCLOAK_ADMIN_URL", srv.URL)

	_, err := CreateUserInKeycloak("user", "user@test.com", "pass", "client")
	if err == nil {
		t.Fatal("esperaba error asignando rol")
	}
}

func TestCreateUserInKeycloak_AdminTokenError(t *testing.T) {
	t.Setenv("KEYCLOAK_ADMIN_URL", "://bad-url")
	_, err := CreateUserInKeycloak("user", "user@test.com", "pass", "client")
	if err == nil {
		t.Fatal("esperaba error obteniendo admin token")
	}
}

func TestGetKeycloakBaseURL(t *testing.T) {
	t.Setenv("KEYCLOAK_ADMIN_URL", "http://admin")
	if got := getKeycloakBaseURL(); got != "http://admin" {
		t.Fatalf("esperaba http://admin, obtuvo %s", got)
	}

	t.Setenv("KEYCLOAK_ADMIN_URL", "")
	t.Setenv("KEYCLOAK_URL", "http://public/")
	if got := getKeycloakBaseURL(); got != "http://public" {
		t.Fatalf("esperaba http://public, obtuvo %s", got)
	}
}
