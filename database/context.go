package database

import "context"

type ctxKey int

const (
	RunnerKey ctxKey = iota
)

func WithRunner(ctx context.Context, runner Runner) context.Context {
	return context.WithValue(ctx, RunnerKey, runner)
}

func GetRunner(ctx context.Context) Runner {
	r, ok := ctx.Value(RunnerKey).(Runner)
	if !ok {
		return nil
	}
	return r
}
