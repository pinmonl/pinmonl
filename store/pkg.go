package store

import (
	"context"
	"database/sql"

	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
)

// PkgOpts defines the filter parameters on Pkg.
type PkgOpts struct {
	ListOpts
	Vendor    string
	VendorURI string
}

// PkgStore defines the services of Pkg store.
type PkgStore interface {
	List(context.Context, *PkgOpts) ([]model.Pkg, error)
	Find(context.Context, *model.Pkg) error
	Create(context.Context, *model.Pkg) error
	Update(context.Context, *model.Pkg) error
	Delete(context.Context, *model.Pkg) error
}

// NewPkgStore creates Pkg store.
func NewPkgStore(s Store) PkgStore {
	return &dbPkgStore{s}
}

type dbPkgStore struct {
	Store
}

// List retrieves Pkg and filters by PkgOpts.
func (s *dbPkgStore) List(ctx context.Context, opts *PkgOpts) ([]model.Pkg, error) {
	e := s.Queryer(ctx)
	br, args := bindPkgOpts(opts)
	rows, err := e.NamedQuery(br.String(), args.Map())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Pkg
	for rows.Next() {
		var m model.Pkg
		err = rows.StructScan(&m)
		if err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

// Find retrieves Pkg by ID.
func (s *dbPkgStore) Find(ctx context.Context, m *model.Pkg) error {
	e := s.Queryer(ctx)
	br, _ := bindPkgOpts(nil)
	br.Where = []string{"id = :pkg_id"}
	br.Limit = 1
	rows, err := e.NamedQuery(br.String(), m)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return sql.ErrNoRows
	}
	var m2 model.Pkg
	err = rows.StructScan(&m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// Create inserts the fields of Pkg with generated ID.
func (s *dbPkgStore) Create(ctx context.Context, m *model.Pkg) error {
	m2 := *m
	m2.ID = newUID()
	m2.CreatedAt = timestamp()
	e := s.Execer(ctx)
	br := database.InsertBuilder{
		Into: pkgTB,
		Fields: map[string]interface{}{
			"id":          ":pkg_id",
			"url":         ":pkg_url",
			"vendor":      ":pkg_vendor",
			"vendor_uri":  ":pkg_vendor_uri",
			"title":       ":pkg_title",
			"description": ":pkg_description",
			"readme":      ":pkg_readme",
			"labels":      ":pkg_labels",
			"image_id":    ":pkg_image_id",
			"created_at":  ":pkg_created_at",
		},
	}
	_, err := e.NamedExec(br.String(), m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// Update updates the fields of Pkg by ID.
func (s *dbPkgStore) Update(ctx context.Context, m *model.Pkg) error {
	m2 := *m
	m2.UpdatedAt = timestamp()
	e := s.Execer(ctx)
	br := database.UpdateBuilder{
		From: pkgTB,
		Fields: map[string]interface{}{
			"url":         ":pkg_url",
			"vendor":      ":pkg_vendor",
			"vendor_uri":  ":pkg_vendor_uri",
			"title":       ":pkg_title",
			"description": ":pkg_description",
			"readme":      ":pkg_readme",
			"labels":      ":pkg_labels",
			"image_id":    ":pkg_image_id",
			"updated_at":  ":pkg_updated_at",
		},
		Where: []string{"id = :pkg_id"},
	}
	_, err := e.NamedExec(br.String(), m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// Delete removes Pkg by ID.
func (s *dbPkgStore) Delete(ctx context.Context, m *model.Pkg) error {
	e := s.Execer(ctx)
	br := database.DeleteBuilder{
		From:  pkgTB,
		Where: []string{"id = :pkg_id"},
	}
	_, err := e.NamedExec(br.String(), m)
	return err
}

func bindPkgOpts(opts *PkgOpts) (database.SelectBuilder, database.QueryVars) {
	br := database.SelectBuilder{
		From: pkgTB,
		Columns: database.NamespacedColumn(
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
		),
	}
	if opts == nil {
		return br, nil
	}

	br = appendListOpts(br, opts.ListOpts)
	args := database.QueryVars{}

	if opts.Vendor != "" {
		br.Where = append(br.Where, "vendor = :vendor")
		args["vendor"] = opts.Vendor
	}
	if opts.VendorURI != "" {
		br.Where = append(br.Where, "vendor_uri = :vendor_uri")
		args["vendor_uri"] = opts.VendorURI
	}

	return br, args
}
