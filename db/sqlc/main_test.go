package db

import (
	"os"
	"testing"

	"github.com/dunky-star/simply-bank/internal/util"
)

var testQueries *Queries
var testStore *Store
var testDB = util.NewTestDB()

func TestMain(m *testing.M) {
	testQueries = New(testDB)
	testStore = NewStore(testDB)

	code := m.Run()
	_ = testDB.Close()
	os.Exit(code)
}
