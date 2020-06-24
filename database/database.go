package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/golang-migrate/migrate/v4"
)

var (
	ErrNoTx = errors.New("Tx does not exist")
)

type DB struct {
	*sql.DB
	Locker
	driver  string
	Migrate *migrate.Migrate
	Builder Builder
}

func NewDB(driver, dsn string) (*DB, error) {
	db, err := Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	m, err := NewMigrate(driver, dsn)
	if err != nil {
		return nil, err
	}

	return &DB{
		DB:      db,
		Locker:  NewDriverLocker(driver),
		Migrate: m,
		driver:  driver,
		Builder: NewBuilderFromBase(squirrel.StatementBuilder),
	}, nil
}

func Open(driver, dsn string) (*sql.DB, error) {
	return sql.Open(driver, dsn)
}

func (d *DB) DriverName() string {
	return d.driver
}

func (d *DB) TxFunc(ctx context.Context, fn func(context.Context) bool) error {
	tx, err := d.Begin()
	if err != nil {
		return err
	}
	ok := fn(WithTx(ctx, tx))
	if !ok {
		return tx.Rollback()
	}
	return tx.Commit()
}

func (d *DB) Begin() (*Tx, error) {
	d.Lock()
	tx, err := d.DB.Begin()
	if err != nil {
		d.Unlock()
		return nil, err
	}
	return &Tx{
		Tx:     tx,
		Locker: d.Locker,
	}, nil
}

func (d *DB) WithTx(ctx context.Context) (context.Context, error) {
	tx, err := d.Begin()
	if err != nil {
		return ctx, err
	}
	return WithTx(ctx, tx), nil
}

func (d *DB) TxFrom(ctx context.Context) (*Tx, error) {
	tx := TxFrom(ctx)
	if tx == nil {
		return nil, ErrNoTx
	}
	return tx, nil
}

func (d *DB) CommitFrom(ctx context.Context) error {
	tx, err := d.TxFrom(ctx)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (d *DB) RollbackFrom(ctx context.Context) error {
	tx, err := d.TxFrom(ctx)
	if err != nil {
		return err
	}
	return tx.Rollback()
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

func (t *Tx) Commit() error {
	err := t.Tx.Commit()
	t.Unlock()
	return err
}

func (t *Tx) Rollback() error {
	err := t.Tx.Rollback()
	t.Unlock()
	return err
}

func (t *Tx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := t.Tx.Query(query, args...)
	return rows, err
}

func (t *Tx) Exec(query string, args ...interface{}) (sql.Result, error) {
	res, err := t.Tx.Exec(query, args...)
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

type Txer interface {
	Begin() (*Tx, error)
	WithTx(context.Context) (context.Context, error)
	TxFrom(context.Context) (*Tx, error)
	TxFunc(context.Context, func(context.Context) bool) error
	CommitFrom(context.Context) error
	RollbackFrom(context.Context) error
}
