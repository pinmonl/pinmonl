package sharing

import (
	"context"

	"github.com/pinmonl/pinmonl/model"
)

type ctxKey int

const (
	userCtxKey ctxKey = iota
	shareCtxKey
)

// WithUser stores User into context.
func WithUser(ctx context.Context, m model.User) context.Context {
	return context.WithValue(ctx, userCtxKey, m)
}

// UserFrom retrieves User from context.
func UserFrom(ctx context.Context) (model.User, bool) {
	m, ok := ctx.Value(userCtxKey).(model.User)
	return m, ok
}

// WithShare stores Share into context.
func WithShare(ctx context.Context, m model.Share) context.Context {
	return context.WithValue(ctx, shareCtxKey, m)
}

// ShareFrom retrieves Share from context.
func ShareFrom(ctx context.Context) (model.Share, bool) {
	m, ok := ctx.Value(shareCtxKey).(model.Share)
	return m, ok
}
