package store

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
)

type Stats struct {
	*Store
}

type StatOpts struct {
	ListOpts
	PkgIDs    []string
	ParentIDs []string
	Kind      field.NullValue
	IsLatest  field.NullBool
}

func NewStats(s *Store) *Stats {
	return &Stats{s}
}

func (s *Stats) table() string {
	return "stats"
}

func (s *Stats) List(ctx context.Context, opts *StatOpts) ([]*model.Stat, error) {
	if opts == nil {
		opts = &StatOpts{}
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
	var list []*model.Stat
	for rows.Next() {
		stat, err := s.scan(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, stat)
	}
	return list, nil
}

func (s *Stats) Count(ctx context.Context, opts *StatOpts) (int64, error) {
	if opts == nil {
		opts = &StatOpts{}
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

func (s *Stats) Find(ctx context.Context, id string) (*model.Stat, error) {
	qb := s.RunnableBuilder(ctx).
		Select(s.columns()...).From(s.table()).
		Where("id = ?", id)
	row := qb.QueryRow()
	stat, err := s.scan(row)
	if err != nil {
		return nil, err
	}
	return stat, nil
}

func (s *Stats) columns() []string {
	return []string{
		"id",
		"pkg_id",
		"parent_id",
		"recorded_at",
		"kind",
		"value",
		"value_type",
		"checksum",
		"weight",
		"is_latest",
		"has_children",
	}
}

func (s *Stats) bindOpts(b squirrel.SelectBuilder, opts *StatOpts) squirrel.SelectBuilder {
	if opts == nil {
		return b
	}

	if len(opts.PkgIDs) > 0 {
		b = b.Where(squirrel.Eq{"pkg_id": opts.PkgIDs})
	}

	if len(opts.ParentIDs) > 0 {
		b = b.Where(squirrel.Eq{"parent_id": opts.ParentIDs})
	}

	if opts.Kind.Valid {
		if k, ok := opts.Kind.Value().(model.StatKind); ok {
			b = b.Where("kind = ?", k)
		}
	}

	if opts.IsLatest.Valid {
		b = b.Where("is_latest = ?", opts.IsLatest.Value())
	}

	return b
}

func (s *Stats) scan(row database.RowScanner) (*model.Stat, error) {
	var stat model.Stat
	err := row.Scan(
		&stat.ID,
		&stat.PkgID,
		&stat.ParentID,
		&stat.RecordedAt,
		&stat.Kind,
		&stat.Value,
		&stat.ValueType,
		&stat.Checksum,
		&stat.Weight,
		&stat.IsLatest,
		&stat.HasChildren)
	if err != nil {
		return nil, err
	}
	return &stat, nil
}

func (s *Stats) Create(ctx context.Context, stat *model.Stat) error {
	stat2 := *stat
	stat2.ID = newID()

	qb := s.RunnableBuilder(ctx).
		Insert(s.table()).
		Columns(
			"id",
			"pkg_id",
			"parent_id",
			"recorded_at",
			"kind",
			"value",
			"value_type",
			"checksum",
			"weight",
			"is_latest",
			"has_children").
		Values(
			stat2.ID,
			stat2.PkgID,
			stat2.ParentID,
			stat2.RecordedAt,
			stat2.Kind,
			stat2.Value,
			stat2.ValueType,
			stat2.Checksum,
			stat2.Weight,
			stat2.IsLatest,
			stat2.HasChildren)
	_, err := qb.Exec()
	if err != nil {
		return err
	}
	*stat = stat2
	return nil
}

func (s *Stats) Update(ctx context.Context, stat *model.Stat) error {
	stat2 := *stat

	qb := s.RunnableBuilder(ctx).
		Update(s.table()).
		Set("pkg_id", stat2.PkgID).
		Set("parent_id", stat2.ParentID).
		Set("recorded_at", stat2.RecordedAt).
		Set("kind", stat2.Kind).
		Set("value", stat2.Value).
		Set("value_type", stat2.ValueType).
		Set("checksum", stat2.Checksum).
		Set("weight", stat2.Weight).
		Set("is_latest", stat2.IsLatest).
		Set("has_children", stat2.HasChildren).
		Where("id = ?", stat2.ID)
	_, err := qb.Exec()
	if err != nil {
		return err
	}
	*stat = stat2
	return nil
}

func (s *Stats) Delete(ctx context.Context, id string) (int64, error) {
	qb := s.RunnableBuilder(ctx).
		Delete(s.table()).
		Where("id = ?", id)
	res, err := qb.Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
