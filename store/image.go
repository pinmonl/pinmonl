package store

import (
	"context"
	"database/sql"

	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
)

// ImageOpts defines the parameters for image filtering.
type ImageOpts struct {
	ListOpts
	Target model.Morphable
}

// ImageStore defines the services of image.
type ImageStore interface {
	List(context.Context, *ImageOpts) ([]model.Image, error)
	Find(context.Context, *model.Image) error
	Create(context.Context, *model.Image) error
	Update(context.Context, *model.Image) error
	Delete(context.Context, *model.Image) error
}

// NewImageStore creates image store.
func NewImageStore(s Store) ImageStore {
	return &dbImageStore{s}
}

type dbImageStore struct {
	Store
}

// List retrieves images by the filter parameters.
func (s *dbImageStore) List(ctx context.Context, opts *ImageOpts) ([]model.Image, error) {
	e := s.Exter(ctx)
	br, args := bindImageOpts(opts)
	br.From = imageTB
	stmt := br.String()
	rows, err := e.NamedQuery(stmt, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Image
	for rows.Next() {
		var m model.Image
		err = rows.StructScan(&m)
		if err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

// Find retrieves image by id.
func (s *dbImageStore) Find(ctx context.Context, m *model.Image) error {
	e := s.Exter(ctx)
	stmt := database.SelectBuilder{
		From:  imageTB,
		Where: []string{"id = :id"},
		Limit: 1,
	}.String()
	rows, err := e.NamedQuery(stmt, m)
	if err != nil {
		return err
	}
	if !rows.Next() {
		return sql.ErrNoRows
	}

	var m2 model.Image
	err = rows.StructScan(&m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// Create inserts the fields of image with generated id.
func (s *dbImageStore) Create(ctx context.Context, m *model.Image) error {
	m2 := *m
	m2.ID = newUID()
	m2.CreatedAt = timestamp()
	e := s.Exter(ctx)
	stmt := database.InsertBuilder{
		Into: imageTB,
		Fields: map[string]interface{}{
			"id":          nil,
			"target_id":   nil,
			"target_name": nil,
			"kind":        nil,
			"sort":        nil,
			"filename":    nil,
			"content":     nil,
			"description": nil,
			"size":        nil,
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

// Update updates the fields of image by id.
func (s *dbImageStore) Update(ctx context.Context, m *model.Image) error {
	m2 := *m
	m2.UpdatedAt = timestamp()
	e := s.Exter(ctx)
	stmt := database.UpdateBuilder{
		From: imageTB,
		Fields: map[string]interface{}{
			"target_id":   nil,
			"target_name": nil,
			"kind":        nil,
			"sort":        nil,
			"filename":    nil,
			"content":     nil,
			"description": nil,
			"size":        nil,
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

// Delete removes image by id.
func (s *dbImageStore) Delete(ctx context.Context, m *model.Image) error {
	e := s.Exter(ctx)
	stmt := database.DeleteBuilder{
		From:  imageTB,
		Where: []string{"id = :id"},
	}.String()
	_, err := e.NamedExec(stmt, m)
	return err
}

func bindImageOpts(opts *ImageOpts) (database.SelectBuilder, map[string]interface{}) {
	br := database.SelectBuilder{}
	if opts == nil {
		return br, nil
	}

	args := make(map[string]interface{})
	if opts.Target != nil {
		br.Where = append(br.Where, "target_id = :target_id")
		br.Where = append(br.Where, "target_name = :target_name")
		args["target_id"] = opts.Target.MorphKey()
		args["target_name"] = opts.Target.MorphName()
	}

	return br, args
}
