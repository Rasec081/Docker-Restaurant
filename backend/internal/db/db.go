//comentario de la base de datos

package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

//variables globales

var DB *sql.DB

func Connect() {
	// es el string de la onexion
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error

	//abrimos la conexion

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error abriendo la conexion a DB:", err)
	}

	//verificamos la conexion con un ping
	err = DB.Ping()
	if err != nil {
		log.Fatal("DB no responde:", err)
	}

	fmt.Println("✅ Conectado a PostgreSQL")
}
