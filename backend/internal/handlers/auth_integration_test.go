package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestAuthIntegration(t *testing.T) {
	setupTestData(t)
	initIntegrationHandlers()

	if os.Getenv("KEYCLOAK_URL") == "" {
		os.Setenv("KEYCLOAK_URL", "http://localhost:8081")
	}
	if os.Getenv("KEYCLOAK_REALM") == "" {
		os.Setenv("KEYCLOAK_REALM", "restaurant-realm")
	}
	if os.Getenv("KEYCLOAK_CLIENT_ID") == "" {
		os.Setenv("KEYCLOAK_CLIENT_ID", "restaurant-client")
	}
	if os.Getenv("KEYCLOAK_CLIENT_SECRET") == "" {
		os.Setenv("KEYCLOAK_CLIENT_SECRET", "")
	}

	router := gin.Default()
	router.POST("/auth/register", Register)
	router.POST("/auth/login", Login)

	username := "intuser" + itoa(int(time.Now().UnixNano()%1000000))
	password := "Pass12345"

	registerBody, _ := json.Marshal(map[string]string{
		"username": username,
		"email":    username + "@test.com",
		"password": password,
	})

	var registerW *httptest.ResponseRecorder
	for i := 0; i < 5; i++ {
		registerReq, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(registerBody))
		registerReq.Header.Set("Content-Type", "application/json")
		registerW = httptest.NewRecorder()
		router.ServeHTTP(registerW, registerReq)
		if registerW.Code == http.StatusCreated {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if registerW.Code != http.StatusCreated {
		t.Fatalf("Esperaba %d, obtuvo %d", http.StatusCreated, registerW.Code)
	}

	loginBody, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})

	var loginW *httptest.ResponseRecorder
	for i := 0; i < 5; i++ {
		loginReq, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(loginBody))
		loginReq.Header.Set("Content-Type", "application/json")
		loginW = httptest.NewRecorder()
		router.ServeHTTP(loginW, loginReq)
		if loginW.Code == http.StatusOK {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if loginW.Code != http.StatusOK {
		t.Fatalf("Esperaba %d, obtuvo %d", http.StatusOK, loginW.Code)
	}
}
