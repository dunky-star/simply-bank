package db

import (
	"os"
	"testing"

	"github.com/dunky-star/simply-bank/internal/util"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	db := util.NewTestDB()
	testQueries = New(db)

	code := m.Run()
	_ = db.Close()
	os.Exit(code)
}
