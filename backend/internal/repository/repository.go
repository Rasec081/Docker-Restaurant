package repository

import (
	"restaurant-backend/internal/models"
)

// UserRepository define operaciones de usuarios
type UserRepository interface {
	GetByID(id int) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id int) error
}

// RestaurantRepository define operaciones de restaurantes
type RestaurantRepository interface {
	GetAll() ([]models.Restaurant, error)
	GetByID(id int) (*models.Restaurant, error)
	Create(restaurant *models.Restaurant) error
	Update(restaurant *models.Restaurant) error
	Delete(id int) error
}

// MenuRepository define operaciones de menús
type MenuRepository interface {
	GetByID(id int) (*models.Menu, error)
	Create(menu *models.Menu) error
	Update(menu *models.Menu) error
	Delete(id int) error
}

// ReservationRepository define operaciones de reservas
type ReservationRepository interface {
	Create(reservation *models.Reservation) error
	Delete(id int) error
	IsTableAvailable(tableID int, fecha string) (bool, error)
}

// OrderRepository define operaciones de órdenes
type OrderRepository interface {
	GetByID(id int) (*models.Order, error)
	Create(order *models.Order) error
}
