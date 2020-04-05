package store

import (
	"context"

	"github.com/pinmonl/pinmonl/database"
)

type ctxKey int

const (
	txCtxKey ctxKey = iota
)

// WithTx passes transaction into context.
func WithTx(ctx context.Context, tx *database.Tx) context.Context {
	return context.WithValue(ctx, txCtxKey, tx)
}

// TxFrom gets transaction from context.
func TxFrom(ctx context.Context) *database.Tx {
	tx, _ := ctx.Value(txCtxKey).(*database.Tx)
	return tx
}
