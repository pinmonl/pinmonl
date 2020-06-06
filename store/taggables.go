package store

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
)

type Taggables struct {
	*Store
}

type TaggableOpts struct {
	ListOpts
	TagIDs  []string
	Targets model.MorphableList
}

func NewTaggables(s *Store) *Taggables {
	return &Taggables{s}
}

func (t *Taggables) table() string {
	return "taggables"
}

func (t *Taggables) List(ctx context.Context, opts *TaggableOpts) ([]*model.Taggable, error) {
	if opts == nil {
		opts = &TaggableOpts{}
	}

	qb := t.RunnableBuilder(ctx).
		Select(t.columns()...).From(t.table())
	qb = t.bindOpts(qb, opts)
	qb = addPagination(qb, opts)
	rows, err := qb.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*model.Taggable
	for rows.Next() {
		taggable, err := t.scan(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, taggable)
	}
	return list, nil
}

func (t *Taggables) Count(ctx context.Context, opts *TaggableOpts) (int64, error) {
	if opts == nil {
		opts = &TaggableOpts{}
	}

	qb := t.RunnableBuilder(ctx).
		Select("count(*)").From(t.table())
	qb = t.bindOpts(qb, opts)
	row := qb.QueryRow()
	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (t *Taggables) Find(ctx context.Context, id string) (*model.Taggable, error) {
	qb := t.RunnableBuilder(ctx).
		Select(t.columns()...).From(t.table()).
		Where("id = ?", id)
	row := qb.QueryRow()
	taggable, err := t.scan(row)
	if err != nil {
		return nil, err
	}
	return taggable, nil
}

func (t *Taggables) columns() []string {
	return []string{
		"id",
		"tag_id",
		"taggable_id",
		"taggable_name",
	}
}

func (t *Taggables) bindOpts(b squirrel.SelectBuilder, opts *TaggableOpts) squirrel.SelectBuilder {
	if opts == nil {
		return b
	}

	if len(opts.TagIDs) > 0 {
		b = b.Where(squirrel.Eq{"tag_id": opts.TagIDs})
	}

	if len(opts.Targets) > 0 && !opts.Targets.IsMixed() {
		b = b.Where("target_name = ?", opts.Targets.MorphName()).
			Where(squirrel.Eq{"target_id": opts.Targets.MorphKeys()})
	}

	return b
}

func (t *Taggables) scan(row database.RowScanner) (*model.Taggable, error) {
	var taggable model.Taggable
	err := row.Scan(
		&taggable.ID,
		&taggable.TagID,
		&taggable.TaggableID,
		&taggable.TaggableName)
	if err != nil {
		return nil, err
	}
	return &taggable, nil
}

func (t *Taggables) ListWithTags(ctx context.Context, opts *TaggableOpts) ([]*model.Taggable, error) {
	if opts == nil {
		opts = &TaggableOpts{}
	}

	tags := &Tags{}
	qb := t.RunnableBuilder(ctx).
		Select(tags.columns()...).
		Columns(
			"taggables.id AS taggable_id",
			"taggable_id",
			"taggable_name").
		From(t.table()).
		LeftJoin(fmt.Sprintf("%s ON %[1]s.id = %s.tag_id", tags.table(), t.table()))
	qb = t.bindOpts(qb, opts)
	qb = addPagination(qb, opts)
	rows, err := qb.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*model.Taggable
	for rows.Next() {
		var (
			mtg model.Taggable
			mt  model.Tag
		)
		err := rows.Scan(
			&mt.ID,
			&mt.Name,
			&mt.UserID,
			&mt.ParentID,
			&mt.Level,
			&mt.Color,
			&mt.BgColor,
			&mt.HasChildren,
			&mt.CreatedAt,
			&mt.UpdatedAt,
			&mtg.ID,
			&mtg.TaggableID,
			&mtg.TaggableName)
		if err != nil {
			return nil, err
		}
		mtg.TagID = mt.ID
		mtg.Tag = &mt
		list = append(list, &mtg)
	}
	return list, nil
}

func (t *Taggables) Create(ctx context.Context, taggable *model.Taggable) error {
	taggable2 := *taggable
	taggable2.ID = newID()

	qb := t.RunnableBuilder(ctx).
		Insert(t.table()).
		Columns(
			"id",
			"tag_id",
			"taggable_id",
			"taggable_name").
		Values(
			taggable2.ID,
			taggable2.TagID,
			taggable2.TaggableID,
			taggable2.TaggableName)
	_, err := qb.Exec()
	if err != nil {
		return err
	}
	*taggable = taggable2
	return nil
}

func (t *Taggables) Update(ctx context.Context, taggable *model.Taggable) error {
	taggable2 := *taggable

	qb := t.RunnableBuilder(ctx).
		Update(t.table()).
		Set("tag_id", taggable2.TagID).
		Set("taggable_id", taggable2.TaggableID).
		Set("taggable_name", taggable2.TaggableName).
		Where("id = ?", taggable2.ID)
	_, err := qb.Exec()
	if err != nil {
		return err
	}
	*taggable = taggable2
	return nil
}

func (t *Taggables) Delete(ctx context.Context, id string) (int64, error) {
	qb := t.RunnableBuilder(ctx).
		Delete(t.table()).
		Where("id = ?", id)
	res, err := qb.Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
