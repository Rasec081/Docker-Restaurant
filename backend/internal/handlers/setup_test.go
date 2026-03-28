package handlers

import (
	"database/sql"
	"net/http"
	"os"
	"testing"

	"restaurant-backend/internal/db"
	"restaurant-backend/internal/repository"
	"restaurant-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	// Cargar .env para tests
	godotenv.Load("../.env.test")

	// Configurar DB de prueba
	db.Connect()
	testDB = db.DB

	// Configurar modo test
	gin.SetMode(gin.TestMode)

	// Ejecutar tests
	code := m.Run()

	// Limpiar después
	cleanupTestData()

	os.Exit(code)
}

func cleanupTestData() {
	if testDB != nil {
		testDB.Exec("DELETE FROM Orders_Items")
		testDB.Exec("DELETE FROM Orders")
		testDB.Exec("DELETE FROM Reservation")
		testDB.Exec("DELETE FROM Menu")
		testDB.Exec("DELETE FROM Tables")
		testDB.Exec("DELETE FROM Restaurant")
		testDB.Exec("DELETE FROM Users")
		testDB.Exec("DELETE FROM Roles")
	}
}

func setupTestData(t *testing.T) {
	// Insertar roles
	testDB.Exec("INSERT INTO Roles (role_id, nombre) VALUES (1, 'Admin'), (2, 'Client') ON CONFLICT DO NOTHING")

	// Insertar usuarios
	testDB.Exec("INSERT INTO Users (user_id, nombre, role_id) VALUES (1, 'admin', 1), (2, 'client', 2) ON CONFLICT DO NOTHING")

	// Insertar restaurante
	testDB.Exec("INSERT INTO Restaurant (restaurant_id, nombre, admin_id, estado) VALUES (1, 'Test Restaurant', 1, 1) ON CONFLICT DO NOTHING")

	// Insertar menu
	testDB.Exec("INSERT INTO Menu (menu_id, dish_name, price, restaurant_id) VALUES (1, 'Test Dish', 10.99, 1) ON CONFLICT DO NOTHING")

	// Insertar tabla
	testDB.Exec("INSERT INTO Tables (table_id, table_number, estado, restaurant_id) VALUES (1, 1, 1, 1) ON CONFLICT DO NOTHING")

	// Sincronizar secuencias para inserts sin ID
	testDB.Exec("SELECT setval('roles_role_id_seq', COALESCE(MAX(role_id), 1), true) FROM Roles")
	testDB.Exec("SELECT setval('users_user_id_seq', COALESCE(MAX(user_id), 1), true) FROM Users")
	testDB.Exec("SELECT setval('restaurant_restaurant_id_seq', COALESCE(MAX(restaurant_id), 1), true) FROM Restaurant")
	testDB.Exec("SELECT setval('menu_menu_id_seq', COALESCE(MAX(menu_id), 1), true) FROM Menu")
	testDB.Exec("SELECT setval('tables_table_id_seq', COALESCE(MAX(table_id), 1), true) FROM Tables")
	testDB.Exec("SELECT setval('orders_orders_id_seq', COALESCE(MAX(orders_id), 1), true) FROM Orders")
	testDB.Exec("SELECT setval('reservation_reservation_id_seq', COALESCE(MAX(reservation_id), 1), true) FROM Reservation")
}

func initIntegrationHandlers() {
	menuRepo := repository.NewPostgresMenuRepository(testDB)
	restaurantRepo := repository.NewPostgresRestaurantRepository(testDB)
	userRepo := repository.NewPostgresUserRepository(testDB)
	reservationRepo := repository.NewPostgresReservationRepository(testDB)
	orderRepo := repository.NewPostgresOrderRepository(testDB)

	InitMenuHandler(menuRepo)
	InitRestaurantHandler(restaurantRepo)
	InitUserHandler(userRepo)
	InitReservationHandler(reservationRepo)
	InitOrderHandler(orderRepo)
	InitAuthHandler(userRepo, services.NewDefaultKeycloakService(), http.DefaultClient)
}
