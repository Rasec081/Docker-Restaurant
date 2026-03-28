package middleware

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func buildJWT(payload map[string]interface{}) string {
	header := map[string]interface{}{"alg": "none", "typ": "JWT"}
	headerBytes, _ := json.Marshal(header)
	payloadBytes, _ := json.Marshal(payload)

	headerEnc := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(headerBytes)
	payloadEnc := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(payloadBytes)

	return headerEnc + "." + payloadEnc + "."
}

func TestRequireRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Missing token", func(t *testing.T) {
		router := gin.New()
		router.GET("/protected", RequireRole("admin"), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest("GET", "/protected", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Invalid format", func(t *testing.T) {
		router := gin.New()
		router.GET("/protected", RequireRole("admin"), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Invalid token", func(t *testing.T) {
		router := gin.New()
		router.GET("/protected", RequireRole("admin"), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer abc")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("No roles", func(t *testing.T) {
		token := buildJWT(map[string]interface{}{
			"sub": "1",
		})

		router := gin.New()
		router.GET("/protected", RequireRole("admin"), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusForbidden, w.Code)
		}
	})

	t.Run("Wrong role", func(t *testing.T) {
		token := buildJWT(map[string]interface{}{
			"realm_access": map[string]interface{}{
				"roles": []interface{}{"client"},
			},
		})

		router := gin.New()
		router.GET("/protected", RequireRole("admin"), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusForbidden, w.Code)
		}
	})

	t.Run("Ok role", func(t *testing.T) {
		token := buildJWT(map[string]interface{}{
			"preferred_username": "admin",
			"realm_access": map[string]interface{}{
				"roles": []interface{}{"admin"},
			},
		})

		router := gin.New()
		router.GET("/protected", RequireRole("admin"), func(c *gin.Context) {
			if username, ok := c.Get("username"); !ok || username != "admin" {
				c.Status(http.StatusInternalServerError)
				return
			}
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusOK, w.Code)
		}
	})

	t.Run("Invalid payload json", func(t *testing.T) {
		// Payload no es JSON valido
		header := buildJWT(map[string]interface{}{})
		parts := strings.Split(header, ".")
		if len(parts) < 2 {
			t.Fatal("token invalido")
		}
		badPayload := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte("not-json"))
		token := parts[0] + "." + badPayload + "."

		router := gin.New()
		router.GET("/protected", RequireRole("admin"), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Roles not array", func(t *testing.T) {
		token := buildJWT(map[string]interface{}{
			"realm_access": map[string]interface{}{
				"roles": "admin",
			},
		})

		router := gin.New()
		router.GET("/protected", RequireRole("admin"), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Fatalf("Esperaba %d, obtuvo %d", http.StatusForbidden, w.Code)
		}
	})
}
