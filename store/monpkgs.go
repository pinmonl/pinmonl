package store

import (
	"context"
	"database/sql"
	"fmt"

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

	joinPkgs bool
}

func NewMonpkgs(s *Store) *Monpkgs {
	return &Monpkgs{s}
}

func (m Monpkgs) table() string {
	return "monpkgs"
}

func (m *Monpkgs) List(ctx context.Context, opts *MonpkgOpts) (model.MonpkgList, error) {
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
	list := make([]*model.Monpkg, 0)
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
		Where("id = ?", id)
	row := qb.QueryRow()
	monpkg, err := m.scan(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return monpkg, nil
}

func (m *Monpkgs) FindOrCreate(ctx context.Context, data *model.Monpkg) (*model.Monpkg, error) {
	found, err := m.List(ctx, &MonpkgOpts{
		MonlIDs: []string{data.MonlID},
		PkgIDs:  []string{data.PkgID},
	})
	if err != nil {
		return nil, err
	}
	if len(found) > 0 {
		return found[0], nil
	}

	monpkg := *data
	err = m.Create(ctx, &monpkg)
	if err != nil {
		return nil, err
	}
	return &monpkg, nil
}

func (m *Monpkgs) ListWithPkg(ctx context.Context, opts *MonpkgOpts) (model.MonpkgList, error) {
	if opts == nil {
		opts = &MonpkgOpts{}
	}
	opts = opts.JoinPkgs()

	qb := m.RunnableBuilder(ctx).
		Select(m.columns()...).
		Columns(Pkgs{}.columns()...).
		From(m.table())
	qb = m.bindOpts(qb, opts)
	qb = addPagination(qb, opts)
	rows, err := qb.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := make([]*model.Monpkg, 0)
	for rows.Next() {
		var (
			mmp model.Monpkg
			mp  model.Pkg
		)
		scanCols := append(m.scanColumns(&mmp), Pkgs{}.scanColumns(&mp)...)
		err := rows.Scan(scanCols...)
		if err != nil {
			return nil, err
		}
		mmp.Pkg = &mp
		list = append(list, &mmp)
	}
	return list, nil
}

func (m Monpkgs) bindOpts(b squirrel.SelectBuilder, opts *MonpkgOpts) squirrel.SelectBuilder {
	if opts == nil {
		return b
	}

	if len(opts.MonlIDs) > 0 {
		b = b.Where(squirrel.Eq{m.table() + ".monl_id": opts.MonlIDs})
	}

	if len(opts.PkgIDs) > 0 {
		b = b.Where(squirrel.Eq{m.table() + ".pkg_id": opts.PkgIDs})
	}

	if opts.joinPkgs {
		b = b.LeftJoin(fmt.Sprintf("%s ON %[1]s.id = %s.pkg_id", Pkgs{}.table(), m.table()))
	}

	return b
}

func (m Monpkgs) columns() []string {
	return []string{
		m.table() + ".id",
		m.table() + ".monl_id",
		m.table() + ".pkg_id",
		m.table() + ".kind",
	}
}

func (m Monpkgs) scanColumns(monpkg *model.Monpkg) []interface{} {
	return []interface{}{
		&monpkg.ID,
		&monpkg.MonlID,
		&monpkg.PkgID,
		&monpkg.Kind,
	}
}

func (m Monpkgs) scan(row database.RowScanner) (*model.Monpkg, error) {
	var monpkg model.Monpkg
	err := row.Scan(m.scanColumns(&monpkg)...)
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
		Where("id = ?", monpkg2.ID)
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
		Where("id = ?", id)
	res, err := qb.Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (o *MonpkgOpts) JoinPkgs() *MonpkgOpts {
	o2 := *o
	o2.joinPkgs = true
	return &o2
}
