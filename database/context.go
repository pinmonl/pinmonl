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
	return ctx.Value(RunnerKey).(Runner)
}
