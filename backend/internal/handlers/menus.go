package handlers

import (
	"net/http"

	"restaurant-backend/internal/db"

	"github.com/gin-gonic/gin"
)

type Menu struct {
	ID     int     `json:"id"`
	Nombre string  `json:"nombrePlato"`
	Precio float64 `json:"precio"`
}

func GetMenu(c *gin.Context) {
	idParam := c.Param("id")

	rows, err := db.DB.Query(
		"SELECT menu_id, dish_name, price FROM Menu WHERE restaurant_id = $1",
		idParam,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer rows.Close()

	var menu []Menu

	for rows.Next() {
		var m Menu
		err := rows.Scan(&m.ID, &m.Nombre, &m.Precio)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		menu = append(menu, m)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, menu)
}
