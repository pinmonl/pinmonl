package store

import (
	"context"
	"database/sql"
	"fmt"

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

	TagNames       []string
	TagNamePattern string
	joinTags       bool
}

func NewSharetags(s *Store) *Sharetags {
	return &Sharetags{s}
}

func (s Sharetags) table() string {
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
	list := make([]*model.Sharetag, 0)
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

func (s *Sharetags) FindOrCreate(ctx context.Context, data *model.Sharetag) (*model.Sharetag, error) {
	found, err := s.List(ctx, &SharetagOpts{
		ShareIDs: []string{data.ShareID},
		TagIDs:   []string{data.TagID},
	})
	if err != nil {
		return nil, err
	}
	if len(found) > 0 {
		return found[0], nil
	}

	sharetag := *data
	err = s.Create(ctx, &sharetag)
	if err != nil {
		return nil, err
	}
	return &sharetag, nil
}

func (s *Sharetags) ListWithTag(ctx context.Context, opts *SharetagOpts) (model.SharetagList, error) {
	if opts == nil {
		opts = &SharetagOpts{}
	}
	opts = opts.JoinTags()

	qb := s.RunnableBuilder(ctx).
		Select(s.columns()...).
		Columns(Tags{}.columns()...).
		From(s.table())
	qb = s.bindOpts(qb, opts)
	qb = addPagination(qb, opts)
	rows, err := qb.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*model.Sharetag, 0)
	for rows.Next() {
		var (
			mst model.Sharetag
			mt  model.Tag
		)
		scanCols := append(s.scanColumns(&mst), Tags{}.scanColumns(&mt)...)
		err := rows.Scan(scanCols...)
		if err != nil {
			return nil, err
		}
		mst.Tag = &mt
		list = append(list, &mst)
	}
	return list, nil
}

func (s Sharetags) bindOpts(b squirrel.SelectBuilder, opts *SharetagOpts) squirrel.SelectBuilder {
	if opts == nil {
		return b
	}

	if len(opts.ShareIDs) > 0 {
		b = b.Where(squirrel.Eq{s.table() + ".share_id": opts.ShareIDs})
	}

	if len(opts.TagIDs) > 0 {
		b = b.Where(squirrel.Eq{s.table() + ".tag_id": opts.TagIDs})
	}

	if opts.Kind.Valid {
		if k, ok := opts.Kind.Value().(model.SharetagKind); ok {
			b = b.Where(s.table()+".kind = ?", k)
		}
	}

	if opts.Level.Valid {
		b = b.Where(s.table()+".level = ?", opts.Level.Value())
	}

	if len(opts.ParentIDs) > 0 {
		b = b.Where(squirrel.Eq{s.table() + ".parent_id": opts.ParentIDs})
	}

	if opts.Status.Valid {
		if sv, ok := opts.Status.Value().(model.Status); ok {
			b = b.Where(s.table()+".status = ?", sv)
		}
	}

	if len(opts.TagNames) > 0 {
		opts = opts.JoinTags()
		b = b.Where(squirrel.Eq{Tags{}.table() + ".name": opts.TagNames})
	}

	if opts.TagNamePattern != "" {
		opts = opts.JoinTags()
		b = b.Where(Tags{}.table()+".name LIKE ?", opts.TagNamePattern)
	}

	if opts.joinTags {
		b = b.LeftJoin(fmt.Sprintf("%s ON %[1]s.id = %s.tag_id", Tags{}.table(), s.table()))
	}

	return b
}

func (s Sharetags) columns() []string {
	return []string{
		s.table() + ".id",
		s.table() + ".share_id",
		s.table() + ".tag_id",
		s.table() + ".kind",
		s.table() + ".parent_id",
		s.table() + ".level",
		s.table() + ".status",
		s.table() + ".has_children",
	}
}

func (s Sharetags) scanColumns(sharetag *model.Sharetag) []interface{} {
	return []interface{}{
		&sharetag.ID,
		&sharetag.ShareID,
		&sharetag.TagID,
		&sharetag.Kind,
		&sharetag.ParentID,
		&sharetag.Level,
		&sharetag.Status,
		&sharetag.HasChildren,
	}
}

func (s Sharetags) scan(row database.RowScanner) (*model.Sharetag, error) {
	var sharetag model.Sharetag
	err := row.Scan(s.scanColumns(&sharetag)...)
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

func (o *SharetagOpts) JoinTags() *SharetagOpts {
	o2 := *o
	o2.joinTags = true
	return &o2
}
