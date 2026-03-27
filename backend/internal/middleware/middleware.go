package middleware

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
