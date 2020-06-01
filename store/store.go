package store

import (
	"context"

	"github.com/pinmonl/pinmonl/database"
)

type Store struct {
	db *database.DB
}

func NewStore(db *database.DB) *Store {
	return &Store{db}
}

func (s *Store) Execer(ctx context.Context) database.Execer {
	return s.Runner(ctx)
}

func (s *Store) Queryer(ctx context.Context) database.Queryer {
	return s.Runner(ctx)
}

func (s *Store) Runner(ctx context.Context) database.Runner {
	runner := database.GetRunner(ctx)
	if runner != nil {
		return runner
	}
	return s.db
}
