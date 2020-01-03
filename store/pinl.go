package store

import (
	"context"
	"database/sql"

	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
)

// PinlOpts defines the paramters for pinl filtering.
type PinlOpts struct {
	ListOpts
	UserID string
}

// PinlStore defines the services of pinl.
type PinlStore interface {
	List(context.Context, *PinlOpts) ([]model.Pinl, error)
	Find(context.Context, *model.Pinl) error
	Create(context.Context, *model.Pinl) error
	Update(context.Context, *model.Pinl) error
	Delete(context.Context, *model.Pinl) error
}

// NewPinlStore creates pinl store.
func NewPinlStore(s Store) PinlStore {
	return &dbPinlStore{s}
}

type dbPinlStore struct {
	Store
}

// List retrieves pinls by the filter parameters.
func (s *dbPinlStore) List(ctx context.Context, opts *PinlOpts) ([]model.Pinl, error) {
	e := s.Exter(ctx)
	br, args := bindPinlOpts(opts)
	br.From = pinlTB
	stmt := br.String()
	rows, err := e.NamedQuery(stmt, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Pinl
	for rows.Next() {
		var m model.Pinl
		err = rows.StructScan(&m)
		if err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

// Find retrieves pinl by id.
func (s *dbPinlStore) Find(ctx context.Context, m *model.Pinl) error {
	e := s.Exter(ctx)
	stmt := database.SelectBuilder{
		From:  pinlTB,
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
	var m2 model.Pinl
	err = rows.StructScan(&m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// Create inserts the fields of pinl with generated id.
func (s *dbPinlStore) Create(ctx context.Context, m *model.Pinl) error {
	m2 := *m
	m2.ID = newUID()
	m2.CreatedAt = timestamp()
	e := s.Exter(ctx)
	stmt := database.InsertBuilder{
		Into: pinlTB,
		Fields: map[string]interface{}{
			"id":          nil,
			"user_id":     nil,
			"url":         nil,
			"title":       nil,
			"description": nil,
			"readme":      nil,
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

// Update updates the fields of pinl by id.
func (s *dbPinlStore) Update(ctx context.Context, m *model.Pinl) error {
	m2 := *m
	m2.UpdatedAt = timestamp()
	e := s.Exter(ctx)
	stmt := database.UpdateBuilder{
		From: pinlTB,
		Fields: map[string]interface{}{
			"user_id":     nil,
			"url":         nil,
			"title":       nil,
			"description": nil,
			"readme":      nil,
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

// Delete removes pinl by id.
func (s *dbPinlStore) Delete(ctx context.Context, m *model.Pinl) error {
	e := s.Exter(ctx)
	stmt := database.DeleteBuilder{
		From:  pinlTB,
		Where: []string{"id = :id"},
	}.String()
	_, err := e.NamedExec(stmt, m)
	return err
}

func bindPinlOpts(opts *PinlOpts) (database.SelectBuilder, map[string]interface{}) {
	br := database.SelectBuilder{}
	if opts == nil {
		return br, nil
	}

	br = bindListOpts(opts.ListOpts)
	args := make(map[string]interface{})
	if opts.UserID != "" {
		br.Where = append(br.Where, "user_id = :user_id")
		args["user_id"] = opts.UserID
	}

	return br, args
}
