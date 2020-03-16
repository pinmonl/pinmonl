package store

import (
	"context"
	"database/sql"

	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
)

// ShareOpts defines the parameters for share filtering.
type ShareOpts struct {
	ListOpts
	UserID string
	Name   string
}

// ShareStore defines the services of share.
type ShareStore interface {
	List(context.Context, *ShareOpts) ([]model.Share, error)
	Count(context.Context, *ShareOpts) (int64, error)
	Find(context.Context, *model.Share) error
	FindByName(context.Context, *model.Share) error
	Create(context.Context, *model.Share) error
	Update(context.Context, *model.Share) error
	Delete(context.Context, *model.Share) error
}

// NewShareStore creates share store.
func NewShareStore(s Store) ShareStore {
	return &dbShareStore{s}
}

type dbShareStore struct {
	Store
}

// List retrieves shares by the filter parameters.
func (s *dbShareStore) List(ctx context.Context, opts *ShareOpts) ([]model.Share, error) {
	e := s.Queryer(ctx)
	br, args := bindShareOpts(opts)
	rows, err := e.NamedQuery(br.String(), args.Map())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Share
	for rows.Next() {
		var m model.Share
		err = rows.StructScan(&m)
		if err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

// Count counts the number of share by the filter parameter.
func (s *dbShareStore) Count(ctx context.Context, opts *ShareOpts) (int64, error) {
	e := s.Queryer(ctx)
	br, args := bindShareOpts(opts)
	br.Columns = []string{"COUNT(*) as count"}
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

// Find retrieves share by id.
func (s *dbShareStore) Find(ctx context.Context, m *model.Share) error {
	e := s.Queryer(ctx)
	br, _ := bindShareOpts(nil)
	br.Where = []string{"id = :share_id"}
	br.Limit = 1
	rows, err := e.NamedQuery(br.String(), m)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return sql.ErrNoRows
	}
	var m2 model.Share
	err = rows.StructScan(&m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// FindByName retieves user's share by name.
func (s *dbShareStore) FindByName(ctx context.Context, m *model.Share) error {
	e := s.Queryer(ctx)
	br, _ := bindShareOpts(nil)
	br.Where = []string{"user_id = :share_user_id", "name = :share_name"}
	br.Limit = 1
	rows, err := e.NamedQuery(br.String(), m)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return sql.ErrNoRows
	}
	var m2 model.Share
	err = rows.StructScan(&m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// Create inserts the fields of share with generated id.
func (s *dbShareStore) Create(ctx context.Context, m *model.Share) error {
	m2 := *m
	m2.ID = newUID()
	m2.CreatedAt = timestamp()
	e := s.Execer(ctx)
	br := database.InsertBuilder{
		Into: shareTB,
		Fields: map[string]interface{}{
			"id":          ":share_id",
			"user_id":     ":share_user_id",
			"name":        ":share_name",
			"description": ":share_description",
			"readme":      ":share_readme",
			"image_id":    ":share_image_id",
			"created_at":  ":share_created_at",
		},
	}
	_, err := e.NamedExec(br.String(), m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// Update updates the fields of share by id.
func (s *dbShareStore) Update(ctx context.Context, m *model.Share) error {
	m2 := *m
	m2.UpdatedAt = timestamp()
	e := s.Execer(ctx)
	br := database.UpdateBuilder{
		From: shareTB,
		Fields: map[string]interface{}{
			"user_id":     ":share_user_id",
			"name":        ":share_name",
			"description": ":share_description",
			"readme":      ":share_readme",
			"image_id":    ":share_image_id",
			"updated_at":  ":share_updated_at",
		},
		Where: []string{"id = :share_id"},
	}
	_, err := e.NamedExec(br.String(), m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// Delete removes share by id.
func (s *dbShareStore) Delete(ctx context.Context, m *model.Share) error {
	e := s.Execer(ctx)
	br := database.DeleteBuilder{
		From:  shareTB,
		Where: []string{"id = :share_id"},
	}
	_, err := e.NamedExec(br.String(), m)
	return err
}

func bindShareOpts(opts *ShareOpts) (database.SelectBuilder, database.QueryVars) {
	br := database.SelectBuilder{
		From: shareTB,
		Columns: database.NamespacedColumn(
			[]string{
				"id AS share_id",
				"user_id AS share_user_id",
				"name AS share_name",
				"description AS share_description",
				"readme AS share_readme",
				"image_id AS share_image_id",
				"created_at AS share_created_at",
				"updated_at AS share_updated_at",
			},
			shareTB,
		),
	}
	if opts == nil {
		return br, nil
	}

	br = appendListOpts(br, opts.ListOpts)
	args := database.QueryVars{}

	if opts.Name != "" {
		br.Where = append(br.Where, "name = :name")
		args["name"] = opts.Name
	}
	if opts.UserID != "" {
		br.Where = append(br.Where, "user_id = :user_id")
		args["user_id"] = opts.UserID
	}

	return br, args
}
