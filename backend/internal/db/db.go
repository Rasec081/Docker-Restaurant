package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

type dbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func getEnvOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func buildConnStr(cfg dbConfig) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
	)
}

func connectWith(
	connStr string,
	maxRetries int,
	sleepFn func(time.Duration),
	openFn func(string, string) (*sql.DB, error),
	pingFn func(*sql.DB) error,
) (*sql.DB, error) {
	var lastErr error
	for i := 1; i <= maxRetries; i++ {
		dbConn, err := openFn("postgres", connStr)
		if err != nil {
			lastErr = err
			log.Printf("Error abriendo conexión (intento %d): %v", i, err)
			sleepFn(2 * time.Second)
			continue
		}

		err = pingFn(dbConn)
		if err == nil {
			return dbConn, nil
		}
		lastErr = err
		log.Printf("Intento %d: DB no disponible (%v), reintentando en 2s...", i, err)
		sleepFn(2 * time.Second)
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("db no disponible")
	}
	return nil, lastErr
}

func Connect() {
	// Determinar qué archivo .env cargar
	envFile := ".env"
	if os.Getenv("GO_ENV") == "test" {
		envFile = ".env.test"
	}

	err := godotenv.Load(envFile)
	if err != nil {
		log.Printf("No se pudo cargar %s, usando variables de entorno", envFile)
	}

	config := dbConfig{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvOrDefault("DB_PORT", "5432"),
		User:     getEnvOrDefault("DB_USER", "postgres"),
		Password: getEnvOrDefault("DB_PASSWORD", "postgres"),
		Name:     getEnvOrDefault("DB_NAME", "restaurantdb"),
	}

	connStr := buildConnStr(config)

	log.Printf("Conectando a DB: host=%s port=%s dbname=%s", config.Host, config.Port, config.Name)

	DB, err = connectWith(connStr, 10, time.Sleep, sql.Open, func(dbConn *sql.DB) error {
		return dbConn.Ping()
	})
	if err == nil {
		log.Println("✅ Conectado a PostgreSQL exitosamente")
		return
	}

	log.Fatal("❌ No se pudo conectar a la BD después de varios intentos")

}

// GetDB retorna la conexión a la BD (para uso en repositorios)
func GetDB() *sql.DB {
	return DB
}

/*package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

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
*/
