package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

var (
	migrationTableName = "migrations"
)

// DB stores sql.DB and the driver name
type DB struct {
	*sqlx.DB
}

// Open creates a db instance
func Open(driverName, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	dbx := sqlx.NewDb(db, driverName)
	return &DB{dbx}, nil
}

// Execer extends sqlx.Execer
type Execer interface {
	sqlx.Execer
	NamedExec(query string, arg interface{}) (sql.Result, error)
}

// Queryer extends sqlx.Queryer
type Queryer interface {
	sqlx.Queryer
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
}

// Ext combines Execer and Queryer
type Ext interface {
	Execer
	Queryer
	DriverName() string
	Rebind(string) string
	BindNamed(string, interface{}) (string, []interface{}, error)
}
