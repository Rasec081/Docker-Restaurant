package handlers

import (
	"net/http"

	"restaurant-backend/internal/db"

	"github.com/gin-gonic/gin"
)

type Order struct {
	ID           int    `json:"id"`
	TableID      *int   `json:"table_id"`
	ClientID     int    `json:"client_id"`
	OrdersType   string `json:"orders_type"`
	RestaurantID int    `json:"restaurant_id"`
}

func GetOrder(c *gin.Context) {
	id := c.Param("id")

	var o Order

	err := db.DB.QueryRow(
		"SELECT orders_id, table_id, client_id, orders_type, restaurant_id FROM Orders WHERE orders_id = $1",
		id,
	).Scan(&o.ID, &o.TableID, &o.ClientID, &o.OrdersType, &o.RestaurantID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Orden no encontrada",
		})
		return
	}

	c.JSON(http.StatusOK, o)
}

func CreateOrder(c *gin.Context) {
	var o Order

	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err := db.DB.Exec(
		"INSERT INTO Orders (table_id, client_id, orders_type, restaurant_id) VALUES ($1, $2, $3, $4)",
		o.TableID, o.ClientID, o.OrdersType, o.RestaurantID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Orden creada correctamente",
	})
}
