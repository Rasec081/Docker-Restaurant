package main

import (
	"log"
	"os"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "restaurant-backend/docs"
	"restaurant-backend/internal/db"
	"restaurant-backend/internal/handlers"
	"restaurant-backend/internal/middleware"
)

func main() {

	/*
		admin debe de tener los 2 roles: admin y client
	*/

	// =========================
	// 1. Cargar variables
	// =========================
	err := godotenv.Load()
	if err != nil {
		log.Println("No se encuentra el .env")
	}

	// =========================
	// 2. Conectar DB
	// =========================
	db.Connect()

	// =========================
	// 3. Router
	// =========================
	router := gin.Default()

	// =========================
	// 4. Public routes
	// =========================
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/auth/register", handlers.Register)
	router.POST("/auth/login", handlers.Login)

	// =========================
	// 5. Grupo protegido (JWT)
	// =========================
	protected := router.Group("/")

	// =========================
	// 6. Subgrupos por rol
	// =========================

	// CLIENT
	client := protected.Group("/")
	client.Use(middleware.RequireRole("client"))

	// ADMIN
	admin := protected.Group("/")
	admin.Use(middleware.RequireRole("admin"))

	// =========================
	// 7. Rutas CLIENT
	// =========================

	client.GET("/restaurants", handlers.GetRestaurants)
	client.GET("/menus/:id", handlers.GetMenu)
	client.GET("/orders/:id", handlers.GetOrder)
	client.GET("/users/me", handlers.GetUserMe)

	client.POST("/reservations", handlers.CreateReservation)
	client.POST("/orders", handlers.CreateOrder)

	// =========================
	// 8. Rutas ADMIN
	// =========================

	admin.POST("/restaurants", handlers.CreateRestaurant)
	admin.POST("/menus", handlers.CreateMenu)

	admin.PUT("/users/:id", handlers.UpdateUser)
	admin.PUT("/menus/:id", handlers.UpdateMenu)

	admin.DELETE("/users/:id", handlers.DeleteUser)
	admin.DELETE("/menus/:id", handlers.DeleteMenu)
	admin.DELETE("/reservations/:id", handlers.DeleteReservation)

	// =========================
	// 9. Puerto
	// =========================
	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server corriendo en puerto", port)
	router.Run(":" + port)
}
