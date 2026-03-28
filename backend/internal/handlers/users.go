package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

// UserHandler maneja operaciones de usuarios
type UserHandler struct {
	UserRepo repository.UserRepository
}

var userHandler *UserHandler

// NewUserHandler crea una nueva instancia
func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{UserRepo: userRepo}
}

// InitUserHandler inicializa el handler global (para rutas y tests)
func InitUserHandler(userRepo repository.UserRepository) {
	userHandler = NewUserHandler(userRepo)
}

func getUserHandler() *UserHandler {
	ensureUserHandler()
	return userHandler
}

/*
funcion del get
*/
// GetUserMe godoc
// @Summary Obtener usuario actual
// @Description Retorna la información del usuario autenticado
// @Tags users
// @Produce json
// @Success 200 {object} models.User
// @Failure 500 {object} map[string]string
// @Router /users/me [get]
func GetUserMe(c *gin.Context) {
	getUserHandler().GetUserMe(c)
}

// GetUserMe godoc
func (h *UserHandler) GetUserMe(c *gin.Context) {
	usernameInterface, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "usuario no autenticado",
		})
		return
	}

	username, ok := usernameInterface.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error procesando username",
		})
		return
	}

	user, err := h.UserRepo.GetByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Usuario no encontrado en DB",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

/*
funcion del update
*/
// UpdateUser godoc
// @Summary Actualizar usuario
// @Description Actualiza los datos de un usuario
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Param body body models.User true "Datos del usuario"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {
	getUserHandler().UpdateUser(c)
}

// UpdateUser godoc
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID inválido",
		})
		return
	}

	var input struct {
		Nombre string `json:"nombre"`
		Rol    int    `json:"rol"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user := models.User{
		ID:     id,
		Nombre: input.Nombre,
		RoleID: input.Rol,
	}

	if err := h.UserRepo.Update(&user); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "no encontrado") || strings.Contains(strings.ToLower(err.Error()), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Usuario no encontrado",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Usuario actualizado correctamente",
	})
}

/*
funcion del delete
*/

/*
como probarlo:

1- haga un build del compose

2- nos metemos como si fuera sql con el comando: docker exec -it restaurant_db psql -U postgres -d restaurantdb

3- agregamos un usuario de prueba con el comando:
INSERT INTO Users (nombre, role_id)
VALUES ('Test Delete', 1);

4- verificamos que se agregó con el comando: SELECT * FROM Users;

5- desde thunder client ejecutamos el delete con la url: http://localhost:8080/users/2 segun el id

6- verificamos que se borró con el comando: SELECT * FROM Users;
*/
// DeleteUser godoc
// @Summary Eliminar usuario
// @Description Elimina un usuario por ID
// @Tags users
// @Produce json
// @Param id path int true "ID del usuario"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	getUserHandler().DeleteUser(c)
}

// DeleteUser godoc
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID inválido",
		})
		return
	}

	if err := h.UserRepo.Delete(id); err != nil {
		if strings.Contains(err.Error(), "violates foreign key") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No se puede eliminar el usuario porque tiene datos asociados",
			})
			return
		}
		if strings.Contains(strings.ToLower(err.Error()), "no encontrado") || strings.Contains(strings.ToLower(err.Error()), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Usuario no encontrado",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Usuario eliminado correctamente",
	})
}
