package database

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/golang-migrate/migrate/v4"
)

type DB struct {
	*sql.DB
	Locker
	drv     string
	Migrate *migrate.Migrate
	Builder Builder
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
		Locker:  NewDriverLocker(opts.Driver),
		Migrate: m,
		drv:     opts.Driver,
		Builder: NewBuilderFromBase(squirrel.StatementBuilder),
	}, nil
}

func Open(driver, dsn string) (*sql.DB, error) {
	return sql.Open(driver, dsn)
}

func (d *DB) DriverName() string {
	return d.drv
}

func (d *DB) TxFunc(fn func(context.Context) error) error {
	tx, err := d.Begin()
	if err != nil {
		return err
	}
	ctx := context.TODO()
	ctx = WithRunner(ctx, tx)
	err = fn(ctx)
	if err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return err2
		}
		return err
	}
	return tx.Commit()
}

func (d *DB) Begin() (*Tx, error) {
	tx, err := d.DB.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{
		Tx:     tx,
		Locker: d.Locker,
	}, nil
}

func (d *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	d.Lock()
	rows, err := d.DB.Query(query, args...)
	d.Unlock()
	return rows, err
}

func (d *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	d.Lock()
	res, err := d.DB.Exec(query, args...)
	d.Unlock()
	return res, err
}

func NewMigrate(driver, dsn string) (*migrate.Migrate, error) {
	srcURL := "pkger://" + getSourceURL(driver)
	dsnURL := driver + "://" + dsn
	return migrate.New(srcURL, dsnURL)
}

func getSourceURL(driver string) string {
	return "/migrations/" + driver
}

type Tx struct {
	*sql.Tx
	Locker
}

func (t *Tx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	t.Lock()
	rows, err := t.Tx.Query(query, args...)
	t.Unlock()
	return rows, err
}

func (t *Tx) Exec(query string, args ...interface{}) (sql.Result, error) {
	t.Lock()
	res, err := t.Tx.Exec(query, args...)
	t.Unlock()
	return res, err
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

type RowScanner interface {
	Scan(dest ...interface{}) error
}

type Builder struct {
	squirrel.StatementBuilderType
}

func NewBuilderFromBase(base squirrel.StatementBuilderType) Builder {
	return Builder{base}
}

func (b Builder) RunWith(runner Runner) Builder {
	return NewBuilderFromBase(b.StatementBuilderType.RunWith(runner))
}
