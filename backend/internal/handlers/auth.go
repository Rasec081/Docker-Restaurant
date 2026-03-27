package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"restaurant-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := services.CreateUserInKeycloak(
		input.Username,
		input.Email,
		input.Password,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"message": "Usuario creado correctamente",
	})
}

func Login(c *gin.Context) {
	var input LoginRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "JSON inválido",
		})
		return
	}

	if input.Username == "" || input.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "username y password son requeridos",
		})
		return
	}

	keycloakURL := os.Getenv("KEYCLOAK_URL")
	realm := os.Getenv("KEYCLOAK_REALM")
	clientID := os.Getenv("KEYCLOAK_CLIENT_ID")
	clientSecret := os.Getenv("KEYCLOAK_CLIENT_SECRET")

	tokenURL := keycloakURL + "/realms/" + realm + "/protocol/openid-connect/token"

	form := url.Values{}
	form.Set("grant_type", "password")
	form.Set("client_id", clientID)
	form.Set("username", input.Username)
	form.Set("password", input.Password)

	if clientSecret != "" {
		form.Set("client_secret", clientSecret)
	}

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(form.Encode()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "No se pudo crear la petición a Keycloak",
		})
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "No se pudo conectar con Keycloak",
		})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "No se pudo leer la respuesta de Keycloak",
		})
		return
	}

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Credenciales inválidas",
			"details": string(body),
		})
		return
	}

	var tokenResponse map[string]interface{}
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Respuesta inválida de Keycloak",
		})
		return
	}

	c.JSON(http.StatusOK, tokenResponse)
}
