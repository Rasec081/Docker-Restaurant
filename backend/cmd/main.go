package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"restaurant-backend/internal/db"
	"restaurant-backend/internal/handlers"
	"restaurant-backend/internal/middleware"
	//"restaurant-backend/internal/handlers"
)

// La funcion que va a correr todo
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No se encuentra el .env")
	}

	// Busca en db y usa la funcion
	db.Connect()

	// router := gin.Default()
	r := gin.Default()

	router := r.Group("/")
	router.Use(middleware.JWTMiddleware())

	// =========================
	// 4. Rutas básicas
	// =========================

	// health check
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 5. Rutas reales (handlers)
	router.GET("/restaurants", handlers.GetRestaurants)

	router.GET("/users/me", handlers.GetUserMe)

	router.PUT("/users/:id", handlers.UpdateUser)

	router.DELETE("/users/:id", handlers.DeleteUser)

	router.GET("/menus/:id", handlers.GetMenu)

	router.POST("/menus", handlers.CreateMenu)

	router.PUT("/menus/:id", handlers.UpdateMenu)

	router.DELETE("/menus/:id", handlers.DeleteMenu)

	router.GET("/orders/:id", handlers.GetOrder)

	// POSTs
	router.POST("/restaurants", handlers.CreateRestaurant)

	router.POST("/reservations", handlers.CreateReservation)

	router.POST("/orders", handlers.CreateOrder)

	// DELETEs
	router.DELETE("/reservations/:id", handlers.DeleteReservation)

	// 6. Puerto dinámico
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8080"
	}

	// 7. Iniciar servidor
	log.Println("Server corriendo en puerto", port)
	r.Run(":" + port)

	/*
		#################################################
		nota: si el usuario tiene una orden pendiente, no
		se puede eliminar, independinetemente de su rol
		#################################################
	*/

}
