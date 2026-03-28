package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB
var Database DBInterface

type DBInterface interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func Connect() {

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	Database = DB

	var err error

	// Intentar conectar varias veces
	for i := 1; i <= 10; i++ {

		DB, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Println("Error abriendo conexión:", err)
		}

		err = DB.Ping()
		if err == nil {
			fmt.Println("Conectado a PostgreSQL")
			return
		}

		log.Printf("Intento %d: DB no disponible, reintentando en 2s...\n", i)
		time.Sleep(2 * time.Second)
	}

	log.Fatal("No se pudo conectar a la DB después de varios intentos")
}
