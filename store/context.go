package store

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type ctxKey int

const (
	txCtxKey ctxKey = iota
)

// WithTx passes transaction into context.
func WithTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, txCtxKey, tx)
}

// TxFrom gets transaction from context.
func TxFrom(ctx context.Context) *sqlx.Tx {
	tx, _ := ctx.Value(txCtxKey).(*sqlx.Tx)
	return tx
}
