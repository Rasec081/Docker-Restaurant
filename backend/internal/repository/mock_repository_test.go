package repository

import (
	"errors"
	"testing"

	"restaurant-backend/internal/models"
)

var errTest = errors.New("test error")

func TestMockUserRepository(t *testing.T) {
	repo := NewMockUserRepository()

	user := &models.User{Nombre: "u1", RoleID: 1}
	if err := repo.Create(user); err != nil {
		t.Fatalf("Create error: %v", err)
	}

	if _, err := repo.GetByID(user.ID); err != nil {
		t.Fatalf("GetByID error: %v", err)
	}
	if _, err := repo.GetByID(999); err == nil {
		t.Fatal("esperaba error en GetByID")
	}
	if _, err := repo.GetByUsername("u1"); err != nil {
		t.Fatalf("GetByUsername error: %v", err)
	}

	user.Nombre = "u2"
	if err := repo.Update(user); err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if err := repo.Update(&models.User{ID: 999, Nombre: "x", RoleID: 1}); err == nil {
		t.Fatal("esperaba error en Update")
	}

	if err := repo.Delete(user.ID); err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if err := repo.Delete(999); err == nil {
		t.Fatal("esperaba error en Delete")
	}

	repo.CreateError = errTest
	if err := repo.Create(&models.User{Nombre: "x", RoleID: 1}); err == nil {
		t.Fatal("esperaba error en Create")
	}
	repo.CreateError = nil

	repo.UpdateError = errTest
	if err := repo.Update(&models.User{ID: 1, Nombre: "x", RoleID: 1}); err == nil {
		t.Fatal("esperaba error en Update")
	}
	repo.UpdateError = nil

	repo.DeleteError = errTest
	if err := repo.Delete(1); err == nil {
		t.Fatal("esperaba error en Delete")
	}
}

func TestMockMenuRepository(t *testing.T) {
	repo := NewMockMenuRepository()

	menu := &models.Menu{Nombre: "Dish", Precio: 10, RestaurantID: 1}
	if err := repo.Create(menu); err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if _, err := repo.GetByID(menu.ID); err != nil {
		t.Fatalf("GetByID error: %v", err)
	}
	if _, err := repo.GetByID(999); err == nil {
		t.Fatal("esperaba error en GetByID")
	}

	menu.Nombre = "Dish2"
	if err := repo.Update(menu); err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if err := repo.Update(&models.Menu{ID: 999, Nombre: "x", Precio: 1, RestaurantID: 1}); err == nil {
		t.Fatal("esperaba error en Update")
	}

	if err := repo.Delete(menu.ID); err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if err := repo.Delete(999); err == nil {
		t.Fatal("esperaba error en Delete")
	}

	repo.CreateError = errTest
	if err := repo.Create(&models.Menu{Nombre: "x", Precio: 1, RestaurantID: 1}); err == nil {
		t.Fatal("esperaba error en Create")
	}
	repo.CreateError = nil

	repo.UpdateError = errTest
	if err := repo.Update(&models.Menu{ID: 1, Nombre: "x", Precio: 1, RestaurantID: 1}); err == nil {
		t.Fatal("esperaba error en Update")
	}
	repo.UpdateError = nil

	repo.DeleteError = errTest
	if err := repo.Delete(1); err == nil {
		t.Fatal("esperaba error en Delete")
	}
}

func TestMockRestaurantRepository(t *testing.T) {
	repo := NewMockRestaurantRepository()

	repo.GetAllError = errTest
	if _, err := repo.GetAll(); err == nil {
		t.Fatal("esperaba error en GetAll")
	}
	repo.GetAllError = nil

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
	if _, err := repo.GetByID(999); err == nil {
		t.Fatal("esperaba error en GetByID")
	}

	rest.Nombre = "Rest2"
	if err := repo.Update(rest); err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if err := repo.Update(&models.Restaurant{ID: 999, Nombre: "x", AdminID: 1, Estado: 1}); err == nil {
		t.Fatal("esperaba error en Update")
	}

	if err := repo.Delete(rest.ID); err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if err := repo.Delete(999); err == nil {
		t.Fatal("esperaba error en Delete")
	}

	repo.CreateError = errTest
	if err := repo.Create(&models.Restaurant{Nombre: "x", AdminID: 1, Estado: 1}); err == nil {
		t.Fatal("esperaba error en Create")
	}
	repo.CreateError = nil

	repo.UpdateError = errTest
	if err := repo.Update(&models.Restaurant{ID: 1, Nombre: "x", AdminID: 1, Estado: 1}); err == nil {
		t.Fatal("esperaba error en Update")
	}
	repo.UpdateError = nil

	repo.DeleteError = errTest
	if err := repo.Delete(1); err == nil {
		t.Fatal("esperaba error en Delete")
	}
}

func TestMockReservationRepository(t *testing.T) {
	repo := NewMockReservationRepository()

	reservation := &models.Reservation{TableID: 1, ClientID: 2, Fecha: "2026-03-27 19:00:00", Estado: 1}
	if err := repo.Create(reservation); err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if _, err := repo.GetByID(reservation.ID); err != nil {
		t.Fatalf("GetByID error: %v", err)
	}
	if _, err := repo.GetByID(999); err == nil {
		t.Fatal("esperaba error en GetByID")
	}

	if err := repo.Delete(reservation.ID); err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if err := repo.Delete(999); err == nil {
		t.Fatal("esperaba error en Delete")
	}

	repo.CreateError = errTest
	if err := repo.Create(&models.Reservation{TableID: 1, ClientID: 2, Fecha: "2026-03-27 19:00:00", Estado: 1}); err == nil {
		t.Fatal("esperaba error en Create")
	}
	repo.CreateError = nil

	repo.DeleteError = errTest
	if err := repo.Delete(1); err == nil {
		t.Fatal("esperaba error en Delete")
	}
}

func TestMockOrderRepository(t *testing.T) {
	repo := NewMockOrderRepository()

	order := &models.Order{ClientID: 1, OrdersType: "Dine in", RestaurantID: 1}
	if err := repo.Create(order); err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if _, err := repo.GetByID(order.ID); err != nil {
		t.Fatalf("GetByID error: %v", err)
	}
	if _, err := repo.GetByID(999); err == nil {
		t.Fatal("esperaba error en GetByID")
	}

	repo.CreateError = errTest
	if err := repo.Create(&models.Order{ClientID: 1, OrdersType: "Dine in", RestaurantID: 1}); err == nil {
		t.Fatal("esperaba error en Create")
	}
}
