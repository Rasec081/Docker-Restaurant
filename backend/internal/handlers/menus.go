package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"restaurant-backend/internal/db"

	"github.com/gin-gonic/gin"
)

type Menu struct {
	ID           int     `json:"id"`
	Nombre       string  `json:"nombre"`
	Precio       float64 `json:"precio"`
	RestaurantID int     `json:"restaurant_id"`
}

// GetMenu godoc
// @Summary Obtener menú por ID
// @Description Retorna un menú específico
// @Tags menus
// @Produce json
// @Param id path int true "ID del menú"
// @Success 200 {object} Menu
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /menus/{id} [get]
func GetMenu(c *gin.Context) {

	id := c.Param("id")

	var menu Menu

	err := db.DB.QueryRow(
		"SELECT menu_id, dish_name, price, restaurant_id FROM Menu WHERE menu_id = $1",
		id,
	).Scan(&menu.ID, &menu.Nombre, &menu.Precio, &menu.RestaurantID)

	if err != nil {

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Menú no encontrado",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, menu)
}

// CreateMenu godoc
// @Summary Crear menú
// @Description Crea un nuevo menú
// @Tags menus
// @Accept json
// @Produce json
// @Param body body Menu true "Datos del menú"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /menus [post]
// Funcion del POST /menus
func CreateMenu(c *gin.Context) {

	var m Menu

	//Leer JSON
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "JSON inválido",
		})
		return
	}

	if m.Nombre == "" || m.Precio <= 0 || m.RestaurantID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos inválidos",
		})
		return
	}

	// Insertar en DB
	var menuID int
	err := db.DB.QueryRow(
		`INSERT INTO Menu (dish_name, price, restaurant_id)
		 VALUES ($1, $2, $3)
		 RETURNING menu_id`,
		m.Nombre, m.Precio, m.RestaurantID,
	).Scan(&menuID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Respuesta
	c.JSON(http.StatusCreated, gin.H{
		"message": "Menú creado exitosamente",
		"id":      menuID,
	})
}

// UpdateMenu godoc
// @Summary Actualizar menú
// @Description Actualiza un menú existente
// @Tags menus
// @Accept json
// @Produce json
// @Param id path int true "ID del menú"
// @Param body body Menu true "Datos actualizados"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /menus/{id} [put]
// PUT /menus/:id
func UpdateMenu(c *gin.Context) {

	id := c.Param("id")

	var input struct {
		Nombre       string  `json:"nombre"`
		Precio       float64 `json:"precio"`
		RestaurantID int     `json:"restaurant_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "JSON inválido",
		})
		return
	}

	if input.Nombre == "" || input.Precio <= 0 || input.RestaurantID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos inválidos",
		})
		return
	}

	query := `
		UPDATE Menu
		SET dish_name = $1, price = $2, restaurant_id = $3
		WHERE menu_id = $4
	`

	result, err := db.DB.Exec(query, input.Nombre, input.Precio, input.RestaurantID, id)

	if err != nil {

		if strings.Contains(err.Error(), "violates foreign key") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "restaurant_id no válido",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Menú no encontrado",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Menú actualizado correctamente",
	})
}

// DeleteMenu godoc
// @Summary Eliminar menú
// @Description Elimina un menú por ID
// @Tags menus
// @Produce json
// @Param id path int true "ID del menú"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /menus/{id} [delete]
// delete /menus/:id
func DeleteMenu(c *gin.Context) {

	id := c.Param("id")

	query := `
		DELETE FROM Menu
		WHERE menu_id = $1
	`

	result, err := db.DB.Exec(query, id)

	if err != nil {
		if strings.Contains(err.Error(), "violates foreign key") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No se puede eliminar el menú porque tiene datos asociados",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Menú no encontrado",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Menú eliminado correctamente",
	})
}
