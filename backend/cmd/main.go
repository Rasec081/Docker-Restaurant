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

	// protected := gin.Default()
	router := gin.Default()

	protected := router.Group("/")
	protected.Use(middleware.JWTMiddleware())

	// =========================
	// 4. Rutas básicas
	// =========================

	// health check
	protected.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 5. Rutas reales (handlers)
	protected.GET("/restaurants", handlers.GetRestaurants)

	protected.GET("/users/me", handlers.GetUserMe)

	protected.GET("/menus/:id", handlers.GetMenu)

	protected.GET("/orders/:id", handlers.GetOrder)

	// POSTs
	router.POST("/auth/register", handlers.Register)

	protected.POST("/restaurants", handlers.CreateRestaurant)

	protected.POST("/reservations", handlers.CreateReservation)

	protected.POST("/orders", handlers.CreateOrder)

	protected.POST("/menus", handlers.CreateMenu)

	// PUTs
	protected.PUT("/users/:id", handlers.UpdateUser)

	protected.PUT("/menus/:id", handlers.UpdateMenu)

	// DELETEs
	protected.DELETE("/reservations/:id", handlers.DeleteReservation)

	protected.DELETE("/menus/:id", handlers.DeleteMenu)

	protected.DELETE("/users/:id", handlers.DeleteUser)

	//auth
	router.POST("/login", handlers.Login)

	// 6. Puerto dinámico
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8080"
	}

	// 7. Iniciar servidor
	log.Println("Server corriendo en puerto", port)
	router.Run(":" + port)

	/*
		#################################################
		nota: si el usuario tiene una orden pendiente, no
		se puede eliminar, independinetemente de su rol
		#################################################
	*/

}
