package repository

import (
	"database/sql"
	"errors"
	"restaurant-backend/internal/models"
)

type PostgresRestaurantRepository struct {
	DB *sql.DB
}

func NewPostgresRestaurantRepository(db *sql.DB) *PostgresRestaurantRepository {
	return &PostgresRestaurantRepository{DB: db}
}

func (r *PostgresRestaurantRepository) GetAll() ([]models.Restaurant, error) {
	rows, err := r.DB.Query("SELECT restaurant_id, nombre, estado, admin_id FROM Restaurant")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var restaurants []models.Restaurant
	for rows.Next() {
		var r models.Restaurant
		err := rows.Scan(&r.ID, &r.Nombre, &r.Estado, &r.AdminID)
		if err != nil {
			return nil, err
		}
		restaurants = append(restaurants, r)
	}
	return restaurants, nil
}

func (r *PostgresRestaurantRepository) GetByID(id int) (*models.Restaurant, error) {
	var restaurant models.Restaurant
	err := r.DB.QueryRow(
		"SELECT restaurant_id, nombre, estado, admin_id FROM Restaurant WHERE restaurant_id = $1",
		id,
	).Scan(&restaurant.ID, &restaurant.Nombre, &restaurant.Estado, &restaurant.AdminID)

	if err == sql.ErrNoRows {
		return nil, errors.New("restaurante no encontrado")
	}
	if err != nil {
		return nil, err
	}
	return &restaurant, nil
}

func (r *PostgresRestaurantRepository) Create(restaurant *models.Restaurant) error {
	err := r.DB.QueryRow(
		"INSERT INTO Restaurant (nombre, admin_id, estado) VALUES ($1, $2, $3) RETURNING restaurant_id",
		restaurant.Nombre, restaurant.AdminID, restaurant.Estado,
	).Scan(&restaurant.ID)
	return err
}

func (r *PostgresRestaurantRepository) Update(restaurant *models.Restaurant) error {
	result, err := r.DB.Exec(
		"UPDATE Restaurant SET nombre = $1, admin_id = $2, estado = $3 WHERE restaurant_id = $4",
		restaurant.Nombre, restaurant.AdminID, restaurant.Estado, restaurant.ID,
	)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("restaurante no encontrado")
	}
	return nil
}

func (r *PostgresRestaurantRepository) Delete(id int) error {
	result, err := r.DB.Exec("DELETE FROM Restaurant WHERE restaurant_id = $1", id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("restaurante no encontrado")
	}
	return nil
}
