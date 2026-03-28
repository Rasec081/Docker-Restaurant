package handlers

import (
	"net/http"
	"strings"

	"restaurant-backend/internal/db"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
	Rol    int    `json:"rol"`
}

/*
funcion del get
*/
// GetUserMe godoc
// @Summary Obtener usuario actual
// @Description Retorna la información del usuario autenticado
// @Tags users
// @Produce json
// @Success 200 {object} User
// @Failure 500 {object} map[string]string
// @Router /users/me [get]
func GetUserMe(c *gin.Context) {

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "usuario no autenticado",
		})
		return
	}

	var user User

	err := db.Database.QueryRow(
		"SELECT user_id, nombre, role_id FROM Users WHERE nombre = $1",
		username,
	).Scan(&user.ID, &user.Nombre, &user.Rol)

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
// @Param body body User true "Datos del usuario"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id} [put]
func UpdateUser(c *gin.Context) {

	// sacamos la URL
	id := c.Param("id") //no se si el id va a venir del jwt o de la url, hay que definirlo

	//el body (JSON)
	var input struct {
		Nombre string `json:"nombre"`
		Rol    int    `json:"rol"`
	}

	//si hay error
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//emepzamos con el update
	query := `
		UPDATE Users
		SET nombre = $1, role_id = $2
		WHERE user_id = $3
	`

	result, err := db.Database.Exec(query, input.Nombre, input.Rol, id)

	//verificacionde si funkó
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//ocupamos saber que si existe ese usuario para dar la respuesta correcta
	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Usuario no encontrado",
		})
		return
	}

	//damos la repsues que si funciono
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

	id := c.Param("id")

	query := `
		DELETE FROM Users
		WHERE user_id = $1
	`

	result, err := db.Database.Exec(query, id)

	//Manejo de errores
	if err != nil {

		// Error de foreign key (como el que te salió)
		if strings.Contains(err.Error(), "violates foreign key") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No se puede eliminar el usuario porque tiene datos asociados",
			})
			return
		}

		// Error general
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Verificar si realmente borró algo
	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Usuario no encontrado",
		})
		return
	}

	//se logro
	c.JSON(http.StatusOK, gin.H{
		"message": "Usuario eliminado correctamente",
	})
}
