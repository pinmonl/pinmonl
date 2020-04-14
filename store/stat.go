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
	ParentID      string
	ParentIDs     []string
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
	rows, err := e.NamedQuery(br.String(), args.Map())
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
	br := database.InsertBuilder{
		Into: statTB,
		Fields: map[string]interface{}{
			"id":          ":stat_id",
			"pkg_id":      ":stat_pkg_id",
			"parent_id":   ":stat_parent_id",
			"recorded_at": ":stat_recorded_at",
			"kind":        ":stat_kind",
			"value":       ":stat_value",
			"digest":      ":stat_digest",
			"is_latest":   ":stat_is_latest",
			"labels":      ":stat_labels",
		},
	}
	_, err := e.NamedExec(br.String(), m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// Update updates the fields of stat by id.
func (s *dbStatStore) Update(ctx context.Context, m *model.Stat) error {
	e := s.Execer(ctx)
	br := database.UpdateBuilder{
		From: statTB,
		Fields: map[string]interface{}{
			"pkg_id":      ":stat_pkg_id",
			"parent_id":   ":stat_parent_id",
			"recorded_at": ":stat_recorded_at",
			"kind":        ":stat_kind",
			"value":       ":stat_value",
			"digest":      ":stat_digest",
			"is_latest":   ":stat_is_latest",
			"labels":      ":stat_labels",
		},
		Where: []string{"id = :stat_id"},
	}
	_, err := e.NamedExec(br.String(), m)
	return err
}

func bindStatOpts(opts *StatOpts) (database.SelectBuilder, database.QueryVars) {
	br := database.SelectBuilder{
		From: statTB,
		Columns: database.NamespacedColumn(
			[]string{
				"id AS stat_id",
				"pkg_id AS stat_pkg_id",
				"parent_id AS stat_parent_id",
				"recorded_at AS stat_recorded_at",
				"kind AS stat_kind",
				"value AS stat_value",
				"digest AS stat_digest",
				"is_latest AS stat_is_latest",
				"labels AS stat_labels",
			},
			statTB,
		),
	}
	if opts == nil {
		return br, nil
	}

	br = appendListOpts(br, opts.ListOpts)
	args := database.QueryVars{}

	if opts.Kind != "" {
		br.Where = append(br.Where, "kind = :kind")
		args["kind"] = opts.Kind
	}

	if opts.PkgID != "" {
		opts.PkgIDs = append(opts.PkgIDs, opts.PkgID)
	}
	if opts.PkgIDs != nil {
		ks, ids := bindQueryIDs("pkg_ids", opts.PkgIDs)
		args.AppendStringMap(ids)
		br.Where = append(br.Where, fmt.Sprintf("pkg_id IN (%s)", strings.Join(ks, ",")))
	}

	if opts.ParentID != "" {
		opts.ParentIDs = append(opts.ParentIDs, opts.ParentID)
	}
	if opts.ParentIDs != nil {
		ks, ids := bindQueryIDs("parent_ids", opts.ParentIDs)
		args.AppendStringMap(ids)
		br.Where = append(br.Where, fmt.Sprintf("parent_id IN (%s)", strings.Join(ks, ",")))
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
