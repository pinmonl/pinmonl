package job

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/pinmonl/pinmonl/exchange"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pubsub"
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

var (
	ErrNoStores          = errors.New("job: stores is missing")
	ErrNoPubsuber        = errors.New("job: pubsuber is missing")
	ErrNoExchangeManager = errors.New("job: exchange manager is missing")
)

type CtxKey int

const (
	StoresCtxKey CtxKey = iota
	ExchangeManagerCtxKey
	PubsuberCtxKey
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

func WithExchangeManager(ctx context.Context, exm *exchange.Manager) context.Context {
	return context.WithValue(ctx, ExchangeManagerCtxKey, exm)
}

func ExchangeManagerFrom(ctx context.Context) *exchange.Manager {
	exm, ok := ctx.Value(ExchangeManagerCtxKey).(*exchange.Manager)
	if !ok {
		return nil
	}
	return exm
}

func WithPubsuber(ctx context.Context, hub pubsub.Pubsuber) context.Context {
	return context.WithValue(ctx, PubsuberCtxKey, hub)
}

func PubsuberFrom(ctx context.Context) pubsub.Pubsuber {
	hub, ok := ctx.Value(PubsuberCtxKey).(pubsub.Pubsuber)
	if !ok {
		return nil
	}
	return hub
}
