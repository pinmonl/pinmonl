package database

import (
	"testing"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func TestSqlite3Driver(t *testing.T) {
	db, err := Open("sqlite3", ":memory:")
	if err != nil {
		t.Error(err)
	}

	err = db.Ping()
	if err != nil {
		t.Error(err)
	}
}
