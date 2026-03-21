package handlers

import (
	"net/http"

	"restaurant-backend/internal/db"

	"github.com/gin-gonic/gin"
)

type Menu struct {
	ID     int     `json:"id"`
	Nombre string  `json:"nombrePlato"`
	Precio float32 `json:"precio"`
}

func GetMenu(c *gin.Context) {
	id := c.Param("id")

	var menu Menu
	err := db.DB.QueryRow(
		"SELECT dish_name, price FROM menus WHERE menu_id = $1",
		id,
	).Scan(&menu.Nombre, &menu.Precio)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Menú no encontrado",
		})
		return
	}

	c.JSON(http.StatusOK, menu)
}
