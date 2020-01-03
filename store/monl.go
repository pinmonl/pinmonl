package store

import (
	"context"
	"database/sql"

	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
)

// MonlOpts defines the parameters for monl filtering.
type MonlOpts struct {
	ListOpts
	URL       string
	Vendor    string
	VendorURI string
}

// MonlStore defines the services of monl.
type MonlStore interface {
	List(context.Context, *MonlOpts) ([]model.Monl, error)
	Find(context.Context, *model.Monl) error
	Create(context.Context, *model.Monl) error
	Update(context.Context, *model.Monl) error
	Delete(context.Context, *model.Monl) error
}

// NewMonlStore creates monl store.
func NewMonlStore(s Store) MonlStore {
	return &dbMonlStore{s}
}

type dbMonlStore struct {
	Store
}

// List retrieves monls by the filter parameters.
func (s *dbMonlStore) List(ctx context.Context, opts *MonlOpts) ([]model.Monl, error) {
	e := s.Exter(ctx)
	br, args := bindMonlOpts(opts)
	br.From = monlTB
	stmt := br.String()
	rows, err := e.NamedQuery(stmt, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Monl
	for rows.Next() {
		var m model.Monl
		err = rows.StructScan(&m)
		if err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

// Find retrieves monl by id.
func (s *dbMonlStore) Find(ctx context.Context, m *model.Monl) error {
	e := s.Exter(ctx)
	stmt := database.SelectBuilder{
		From:  monlTB,
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
	var m2 model.Monl
	err = rows.StructScan(&m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// Create inserts the fields of monl with generated id.
func (s *dbMonlStore) Create(ctx context.Context, m *model.Monl) error {
	m2 := *m
	m2.ID = newUID()
	m2.CreatedAt = timestamp()
	e := s.Exter(ctx)
	stmt := database.InsertBuilder{
		Into: monlTB,
		Fields: map[string]interface{}{
			"id":          nil,
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

// Update updates the fields of monl by id.
func (s *dbMonlStore) Update(ctx context.Context, m *model.Monl) error {
	m2 := *m
	m2.UpdatedAt = timestamp()
	e := s.Exter(ctx)
	stmt := database.UpdateBuilder{
		From: monlTB,
		Fields: map[string]interface{}{
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

// Delete removes monl by id.
func (s *dbMonlStore) Delete(ctx context.Context, m *model.Monl) error {
	e := s.Exter(ctx)
	stmt := database.DeleteBuilder{
		From:  monlTB,
		Where: []string{"id = :id"},
	}.String()
	_, err := e.NamedExec(stmt, m)
	return err
}

func bindMonlOpts(opts *MonlOpts) (database.SelectBuilder, map[string]interface{}) {
	br := database.SelectBuilder{}
	if opts != nil {
		return br, nil
	}

	br = bindListOpts(opts.ListOpts)
	args := make(map[string]interface{})
	if opts.URL != "" {
		br.Where = append(br.Where, "url = :url")
		args["url"] = opts.URL
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
