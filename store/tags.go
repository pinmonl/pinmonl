package store

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
)

type Tags struct {
	*Store
}

type TagOpts struct {
	ListOpts
	UserID      string
	UserIDs     []string
	Name        string
	NamePattern string
	ParentIDs   []string
	Level       field.NullInt64
}

func NewTags(s *Store) *Tags {
	return &Tags{s}
}

func (t *Tags) table() string {
	return "tags"
}

func (t *Tags) List(ctx context.Context, opts *TagOpts) (model.TagList, error) {
	if opts == nil {
		opts = &TagOpts{}
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
	var list []*model.Tag
	for rows.Next() {
		tag, err := t.scan(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, tag)
	}
	return list, nil
}

func (t *Tags) Count(ctx context.Context, opts *TagOpts) (int64, error) {
	if opts == nil {
		opts = &TagOpts{}
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

func (t *Tags) Find(ctx context.Context, id string) (*model.Tag, error) {
	qb := t.RunnableBuilder(ctx).
		Select(t.columns()...).From(t.table()).
		Where("id = ?", id)
	row := qb.QueryRow()
	tag, err := t.scan(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func (t *Tags) columns() []string {
	return []string{
		"id",
		"name",
		"user_id",
		"parent_id",
		"level",
		"color",
		"bg_color",
		"has_children",
		"created_at",
		"updated_at",
	}
}

func (t *Tags) bindOpts(b squirrel.SelectBuilder, opts *TagOpts) squirrel.SelectBuilder {
	if opts == nil {
		return b
	}

	if opts.UserID != "" {
		opts.UserIDs = append(opts.UserIDs, opts.UserID)
	}
	if len(opts.UserIDs) > 0 {
		b = b.Where(squirrel.Eq{"user_id": opts.UserIDs})
	}

	if opts.Name != "" {
		b = b.Where("name = ?", opts.Name)
	}

	if opts.NamePattern != "" {
		b = b.Where("name LIKE ?", opts.NamePattern)
	}

	if len(opts.ParentIDs) > 0 {
		b = b.Where(squirrel.Eq{"parent_id": opts.ParentIDs})
	}

	if opts.Level.Valid {
		b = b.Where("level = ?", opts.Level.Value())
	}

	return b
}

func (t *Tags) scan(row database.RowScanner) (*model.Tag, error) {
	var tag model.Tag
	err := row.Scan(
		&tag.ID,
		&tag.Name,
		&tag.UserID,
		&tag.ParentID,
		&tag.Level,
		&tag.Color,
		&tag.BgColor,
		&tag.HasChildren,
		&tag.CreatedAt,
		&tag.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (t *Tags) Create(ctx context.Context, tag *model.Tag) error {
	tag2 := *tag
	tag2.ID = newID()
	tag2.CreatedAt = timestamp()
	tag2.UpdatedAt = timestamp()

	qb := t.RunnableBuilder(ctx).
		Insert(t.table()).
		Columns(
			"id",
			"name",
			"user_id",
			"parent_id",
			"level",
			"color",
			"bg_color",
			"has_children",
			"created_at",
			"updated_at").
		Values(
			tag2.ID,
			tag2.Name,
			tag2.UserID,
			tag2.ParentID,
			tag2.Level,
			tag2.Color,
			tag2.BgColor,
			tag2.HasChildren,
			tag2.CreatedAt,
			tag2.UpdatedAt)
	_, err := qb.Exec()
	if err != nil {
		return err
	}
	*tag = tag2
	return nil
}

func (t *Tags) Update(ctx context.Context, tag *model.Tag) error {
	tag2 := *tag
	tag2.UpdatedAt = timestamp()

	qb := t.RunnableBuilder(ctx).
		Update(t.table()).
		Set("name", tag2.Name).
		Set("user_id", tag2.UserID).
		Set("parent_id", tag2.ParentID).
		Set("level", tag2.Level).
		Set("color", tag2.Color).
		Set("bg_color", tag2.BgColor).
		Set("has_children", tag2.HasChildren).
		Set("updated_at", tag2.UpdatedAt).
		Where("id = ?", tag2.ID)
	_, err := qb.Exec()
	if err != nil {
		return err
	}
	*tag = tag2
	return nil
}

func (t *Tags) Delete(ctx context.Context, id string) (int64, error) {
	qb := t.RunnableBuilder(ctx).
		Delete(t.table()).
		Where("id = ?", id)
	res, err := qb.Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
