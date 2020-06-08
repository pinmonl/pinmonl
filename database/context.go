package database

import "context"

type ctxKey int

const (
	TxKey ctxKey = iota
)

func WithTx(ctx context.Context, tx *Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}

func TxFrom(ctx context.Context) *Tx {
	tx, ok := ctx.Value(TxKey).(*Tx)
	if !ok {
		return nil
	}
	return tx
}
