package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
)

// MonlOpts defines the parameters for monl filtering.
type MonlOpts struct {
	ListOpts
	URL           string
	UpdatedBefore time.Time
}

// MonlStore defines the services of monl.
type MonlStore interface {
	List(context.Context, *MonlOpts) ([]model.Monl, error)
	Count(context.Context, *MonlOpts) (int64, error)
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
	e := s.Queryer(ctx)
	br, args := bindMonlOpts(opts)
	rows, err := e.NamedQuery(br.String(), args.Map())
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

// Count counts the number of monls with the filter parameters.
func (s *dbMonlStore) Count(ctx context.Context, opts *MonlOpts) (int64, error) {
	e := s.Queryer(ctx)
	br, args := bindMonlOpts(opts)
	br.Columns = []string{"COUNT(*) AS count"}
	rows, err := e.NamedQuery(br.String(), args.Map())
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, sql.ErrNoRows
	}
	var count int64
	err = rows.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Find retrieves monl by id.
func (s *dbMonlStore) Find(ctx context.Context, m *model.Monl) error {
	e := s.Queryer(ctx)
	br, _ := bindMonlOpts(nil)
	br.Where = []string{"id = :monl_id"}
	br.Limit = 1
	rows, err := e.NamedQuery(br.String(), m)
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
	e := s.Execer(ctx)
	br := database.InsertBuilder{
		Into: monlTB,
		Fields: map[string]interface{}{
			"id":          ":monl_id",
			"url":         ":monl_url",
			"title":       ":monl_title",
			"description": ":monl_description",
			"readme":      ":monl_readme",
			"image_id":    ":monl_image_id",
			"created_at":  ":monl_created_at",
		},
	}
	_, err := e.NamedExec(br.String(), m2)
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
	e := s.Execer(ctx)
	br := database.UpdateBuilder{
		From: monlTB,
		Fields: map[string]interface{}{
			"url":         ":monl_url",
			"title":       ":monl_title",
			"description": ":monl_description",
			"readme":      ":monl_readme",
			"image_id":    ":monl_image_id",
			"updated_at":  ":monl_updated_at",
		},
		Where: []string{"id = :monl_id"},
	}
	_, err := e.NamedExec(br.String(), m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// Delete removes monl by id.
func (s *dbMonlStore) Delete(ctx context.Context, m *model.Monl) error {
	e := s.Execer(ctx)
	br := database.DeleteBuilder{
		From:  monlTB,
		Where: []string{"id = :monl_id"},
	}
	_, err := e.NamedExec(br.String(), m)
	return err
}

func bindMonlOpts(opts *MonlOpts) (database.SelectBuilder, database.QueryVars) {
	br := database.SelectBuilder{
		From: monlTB,
		Columns: database.NamespacedColumn(
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
		),
	}
	if opts == nil {
		return br, nil
	}

	br = appendListOpts(br, opts.ListOpts)
	args := database.QueryVars{}

	if opts.URL != "" {
		br.Where = append(br.Where, "url = :url")
		args["url"] = opts.URL
	}

	if !opts.UpdatedBefore.IsZero() {
		br.Where = append(br.Where, "(updated_at <= :updated_before OR updated_at IS NULL)")
		args["updated_before"] = opts.UpdatedBefore
	}

	return br, args
}
