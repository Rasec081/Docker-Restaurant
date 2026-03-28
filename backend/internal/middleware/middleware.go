package middleware

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ========================
// Middleware de roles
// ========================
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token requerido",
			})
			c.Abort()
			return
		}

		// formato: Bearer token
		tokenParts := strings.Split(authHeader, " ")

		if len(tokenParts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "formato de token inválido",
			})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// 🔥 DECODIFICAR JWT (sin validar firma por ahora)
		payload := strings.Split(token, ".")

		if len(payload) < 2 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token inválido",
			})
			c.Abort()
			return
		}

		// base64 decode
		decoded, err := decodeSegment(payload[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "error decodificando token",
			})
			c.Abort()
			return
		}

		// convertir a map
		var tokenData map[string]interface{}
		if err := json.Unmarshal(decoded, &tokenData); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "error parseando token",
			})
			c.Abort()
			return
		}

		// 🔥 EXTRAER ROLES
		realmAccess, ok := tokenData["realm_access"].(map[string]interface{})
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "no hay roles en token",
			})
			c.Abort()
			return
		}

		roles, ok := realmAccess["roles"].([]interface{})
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "roles inválidos",
			})
			c.Abort()
			return
		}

		// 🔍 BUSCAR SI TIENE EL ROL
		hasRole := false
		for _, r := range roles {
			if r.(string) == requiredRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "no tienes permisos",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func decodeSegment(seg string) ([]byte, error) {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}
	return base64.URLEncoding.DecodeString(seg)
}

//me diejron que esto lo borrara pero si no funciona me sirve para probar el middleware de roles, lo dejo comentado por si acaso
/*package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func loadJWKWithRetry(jwksURL string) (*keyfunc.JWKS, error) {
	var jwks *keyfunc.JWKS
	var err error

	for i := 0; i < 10; i++ {
		jwks, err = keyfunc.Get(jwksURL, keyfunc.Options{})
		if err == nil {
			log.Println("JWK cargado correctamente")
			return jwks, nil
		}

		log.Println("Intento", i+1, ": Keycloak no listo, reintentando...")
		time.Sleep(3 * time.Second)
	}

	return nil, err
}

func JWTMiddleware() gin.HandlerFunc {

	jwksURL := "http://keycloak:8080/realms/restaurant-realm/protocol/openid-connect/certs"

	jwks, err := loadJWKWithRetry(jwksURL)
	if err != nil {
		log.Println("Error cargando JWK:", err)
	}

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, jwks.Keyfunc)

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		c.Set("user_id", claims["sub"])
		c.Set("email", claims["email"])

		c.Next()
	}
}
*/
