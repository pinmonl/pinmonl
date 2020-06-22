package job

import (
	"context"
	"fmt"
	"time"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/store"
)

type Job interface {
	fmt.Stringer

	// Describe reports the descriptor of the job,
	// which is the unique identifer in the queue.
	Describe() []string

	// Target reports the related model.
	Target() model.Morphable

	// RunAt reports the time of the job that should
	// be executed.
	RunAt() time.Time

	// PreRun starts fetching the data for Run.
	PreRun(context.Context) error

	// Run starts the process which involves changes of the
	// database entries.
	Run(context.Context) ([]Job, error)
}

type CtxKey int

const (
	StoresCtxKey CtxKey = iota
)

func WithStores(ctx context.Context, stores *store.Stores) context.Context {
	return context.WithValue(ctx, StoresCtxKey, stores)
}

func StoresFrom(ctx context.Context) *store.Stores {
	stores, ok := ctx.Value(StoresCtxKey).(*store.Stores)
	if !ok {
		return nil
	}
	return stores
}
