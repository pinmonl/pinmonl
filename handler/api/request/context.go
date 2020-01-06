package request

import (
	"context"

	"github.com/pinmonl/pinmonl/logx"
	"github.com/pinmonl/pinmonl/model"
)

type ctxKey int

const (
	loggerCtxKey ctxKey = iota
	userCtxKey
	pinlCtxKey
	tagCtxKey
	shareCtxKey
	imageCtxKey
)

// WithLogger passes logger into context.
func WithLogger(ctx context.Context, l logx.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, l)
}

// LoggerFrom gets logger from context.
func LoggerFrom(ctx context.Context) (logx.Logger, bool) {
	l, ok := ctx.Value(loggerCtxKey).(logx.Logger)
	return l, ok
}

// MustLogger must get the logger from context.
func MustLogger(ctx context.Context) logx.Logger {
	l, ok := LoggerFrom(ctx)
	if !ok {
		panic("logger cannot find")
	}
	return l
}

// WithUser passes user into context.
func WithUser(ctx context.Context, m model.User) context.Context {
	return context.WithValue(ctx, userCtxKey, m)
}

// UserFrom gets user from context.
func UserFrom(ctx context.Context) (model.User, bool) {
	m, ok := ctx.Value(userCtxKey).(model.User)
	return m, ok
}

// WithPinl passes pinl into context.
func WithPinl(ctx context.Context, m model.Pinl) context.Context {
	return context.WithValue(ctx, pinlCtxKey, m)
}

// PinlFrom gets pinl from context.
func PinlFrom(ctx context.Context) (model.Pinl, bool) {
	m, ok := ctx.Value(pinlCtxKey).(model.Pinl)
	return m, ok
}

// WithTag passes tag into context.
func WithTag(ctx context.Context, m model.Tag) context.Context {
	return context.WithValue(ctx, tagCtxKey, m)
}

// TagFrom gets tag from context.
func TagFrom(ctx context.Context) (model.Tag, bool) {
	m, ok := ctx.Value(tagCtxKey).(model.Tag)
	return m, ok
}

// WithShare passes share into context.
func WithShare(ctx context.Context, m model.Share) context.Context {
	return context.WithValue(ctx, shareCtxKey, m)
}

// ShareFrom gets tag from context.
func ShareFrom(ctx context.Context) (model.Share, bool) {
	m, ok := ctx.Value(shareCtxKey).(model.Share)
	return m, ok
}

// WithImage passes share into context.
func WithImage(ctx context.Context, m model.Image) context.Context {
	return context.WithValue(ctx, imageCtxKey, m)
}

// ImageFrom gets tag from context.
func ImageFrom(ctx context.Context) (model.Image, bool) {
	m, ok := ctx.Value(imageCtxKey).(model.Image)
	return m, ok
}
