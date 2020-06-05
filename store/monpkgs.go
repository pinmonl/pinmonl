package store

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
)

type Monpkgs struct {
	*Store
}

type MonpkgOpts struct {
	ListOpts
	MonlIDs []string
	PkgIDs  []string
}

func NewMonpkgs(s *Store) *Monpkgs {
	return &Monpkgs{s}
}

func (m *Monpkgs) table() string {
	return "monpkgs"
}

func (m *Monpkgs) List(ctx context.Context, opts *MonpkgOpts) ([]*model.Monpkg, error) {
	if opts == nil {
		opts = &MonpkgOpts{}
	}

	qb := m.RunnableBuilder(ctx).
		Select(m.columns()...).From(m.table())
	qb = m.bindOpts(qb, opts)
	qb = addPagination(qb, opts)
	rows, err := qb.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []*model.Monpkg
	for rows.Next() {
		monpkg, err := m.scan(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, monpkg)
	}
	return list, nil
}

func (m *Monpkgs) Count(ctx context.Context, opts *MonpkgOpts) (int64, error) {
	if opts == nil {
		opts = &MonpkgOpts{}
	}

	qb := m.RunnableBuilder(ctx).
		Select("count(*)").From(m.table())
	qb = m.bindOpts(qb, opts)
	row := qb.QueryRow()
	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *Monpkgs) Find(ctx context.Context, id string) (*model.Monpkg, error) {
	qb := m.RunnableBuilder(ctx).
		Select(m.columns()...).From(m.table()).
		Where("id = ?", id).
		Limit(1)
	row := qb.QueryRow()
	monpkg, err := m.scan(row)
	if err != nil {
		return nil, err
	}
	return monpkg, nil
}

func (m *Monpkgs) columns() []string {
	return []string{
		"id",
		"monl_id",
		"pkg_id",
		"kind",
	}
}

func (m *Monpkgs) bindOpts(b squirrel.SelectBuilder, opts *MonpkgOpts) squirrel.SelectBuilder {
	if opts == nil {
		return b
	}

	if len(opts.MonlIDs) > 0 {
		b = b.Where(squirrel.Eq{"monl_id": opts.MonlIDs})
	}

	if len(opts.PkgIDs) > 0 {
		b = b.Where(squirrel.Eq{"pkg_id": opts.PkgIDs})
	}

	return b
}

func (m *Monpkgs) scan(row database.RowScanner) (*model.Monpkg, error) {
	var monpkg model.Monpkg
	err := row.Scan(
		&monpkg.ID,
		&monpkg.MonlID,
		&monpkg.PkgID,
		&monpkg.Kind)
	if err != nil {
		return nil, err
	}
	return &monpkg, nil
}

func (m *Monpkgs) Create(ctx context.Context, monpkg *model.Monpkg) error {
	monpkg2 := *monpkg
	monpkg2.ID = newID()

	qb := m.RunnableBuilder(ctx).
		Insert(m.table()).
		Columns(
			"id",
			"monl_id",
			"pkg_id",
			"kind").
		Values(
			monpkg2.ID,
			monpkg2.MonlID,
			monpkg2.PkgID,
			monpkg2.Kind)
	_, err := qb.Exec()
	if err != nil {
		return err
	}
	*monpkg = monpkg2
	return nil
}

func (m *Monpkgs) Update(ctx context.Context, monpkg *model.Monpkg) error {
	monpkg2 := *monpkg

	qb := m.RunnableBuilder(ctx).
		Update(m.table()).
		Set("monl_id", monpkg2.MonlID).
		Set("pkg_id", monpkg2.PkgID).
		Set("kind", monpkg2.Kind).
		Where("id = ?", monpkg2.ID).
		Limit(1)
	_, err := qb.Exec()
	if err != nil {
		return err
	}
	*monpkg = monpkg2
	return nil
}

func (m *Monpkgs) Delete(ctx context.Context, id string) (int64, error) {
	qb := m.RunnableBuilder(ctx).
		Delete(m.table()).
		Where("id = ?", id).
		Limit(1)
	res, err := qb.Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
