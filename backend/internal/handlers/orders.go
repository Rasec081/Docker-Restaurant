package handlers

import (
	"net/http"
	"strconv"

	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

// OrderHandler maneja operaciones de ordenes
type OrderHandler struct {
	OrderRepo repository.OrderRepository
}

var orderHandler *OrderHandler

// NewOrderHandler crea una nueva instancia
func NewOrderHandler(orderRepo repository.OrderRepository) *OrderHandler {
	return &OrderHandler{OrderRepo: orderRepo}
}

// InitOrderHandler inicializa el handler global (para rutas y tests)
func InitOrderHandler(orderRepo repository.OrderRepository) {
	orderHandler = NewOrderHandler(orderRepo)
}

func getOrderHandler() *OrderHandler {
	ensureOrderHandler()
	return orderHandler
}

// GetOrder godoc
// @Summary Obtener orden por ID
// @Description Retorna una orden específica
// @Tags orders
// @Produce json
// @Param id path int true "ID de la orden"
// @Success 200 {object} models.Order
// @Failure 404 {object} map[string]string
// @Router /orders/{id} [get]
func GetOrder(c *gin.Context) {
	getOrderHandler().GetOrder(c)
}

// GetOrder godoc
func (h *OrderHandler) GetOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID inválido",
		})
		return
	}

	order, err := h.OrderRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Orden no encontrada",
		})
		return
	}

	c.JSON(http.StatusOK, order)
}

// CreateOrder godoc
// @Summary Crear orden
// @Description Crea una nueva orden
// @Tags orders
// @Accept json
// @Produce json
// @Param body body models.Order true "Datos de la orden"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders [post]
func CreateOrder(c *gin.Context) {
	getOrderHandler().CreateOrder(c)
}

// CreateOrder godoc
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var o models.Order

	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.OrderRepo.Create(&o); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Orden creada correctamente",
	})
}
