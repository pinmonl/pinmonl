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
	rows, err := e.NamedQuery(br.String(), args.Map())
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
	br := database.InsertBuilder{
		Into: substatTB,
		Fields: map[string]interface{}{
			"id":      ":substat_id",
			"stat_id": ":substat_stat_id",
			"kind":    ":substat_kind",
			"labels":  ":substat_labels",
		},
	}
	_, err := e.NamedExec(br.String(), m2)
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
	br := database.UpdateBuilder{
		From: substatTB,
		Fields: map[string]interface{}{
			"stat_id": ":substat_stat_id",
			"kind":    ":substat_kind",
			"labels":  ":substat_labels",
		},
		Where: []string{"id = :substat_id"},
	}
	_, err := e.NamedExec(br.String(), m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// Delete removes substat by ID.
func (s *dbSubstatStore) Delete(ctx context.Context, m *model.Substat) error {
	e := s.Execer(ctx)
	br := database.DeleteBuilder{
		From:  substatTB,
		Where: []string{"id = :substat_id"},
	}
	_, err := e.NamedExec(br.String(), m)
	return err
}

func bindSubstatOpts(opts *SubstatOpts) (database.SelectBuilder, database.QueryVars) {
	br := database.SelectBuilder{
		From: substatTB,
		Columns: database.NamespacedColumn(
			[]string{
				"id AS substat_id",
				"stat_id AS substat_stat_id",
				"kind AS substat_kind",
				"labels AS substat_labels",
			},
			substatTB,
		),
	}
	if opts == nil {
		return br, nil
	}

	br = appendListOpts(br, opts.ListOpts)
	args := database.QueryVars{}

	if opts.Stats != nil {
		ids := make([]string, len(opts.Stats))
		for i, s := range opts.Stats {
			ids[i] = s.ID
		}
		opts.StatIDs = ids
	}
	if opts.StatIDs != nil {
		ks, ids := bindQueryIDs("stat_ids", opts.StatIDs)
		args.AppendStringMap(ids)
		br.Where = append(br.Where, fmt.Sprintf("stat_id IN (%s)", strings.Join(ks, ",")))
	}

	return br, args
}
