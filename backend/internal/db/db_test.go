package db

import (
	"database/sql"
	"errors"
	"testing"
	"time"
)

func TestConnectAndGetDB(t *testing.T) {
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DB_NAME", "")
	t.Setenv("GO_ENV", "test")
	Connect()
	if DB == nil {
		t.Fatal("DB no inicializada")
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	t.Setenv("DB_HOST", "custom")
	if got := getEnvOrDefault("DB_HOST", "fallback"); got != "custom" {
		t.Fatalf("esperaba custom, obtuvo %s", got)
	}

	t.Setenv("DB_HOST", "")
	if got := getEnvOrDefault("DB_HOST", "fallback"); got != "fallback" {
		t.Fatalf("esperaba fallback, obtuvo %s", got)
	}
}

func TestBuildConnStr(t *testing.T) {
	cfg := dbConfig{Host: "h", Port: "p", User: "u", Password: "pw", Name: "db"}
	connStr := buildConnStr(cfg)
	if connStr == "" {
		t.Fatal("connStr vacio")
	}
}

func TestConnectWithSuccess(t *testing.T) {
	openFn := func(_ string, _ string) (*sql.DB, error) {
		return &sql.DB{}, nil
	}
	pingFn := func(_ *sql.DB) error { return nil }
	sleepFn := func(time.Duration) {}

	if _, err := connectWith("conn", 1, sleepFn, openFn, pingFn); err != nil {
		t.Fatalf("no esperaba error, obtuvo %v", err)
	}
}

func TestConnectWithRetryThenSuccess(t *testing.T) {
	openFn := func(_ string, _ string) (*sql.DB, error) {
		return &sql.DB{}, nil
	}
	attempts := 0
	pingFn := func(_ *sql.DB) error {
		attempts++
		if attempts < 2 {
			return errors.New("ping error")
		}
		return nil
	}
	sleepFn := func(time.Duration) {}

	if _, err := connectWith("conn", 3, sleepFn, openFn, pingFn); err != nil {
		t.Fatalf("no esperaba error, obtuvo %v", err)
	}
}

func TestConnectWithOpenError(t *testing.T) {
	openFn := func(_ string, _ string) (*sql.DB, error) {
		return nil, errors.New("open error")
	}
	pingFn := func(_ *sql.DB) error { return nil }
	sleepFn := func(time.Duration) {}

	if _, err := connectWith("conn", 2, sleepFn, openFn, pingFn); err == nil {
		t.Fatal("esperaba error")
	}
}
