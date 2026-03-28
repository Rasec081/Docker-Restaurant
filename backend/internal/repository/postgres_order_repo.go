package repository

import (
	"database/sql"
	"errors"
	"restaurant-backend/internal/models"
)

type PostgresOrderRepository struct {
	DB *sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{DB: db}
}

func (r *PostgresOrderRepository) GetByID(id int) (*models.Order, error) {
	var order models.Order
	err := r.DB.QueryRow(
		"SELECT orders_id, table_id, client_id, orders_type, restaurant_id FROM Orders WHERE orders_id = $1",
		id,
	).Scan(&order.ID, &order.TableID, &order.ClientID, &order.OrdersType, &order.RestaurantID)

	if err == sql.ErrNoRows {
		return nil, errors.New("orden no encontrada")
	}
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *PostgresOrderRepository) Create(order *models.Order) error {
	err := r.DB.QueryRow(
		"INSERT INTO Orders (table_id, client_id, orders_type, restaurant_id) VALUES ($1, $2, $3, $4) RETURNING orders_id",
		order.TableID, order.ClientID, order.OrdersType, order.RestaurantID,
	).Scan(&order.ID)
	return err
}
