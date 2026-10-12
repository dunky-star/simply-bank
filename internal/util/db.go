package util

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const dbDriver = "postgres"

// OpenDBFromEnv loads .env at envPath (if present), builds the Postgres URL
// from POSTGRES_* vars, and returns a pinged *sql.DB.
func OpenDBFromEnv(envPath string) *sql.DB {
	if err := godotenv.Load(envPath); err != nil {
		log.Println("util: .env not found, using real env vars:", err)
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
		log.Fatal("util: cannot open database: ", err)
	}
	if err := conn.Ping(); err != nil {
		log.Fatal("util: cannot reach database: ", err)
	}
	return conn
}

// NewDB returns an open *sql.DB for the app, loading .env from the working dir if present.
func NewDB() *sql.DB {
	return OpenDBFromEnv(".env")
}
