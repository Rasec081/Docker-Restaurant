package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func getKeycloakBaseURL() string {
	if baseURL := os.Getenv("KEYCLOAK_ADMIN_URL"); baseURL != "" {
		return strings.TrimRight(baseURL, "/")
	}
	if baseURL := os.Getenv("KEYCLOAK_URL"); baseURL != "" {
		return strings.TrimRight(baseURL, "/")
	}
	return "http://keycloak:8080"
}

// ========================
// Obtener token admin
// ========================
func GetAdminToken() (string, error) {

	data := url.Values{}
	data.Set("client_id", "admin-cli")
	data.Set("username", "admin")
	data.Set("password", "admin")
	data.Set("grant_type", "password")

	client := &http.Client{}
	endpoint := getKeycloakBaseURL() + "/realms/master/protocol/openid-connect/token"
	var lastErr error

	for i := 0; i < 5; i++ {
		req, err := http.NewRequest(
			"POST",
			endpoint,
			bytes.NewBufferString(data.Encode()),
		)
		if err != nil {
			return "", err
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(2 * time.Second)
			continue
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		token, ok := result["access_token"].(string)
		if !ok {
			return "", fmt.Errorf("no se pudo obtener el token")
		}

		return token, nil
	}

	if lastErr != nil {
		return "", lastErr
	}
	return "", fmt.Errorf("no se pudo obtener el token")
}

// ========================
// Crear usuario + rol
// ========================
func CreateUserInKeycloak(username, email, password, roleName string) (string, error) {

	// 1. Obtener admin token
	adminToken, err := GetAdminToken()
	if err != nil {
		return "", fmt.Errorf("error obteniendo admin token: %v", err)
	}

	client := &http.Client{}

	// ========================
	// 2. Crear usuario
	// ========================
	createURL := getKeycloakBaseURL() + "/admin/realms/restaurant-realm/users"

	body := map[string]interface{}{
		"username":        username,
		"email":           email,
		"enabled":         true,
		"emailVerified":   true,
		"firstName":       username,
		"lastName":        "User",
		"requiredActions": []string{},
	}

	jsonData, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", createURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 201 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("error creando usuario: %s", string(bodyBytes))
	}

	// ========================
	// 3. Obtener userId
	// ========================
	searchURL := fmt.Sprintf(
		getKeycloakBaseURL()+"/admin/realms/restaurant-realm/users?username=%s",
		username,
	)

	reqSearch, _ := http.NewRequest("GET", searchURL, nil)
	reqSearch.Header.Set("Authorization", "Bearer "+adminToken)

	respSearch, err := client.Do(reqSearch)
	if err != nil {
		return "", err
	}
	defer respSearch.Body.Close()

	var users []map[string]interface{}
	json.NewDecoder(respSearch.Body).Decode(&users)

	if len(users) == 0 {
		return "", fmt.Errorf("usuario no encontrado en keycloak")
	}

	userID := users[0]["id"].(string)

	// ========================
	// 4. Setear password
	// ========================
	passwordURL := fmt.Sprintf(
		getKeycloakBaseURL()+"/admin/realms/restaurant-realm/users/%s/reset-password",
		userID,
	)

	passwordData := map[string]interface{}{
		"type":      "password",
		"value":     password,
		"temporary": false,
	}

	passwordJSON, _ := json.Marshal(passwordData)

	reqPass, _ := http.NewRequest("PUT", passwordURL, bytes.NewBuffer(passwordJSON))
	reqPass.Header.Set("Content-Type", "application/json")
	reqPass.Header.Set("Authorization", "Bearer "+adminToken)

	respPass, err := client.Do(reqPass)
	if err != nil {
		return "", err
	}

	if respPass.StatusCode != 204 {
		bodyBytes, _ := io.ReadAll(respPass.Body)
		return "", fmt.Errorf("error seteando password: %s", string(bodyBytes))
	}

	// ========================
	// 5. Obtener rol (FIX)
	// ========================
	rolesURL := getKeycloakBaseURL() + "/admin/realms/restaurant-realm/roles"

	reqRoles, _ := http.NewRequest("GET", rolesURL, nil)
	reqRoles.Header.Set("Authorization", "Bearer "+adminToken)

	respRoles, err := client.Do(reqRoles)
	if err != nil {
		return "", err
	}
	defer respRoles.Body.Close()

	var roles []map[string]interface{}
	json.NewDecoder(respRoles.Body).Decode(&roles)

	var foundRole map[string]interface{}

	for _, r := range roles {
		if r["name"] == roleName {
			foundRole = r
			break
		}
	}

	if foundRole == nil {
		return "", fmt.Errorf("rol no encontrado: %s", roleName)
	}
	// ========================
	// 6. Asignar rol al usuario
	// ========================
	assignURL := fmt.Sprintf(
		getKeycloakBaseURL()+"/admin/realms/restaurant-realm/users/%s/role-mappings/realm",
		userID,
	)

	cleanRole := map[string]interface{}{
		"id":   foundRole["id"],
		"name": foundRole["name"],
	}

	roleArray := []map[string]interface{}{cleanRole}
	roleJSON, _ := json.Marshal(roleArray)

	reqAssign, _ := http.NewRequest("POST", assignURL, bytes.NewBuffer(roleJSON))
	reqAssign.Header.Set("Content-Type", "application/json")
	reqAssign.Header.Set("Authorization", "Bearer "+adminToken)

	respAssign, err := client.Do(reqAssign)
	if err != nil {
		return "", err
	}

	if respAssign.StatusCode != 204 {
		bodyBytes, _ := io.ReadAll(respAssign.Body)
		return "", fmt.Errorf("error asignando rol: %s", string(bodyBytes))
	}

	return userID, nil
}

// DefaultKeycloakService adapta las funciones existentes para inyeccion
type DefaultKeycloakService struct{}

func NewDefaultKeycloakService() *DefaultKeycloakService {
	return &DefaultKeycloakService{}
}

func (s *DefaultKeycloakService) CreateUser(username, email, password, roleName string) (string, error) {
	return CreateUserInKeycloak(username, email, password, roleName)
}
