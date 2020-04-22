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
	Create(context.Context, *model.Monpkg) error
	Delete(context.Context, *model.Monpkg) error
	Associate(context.Context, model.Monl, model.Pkg) error
	AssociateMany(context.Context, model.Monl, []model.Pkg) error
	Dissociate(context.Context, model.Monl, model.Pkg) error
	DissociateMany(context.Context, model.Monl, []model.Pkg) error
	ReAssociateMany(context.Context, model.Monl, []model.Pkg) error
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

	list := make(map[string][]model.Monl)
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

	list := make(map[string][]model.Pkg)
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

func (s *dbMonpkgStore) Create(ctx context.Context, m *model.Monpkg) error {
	m2 := *m
	e := s.Execer(ctx)
	br := database.InsertBuilder{
		Into: monpkgTB,
		Fields: map[string]interface{}{
			"monl_id": ":monpkg_monl_id",
			"pkg_id":  ":monpkg_pkg_id",
			"tie":     ":monpkg_tie",
		},
	}
	_, err := e.NamedExec(br.String(), m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

func (s *dbMonpkgStore) Delete(ctx context.Context, m *model.Monpkg) error {
	e := s.Execer(ctx)
	br := database.DeleteBuilder{
		From: monpkgTB,
		Where: []string{
			"monl_id = :monpkg_monl_id",
			"pkg_id = :monpkg_pkg_id",
		},
	}
	_, err := e.NamedExec(br.String(), m)
	return err
}

func (s *dbMonpkgStore) Associate(ctx context.Context, monl model.Monl, pkg model.Pkg) error {
	return s.Create(ctx, &model.Monpkg{
		MonlID: monl.ID,
		PkgID:  pkg.ID,
	})
}

func (s *dbMonpkgStore) AssociateMany(ctx context.Context, monl model.Monl, pkgs []model.Pkg) error {
	for _, pkg := range pkgs {
		err := s.Associate(ctx, monl, pkg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *dbMonpkgStore) Dissociate(ctx context.Context, monl model.Monl, pkg model.Pkg) error {
	return s.Delete(ctx, &model.Monpkg{
		MonlID: monl.ID,
		PkgID:  pkg.ID,
	})
}

func (s *dbMonpkgStore) DissociateMany(ctx context.Context, monl model.Monl, pkgs []model.Pkg) error {
	for _, pkg := range pkgs {
		err := s.Dissociate(ctx, monl, pkg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *dbMonpkgStore) ReAssociateMany(ctx context.Context, monl model.Monl, pkgs []model.Pkg) error {
	err := s.clear(ctx, monl)
	if err != nil {
		return err
	}
	return s.AssociateMany(ctx, monl, pkgs)
}

func (s *dbMonpkgStore) clear(ctx context.Context, monl model.Monl) error {
	e := s.Execer(ctx)
	br := database.DeleteBuilder{
		From:  monpkgTB,
		Where: []string{"monl_id = :monpkg_monl_id"},
	}
	_, err := e.NamedExec(br.String(), monl)
	return err
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
		br.Join = append(br.Join, fmt.Sprintf("INNER JOIN %s ON %[1]s.id = %s.monl_id", monlTB, monpkgTB))
	}
	if opts.joinPkg {
		br.Columns = append(br.Columns, database.NamespacedColumn(
			[]string{
				"id AS pkg_id",
				"url AS pkg_url",
				"provider AS pkg_provider",
				"provider_host AS pkg_provider_host",
				"provider_uri AS pkg_provider_uri",
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
		br.Join = append(br.Join, fmt.Sprintf("INNER JOIN %s ON %[1]s.id = %s.pkg_id", pkgTB, monpkgTB))
	}

	return br, args
}
