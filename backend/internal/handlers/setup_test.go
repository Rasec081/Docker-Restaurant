package handlers

import (
	"database/sql"
	"os"
	"testing"

	"restaurant-backend/internal/db"

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
}
