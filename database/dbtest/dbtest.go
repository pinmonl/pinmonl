package dbtest

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Masterminds/squirrel"
	"github.com/pinmonl/pinmonl/database"
)

func New() (*database.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	return &database.DB{
		DB:      db,
		Locker:  &database.NopLocker{},
		Builder: database.NewBuilderFromBase(squirrel.StatementBuilder),
	}, mock, err
}
