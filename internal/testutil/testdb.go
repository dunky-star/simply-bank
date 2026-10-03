// Package testutil provides shared helpers for database-backed tests.
package testutil

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // registers the "postgres" driver
)

const dbDriver = "postgres"

// NewTestDB returns an open *sql.DB for tests, loading .env if present.
func NewTestDB() *sql.DB {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("testutil: .env not found, using real env vars:", err)
	}

	dbSource := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		url.QueryEscape(os.Getenv("POSTGRES_PASSWORD")),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("testutil: cannot open database: ", err)
	}
	if err := conn.Ping(); err != nil {
		log.Fatal("testutil: cannot reach database: ", err)
	}
	return conn
}
