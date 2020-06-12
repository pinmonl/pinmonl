package store

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
)

type Sharepins struct {
	*Store
}

type SharepinOpts struct {
	ListOpts
	ShareIDs []string
	PinlIDs  []string
	Status   field.NullValue
}

func NewSharepins(s *Store) *Sharepins {
	return &Sharepins{s}
}

func (s *Sharepins) table() string {
	return "sharepins"
}

func (s *Sharepins) List(ctx context.Context, opts *SharepinOpts) ([]*model.Sharepin, error) {
	if opts == nil {
		opts = &SharepinOpts{}
	}

	qb := s.RunnableBuilder(ctx).
		Select(s.columns()...).From(s.table())
	qb = s.bindOpts(qb, opts)
	qb = addPagination(qb, opts)
	rows, err := qb.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*model.Sharepin
	for rows.Next() {
		sharepin, err := s.scan(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, sharepin)
	}
	return list, nil
}

func (s *Sharepins) Count(ctx context.Context, opts *SharepinOpts) (int64, error) {
	if opts == nil {
		opts = &SharepinOpts{}
	}

	qb := s.RunnableBuilder(ctx).
		Select("count(*)").From(s.table())
	row := qb.QueryRow()
	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Sharepins) Find(ctx context.Context, id string) (*model.Sharepin, error) {
	qb := s.RunnableBuilder(ctx).
		Select(s.columns()...).From(s.table()).
		Where("id = ?", id)
	row := qb.QueryRow()
	sharepin, err := s.scan(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return sharepin, nil
}

func (s *Sharepins) columns() []string {
	return []string{
		"id",
		"share_id",
		"pinl_id",
		"status",
	}
}

func (s *Sharepins) bindOpts(b squirrel.SelectBuilder, opts *SharepinOpts) squirrel.SelectBuilder {
	if opts == nil {
		return b
	}

	if len(opts.ShareIDs) > 0 {
		b = b.Where(squirrel.Eq{"share_id": opts.ShareIDs})
	}

	if len(opts.PinlIDs) > 0 {
		b = b.Where(squirrel.Eq{"pinl_id": opts.PinlIDs})
	}

	if opts.Status.Valid {
		if s, ok := opts.Status.Value().(model.ShareStatus); ok {
			b = b.Where("status = ?", s)
		}
	}

	return b
}

func (s *Sharepins) scan(row database.RowScanner) (*model.Sharepin, error) {
	var sharepin model.Sharepin
	err := row.Scan(
		&sharepin.ID,
		&sharepin.ShareID,
		&sharepin.PinlID,
		&sharepin.Status)
	if err != nil {
		return nil, err
	}
	return &sharepin, nil
}

func (s *Sharepins) Create(ctx context.Context, sharepin *model.Sharepin) error {
	sharepin2 := *sharepin
	sharepin2.ID = newID()

	qb := s.RunnableBuilder(ctx).
		Insert(s.table()).
		Columns(
			"id",
			"share_id",
			"pinl_id",
			"status").
		Values(
			sharepin2.ID,
			sharepin2.ShareID,
			sharepin2.PinlID,
			sharepin2.Status)
	_, err := qb.Exec()
	if err != nil {
		return err
	}
	*sharepin = sharepin2
	return nil
}

func (s *Sharepins) Update(ctx context.Context, sharepin *model.Sharepin) error {
	sharepin2 := *sharepin

	qb := s.RunnableBuilder(ctx).
		Update(s.table()).
		Set("share_id", sharepin2.ShareID).
		Set("pinl_id", sharepin2.PinlID).
		Set("status", sharepin2.Status).
		Where("id = ?", sharepin2.ID)
	_, err := qb.Exec()
	if err != nil {
		return err
	}
	*sharepin = sharepin2
	return nil
}

func (s *Sharepins) Delete(ctx context.Context, id string) (int64, error) {
	qb := s.RunnableBuilder(ctx).
		Delete(s.table()).
		Where("id = ?", id)
	res, err := qb.Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
