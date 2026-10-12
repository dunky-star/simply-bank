package util

import (
	"database/sql"
	"log"
)

// NewTestDB returns an open *sql.DB for tests, loading config from the project root.
// Tests run from db/sqlc, so the project root is ../..
func NewTestDB() *sql.DB {
	config, err := LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load test config: ", err)
	}
	conn, err := OpenDB(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to test db: ", err)
	}
	return conn
}
