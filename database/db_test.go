package database

import (
	"os"
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

func TestPostgresDriver(t *testing.T) {
	dsn := os.Getenv("POSTGRES_DSN")
	db, err := Open("postgres", dsn)
	if err != nil {
		t.Error(err)
	}

	err = db.Ping()
	if err != nil {
		t.Error(err)
	}
}
