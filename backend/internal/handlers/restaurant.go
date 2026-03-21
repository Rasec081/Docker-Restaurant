package handlers

import (
	"net/http"

	"restaurant-backend/internal/db"

	"github.com/gin-gonic/gin"
)

type Restaurant struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
	Estado int    `json:"estado"`
}

func GetRestaurants(c *gin.Context) {

	rows, err := db.DB.Query("SELECT restaurant_id, nombre, estado FROM restaurant")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer rows.Close()

	var restaurants []Restaurant

	for rows.Next() {
		var r Restaurant
		err := rows.Scan(&r.ID, &r.Nombre, &r.Estado)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		restaurants = append(restaurants, r)
	}

	c.JSON(http.StatusOK, restaurants)
}
