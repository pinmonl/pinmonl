package store

import (
	"context"
	"database/sql"

	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
)

// UserOpts defines the parameters for user filtering.
type UserOpts struct {
	ListOpts
	Login string
	Hash  string
}

// UserStore defines the services of user.
type UserStore interface {
	List(context.Context, *UserOpts) ([]model.User, error)
	Find(context.Context, *model.User) error
	FindLogin(context.Context, *model.User) error
	Create(context.Context, *model.User) error
	Update(context.Context, *model.User) error
	Delete(context.Context, *model.User) error
}

// NewUserStore creates user store.
func NewUserStore(s Store) UserStore {
	return &dbUserStore{s}
}

type dbUserStore struct {
	Store
}

// List retrieves users by the filter parameters.
func (s *dbUserStore) List(ctx context.Context, opts *UserOpts) ([]model.User, error) {
	e := s.Queryer(ctx)
	br, args := bindUserOpts(opts)
	rows, err := e.NamedQuery(br.String(), args.Map())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ms []model.User
	for rows.Next() {
		var m model.User
		err = rows.StructScan(&m)
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}
	return ms, nil
}

// Find retrieves user by id.
func (s *dbUserStore) Find(ctx context.Context, m *model.User) error {
	e := s.Queryer(ctx)
	br, _ := bindUserOpts(nil)
	br.Where = []string{"id = :user_id"}
	br.Limit = 1
	rows, err := e.NamedQuery(br.String(), m)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return sql.ErrNoRows
	}
	var m2 model.User
	err = rows.StructScan(&m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// FindLogin retrieves user by login.
func (s *dbUserStore) FindLogin(ctx context.Context, m *model.User) error {
	e := s.Queryer(ctx)
	br, _ := bindUserOpts(nil)
	br.Where = []string{"login = :user_login"}
	br.Limit = 1
	rows, err := e.NamedQuery(br.String(), m)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return sql.ErrNoRows
	}
	var m2 model.User
	err = rows.StructScan(&m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// Create inserts the fields of user with generated id.
func (s *dbUserStore) Create(ctx context.Context, m *model.User) error {
	m2 := *m
	m2.ID = newUID()
	m2.CreatedAt = timestamp()
	e := s.Execer(ctx)
	br := database.InsertBuilder{
		Into: userTB,
		Fields: map[string]interface{}{
			"id":         ":user_id",
			"login":      ":user_login",
			"password":   ":user_password",
			"name":       ":user_name",
			"image_id":   ":user_image_id",
			"role":       ":user_role",
			"hash":       ":user_hash",
			"created_at": ":user_created_at",
			"last_log":   ":user_last_log",
		},
	}
	_, err := e.NamedExec(br.String(), m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// Update updates the fields of user by id.
func (s *dbUserStore) Update(ctx context.Context, m *model.User) error {
	m2 := *m
	m2.UpdatedAt = timestamp()
	e := s.Execer(ctx)
	br := database.UpdateBuilder{
		From: userTB,
		Fields: map[string]interface{}{
			"login":      ":user_login",
			"password":   ":user_password",
			"name":       ":user_name",
			"image_id":   ":user_image_id",
			"role":       ":user_role",
			"hash":       ":user_hash",
			"updated_at": ":user_updated_at",
			"last_log":   ":user_last_log",
		},
		Where: []string{"id = :user_id"},
	}
	_, err := e.NamedExec(br.String(), m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// Delete removes user by id.
func (s *dbUserStore) Delete(ctx context.Context, m *model.User) error {
	e := s.Execer(ctx)
	br := database.DeleteBuilder{
		From:  userTB,
		Where: []string{"id = :user_id"},
	}
	_, err := e.NamedExec(br.String(), m)
	return err
}

func bindUserOpts(opts *UserOpts) (database.SelectBuilder, database.QueryVars) {
	br := database.SelectBuilder{
		From: userTB,
		Columns: database.NamespacedColumn(
			[]string{
				"id AS user_id",
				"login AS user_login",
				"password AS user_password",
				"name AS user_name",
				"image_id AS user_image_id",
				"role AS user_role",
				"hash AS user_hash",
				"created_at AS user_created_at",
				"updated_at AS user_updated_at",
				"last_log AS user_last_log",
			},
			userTB,
		),
	}
	if opts == nil {
		return br, nil
	}

	br = appendListOpts(br, opts.ListOpts)
	args := database.QueryVars{}

	if opts.Login != "" {
		br.Where = append(br.Where, "login = :login")
		args["login"] = opts.Login
	}
	if opts.Hash != "" {
		br.Where = append(br.Where, "hash = :hash")
		args["hash"] = opts.Hash
	}

	return br, args
}
