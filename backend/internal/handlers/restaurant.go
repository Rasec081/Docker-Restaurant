package handlers

import (
	"net/http"
	"restaurant-backend/internal/db"

	"github.com/gin-gonic/gin"
)

type Restaurant struct {
	ID      int    `json:"id"`
	Nombre  string `json:"nombre"`
	Estado  int    `json:"estado"`
	AdminID int    `json:"admin_id"`
}

// GetRestaurants godoc
// @Summary Obtener todos los restaurantes
// @Description Retorna una lista de restaurantes
// @Tags restaurants
// @Produce json
// @Success 200 {array} Restaurant
// @Failure 500 {object} map[string]string
// @Router /restaurants [get]
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

// CreateRestaurant godoc
// @Summary Crear restaurante
// @Description Crea un nuevo restaurante
// @Tags restaurants
// @Accept json
// @Produce json
// @Param body body Restaurant true "Datos del restaurante"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /restaurants [post]
func CreateRestaurant(c *gin.Context) {
	var r Restaurant

	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err := db.DB.Exec(
		"INSERT INTO restaurant (nombre, admin_id, estado) VALUES ($1, $2, $3)",
		r.Nombre, r.AdminID, r.Estado,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Restaurant creado correctamente",
	})
}
