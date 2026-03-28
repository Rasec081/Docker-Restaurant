package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// KeycloakService define operaciones de autenticacion externa
type KeycloakService interface {
	CreateUser(username, email, password, role string) (string, error)
}

// AuthHandler maneja operaciones de auth
type AuthHandler struct {
	UserRepo   repository.UserRepository
	Keycloak   KeycloakService
	HTTPClient *http.Client
}

var authHandler *AuthHandler

// NewAuthHandler crea una nueva instancia
func NewAuthHandler(userRepo repository.UserRepository, keycloakSvc KeycloakService, httpClient *http.Client) *AuthHandler {
	return &AuthHandler{UserRepo: userRepo, Keycloak: keycloakSvc, HTTPClient: httpClient}
}

// InitAuthHandler inicializa el handler global (para rutas y tests)
func InitAuthHandler(userRepo repository.UserRepository, keycloakSvc KeycloakService, httpClient *http.Client) {
	authHandler = NewAuthHandler(userRepo, keycloakSvc, httpClient)
}

func getAuthHandler() *AuthHandler {
	ensureAuthHandler()
	return authHandler
}

// Register godoc
// @Summary Registrar usuario
// @Description Crea un usuario en Keycloak
// @Tags auth
// @Accept json
// @Produce json
// @Param body body object true "Datos de registro"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
func Register(c *gin.Context) {
	getAuthHandler().Register(c)
}

// Register godoc
func (h *AuthHandler) Register(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if input.Username == "" || input.Email == "" || input.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username, email y password son requeridos"})
		return
	}

	// 1. Crear en Keycloak
	_, err := h.Keycloak.CreateUser(
		input.Username,
		input.Email,
		input.Password,
		"client",
	)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 2. Crear en DB
	user := models.User{
		Nombre: input.Username,
		RoleID: 2,
	}
	if err := h.UserRepo.Create(&user); err != nil {
		c.JSON(500, gin.H{
			"error": "usuario creado en Keycloak pero falló en DB",
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "Usuario creado correctamente",
		"user_id": user.ID,
	})
}

// Login godoc
// @Summary Login de usuario
// @Description Autentica un usuario contra Keycloak y devuelve el token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Credenciales de usuario"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /login [post]
func Login(c *gin.Context) {
	getAuthHandler().Login(c)
}

// Login godoc
func (h *AuthHandler) Login(c *gin.Context) {

	// =========================
	// 1. Leer JSON
	// =========================
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

	// =========================
	// 2. Variables de entorno
	// =========================
	keycloakURL := os.Getenv("KEYCLOAK_URL")
	realm := os.Getenv("KEYCLOAK_REALM")
	clientID := os.Getenv("KEYCLOAK_CLIENT_ID")
	clientSecret := os.Getenv("KEYCLOAK_CLIENT_SECRET")

	tokenURL := keycloakURL + "/realms/" + realm + "/protocol/openid-connect/token"

	// =========================
	// 3. Crear body (form-data)
	// =========================
	form := url.Values{}
	form.Set("grant_type", "password")
	form.Set("client_id", clientID)
	form.Set("username", input.Username)
	form.Set("password", input.Password)

	if clientSecret != "" {
		form.Set("client_secret", clientSecret)
	}

	// =========================
	// 4. Crear request
	// =========================
	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(form.Encode()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "No se pudo crear la petición a Keycloak",
		})
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := h.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}

	// =========================
	// 5. Retry (IMPORTANTE)
	// =========================
	var resp *http.Response

	for i := 0; i < 5; i++ {
		resp, err = client.Do(req)
		if err == nil {
			break
		}
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "No se pudo conectar con Keycloak",
		})
		return
	}
	defer resp.Body.Close()

	// =========================
	// 6. Leer respuesta
	// =========================
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "No se pudo leer la respuesta de Keycloak",
		})
		return
	}

	// =========================
	// 7. Manejar error login
	// =========================
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Credenciales inválidas",
			"details": string(body),
		})
		return
	}

	// =========================
	// 8. Parsear token
	// =========================
	var tokenResponse map[string]interface{}
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Respuesta inválida de Keycloak",
		})
		return
	}

	// =========================
	// 9. Respuesta final
	// =========================
	c.JSON(http.StatusOK, tokenResponse)
}
