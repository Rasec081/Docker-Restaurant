package handlers

import (
	"net/http"

	"restaurant-backend/internal/db"

	"github.com/gin-gonic/gin"
)

type Order struct {
	ID          int     `json:"id"`
	Table       string  `json:"mesa"`
	Cliente     float64 `json:"cliente"`
	Orders_type string  `json:"tipoOrden"`
}

func GetOrder(c *gin.Context) {
	idParam := c.Param("id")

	var order Order

	err := db.DB.QueryRow(
		"SELECT orders_id, table_id, client_id, orders_type FROM Orders WHERE orders_id = $1", // Hacer un join para que se muestre el nombre del cliente
		idParam,
	).Scan(&order.ID, &order.Table, &order.Cliente, &order.Orders_type)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, order)
}
