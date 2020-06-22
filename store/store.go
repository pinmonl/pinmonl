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

type Stores struct {
	Store *Store

	Images    *Images
	Jobs      *Jobs
	Monls     *Monls
	Monpkgs   *Monpkgs
	Pinls     *Pinls
	Pinpkgs   *Pinpkgs
	Pkgs      *Pkgs
	Sharepins *Sharepins
	Shares    *Shares
	Sharetags *Sharetags
	Stats     *Stats
	Taggables *Taggables
	Tags      *Tags
	Users     *Users
}

func NewStores(db *database.DB) *Stores {
	s := NewStore(db)
	return &Stores{
		Store: s,

		Images:    NewImages(s),
		Jobs:      NewJobs(s),
		Monls:     NewMonls(s),
		Monpkgs:   NewMonpkgs(s),
		Pinls:     NewPinls(s),
		Pinpkgs:   NewPinpkgs(s),
		Pkgs:      NewPkgs(s),
		Sharepins: NewSharepins(s),
		Shares:    NewShares(s),
		Sharetags: NewSharetags(s),
		Stats:     NewStats(s),
		Taggables: NewTaggables(s),
		Tags:      NewTags(s),
		Users:     NewUsers(s),
	}
}
