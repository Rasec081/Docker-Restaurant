package middleware

import (
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func JWTMiddleware() gin.HandlerFunc {

	jwksURL := "http://localhost:8080/realms/restaurant-realm/protocol/openid-connect/certs"

	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{})
	if err != nil {
		panic("No se pudo cargar el JWK")
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
