package repository

import (
	"database/sql"
	"errors"
	"restaurant-backend/internal/models"
)

type PostgresReservationRepository struct {
	DB *sql.DB
}

func NewPostgresReservationRepository(db *sql.DB) *PostgresReservationRepository {
	return &PostgresReservationRepository{DB: db}
}

func (r *PostgresReservationRepository) Create(reservation *models.Reservation) error {
	err := r.DB.QueryRow(
		"INSERT INTO Reservation (table_id, client_id, fecha, estado) VALUES ($1, $2, $3, $4) RETURNING reservation_id",
		reservation.TableID, reservation.ClientID, reservation.Fecha, reservation.Estado,
	).Scan(&reservation.ID)
	return err
}

func (r *PostgresReservationRepository) Delete(id int) error {
	result, err := r.DB.Exec("DELETE FROM Reservation WHERE reservation_id = $1", id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("reserva no encontrada")
	}
	return nil
}

func (r *PostgresReservationRepository) GetByID(id int) (*models.Reservation, error) {
	var reservation models.Reservation
	err := r.DB.QueryRow(
		"SELECT reservation_id, table_id, client_id, fecha, estado FROM Reservation WHERE reservation_id = $1",
		id,
	).Scan(&reservation.ID, &reservation.TableID, &reservation.ClientID, &reservation.Fecha, &reservation.Estado)

	if err == sql.ErrNoRows {
		return nil, errors.New("reserva no encontrada")
	}
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}
