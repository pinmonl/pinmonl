package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
)

// PkgOpts defines the filter parameters on Pkg.
type PkgOpts struct {
	ListOpts
	MonlURL   string
	MonlID    string
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
	br.From = pkgTB
	stmt := br.String()
	rows, err := e.NamedQuery(stmt, args)
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
	stmt := database.SelectBuilder{
		From:  pkgTB,
		Where: []string{"id = :id"},
		Limit: 1,
	}.String()
	rows, err := e.NamedQuery(stmt, m)
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
	stmt := database.InsertBuilder{
		Into: pkgTB,
		Fields: map[string]interface{}{
			"id":          nil,
			"monl_id":     nil,
			"url":         nil,
			"vendor":      nil,
			"vendor_uri":  nil,
			"title":       nil,
			"description": nil,
			"readme":      nil,
			"labels":      nil,
			"image_id":    nil,
			"created_at":  nil,
		},
	}.String()
	_, err := e.NamedExec(stmt, m2)
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
	stmt := database.UpdateBuilder{
		From: pkgTB,
		Fields: map[string]interface{}{
			"monl_id":     nil,
			"url":         nil,
			"vendor":      nil,
			"vendor_uri":  nil,
			"title":       nil,
			"description": nil,
			"readme":      nil,
			"labels":      nil,
			"image_id":    nil,
			"updated_at":  nil,
		},
		Where: []string{"id = :id"},
	}.String()
	_, err := e.NamedExec(stmt, m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// Delete removes Pkg by ID.
func (s *dbPkgStore) Delete(ctx context.Context, m *model.Pkg) error {
	e := s.Execer(ctx)
	stmt := database.DeleteBuilder{
		From:  pkgTB,
		Where: []string{"id = :id"},
	}.String()
	_, err := e.NamedExec(stmt, m)
	return err
}

func bindPkgOpts(opts *PkgOpts) (database.SelectBuilder, map[string]interface{}) {
	br := database.SelectBuilder{}
	if opts == nil {
		return br, nil
	}

	br = bindListOpts(opts.ListOpts)
	args := make(map[string]interface{})
	if opts.MonlID != "" {
		br.Where = append(br.Where, "monl_id = :monl_id")
		args["monl_id"] = opts.MonlID
	}

	if opts.MonlURL != "" {
		sq := database.SelectBuilder{
			Columns: []string{"1"},
			From:    monlTB,
			Where: []string{
				fmt.Sprintf("%s.monl_id = id", pkgTB),
				"url = :monl_url",
			},
		}
		args["monl_url"] = opts.MonlURL
		br.Where = append(br.Where, fmt.Sprintf("EXISTS (%s)", sq.String()))
	}

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
