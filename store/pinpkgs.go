package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
)

type Pinpkgs struct {
	*Store
}

type PinpkgOpts struct {
	ListOpts
	PinlIDs []string
	PkgIDs  []string

	joinPinls bool
	joinPkgs  bool
}

func NewPinpkgs(s *Store) *Pinpkgs {
	return &Pinpkgs{s}
}

func (p Pinpkgs) table() string {
	return "pinpkgs"
}

func (p *Pinpkgs) List(ctx context.Context, opts *PinpkgOpts) (model.PinpkgList, error) {
	if opts == nil {
		opts = &PinpkgOpts{}
	}

	qb := p.RunnableBuilder(ctx).
		Select(p.columns()...).From(p.table())
	qb = p.bindOpts(qb, opts)
	qb = addPagination(qb, opts)
	rows, err := qb.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := make([]*model.Pinpkg, 0)
	for rows.Next() {
		pinpkg, err := p.scan(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, pinpkg)
	}
	return list, nil
}

func (p *Pinpkgs) Count(ctx context.Context, opts *PinpkgOpts) (int64, error) {
	if opts == nil {
		opts = &PinpkgOpts{}
	}

	qb := p.RunnableBuilder(ctx).
		Select("count(*)").From(p.table())
	qb = p.bindOpts(qb, opts)
	row := qb.QueryRow()
	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (p *Pinpkgs) Find(ctx context.Context, id string) (*model.Pinpkg, error) {
	qb := p.RunnableBuilder(ctx).
		Select(p.columns()...).From(p.table()).
		Where("id = ?", id)
	row := qb.QueryRow()
	pinpkg, err := p.scan(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return pinpkg, nil
}

func (p *Pinpkgs) FindOrCreate(ctx context.Context, data *model.Pinpkg) (*model.Pinpkg, error) {
	found, err := p.List(ctx, &PinpkgOpts{
		PinlIDs: []string{data.PinlID},
		PkgIDs:  []string{data.PkgID},
	})
	if err != nil {
		return nil, err
	}
	if len(found) > 0 {
		return found[0], nil
	}

	pinpkg := *data
	err = p.Create(ctx, &pinpkg)
	if err != nil {
		return nil, err
	}
	return &pinpkg, nil
}

func (p *Pinpkgs) ListWithPinl(ctx context.Context, opts *PinpkgOpts) (model.PinpkgList, error) {
	if opts == nil {
		opts = &PinpkgOpts{}
	}
	opts = opts.JoinPinls()

	qb := p.RunnableBuilder(ctx).
		Select(p.columns()...).
		Columns(Pinls{}.columns()...).
		From(p.table())
	qb = p.bindOpts(qb, opts)
	qb = addPagination(qb, opts)
	rows, err := qb.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := make([]*model.Pinpkg, 0)
	for rows.Next() {
		var (
			mpp model.Pinpkg
			mp  model.Pinl
		)
		scanCols := append(p.scanColumns(&mpp), Pinls{}.scanColumns(&mp)...)
		err := rows.Scan(scanCols...)
		if err != nil {
			return nil, err
		}
		mpp.Pinl = &mp
		list = append(list, &mpp)
	}
	return list, nil
}

func (p *Pinpkgs) ListWithPkg(ctx context.Context, opts *PinpkgOpts) (model.PinpkgList, error) {
	if opts == nil {
		opts = &PinpkgOpts{}
	}
	opts = opts.JoinPkgs()

	qb := p.RunnableBuilder(ctx).
		Select(p.columns()...).
		Columns(Pkgs{}.columns()...).
		From(p.table())
	qb = p.bindOpts(qb, opts)
	qb = addPagination(qb, opts)
	rows, err := qb.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := make([]*model.Pinpkg, 0)
	for rows.Next() {
		var (
			mpp model.Pinpkg
			mp  model.Pkg
		)
		scanCols := append(p.scanColumns(&mpp), Pkgs{}.scanColumns(&mp)...)
		err := rows.Scan(scanCols...)
		if err != nil {
			return nil, err
		}
		mpp.Pkg = &mp
		list = append(list, &mpp)
	}
	return list, nil
}

func (p Pinpkgs) bindOpts(b squirrel.SelectBuilder, opts *PinpkgOpts) squirrel.SelectBuilder {
	if opts == nil {
		return b
	}

	if len(opts.PinlIDs) > 0 {
		b = b.Where(squirrel.Eq{"pinl_id": opts.PinlIDs})
	}

	if len(opts.PkgIDs) > 0 {
		b = b.Where(squirrel.Eq{"pkg_id": opts.PkgIDs})
	}

	if opts.joinPinls {
		b = b.Join(fmt.Sprintf("%s ON %[1]s.id = %s.pinl_id", Pinls{}.table(), p.table()))
	}

	if opts.joinPkgs {
		b = b.Join(fmt.Sprintf("%s ON %[1]s.id = %s.pkg_id", Pkgs{}.table(), p.table()))
	}

	return b
}

func (p Pinpkgs) columns() []string {
	return []string{
		p.table() + ".id",
		p.table() + ".pinl_id",
		p.table() + ".pkg_id",
	}
}

func (p Pinpkgs) scanColumns(pinpkg *model.Pinpkg) []interface{} {
	return []interface{}{
		&pinpkg.ID,
		&pinpkg.PinlID,
		&pinpkg.PkgID,
	}
}

func (p Pinpkgs) scan(row database.RowScanner) (*model.Pinpkg, error) {
	var pinpkg model.Pinpkg
	err := row.Scan(p.scanColumns(&pinpkg)...)
	if err != nil {
		return nil, err
	}
	return &pinpkg, nil
}

func (p *Pinpkgs) Create(ctx context.Context, pinpkg *model.Pinpkg) error {
	pinpkg2 := *pinpkg
	pinpkg2.ID = newID()

	qb := p.RunnableBuilder(ctx).
		Insert(p.table()).
		Columns(
			"id",
			"pinl_id",
			"pkg_id").
		Values(
			pinpkg2.ID,
			pinpkg2.PinlID,
			pinpkg2.PkgID)
	_, err := qb.Exec()
	if err != nil {
		return err
	}
	*pinpkg = pinpkg2
	return nil
}

func (p *Pinpkgs) Update(ctx context.Context, pinpkg *model.Pinpkg) error {
	pinpkg2 := *pinpkg

	qb := p.RunnableBuilder(ctx).
		Update(p.table()).
		Set("pinl_id", pinpkg2.PinlID).
		Set("pkg_id", pinpkg2.PkgID).
		Where("id = ?", pinpkg2.ID)
	_, err := qb.Exec()
	if err != nil {
		return err
	}
	*pinpkg = pinpkg2
	return nil
}

func (p *Pinpkgs) Delete(ctx context.Context, id string) (int64, error) {
	qb := p.RunnableBuilder(ctx).
		Delete(p.table()).
		Where("id = ?", id)
	res, err := qb.Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (o *PinpkgOpts) JoinPinls() *PinpkgOpts {
	o2 := *o
	o2.joinPinls = true
	return &o2
}

func (o *PinpkgOpts) JoinPkgs() *PinpkgOpts {
	o2 := *o
	o2.joinPkgs = true
	return &o2
}
