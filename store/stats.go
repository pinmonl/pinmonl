package store

import (
	"context"
	"database/sql"

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
	PkgIDs        []string
	ParentIDs     []string
	Kind          field.NullValue
	IsLatest      field.NullBool
	Kinds         []model.StatKind
	KindsExcluded []model.StatKind

	Orders []StatOrder
}

type StatOrder int

const (
	StatOrderByRecordDesc StatOrder = iota
)

func NewStats(s *Store) *Stats {
	return &Stats{s}
}

func (s Stats) table() string {
	return "stats"
}

func (s *Stats) List(ctx context.Context, opts *StatOpts) (model.StatList, error) {
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
	list := make([]*model.Stat, 0)
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

	o2 := *opts
	o2.Orders = nil

	qb := s.RunnableBuilder(ctx).
		Select("count(*)").From(s.table())
	qb = s.bindOpts(qb, &o2)
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
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return stat, nil
}

func (s *Stats) FindMany(ctx context.Context, ids []string) (model.StatList, error) {
	qb := s.RunnableBuilder(ctx).
		Select(s.columns()...).From(s.table()).
		Where(squirrel.Eq{"id": ids})
	rows, err := qb.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := make([]*model.Stat, 0)
	for rows.Next() {
		stat, err := s.scan(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, stat)
	}
	return list, nil
}

func (s Stats) bindOpts(b squirrel.SelectBuilder, opts *StatOpts) squirrel.SelectBuilder {
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
			opts.Kinds = append(opts.Kinds, k)
		}
	}
	if len(opts.Kinds) > 0 {
		b = b.Where(squirrel.Eq{"kind": opts.Kinds})
	}

	if len(opts.KindsExcluded) > 0 {
		b = b.Where(squirrel.NotEq{"kind": opts.KindsExcluded})
	}

	if opts.IsLatest.Valid {
		b = b.Where("is_latest = ?", opts.IsLatest.Value())
	}

	for _, order := range opts.Orders {
		switch order {
		case StatOrderByRecordDesc:
			b = b.OrderBy("recorded_at DESC")
		}
	}

	return b
}

func (s Stats) columns() []string {
	return []string{
		s.table() + ".id",
		s.table() + ".pkg_id",
		s.table() + ".parent_id",
		s.table() + ".recorded_at",
		s.table() + ".kind",
		s.table() + ".name",
		s.table() + ".value",
		s.table() + ".value_type",
		s.table() + ".checksum",
		s.table() + ".weight",
		s.table() + ".is_latest",
		s.table() + ".has_children",
	}
}

func (s Stats) scanColumns(stat *model.Stat) []interface{} {
	return []interface{}{
		&stat.ID,
		&stat.PkgID,
		&stat.ParentID,
		&stat.RecordedAt,
		&stat.Kind,
		&stat.Name,
		&stat.Value,
		&stat.ValueType,
		&stat.Checksum,
		&stat.Weight,
		&stat.IsLatest,
		&stat.HasChildren,
	}
}

func (s Stats) scan(row database.RowScanner) (*model.Stat, error) {
	var stat model.Stat
	err := row.Scan(s.scanColumns(&stat)...)
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
			"name",
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
			stat2.Name,
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
		Set("name", stat2.Name).
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
