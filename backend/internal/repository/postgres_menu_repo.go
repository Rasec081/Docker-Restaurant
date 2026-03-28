package repository

import (
	"database/sql"
	"errors"
	"restaurant-backend/internal/models"
)

type PostgresMenuRepository struct {
	DB *sql.DB
}

func NewPostgresMenuRepository(db *sql.DB) *PostgresMenuRepository {
	return &PostgresMenuRepository{DB: db}
}

func (r *PostgresMenuRepository) GetByID(id int) (*models.Menu, error) {
	var menu models.Menu
	err := r.DB.QueryRow(
		"SELECT menu_id, dish_name, price, restaurant_id FROM Menu WHERE menu_id = $1",
		id,
	).Scan(&menu.ID, &menu.Nombre, &menu.Precio, &menu.RestaurantID)

	if err == sql.ErrNoRows {
		return nil, errors.New("menú no encontrado")
	}
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *PostgresMenuRepository) Create(menu *models.Menu) error {
	err := r.DB.QueryRow(
		"INSERT INTO Menu (dish_name, price, restaurant_id) VALUES ($1, $2, $3) RETURNING menu_id",
		menu.Nombre, menu.Precio, menu.RestaurantID,
	).Scan(&menu.ID)
	return err
}

func (r *PostgresMenuRepository) Update(menu *models.Menu) error {
	result, err := r.DB.Exec(
		"UPDATE Menu SET dish_name = $1, price = $2, restaurant_id = $3 WHERE menu_id = $4",
		menu.Nombre, menu.Precio, menu.RestaurantID, menu.ID,
	)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("menú no encontrado")
	}
	return nil
}

func (r *PostgresMenuRepository) Delete(id int) error {
	result, err := r.DB.Exec("DELETE FROM Menu WHERE menu_id = $1", id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("menú no encontrado")
	}
	return nil
}
