package handlers

import (
	"net/http"
	"restaurant-backend/internal/db"

	"github.com/gin-gonic/gin"
)

type Reservation struct {
	ID       int    `json:"id"`
	TableID  int    `json:"table_id"`
	ClientID int    `json:"client_id"`
	Fecha    string `json:"fecha"`
	Estado   int    `json:"estado"`
}

func CreateReservation(c *gin.Context) {
	var r Reservation

	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var id int
	err := db.DB.QueryRow(
		"INSERT INTO Reservation (table_id, client_id, fecha, estado) VALUES ($1, $2, $3, $4) RETURNING reservation_id",
		r.TableID, r.ClientID, r.Fecha, r.Estado,
	).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Reserva creada correctamente",
		"reservation_id": id,
	})
}

func DeleteReservation(c *gin.Context) {
	id := c.Param("id")

	result, err := db.DB.Exec(
		"DELETE FROM Reservation WHERE reservation_id = $1",
		id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	rows, _ := result.RowsAffected()

	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Reserva no encontrada",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reserva cancelada correctamente",
	})
}
