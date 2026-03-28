package handlers

import (
	"net/http"

	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

// RestaurantHandler maneja operaciones de restaurantes
type RestaurantHandler struct {
	RestaurantRepo repository.RestaurantRepository
}

var restaurantHandler *RestaurantHandler

// NewRestaurantHandler crea una nueva instancia
func NewRestaurantHandler(restaurantRepo repository.RestaurantRepository) *RestaurantHandler {
	return &RestaurantHandler{RestaurantRepo: restaurantRepo}
}

// InitRestaurantHandler inicializa el handler global (para rutas y tests)
func InitRestaurantHandler(restaurantRepo repository.RestaurantRepository) {
	restaurantHandler = NewRestaurantHandler(restaurantRepo)
}

func getRestaurantHandler() *RestaurantHandler {
	ensureRestaurantHandler()
	return restaurantHandler
}

// GetRestaurants godoc
// @Summary Obtener todos los restaurantes
// @Description Retorna una lista de restaurantes
// @Tags restaurants
// @Produce json
// @Success 200 {array} models.Restaurant
// @Failure 500 {object} map[string]string
// @Router /restaurants [get]
func GetRestaurants(c *gin.Context) {
	getRestaurantHandler().GetRestaurants(c)
}

// GetRestaurants godoc
func (h *RestaurantHandler) GetRestaurants(c *gin.Context) {
	restaurants, err := h.RestaurantRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, restaurants)
}

// CreateRestaurant godoc
// @Summary Crear restaurante
// @Description Crea un nuevo restaurante
// @Tags restaurants
// @Accept json
// @Produce json
// @Param body body models.Restaurant true "Datos del restaurante"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /restaurants [post]
func CreateRestaurant(c *gin.Context) {
	getRestaurantHandler().CreateRestaurant(c)
}

// CreateRestaurant godoc
func (h *RestaurantHandler) CreateRestaurant(c *gin.Context) {
	var r models.Restaurant

	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if r.Nombre == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "nombre requerido",
		})
		return
	}

	if err := h.RestaurantRepo.Create(&r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Restaurant creado correctamente",
	})
}
