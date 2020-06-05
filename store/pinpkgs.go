package store

import (
	"context"

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
}

func NewPinpkgs(s *Store) *Pinpkgs {
	return &Pinpkgs{s}
}

func (p *Pinpkgs) table() string {
	return "pinpkgs"
}

func (p *Pinpkgs) List(ctx context.Context, opts *PinpkgOpts) ([]*model.Pinpkg, error) {
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
	var list []*model.Pinpkg
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
		Where("id = ?", id).
		Limit(1)
	row := qb.QueryRow()
	pinpkg, err := p.scan(row)
	if err != nil {
		return nil, err
	}
	return pinpkg, nil
}

func (p *Pinpkgs) columns() []string {
	return []string{
		"id",
		"pinl_id",
		"pkg_id",
	}
}

func (p *Pinpkgs) bindOpts(b squirrel.SelectBuilder, opts *PinpkgOpts) squirrel.SelectBuilder {
	if opts == nil {
		return b
	}

	if len(opts.PinlIDs) > 0 {
		b = b.Where(squirrel.Eq{"pinl_id": opts.PinlIDs})
	}

	if len(opts.PkgIDs) > 0 {
		b = b.Where(squirrel.Eq{"pkg_id": opts.PkgIDs})
	}

	return b
}

func (p *Pinpkgs) scan(row database.RowScanner) (*model.Pinpkg, error) {
	var pinpkg model.Pinpkg
	err := row.Scan(
		&pinpkg.ID,
		&pinpkg.PinlID,
		&pinpkg.PkgID)
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
		Where("id = ?", pinpkg2.ID).
		Limit(1)
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
		Where("id = ?", id).
		Limit(1)
	res, err := qb.Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
