package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
)

// PinmonOpts defines the parameters for pinmon fitlering.
type PinmonOpts struct {
	ListOpts
	UserID  string
	MonlID  string
	PinlID  string
	MonlIDs []string
	PinlIDs []string
}

// PinmonStore defines the services of pinmon.
type PinmonStore interface {
	List(context.Context, *PinmonOpts) ([]model.Pinmon, error)
	Create(context.Context, *model.Pinmon) error
	Delete(context.Context, *model.Pinmon) error
}

// NewPinmonStore creates pinmon store.
func NewPinmonStore(s Store) PinmonStore {
	return &dbPinmonStore{s}
}

type dbPinmonStore struct {
	Store
}

// List retrieves pinmons by the filter parameters.
func (s *dbPinmonStore) List(ctx context.Context, opts *PinmonOpts) ([]model.Pinmon, error) {
	e := s.Queryer(ctx)
	br, args := bindPinmonOpts(opts)
	br.From = pinmonTB
	stmt := br.String()
	rows, err := e.NamedQuery(stmt, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Pinmon
	for rows.Next() {
		var m model.Pinmon
		err = rows.StructScan(&m)
		if err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

// Create inserts the fields of pinmon.
func (s *dbPinmonStore) Create(ctx context.Context, m *model.Pinmon) error {
	e := s.Execer(ctx)
	stmt := database.InsertBuilder{
		Into: pinmonTB,
		Fields: map[string]interface{}{
			"user_id": nil,
			"pinl_id": nil,
			"monl_id": nil,
			"sort":    nil,
		},
	}.String()
	_, err := e.NamedExec(stmt, m)
	return err
}

// Delete removes pinmon by the relationship.
func (s *dbPinmonStore) Delete(ctx context.Context, m *model.Pinmon) error {
	e := s.Execer(ctx)
	stmt := database.DeleteBuilder{
		From: pinmonTB,
		Where: []string{
			"user_id = :user_id",
			"monl_id = :monl_id",
			"pinl_id = :pinl_id",
		},
	}.String()
	_, err := e.NamedExec(stmt, m)
	return err
}

func bindPinmonOpts(opts *PinmonOpts) (database.SelectBuilder, map[string]interface{}) {
	br := database.SelectBuilder{}
	if opts == nil {
		return br, nil
	}

	args := make(map[string]interface{})
	if opts.UserID != "" {
		br.Where = append(br.Where, "user_id = :user_id")
		args["user_id"] = opts.UserID
	}
	if opts.MonlID != "" {
		opts.MonlIDs = append(opts.MonlIDs, opts.MonlID)
	}
	if opts.PinlID != "" {
		opts.PinlIDs = append(opts.PinlIDs, opts.PinlID)
	}
	if opts.MonlIDs != nil {
		inKeys := make([]string, 0)
		for i, id := range opts.MonlIDs {
			key := fmt.Sprintf("monl_id%d", i)
			inKeys = append(inKeys, fmt.Sprintf(":%s", key))
			args[key] = id
		}
		br.Where = append(br.Where, "monl_id IN ("+strings.Join(inKeys, ", ")+")")
	}
	if opts.PinlIDs != nil {
		inKeys := make([]string, 0)
		for i, id := range opts.PinlIDs {
			key := fmt.Sprintf("pinl_id%d", i)
			inKeys = append(inKeys, fmt.Sprintf(":%s", key))
			args[key] = id
		}
		br.Where = append(br.Where, "pinl_id IN ("+strings.Join(inKeys, ", ")+")")
	}

	return br, args
}
