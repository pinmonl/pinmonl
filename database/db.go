package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// DB stores sql.DB and the driver name.
type DB struct {
	*sqlx.DB
}

// Open creates db.
func Open(driverName, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	dbx := sqlx.NewDb(db, driverName)
	return &DB{dbx}, nil
}

// Binder provides database bind var functions.
type Binder interface {
	Rebind(string) string
	BindNamed(string, interface{}) (string, []interface{}, error)
}

// Execer extends sqlx.Execer.
type Execer interface {
	sqlx.Execer
	Binder
	NamedExec(query string, arg interface{}) (sql.Result, error)
}

// Queryer extends sqlx.Queryer.
type Queryer interface {
	sqlx.Queryer
	Binder
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
}

// Ext combines Execer and Queryer.
type Ext interface {
	sqlx.Execer
	sqlx.Queryer
	Binder
	DriverName() string
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
}
