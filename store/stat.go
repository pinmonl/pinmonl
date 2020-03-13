package store

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
)

// StatOpts defines the parameters for stat filtering.
type StatOpts struct {
	ListOpts
	Kind          string
	WithLatest    bool
	WithoutLatest bool
	PkgID         string
	PkgIDs        []string
	After         time.Time
	Before        time.Time
}

// StatStore defines the services of stat.
type StatStore interface {
	List(context.Context, *StatOpts) ([]model.Stat, error)
	Create(context.Context, *model.Stat) error
	Update(context.Context, *model.Stat) error
}

// NewStatStore creates stat store.
func NewStatStore(s Store) StatStore {
	return &dbStatStore{s}
}

type dbStatStore struct {
	Store
}

// List retrieves stats by the filter parameters.
func (s *dbStatStore) List(ctx context.Context, opts *StatOpts) ([]model.Stat, error) {
	e := s.Queryer(ctx)
	br, args := bindStatOpts(opts)
	rows, err := e.NamedQuery(br.String(), args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Stat
	for rows.Next() {
		var m model.Stat
		err = rows.StructScan(&m)
		if err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

// Create inserts the fields of stat with generated id.
func (s *dbStatStore) Create(ctx context.Context, m *model.Stat) error {
	m2 := *m
	m2.ID = newUID()
	e := s.Execer(ctx)
	stmt := database.InsertBuilder{
		Into: statTB,
		Fields: map[string]interface{}{
			"id":          nil,
			"pkg_id":      nil,
			"recorded_at": nil,
			"kind":        nil,
			"value":       nil,
			"is_latest":   nil,
			"labels":      nil,
		},
	}.String()
	_, err := e.NamedExec(stmt, m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// Update updates the fields of stat by id.
func (s *dbStatStore) Update(ctx context.Context, m *model.Stat) error {
	e := s.Execer(ctx)
	stmt := database.UpdateBuilder{
		From: statTB,
		Fields: map[string]interface{}{
			"pkg_id":      nil,
			"recorded_at": nil,
			"kind":        nil,
			"value":       nil,
			"is_latest":   nil,
			"labels":      nil,
		},
		Where: []string{"id = :id"},
	}.String()
	_, err := e.NamedExec(stmt, m)
	return err
}

func bindStatOpts(opts *StatOpts) (database.SelectBuilder, map[string]interface{}) {
	br := database.SelectBuilder{
		From: statTB,
		Columns: database.NamespacedColumn(
			[]string{"id", "pkg_id", "recorded_at", "kind", "value", "is_latest", "labels"},
			statTB,
		),
	}
	if opts == nil {
		return br, nil
	}

	br = appendListOpts(br, opts.ListOpts)
	args := make(map[string]interface{})

	if opts.Kind != "" {
		br.Where = append(br.Where, "kind = :kind")
		args["kind"] = opts.Kind
	}

	if opts.PkgID != "" {
		opts.PkgIDs = append(opts.PkgIDs, opts.PkgID)
	}
	if opts.PkgIDs != nil {
		ks, ids := bindQueryIDs("pkg_ids", opts.PkgIDs)
		for k, id := range ids {
			args[k] = id
		}
		br.Where = append(br.Where, fmt.Sprintf("pkg_id IN (%s)", strings.Join(ks, ",")))
	}

	if opts.WithLatest {
		br.Where = append(br.Where, "is_latest = :is_latest")
		args["is_latest"] = true
	} else if opts.WithoutLatest {
		br.Where = append(br.Where, "is_latest = :is_latest")
		args["is_latest"] = false
	}

	if !opts.After.IsZero() {
		br.Where = append(br.Where, "recorded_at >= :date_after")
		args["date_after"] = opts.After
	}
	if !opts.Before.IsZero() {
		br.Where = append(br.Where, "recorded_at <= :date_before")
		args["date_before"] = opts.Before
	}

	return br, args
}
