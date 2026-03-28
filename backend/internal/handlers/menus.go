package handlers

import (
	"net/http"
	"strconv"

	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

// MenuHandler maneja las operaciones de menús
type MenuHandler struct {
	MenuRepo repository.MenuRepository
}

var menuHandler *MenuHandler

// NewMenuHandler crea una nueva instancia de MenuHandler
func NewMenuHandler(menuRepo repository.MenuRepository) *MenuHandler {
	return &MenuHandler{MenuRepo: menuRepo}
}

// InitMenuHandler inicializa el handler global (para rutas y tests)
func InitMenuHandler(menuRepo repository.MenuRepository) {
	menuHandler = NewMenuHandler(menuRepo)
}

func getMenuHandler() *MenuHandler {
	ensureMenuHandler()
	return menuHandler
}

// GetMenu expone el handler global
func GetMenu(c *gin.Context) {
	getMenuHandler().GetMenu(c)
}

// GetMenu godoc
// @Summary Obtener menú por ID
// @Description Retorna un menú específico
// @Tags menus
// @Produce json
// @Param id path int true "ID del menú"
// @Success 200 {object} models.Menu
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /menus/{id} [get]
func (h *MenuHandler) GetMenu(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	menu, err := h.MenuRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menú no encontrado"})
		return
	}

	c.JSON(http.StatusOK, menu)
}

// CreateMenu expone el handler global
func CreateMenu(c *gin.Context) {
	getMenuHandler().CreateMenu(c)
}

// CreateMenu godoc
// @Summary Crear menú
// @Description Crea un nuevo menú
// @Tags menus
// @Accept json
// @Produce json
// @Param body body models.Menu true "Datos del menú"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /menus [post]
func (h *MenuHandler) CreateMenu(c *gin.Context) {
	var menu models.Menu

	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	if menu.Nombre == "" || menu.Precio <= 0 || menu.RestaurantID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	if err := h.MenuRepo.Create(&menu); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Menú creado exitosamente",
		"id":      menu.ID,
	})
}

// UpdateMenu expone el handler global
func UpdateMenu(c *gin.Context) {
	getMenuHandler().UpdateMenu(c)
}

// UpdateMenu godoc
// @Summary Actualizar menú
// @Description Actualiza un menú existente
// @Tags menus
// @Accept json
// @Produce json
// @Param id path int true "ID del menú"
// @Param body body models.Menu true "Datos actualizados"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /menus/{id} [put]
func (h *MenuHandler) UpdateMenu(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var menu models.Menu
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	menu.ID = id
	if err := h.MenuRepo.Update(&menu); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menú actualizado correctamente"})
}

// DeleteMenu expone el handler global
func DeleteMenu(c *gin.Context) {
	getMenuHandler().DeleteMenu(c)
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
// @Router /menus/{id} [delete]
func (h *MenuHandler) DeleteMenu(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.MenuRepo.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menú eliminado correctamente"})
}
