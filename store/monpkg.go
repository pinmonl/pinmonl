package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
)

// MonpkgOpts defines the parameters for monlpkg filtering.
type MonpkgOpts struct {
	ListOpts
	MonlURL string
	Monls   []model.Monl
	MonlIDs []string
	Pkgs    []model.Pkg
	PkgIDs  []string

	joinPkg  bool
	joinMonl bool
}

// MonpkgStore defines the service of monpkg.
type MonpkgStore interface {
	List(context.Context, *MonpkgOpts) ([]model.Monpkg, error)
	ListMonls(context.Context, *MonpkgOpts) (map[string][]model.Monl, error)
	ListPkgs(context.Context, *MonpkgOpts) (map[string][]model.Pkg, error)
}

// NewMonpkgStore creates monpkg store.
func NewMonpkgStore(s Store) MonpkgStore {
	return &dbMonpkgStore{s}
}

type dbMonpkgStore struct {
	Store
}

// List retrieves monpkgs with options.
func (s *dbMonpkgStore) List(ctx context.Context, opts *MonpkgOpts) ([]model.Monpkg, error) {
	e := s.Queryer(ctx)
	br, args := bindMonpkgOpts(opts)
	rows, err := e.NamedQuery(br.String(), args.Map())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Monpkg
	for rows.Next() {
		var m model.Monpkg
		err = rows.StructScan(&m)
		if err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

// ListMonls retrieves monls with options.
func (s *dbMonpkgStore) ListMonls(ctx context.Context, opts *MonpkgOpts) (map[string][]model.Monl, error) {
	if opts == nil {
		opts = &MonpkgOpts{}
	}
	opts.joinMonl = true

	e := s.Queryer(ctx)
	br, args := bindMonpkgOpts(opts)
	rows, err := e.NamedQuery(br.String(), args.Map())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list map[string][]model.Monl
	for rows.Next() {
		var m model.Monpkg
		err = rows.StructScan(&m)
		if err != nil {
			return nil, err
		}
		k := m.PkgID
		list[k] = append(list[k], *m.Monl)
	}
	return list, nil
}

// ListMonls retrieves pkgs with options.
func (s *dbMonpkgStore) ListPkgs(ctx context.Context, opts *MonpkgOpts) (map[string][]model.Pkg, error) {
	if opts == nil {
		opts = &MonpkgOpts{}
	}
	opts.joinPkg = true

	e := s.Queryer(ctx)
	br, args := bindMonpkgOpts(opts)
	rows, err := e.NamedQuery(br.String(), args.Map())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list map[string][]model.Pkg
	for rows.Next() {
		var m model.Monpkg
		err = rows.StructScan(&m)
		if err != nil {
			return nil, err
		}
		k := m.MonlID
		list[k] = append(list[k], *m.Pkg)
	}
	return list, nil
}

func bindMonpkgOpts(opts *MonpkgOpts) (database.SelectBuilder, database.QueryVars) {
	br := database.SelectBuilder{
		From: monpkgTB,
		Columns: database.NamespacedColumn(
			[]string{
				"monl_id AS monpkg_monl_id",
				"pkg_id AS monpkg_pkg_id",
				"tie AS monpkg_tie",
			},
			monpkgTB,
		),
	}
	if opts == nil {
		return br, nil
	}

	br = appendListOpts(br, opts.ListOpts)
	args := database.QueryVars{}

	if opts.MonlURL != "" {
		br.Where = append(br.Where, fmt.Sprintf("%s.url = :monl_url", monlTB))
		args.Set("monl_url", opts.MonlURL)
		opts.joinMonl = true
	}
	if opts.Monls != nil {
		opts.MonlIDs = append(opts.MonlIDs, (model.MonlList)(opts.Monls).Keys()...)
	}
	if opts.MonlIDs != nil {
		ks, ids := bindQueryIDs("monl_ids", opts.MonlIDs)
		args.AppendStringMap(ids)
		br.Where = append(br.Where, fmt.Sprintf("%s.monl_id IN (%s)", monpkgTB, strings.Join(ks, ",")))
	}

	if opts.Pkgs != nil {
		opts.PkgIDs = append(opts.PkgIDs, (model.PkgList)(opts.Pkgs).Keys()...)
	}
	if opts.PkgIDs != nil {
		ks, ids := bindQueryIDs("pkgs_ids", opts.PkgIDs)
		args.AppendStringMap(ids)
		br.Where = append(br.Where, fmt.Sprintf("%s.pkg_id IN (%s)", monpkgTB, strings.Join(ks, ",")))
	}

	if opts.joinMonl {
		br.Columns = append(br.Columns, database.NamespacedColumn(
			[]string{
				"id AS monl_id",
				"url AS monl_url",
				"title AS monl_title",
				"description AS monl_description",
				"readme AS monl_readme",
				"image_id AS monl_image_id",
				"created_at AS monl_created_at",
				"updated_at AS monl_updated_at",
			},
			monlTB,
		)...)
		br.Join = append(br.Join, fmt.Sprintf("%s ON %[1]s.id = %s.monl_id", monlTB, monpkgTB))
	}
	if opts.joinPkg {
		br.Columns = append(br.Columns, database.NamespacedColumn(
			[]string{
				"id AS pkg_id",
				"url AS pkg_url",
				"vendor AS pkg_vendor",
				"vendor_uri AS pkg_vendor_uri",
				"title AS pkg_title",
				"description AS pkg_description",
				"readme AS pkg_readme",
				"image_id AS pkg_image_id",
				"labels AS pkg_labels",
				"created_at AS pkg_created_at",
				"updated_at AS pkg_updated_at",
			},
			pkgTB,
		)...)
		br.Join = append(br.Join, fmt.Sprintf("%s ON %[1]s.id = %s.pkg_id", pkgTB, monpkgTB))
	}

	return br, args
}
