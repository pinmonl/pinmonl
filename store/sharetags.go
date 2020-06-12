package store

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
)

type Sharetags struct {
	*Store
}

type SharetagOpts struct {
	ListOpts
	ShareIDs  []string
	TagIDs    []string
	Kind      field.NullValue
	Level     field.NullInt64
	ParentIDs []string
	Status    field.NullValue
}

func NewSharetags(s *Store) *Sharetags {
	return &Sharetags{s}
}

func (s *Sharetags) table() string {
	return "sharetags"
}

func (s *Sharetags) List(ctx context.Context, opts *SharetagOpts) (model.SharetagList, error) {
	if opts == nil {
		opts = &SharetagOpts{}
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
	var list []*model.Sharetag
	for rows.Next() {
		sharetag, err := s.scan(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, sharetag)
	}
	return list, nil
}

func (s *Sharetags) Count(ctx context.Context, opts *SharetagOpts) (int64, error) {
	if opts == nil {
		opts = &SharetagOpts{}
	}

	qb := s.RunnableBuilder(ctx).
		Select("count(*)").From(s.table())
	qb = s.bindOpts(qb, opts)
	row := qb.QueryRow()
	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Sharetags) Find(ctx context.Context, id string) (*model.Sharetag, error) {
	qb := s.RunnableBuilder(ctx).
		Select(s.columns()...).From(s.table()).
		Where("id = ?", id)
	row := qb.QueryRow()
	sharetag, err := s.scan(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return sharetag, nil
}

func (s *Sharetags) columns() []string {
	return []string{
		"id",
		"share_id",
		"tag_id",
		"kind",
		"parent_id",
		"level",
		"status",
		"has_children",
	}
}

func (s *Sharetags) bindOpts(b squirrel.SelectBuilder, opts *SharetagOpts) squirrel.SelectBuilder {
	if opts == nil {
		return b
	}

	if len(opts.ShareIDs) > 0 {
		b = b.Where(squirrel.Eq{"share_id": opts.ShareIDs})
	}

	if len(opts.TagIDs) > 0 {
		b = b.Where(squirrel.Eq{"tag_id": opts.TagIDs})
	}

	if opts.Kind.Valid {
		if k, ok := opts.Kind.Value().(model.SharetagKind); ok {
			b = b.Where("kind = ?", k)
		}
	}

	if opts.Level.Valid {
		b = b.Where("level = ?", opts.Level.Value())
	}

	if len(opts.ParentIDs) > 0 {
		b = b.Where(squirrel.Eq{"parent_id": opts.ParentIDs})
	}

	if opts.Status.Valid {
		if s, ok := opts.Status.Value().(model.ShareStatus); ok {
			b = b.Where("status = ?", s)
		}
	}

	return b
}

func (s *Sharetags) scan(row database.RowScanner) (*model.Sharetag, error) {
	var sharetag model.Sharetag
	err := row.Scan(
		&sharetag.ID,
		&sharetag.ShareID,
		&sharetag.TagID,
		&sharetag.Kind,
		&sharetag.ParentID,
		&sharetag.Level,
		&sharetag.Status,
		&sharetag.HasChildren)
	if err != nil {
		return nil, err
	}
	return &sharetag, nil
}

func (s *Sharetags) Create(ctx context.Context, sharetag *model.Sharetag) error {
	sharetag2 := *sharetag
	sharetag2.ID = newID()

	qb := s.RunnableBuilder(ctx).
		Insert(s.table()).
		Columns(
			"id",
			"share_id",
			"tag_id",
			"kind",
			"parent_id",
			"level",
			"status",
			"has_children").
		Values(
			sharetag2.ID,
			sharetag2.ShareID,
			sharetag2.TagID,
			sharetag2.Kind,
			sharetag2.ParentID,
			sharetag2.Level,
			sharetag2.Status,
			sharetag2.HasChildren)
	_, err := qb.Exec()
	if err != nil {
		return err
	}
	*sharetag = sharetag2
	return nil
}

func (s *Sharetags) Update(ctx context.Context, sharetag *model.Sharetag) error {
	sharetag2 := *sharetag

	qb := s.RunnableBuilder(ctx).
		Update(s.table()).
		Set("share_id", sharetag2.ShareID).
		Set("tag_id", sharetag2.TagID).
		Set("kind", sharetag2.Kind).
		Set("parent_id", sharetag2.ParentID).
		Set("level", sharetag2.Level).
		Set("status", sharetag2.Status).
		Set("has_children", sharetag2.HasChildren).
		Where("id = ?", sharetag2.ID)
	_, err := qb.Exec()
	if err != nil {
		return err
	}
	*sharetag = sharetag2
	return nil
}

func (s *Sharetags) Delete(ctx context.Context, id string) (int64, error) {
	qb := s.RunnableBuilder(ctx).
		Delete(s.table()).
		Where("id = ?", id)
	res, err := qb.Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (s *Sharetags) DeleteByShare(ctx context.Context, shareID string) (int64, error) {
	qb := s.RunnableBuilder(ctx).
		Delete(s.table()).
		Where("share_id = ?", shareID)
	res, err := qb.Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
