package util

import "database/sql"

// NewTestDB returns an open *sql.DB for tests, loading .env from the project root.
func NewTestDB() *sql.DB {
	return OpenDBFromEnv("../../.env")
}
