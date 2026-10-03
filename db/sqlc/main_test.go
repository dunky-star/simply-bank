package db

import (
	"os"
	"testing"

	"github.com/dunky-star/simply-bank/internal/testutil"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn := testutil.NewTestDB()
	testQueries = New(conn)

	code := m.Run()
	_ = conn.Close()
	os.Exit(code)
}
