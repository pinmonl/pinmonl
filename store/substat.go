package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
)

// SubstatOpts defines the parameters for substat filtering.
type SubstatOpts struct {
	ListOpts
	Stats   []model.Stat
	StatIDs []string
}

// SubstatStore defines the services of substat.
type SubstatStore interface {
	List(context.Context, *SubstatOpts) ([]model.Substat, error)
	Create(context.Context, *model.Substat) error
	Update(context.Context, *model.Substat) error
	Delete(context.Context, *model.Substat) error
}

// NewSubstatStore creates substat store.
func NewSubstatStore(s Store) SubstatStore {
	return &dbSubstatStore{s}
}

type dbSubstatStore struct {
	Store
}

// List retrieves substats by the filter parameters.
func (s *dbSubstatStore) List(ctx context.Context, opts *SubstatOpts) ([]model.Substat, error) {
	e := s.Queryer(ctx)
	br, args := bindSubstatOpts(opts)
	rows, err := e.NamedQuery(br.String(), args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Substat
	for rows.Next() {
		var m model.Substat
		err = rows.StructScan(&m)
		if err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

// Create inserts the fields of substat.
func (s *dbSubstatStore) Create(ctx context.Context, m *model.Substat) error {
	m2 := *m
	m2.ID = newUID()
	e := s.Execer(ctx)
	stmt := database.InsertBuilder{
		Into: substatTB,
		Fields: map[string]interface{}{
			"id":      nil,
			"stat_id": nil,
			"kind":    nil,
			"labels":  nil,
		},
	}.String()
	_, err := e.NamedExec(stmt, m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// Update updates the field of substat by ID.
func (s *dbSubstatStore) Update(ctx context.Context, m *model.Substat) error {
	m2 := *m
	e := s.Execer(ctx)
	stmt := database.UpdateBuilder{
		From: substatTB,
		Fields: map[string]interface{}{
			"stat_id": nil,
			"kind":    nil,
			"labels":  nil,
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

// Delete removes substat by ID.
func (s *dbSubstatStore) Delete(ctx context.Context, m *model.Substat) error {
	e := s.Execer(ctx)
	stmt := database.DeleteBuilder{
		From:  substatTB,
		Where: []string{"id = :id"},
	}.String()
	_, err := e.NamedExec(stmt, m)
	return err
}

func bindSubstatOpts(opts *SubstatOpts) (database.SelectBuilder, map[string]interface{}) {
	br := database.SelectBuilder{
		From: substatTB,
		Columns: database.NamespacedColumn(
			[]string{"id", "stat_id", "kind", "labels"},
			substatTB,
		),
	}
	if opts == nil {
		return br, nil
	}

	br = appendListOpts(br, opts.ListOpts)
	args := make(map[string]interface{})

	if opts.Stats != nil {
		ids := make([]string, len(opts.Stats))
		for i, s := range opts.Stats {
			ids[i] = s.ID
		}
		opts.StatIDs = ids
	}
	if opts.StatIDs != nil {
		ks, ids := bindQueryIDs("stat_ids", opts.StatIDs)
		for k, id := range ids {
			args[k] = id
		}
		br.Where = append(br.Where, fmt.Sprintf("stat_id IN (%s)", strings.Join(ks, ",")))
	}

	return br, args
}
