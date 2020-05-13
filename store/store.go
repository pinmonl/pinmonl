package store

import (
	"context"
	"fmt"

	"github.com/pinmonl/pinmonl/database"
)

// Store provides the general interface of store.
type Store interface {
	BeginTx(context.Context) (context.Context, error)
	EndTx(context.Context) (context.Context, error)
	Rollback(context.Context) error
	Commit(context.Context) error
	Queryer(context.Context) database.Queryer
	Execer(context.Context) database.Execer
	Ext(context.Context) database.Ext
}

// NewStore creates the store with database.
func NewStore(db *database.DB) Store {
	return &dbStore{db}
}

type dbStore struct {
	db *database.DB
}

// BeginTx starts a Tx and passes it into context.
func (s *dbStore) BeginTx(ctx context.Context) (context.Context, error) {
	tx := TxFrom(ctx)
	if tx != nil {
		return ctx, fmt.Errorf("store: transaction has already been started")
	}
	tx, _ = s.db.Beginx()
	ctx = WithTx(ctx, tx)
	return ctx, nil
}

// EndTx removes transaction from context.
func (s *dbStore) EndTx(ctx context.Context) (context.Context, error) {
	tx := TxFrom(ctx)
	if tx == nil {
		return ctx, fmt.Errorf("store: transaction does not exist")
	}
	ctx = WithTx(ctx, nil)
	return ctx, nil
}

// Commit commits the transaction in context if exists.
func (s *dbStore) Commit(ctx context.Context) error {
	tx := TxFrom(ctx)
	if tx != nil {
		return tx.Commit()
	}
	return nil
}

// Rollback rollbacks the transaction in context if exists.
func (s *dbStore) Rollback(ctx context.Context) error {
	tx := TxFrom(ctx)
	if tx != nil {
		return tx.Rollback()
	}
	return nil
}

func (s *dbStore) Queryer(ctx context.Context) database.Queryer {
	tx := TxFrom(ctx)
	if tx != nil {
		return tx
	}
	return s.db
}

func (s *dbStore) Execer(ctx context.Context) database.Execer {
	tx := TxFrom(ctx)
	if tx != nil {
		return tx
	}
	return s.db
}

// Ext returns the Tx from context if exists, otherwise returns the DB.
func (s *dbStore) Ext(ctx context.Context) database.Ext {
	return s.db
}
