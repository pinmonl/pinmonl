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
	e := s.Queryer(ctx)
	br, args := bindImageOpts(opts)
	rows, err := e.NamedQuery(br.String(), args.Map())
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
	e := s.Queryer(ctx)
	br, _ := bindImageOpts(nil)
	br.Where = []string{"id = :image_id"}
	br.Limit = 1
	rows, err := e.NamedQuery(br.String(), m)
	if err != nil {
		return err
	}
	if !rows.Next() {
		return sql.ErrNoRows
	}
	defer rows.Close()

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
	e := s.Execer(ctx)
	br := database.InsertBuilder{
		Into: imageTB,
		Fields: map[string]interface{}{
			"id":           ":image_id",
			"target_id":    ":image_target_id",
			"target_name":  ":image_target_name",
			"content_type": ":image_content_type",
			"sort":         ":image_sort",
			"filename":     ":image_filename",
			"content":      ":image_content",
			"description":  ":image_description",
			"size":         ":image_size",
			"created_at":   ":image_created_at",
		},
	}
	_, err := e.NamedExec(br.String(), m2)
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
	e := s.Execer(ctx)
	br := database.UpdateBuilder{
		From: imageTB,
		Fields: map[string]interface{}{
			"target_id":    ":image_target_id",
			"target_name":  ":image_target_name",
			"content_type": ":image_content_type",
			"sort":         ":image_sort",
			"filename":     ":image_filename",
			"content":      ":image_content",
			"description":  ":image_description",
			"size":         ":image_size",
			"updated_at":   ":image_updated_at",
		},
		Where: []string{"id = :image_id"},
	}
	_, err := e.NamedExec(br.String(), m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// Delete removes image by id.
func (s *dbImageStore) Delete(ctx context.Context, m *model.Image) error {
	e := s.Execer(ctx)
	br := database.DeleteBuilder{
		From:  imageTB,
		Where: []string{"id = :image_id"},
	}
	_, err := e.NamedExec(br.String(), m)
	return err
}

func bindImageOpts(opts *ImageOpts) (database.SelectBuilder, database.QueryVars) {
	br := database.SelectBuilder{
		From: imageTB,
		Columns: database.NamespacedColumn(
			[]string{
				"id AS image_id",
				"target_id AS image_target_id",
				"target_name AS image_target_name",
				"content_type AS image_content_type",
				"sort AS image_sort",
				"filename AS image_filename",
				"content AS image_content",
				"description AS image_description",
				"size AS image_size",
				"created_at AS image_created_at",
				"updated_at AS image_updated_at",
			},
			imageTB,
		),
	}
	if opts == nil {
		return br, nil
	}

	br = appendListOpts(br, opts.ListOpts)
	args := database.QueryVars{}

	if opts.Target != nil {
		br.Where = append(br.Where, "target_id = :target_id")
		br.Where = append(br.Where, "target_name = :target_name")
		args["target_id"] = opts.Target.MorphKey()
		args["target_name"] = opts.Target.MorphName()
	}

	return br, args
}
