package store

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pkgs/pkguri"
)

type Pkgs struct {
	*Store
}

type PkgOpts struct {
	ListOpts
	Provider     string
	ProviderHost string
	ProviderURI  string
}

func NewPkgs(s *Store) *Pkgs {
	return &Pkgs{s}
}

func (p Pkgs) table() string {
	return "pkgs"
}

func (p *Pkgs) List(ctx context.Context, opts *PkgOpts) (model.PkgList, error) {
	if opts == nil {
		opts = &PkgOpts{}
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
	list := make([]*model.Pkg, 0)
	for rows.Next() {
		pkg, err := p.scan(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, pkg)
	}
	return list, nil
}

func (p *Pkgs) Count(ctx context.Context, opts *PkgOpts) (int64, error) {
	if opts == nil {
		opts = &PkgOpts{}
	}

	qb := p.RunnableBuilder(ctx).
		Select("count(*)").From(p.table())
	row := qb.QueryRow()
	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (p *Pkgs) Find(ctx context.Context, id string) (*model.Pkg, error) {
	qb := p.RunnableBuilder(ctx).
		Select(p.columns()...).From(p.table()).
		Where("id = ?", id)
	row := qb.QueryRow()
	pkg, err := p.scan(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return pkg, nil
}

func (p *Pkgs) FindURI(ctx context.Context, pu *pkguri.PkgURI) (*model.Pkg, error) {
	qb := p.RunnableBuilder(ctx).
		Select(p.columns()...).From(p.table()).
		Where("provider = ?", pu.Provider).
		Where("provider_host = ?", pu.Host).
		Where("provider_uri = ?", pu.URI).
		Where("provider_proto = ?", pu.Proto)
	row := qb.QueryRow()
	pkg, err := p.scan(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return pkg, nil
}

func (p Pkgs) bindOpts(b squirrel.SelectBuilder, opts *PkgOpts) squirrel.SelectBuilder {
	if opts == nil {
		return b
	}

	if opts.Provider != "" {
		b = b.Where("provider = ?", opts.Provider)
	}

	if opts.ProviderHost != "" {
		b = b.Where("provider_host = ?", opts.ProviderHost)
	}

	if opts.ProviderURI != "" {
		b = b.Where("provider_uri = ?", opts.ProviderURI)
	}

	return b
}

func (p Pkgs) columns() []string {
	return []string{
		p.table() + ".id",
		p.table() + ".url",
		p.table() + ".provider",
		p.table() + ".provider_host",
		p.table() + ".provider_uri",
		p.table() + ".provider_proto",
		p.table() + ".fetched_at",
		p.table() + ".created_at",
		p.table() + ".updated_at",
	}
}

func (p Pkgs) scanColumns(pkg *model.Pkg) []interface{} {
	return []interface{}{
		&pkg.ID,
		&pkg.URL,
		&pkg.Provider,
		&pkg.ProviderHost,
		&pkg.ProviderURI,
		&pkg.ProviderProto,
		&pkg.FetchedAt,
		&pkg.CreatedAt,
		&pkg.UpdatedAt,
	}
}

func (p Pkgs) scan(row database.RowScanner) (*model.Pkg, error) {
	var pkg model.Pkg
	err := row.Scan(p.scanColumns(&pkg)...)
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}

func (p *Pkgs) Create(ctx context.Context, pkg *model.Pkg) error {
	pkg2 := *pkg
	pkg2.ID = newID()
	pkg2.CreatedAt = timestamp()
	pkg2.UpdatedAt = timestamp()

	qb := p.RunnableBuilder(ctx).
		Insert(p.table()).
		Columns(
			"id",
			"url",
			"provider",
			"provider_host",
			"provider_uri",
			"provider_proto",
			"fetched_at",
			"created_at",
			"updated_at").
		Values(
			pkg2.ID,
			pkg2.URL,
			pkg2.Provider,
			pkg2.ProviderHost,
			pkg2.ProviderURI,
			pkg2.ProviderProto,
			pkg2.FetchedAt,
			pkg2.CreatedAt,
			pkg2.UpdatedAt)
	_, err := qb.Exec()
	if err != nil {
		return err
	}
	*pkg = pkg2
	return nil
}

func (p *Pkgs) Update(ctx context.Context, pkg *model.Pkg) error {
	pkg2 := *pkg
	pkg2.UpdatedAt = timestamp()

	qb := p.RunnableBuilder(ctx).
		Update(p.table()).
		Set("url", pkg2.URL).
		Set("provider", pkg2.Provider).
		Set("provider_host", pkg2.ProviderHost).
		Set("provider_uri", pkg2.ProviderURI).
		Set("provider_proto", pkg2.ProviderProto).
		Set("fetched_at", pkg2.FetchedAt).
		Set("updated_at", pkg2.UpdatedAt).
		Where("id = ?", pkg2.ID)
	_, err := qb.Exec()
	if err != nil {
		return err
	}
	*pkg = pkg2
	return nil
}

func (p *Pkgs) Delete(ctx context.Context, id string) (int64, error) {
	qb := p.RunnableBuilder(ctx).
		Delete(p.table()).
		Where("id = ?", id)
	res, err := qb.Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
