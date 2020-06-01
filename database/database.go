package database

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/golang-migrate/migrate/v4"
)

type DB struct {
	*sql.DB
	Migrate *migrate.Migrate
	Driver  string
	Builder StatementBuilder
}

type DBOpts struct {
	Driver string
	DSN    string
}

func NewDB(opts *DBOpts) (*DB, error) {
	db, err := Open(opts.Driver, opts.DSN)
	if err != nil {
		return nil, err
	}

	m, err := NewMigrate(opts.Driver, opts.DSN)
	if err != nil {
		return nil, err
	}

	return &DB{
		DB:      db,
		Migrate: m,
		Driver:  opts.Driver,
		Builder: StatementBuilder(squirrel.StatementBuilder),
	}, nil
}

type StatementBuilder squirrel.StatementBuilderType

func Open(driver, dsn string) (*sql.DB, error) {
	return sql.Open(driver, dsn)
}

func NewMigrate(driver, dsn string) (*migrate.Migrate, error) {
	srcURL := "pkger://" + getSourceURL(driver)
	dsnURL := driver + "://" + dsn
	return migrate.New(srcURL, dsnURL)
}

func getSourceURL(driver string) string {
	return "/migrations/" + driver
}

type Queryer interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type Execer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type Runner interface {
	Queryer
	Execer
}
