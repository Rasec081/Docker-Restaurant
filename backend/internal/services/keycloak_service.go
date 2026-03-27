package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func GetAdminToken() (string, error) {

	data := url.Values{}
	data.Set("client_id", "admin-cli")
	data.Set("username", "admin")
	data.Set("password", "admin")
	data.Set("grant_type", "password")

	req, err := http.NewRequest(
		"POST",
		"http://keycloak:8080/realms/master/protocol/openid-connect/token",
		bytes.NewBufferString(data.Encode()),
	)

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
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

func CreateUserInKeycloak(username, email, password string) error {

	adminToken, err := GetAdminToken()
	if err != nil {
		return fmt.Errorf("error obteniendo admin token: %v", err)
	}

	// 🔥 Extraer nombre y apellido desde el email (simple)
	firstName := username
	lastName := "User"

	// ========================
	// 1. Crear usuario
	// ========================
	createURL := "http://keycloak:8080/admin/realms/restaurant-realm/users"

	body := map[string]interface{}{
		"username":        username,
		"email":           email,
		"enabled":         true,
		"emailVerified":   true,
		"firstName":       firstName,
		"lastName":        lastName,
		"requiredActions": []string{},
	}

	jsonData, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", createURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 201 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("keycloak error: %s", string(bodyBytes))
	}

	// ========================
	// 2. Obtener userId
	// ========================
	searchURL := fmt.Sprintf(
		"http://keycloak:8080/admin/realms/restaurant-realm/users?username=%s",
		username,
	)

	reqSearch, _ := http.NewRequest("GET", searchURL, nil)
	reqSearch.Header.Set("Authorization", "Bearer "+adminToken)

	respSearch, err := client.Do(reqSearch)
	if err != nil {
		return err
	}
	defer respSearch.Body.Close()

	var users []map[string]interface{}
	json.NewDecoder(respSearch.Body).Decode(&users)

	if len(users) == 0 {
		return fmt.Errorf("usuario no encontrado en keycloak")
	}

	userID := users[0]["id"].(string)

	// ========================
	// 3. Setear password
	// ========================
	passwordURL := fmt.Sprintf(
		"http://keycloak:8080/admin/realms/restaurant-realm/users/%s/reset-password",
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
		return err
	}

	if respPass.StatusCode != 204 {
		bodyBytes, _ := io.ReadAll(respPass.Body)
		return fmt.Errorf("error seteando password: %s", string(bodyBytes))
	}

	return nil
}
