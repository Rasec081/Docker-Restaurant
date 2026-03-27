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

	url := "http://keycloak:8080/admin/realms/restaurant-realm/users"

	adminToken, err := GetAdminToken()
	if err != nil {
		return fmt.Errorf("error obteniendo admin token: %v", err)
	}

	body := map[string]interface{}{
		"username": username,
		"email":    email,
		"enabled":  true,
		"credentials": []map[string]interface{}{
			{
				"type":      "password",
				"value":     password,
				"temporary": false,
			},
		},
	}

	jsonData, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

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

	return nil
}
