package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

// ReservationHandler maneja operaciones de reservas
type ReservationHandler struct {
	ReservationRepo repository.ReservationRepository
}

var reservationHandler *ReservationHandler

// NewReservationHandler crea una nueva instancia
func NewReservationHandler(reservationRepo repository.ReservationRepository) *ReservationHandler {
	return &ReservationHandler{ReservationRepo: reservationRepo}
}

// InitReservationHandler inicializa el handler global (para rutas y tests)
func InitReservationHandler(reservationRepo repository.ReservationRepository) {
	reservationHandler = NewReservationHandler(reservationRepo)
}

func getReservationHandler() *ReservationHandler {
	ensureReservationHandler()
	return reservationHandler
}

// CreateReservation godoc
// @Summary Crear una reserva
// @Description Crea una nueva reserva en el sistema
// @Tags reservations
// @Accept json
// @Produce json
// @Param reservation body models.Reservation true "Datos de la reserva"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /reservations [post]
func CreateReservation(c *gin.Context) {
	getReservationHandler().CreateReservation(c)
}

// CreateReservation godoc
func (h *ReservationHandler) CreateReservation(c *gin.Context) {
	var r models.Reservation

	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	available, err := h.ReservationRepo.IsTableAvailable(r.TableID, r.Fecha)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if !available {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Mesa no disponible",
		})
		return
	}

	if err := h.ReservationRepo.Create(&r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Reserva creada correctamente",
		"reservation_id": r.ID,
	})
}

// DeleteReservation godoc
// @Summary Eliminar reserva
// @Description Cancela una reserva por ID
// @Tags reservations
// @Produce json
// @Param id path int true "ID de la reserva"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reservations/{id} [delete]
func DeleteReservation(c *gin.Context) {
	getReservationHandler().DeleteReservation(c)
}

// DeleteReservation godoc
func (h *ReservationHandler) DeleteReservation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID inválido",
		})
		return
	}

	if err := h.ReservationRepo.Delete(id); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "no encontrada") || strings.Contains(strings.ToLower(err.Error()), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Reserva no encontrada",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reserva cancelada correctamente",
	})
}
