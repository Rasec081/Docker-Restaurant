package repository

import (
	"database/sql"
	"errors"
	"restaurant-backend/internal/models"
)

type PostgresUserRepository struct {
	DB *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{DB: db}
}

func (r *PostgresUserRepository) GetByID(id int) (*models.User, error) {
	var user models.User
	err := r.DB.QueryRow(
		"SELECT user_id, nombre, role_id FROM Users WHERE user_id = $1",
		id,
	).Scan(&user.ID, &user.Nombre, &user.RoleID)

	if err == sql.ErrNoRows {
		return nil, errors.New("usuario no encontrado")
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.DB.QueryRow(
		"SELECT user_id, nombre, role_id FROM Users WHERE nombre = $1",
		username,
	).Scan(&user.ID, &user.Nombre, &user.RoleID)

	if err == sql.ErrNoRows {
		return nil, errors.New("usuario no encontrado")
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) Create(user *models.User) error {
	err := r.DB.QueryRow(
		"INSERT INTO Users (nombre, role_id) VALUES ($1, $2) RETURNING user_id",
		user.Nombre, user.RoleID,
	).Scan(&user.ID)
	return err
}

func (r *PostgresUserRepository) Update(user *models.User) error {
	result, err := r.DB.Exec(
		"UPDATE Users SET nombre = $1, role_id = $2 WHERE user_id = $3",
		user.Nombre, user.RoleID, user.ID,
	)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("usuario no encontrado")
	}
	return nil
}

func (r *PostgresUserRepository) Delete(id int) error {
	result, err := r.DB.Exec("DELETE FROM Users WHERE user_id = $1", id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("usuario no encontrado")
	}
	return nil
}
