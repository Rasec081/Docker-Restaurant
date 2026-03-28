package handlers

import (
	"net/http"

	"restaurant-backend/internal/db"
	"restaurant-backend/internal/repository"
	"restaurant-backend/internal/services"
)

func ensureMenuHandler() {
	if menuHandler != nil {
		return
	}
	InitMenuHandler(repository.NewPostgresMenuRepository(db.DB))
}

func ensureOrderHandler() {
	if orderHandler != nil {
		return
	}
	InitOrderHandler(repository.NewPostgresOrderRepository(db.DB))
}

func ensureReservationHandler() {
	if reservationHandler != nil {
		return
	}
	InitReservationHandler(repository.NewPostgresReservationRepository(db.DB))
}

func ensureRestaurantHandler() {
	if restaurantHandler != nil {
		return
	}
	InitRestaurantHandler(repository.NewPostgresRestaurantRepository(db.DB))
}

func ensureUserHandler() {
	if userHandler != nil {
		return
	}
	InitUserHandler(repository.NewPostgresUserRepository(db.DB))
}

func ensureAuthHandler() {
	if authHandler != nil {
		return
	}
	InitAuthHandler(
		repository.NewPostgresUserRepository(db.DB),
		services.NewDefaultKeycloakService(),
		http.DefaultClient,
	)
}
