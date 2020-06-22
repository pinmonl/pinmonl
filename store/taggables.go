package store

import (
	"context"
	"database/sql"
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
	TagIDs     []string
	Targets    model.MorphableList
	TargetIDs  []string
	TargetName string

	joinTags bool
}

func NewTaggables(s *Store) *Taggables {
	return &Taggables{s}
}

func (t Taggables) table() string {
	return "taggables"
}

func (t *Taggables) List(ctx context.Context, opts *TaggableOpts) (model.TaggableList, error) {
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
	list := make([]*model.Taggable, 0)
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
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return taggable, nil
}

func (t *Taggables) FindOrCreate(ctx context.Context, data *model.Taggable) (*model.Taggable, error) {
	found, err := t.List(ctx, &TaggableOpts{
		TagIDs:     []string{data.TagID},
		TargetIDs:  []string{data.TargetID},
		TargetName: data.TargetName,
	})
	if err != nil {
		return nil, err
	}
	if len(found) > 0 {
		return found[0], nil
	}

	taggable := *data
	err = t.Create(ctx, &taggable)
	if err != nil {
		return nil, err
	}
	return &taggable, nil
}

func (t Taggables) bindOpts(b squirrel.SelectBuilder, opts *TaggableOpts) squirrel.SelectBuilder {
	if opts == nil {
		return b
	}

	if len(opts.TagIDs) > 0 {
		b = b.Where(squirrel.Eq{t.table() + ".tag_id": opts.TagIDs})
	}

	if len(opts.Targets) > 0 && !opts.Targets.IsMixed() {
		opts.TargetName = opts.Targets.MorphName()
		opts.TargetIDs = opts.Targets.MorphKeys()
	}
	if opts.TargetName != "" {
		b = b.Where(t.table()+".target_name = ?", opts.TargetName)
	}
	if len(opts.TargetIDs) > 0 {
		b = b.Where(squirrel.Eq{t.table() + ".target_id": opts.TargetIDs})
	}

	return b
}

func (t Taggables) columns() []string {
	return []string{
		t.table() + ".id",
		t.table() + ".tag_id",
		t.table() + ".target_id",
		t.table() + ".target_name",
	}
}

func (t Taggables) scanColumns(taggable *model.Taggable) []interface{} {
	return []interface{}{
		&taggable.ID,
		&taggable.TagID,
		&taggable.TargetID,
		&taggable.TargetName,
	}
}

func (t Taggables) scan(row database.RowScanner) (*model.Taggable, error) {
	var taggable model.Taggable
	err := row.Scan(t.scanColumns(&taggable)...)
	if err != nil {
		return nil, err
	}
	return &taggable, nil
}

func (t *Taggables) ListWithTag(ctx context.Context, opts *TaggableOpts) (model.TaggableList, error) {
	if opts == nil {
		opts = &TaggableOpts{}
	}

	qb := t.RunnableBuilder(ctx).
		Select(t.columns()...).
		Columns(Tags{}.columns()...).
		From(t.table()).
		LeftJoin(fmt.Sprintf("%s ON %[1]s.id = %s.tag_id", Tags{}.table(), t.table()))
	qb = t.bindOpts(qb, opts)
	qb = addPagination(qb, opts)
	rows, err := qb.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := make([]*model.Taggable, 0)
	for rows.Next() {
		var (
			mtg model.Taggable
			mt  model.Tag
		)
		scanCols := append(t.scanColumns(&mtg), Tags{}.scanColumns(&mt)...)
		err := rows.Scan(scanCols...)
		if err != nil {
			return nil, err
		}
		mtg.Tag = &mt
		list = append(list, &mtg)
	}
	return list, nil
}

func (t *Taggables) listWithTarget(
	ctx context.Context,
	columns []string,
	tableName, targetName string,
	opts *TaggableOpts,
	rowScan func(row database.RowScanner) (*model.Taggable, error),
) (model.TaggableList, error) {
	if opts == nil {
		opts = &TaggableOpts{}
	}

	qb := t.RunnableBuilder(ctx).
		Select(t.columns()...).
		Columns(columns...).
		From(t.table()).
		LeftJoin(fmt.Sprintf("%s ON %[1]s.id = %s.target_id", tableName, t.table())).
		Where(t.table()+".target_name = ?", targetName)
	qb = t.bindOpts(qb, opts)
	qb = addPagination(qb, opts)
	rows, err := qb.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := make([]*model.Taggable, 0)
	for rows.Next() {
		taggable, err := rowScan(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, taggable)
	}
	return list, nil
}

func (t *Taggables) ListWithPinl(ctx context.Context, opts *TaggableOpts) (model.TaggableList, error) {
	scanner := func(row database.RowScanner) (*model.Taggable, error) {
		var (
			mtg model.Taggable
			mp  model.Pinl
		)
		scanCols := append(t.scanColumns(&mtg), Pinls{}.scanColumns(&mp)...)
		err := row.Scan(scanCols...)
		if err != nil {
			return nil, err
		}
		mtg.Pinl = &mp
		return &mtg, nil
	}

	return t.listWithTarget(ctx, Pinls{}.columns(), Pinls{}.table(), model.Pinl{}.MorphName(), opts, scanner)
}

func (t *Taggables) Create(ctx context.Context, taggable *model.Taggable) error {
	taggable2 := *taggable
	taggable2.ID = newID()

	qb := t.RunnableBuilder(ctx).
		Insert(t.table()).
		Columns(
			"id",
			"tag_id",
			"target_id",
			"target_name").
		Values(
			taggable2.ID,
			taggable2.TagID,
			taggable2.TargetID,
			taggable2.TargetName)
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
		Set("target_id", taggable2.TargetID).
		Set("target_name", taggable2.TargetName).
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

func (t *Taggables) DeleteByTarget(ctx context.Context, target model.Morphable) (int64, error) {
	qb := t.RunnableBuilder(ctx).
		Delete(t.table()).
		Where("target_name = ?", target.MorphName()).
		Where("target_id = ?", target.MorphKey())
	res, err := qb.Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (t *Taggables) DeleteByTag(ctx context.Context, tag *model.Tag) (int64, error) {
	qb := t.RunnableBuilder(ctx).
		Delete(t.table()).
		Where("tag_id = ?", tag.ID)
	res, err := qb.Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (o *TaggableOpts) JoinTags() *TaggableOpts {
	o2 := *o
	o2.joinTags = true
	return &o2
}
