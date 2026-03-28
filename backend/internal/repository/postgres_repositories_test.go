package repository

import (
	"database/sql"
	"testing"
	"time"

	"restaurant-backend/internal/db"
	"restaurant-backend/internal/models"
)

var repoTestDB *sql.DB

func setupRepoTestDB(t *testing.T) {
	if repoTestDB == nil {
		t.Setenv("GO_ENV", "test")
		db.Connect()
		repoTestDB = db.DB
	}
	cleanupRepoTestData(t)
	t.Cleanup(func() { cleanupRepoTestData(t) })
}

func cleanupRepoTestData(t *testing.T) {
	if repoTestDB == nil {
		return
	}
	repoTestDB.Exec("DELETE FROM Orders_Items")
	repoTestDB.Exec("DELETE FROM Orders")
	repoTestDB.Exec("DELETE FROM Reservation")
	repoTestDB.Exec("DELETE FROM Menu")
	repoTestDB.Exec("DELETE FROM Tables")
	repoTestDB.Exec("DELETE FROM Restaurant")
	repoTestDB.Exec("DELETE FROM Users")
	repoTestDB.Exec("DELETE FROM Roles")
}

func seedBaseData(t *testing.T) {
	repoTestDB.Exec("INSERT INTO Roles (role_id, nombre) VALUES (1, 'Admin'), (2, 'Client') ON CONFLICT DO NOTHING")
	repoTestDB.Exec("INSERT INTO Users (user_id, nombre, role_id) VALUES (1, 'admin', 1), (2, 'client', 2) ON CONFLICT DO NOTHING")
	repoTestDB.Exec("INSERT INTO Restaurant (restaurant_id, nombre, admin_id, estado) VALUES (1, 'Test Restaurant', 1, 1) ON CONFLICT DO NOTHING")
	repoTestDB.Exec("INSERT INTO Tables (table_id, table_number, estado, restaurant_id) VALUES (1, 1, 1, 1) ON CONFLICT DO NOTHING")
}

func TestPostgresUserRepositoryCRUD(t *testing.T) {
	setupRepoTestDB(t)
	seedBaseData(t)

	repo := NewPostgresUserRepository(repoTestDB)

	user := &models.User{Nombre: "user1", RoleID: 2}
	if err := repo.Create(user); err != nil {
		t.Fatalf("Create error: %v", err)
	}

	if _, err := repo.GetByID(user.ID); err != nil {
		t.Fatalf("GetByID error: %v", err)
	}
	if _, err := repo.GetByID(9999); err == nil {
		t.Fatal("esperaba error al buscar no existente")
	}
	if _, err := repo.GetByUsername("user1"); err != nil {
		t.Fatalf("GetByUsername error: %v", err)
	}
	if _, err := repo.GetByUsername("missing"); err == nil {
		t.Fatal("esperaba error en GetByUsername")
	}

	user.Nombre = "user2"
	user.RoleID = 1
	if err := repo.Update(user); err != nil {
		t.Fatalf("Update error: %v", err)
	}

	if err := repo.Delete(user.ID); err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if err := repo.Delete(9999); err == nil {
		t.Fatal("esperaba error al borrar no existente")
	}
}

func TestPostgresMenuRepositoryCRUD(t *testing.T) {
	setupRepoTestDB(t)
	seedBaseData(t)

	repo := NewPostgresMenuRepository(repoTestDB)

	menu := &models.Menu{Nombre: "Dish", Precio: 12.5, RestaurantID: 1}
	if err := repo.Create(menu); err != nil {
		t.Fatalf("Create error: %v", err)
	}

	if _, err := repo.GetByID(menu.ID); err != nil {
		t.Fatalf("GetByID error: %v", err)
	}

	menu.Nombre = "Dish2"
	menu.Precio = 15
	if err := repo.Update(menu); err != nil {
		t.Fatalf("Update error: %v", err)
	}

	if err := repo.Update(&models.Menu{ID: 9999, Nombre: "X", Precio: 1, RestaurantID: 1}); err == nil {
		t.Fatal("esperaba error al actualizar no existente")
	}

	if err := repo.Delete(menu.ID); err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if err := repo.Delete(9999); err == nil {
		t.Fatal("esperaba error al borrar no existente")
	}
}

func TestPostgresRestaurantRepositoryCRUD(t *testing.T) {
	setupRepoTestDB(t)
	seedBaseData(t)

	repo := NewPostgresRestaurantRepository(repoTestDB)

	rest := &models.Restaurant{Nombre: "Rest", AdminID: 1, Estado: 1}
	if err := repo.Create(rest); err != nil {
		t.Fatalf("Create error: %v", err)
	}

	if _, err := repo.GetAll(); err != nil {
		t.Fatalf("GetAll error: %v", err)
	}

	if _, err := repo.GetByID(rest.ID); err != nil {
		t.Fatalf("GetByID error: %v", err)
	}
	if _, err := repo.GetByID(9999); err == nil {
		t.Fatal("esperaba error en GetByID")
	}

	rest.Nombre = "Rest2"
	if err := repo.Update(rest); err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if err := repo.Update(&models.Restaurant{ID: 9999, Nombre: "X", AdminID: 1, Estado: 1}); err == nil {
		t.Fatal("esperaba error al actualizar no existente")
	}

	if err := repo.Delete(rest.ID); err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if err := repo.Delete(9999); err == nil {
		t.Fatal("esperaba error al borrar no existente")
	}
}

func TestPostgresReservationRepositoryCRUD(t *testing.T) {
	setupRepoTestDB(t)
	seedBaseData(t)

	repo := NewPostgresReservationRepository(repoTestDB)

	reservation := &models.Reservation{TableID: 1, ClientID: 2, Fecha: time.Now().Format("2006-01-02 15:04:05"), Estado: 1}
	if err := repo.Create(reservation); err != nil {
		t.Fatalf("Create error: %v", err)
	}

	if _, err := repo.GetByID(reservation.ID); err != nil {
		t.Fatalf("GetByID error: %v", err)
	}

	if err := repo.Delete(reservation.ID); err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if err := repo.Delete(9999); err == nil {
		t.Fatal("esperaba error al borrar no existente")
	}
}

func TestPostgresOrderRepositoryCRUD(t *testing.T) {
	setupRepoTestDB(t)
	seedBaseData(t)

	repo := NewPostgresOrderRepository(repoTestDB)

	order := &models.Order{TableID: nil, ClientID: 2, OrdersType: "Dine in", RestaurantID: 1}
	if err := repo.Create(order); err != nil {
		t.Fatalf("Create error: %v", err)
	}

	if _, err := repo.GetByID(order.ID); err != nil {
		t.Fatalf("GetByID error: %v", err)
	}
	if _, err := repo.GetByID(9999); err == nil {
		t.Fatal("esperaba error al buscar no existente")
	}
}

func TestDBMenuRepositoryCRUD(t *testing.T) {
	setupRepoTestDB(t)
	seedBaseData(t)

	repo := &DBMenuRepository{DB: repoTestDB}

	menu := &models.Menu{Nombre: "Dish", Precio: 12.5, RestaurantID: 1}
	if err := repo.Create(menu); err != nil {
		t.Fatalf("Create error: %v", err)
	}

	if _, err := repo.GetByID(menu.ID); err != nil {
		t.Fatalf("GetByID error: %v", err)
	}
	if _, err := repo.GetByID(9999); err == nil {
		t.Fatal("esperaba error al buscar no existente")
	}

	menu.Nombre = "Dish2"
	if err := repo.Update(menu); err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if err := repo.Update(&models.Menu{ID: 9999, Nombre: "X", Precio: 1, RestaurantID: 1}); err == nil {
		t.Fatal("esperaba error al actualizar no existente")
	}

	if err := repo.Delete(menu.ID); err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if err := repo.Delete(9999); err == nil {
		t.Fatal("esperaba error al borrar no existente")
	}
}
