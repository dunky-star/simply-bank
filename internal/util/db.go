package util

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// OpenDB opens and pings a *sql.DB for the given driver and datasource.
// It returns a wrapped error so callers can decide how to handle failure.
func OpenDB(driver, source string) (*sql.DB, error) {
	conn, err := sql.Open(driver, source)
	if err != nil {
		return nil, fmt.Errorf("open %s db: %w", driver, err)
	}
	if err := conn.Ping(); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("ping %s db: %w", driver, err)
	}
	return conn, nil
}

// NewDB returns an open, pinged *sql.DB built from the Viper-loaded config.
// Fatal here is intentional: the app cannot run without a database.
func NewDB(config Config) *sql.DB {
	conn, err := OpenDB(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	return conn
}
