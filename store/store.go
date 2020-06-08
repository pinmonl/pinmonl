package store

import (
	"context"

	"github.com/Masterminds/squirrel"
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
	tx := database.TxFrom(ctx)
	if tx != nil {
		return tx
	}
	return s.db
}

func (s *Store) Builder() database.Builder {
	return s.db.Builder
}

func (s *Store) RunnableBuilder(ctx context.Context) database.Builder {
	r := s.Runner(ctx)
	b := s.db.Builder.RunWith(r)
	return b
}

type ListOpts struct {
	Limit  int
	Offset int
}

func (o ListOpts) LimitUint64() uint64 {
	return uint64(o.Limit)
}

func (o ListOpts) OffsetUint64() uint64 {
	return uint64(o.Offset)
}

type Paginator interface {
	LimitUint64() uint64
	OffsetUint64() uint64
}

func addPagination(b squirrel.SelectBuilder, pt Paginator) squirrel.SelectBuilder {
	if pt == nil {
		return b
	}
	if pt.LimitUint64() > 0 {
		return b.
			Limit(pt.LimitUint64()).
			Offset(pt.OffsetUint64())
	}
	return b
}
