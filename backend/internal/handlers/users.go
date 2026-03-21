package handlers

import (
	"net/http"

	"restaurant-backend/internal/db"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
	Rol    int    `json:"rol"`
}

func GetUserMe(c *gin.Context) {
	row := db.DB.QueryRow("SELECT user_id, nombre, role_id FROM Users WHERE user_id = 1") // Despues el where tiene que venir del jwk

	var user User
	err := row.Scan(&user.ID, &user.Nombre, &user.Rol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
