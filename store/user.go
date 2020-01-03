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
	Email string
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
	e := s.Exter(ctx)
	br, args := bindUserOpts(opts)
	br.From = userTB
	stmt := br.String()
	rows, err := e.NamedQuery(stmt, args)
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
	e := s.Exter(ctx)
	stmt := database.SelectBuilder{
		From:  userTB,
		Where: []string{"id = :id"},
	}.String()
	rows, err := e.NamedQuery(stmt, m)
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
	e := s.Exter(ctx)
	stmt := database.SelectBuilder{
		From:  userTB,
		Where: []string{"login = :login"},
	}.String()
	rows, err := e.NamedQuery(stmt, m)
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
	e := s.Exter(ctx)
	stmt := database.InsertBuilder{
		Into: userTB,
		Fields: map[string]interface{}{
			"id":         nil,
			"login":      nil,
			"password":   nil,
			"name":       nil,
			"email":      nil,
			"created_at": nil,
		},
	}.String()
	_, err := e.NamedExec(stmt, m2)
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
	e := s.Exter(ctx)
	stmt := database.UpdateBuilder{
		From: userTB,
		Fields: map[string]interface{}{
			"login":      nil,
			"password":   nil,
			"name":       nil,
			"email":      nil,
			"updated_at": nil,
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

// Delete removes user by id.
func (s *dbUserStore) Delete(ctx context.Context, m *model.User) error {
	e := s.Exter(ctx)
	stmt := database.DeleteBuilder{
		From:  userTB,
		Where: []string{"id = :id"},
	}.String()
	_, err := e.NamedExec(stmt, m)
	return err
}

func bindUserOpts(opts *UserOpts) (database.SelectBuilder, map[string]interface{}) {
	br := database.SelectBuilder{}
	if opts == nil {
		return br, nil
	}

	br = bindListOpts(opts.ListOpts)
	args := make(map[string]interface{})
	if opts.Login != "" {
		br.Where = append(br.Where, "login = :login")
		args["login"] = opts.Login
	}
	if opts.Email != "" {
		br.Where = append(br.Where, "email = :email")
		args["email"] = opts.Email
	}

	return br, args
}
