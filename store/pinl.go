package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
)

// PinlOpts defines the paramters for pinl filtering.
type PinlOpts struct {
	ListOpts
	ID         string
	UserID     string
	MustTagIDs []string
	AnyTagIDs  []string
}

// PinlStore defines the services of pinl.
type PinlStore interface {
	List(context.Context, *PinlOpts) ([]model.Pinl, error)
	Count(context.Context, *PinlOpts) (int64, error)
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
	e := s.Queryer(ctx)
	br, args := bindPinlOpts(opts)
	rows, err := e.NamedQuery(br.String(), args)
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

// Count counts the number of pinls with the filter parameters.
func (s *dbPinlStore) Count(ctx context.Context, opts *PinlOpts) (int64, error) {
	e := s.Queryer(ctx)
	br, args := bindPinlOpts(opts)
	br.Columns = []string{"COUNT(*) as count"}
	rows, err := e.NamedQuery(br.String(), args)
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

// Find retrieves pinl by id.
func (s *dbPinlStore) Find(ctx context.Context, m *model.Pinl) error {
	e := s.Queryer(ctx)
	br, _ := bindPinlOpts(nil)
	br.Where = []string{"id = :id"}
	br.Limit = 1
	rows, err := e.NamedQuery(br.String(), m)
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
	e := s.Execer(ctx)
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
	e := s.Execer(ctx)
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
	e := s.Execer(ctx)
	stmt := database.DeleteBuilder{
		From:  pinlTB,
		Where: []string{"id = :id"},
	}.String()
	_, err := e.NamedExec(stmt, m)
	return err
}

func bindPinlOpts(opts *PinlOpts) (database.SelectBuilder, map[string]interface{}) {
	br := database.SelectBuilder{
		From: pinlTB,
		Columns: database.NamespacedColumn(
			[]string{"id", "user_id", "url", "title", "description", "readme", "image_id", "created_at", "updated_at"},
			pinlTB,
		),
	}
	if opts == nil {
		return br, nil
	}

	br = appendListOpts(br, opts.ListOpts)
	args := make(map[string]interface{})

	if opts.UserID != "" {
		br.Where = append(br.Where, "user_id = :user_id")
		args["user_id"] = opts.UserID
	}

	if opts.ID != "" {
		br.Where = append(br.Where, "id = :id")
		args["id"] = opts.ID
	}

	if opts.MustTagIDs != nil {
		sq := database.SelectBuilder{
			Columns: []string{"1"},
			From:    taggableTB,
			Where: []string{
				"target_name = :must_tag_target_name",
				fmt.Sprintf("target_id = %s.id", pinlTB),
			},
			GroupBy: []string{"target_id"},
			Having:  []string{"COUNT(DISTINCT tag_id) = :must_tag_count"},
		}
		args["must_tag_target_name"] = model.Pinl{}.MorphName()
		args["must_tag_count"] = len(opts.MustTagIDs)

		ks := make([]string, len(opts.MustTagIDs))
		for i, t := range opts.MustTagIDs {
			k := fmt.Sprintf("must_tag_id%d", i)
			args[k] = t
			ks[i] = ":" + k
		}
		sq.Where = append(sq.Where, fmt.Sprintf("tag_id IN (%s)", strings.Join(ks, ", ")))
		br.Where = append(br.Where, fmt.Sprintf("EXISTS (%s)", sq.String()))
	}

	if opts.AnyTagIDs != nil {
		sq := database.SelectBuilder{
			Columns: []string{"1"},
			From:    taggableTB,
			Where: []string{
				"target_name = :any_tag_target_name",
				fmt.Sprintf("target_id = %s.id", pinlTB),
			},
			GroupBy: []string{"target_id"},
		}
		args["any_tag_target_name"] = model.Pinl{}.MorphName()

		ks := make([]string, len(opts.AnyTagIDs))
		for i, t := range opts.AnyTagIDs {
			k := fmt.Sprintf("any_tag_id%d", i)
			args[k] = t
			ks[i] = ":" + k
		}
		sq.Where = append(sq.Where, fmt.Sprintf("tag_id IN (%s)", strings.Join(ks, ", ")))
		br.Where = append(br.Where, fmt.Sprintf("EXISTS (%s)", sq.String()))
	}

	return br, args
}
