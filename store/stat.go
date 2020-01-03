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
	MonlID        string
	MonlIDs       []string
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
	e := s.Exter(ctx)
	br, args := bindStatOpts(opts)
	br.From = statTB
	stmt := br.String()
	rows, err := e.NamedQuery(stmt, args)
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
	e := s.Exter(ctx)
	stmt := database.InsertBuilder{
		Into: statTB,
		Fields: map[string]interface{}{
			"id":          nil,
			"monl_id":     nil,
			"recorded_at": nil,
			"kind":        nil,
			"value":       nil,
			"is_latest":   nil,
			"manifest":    nil,
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
	e := s.Exter(ctx)
	stmt := database.UpdateBuilder{
		From: statTB,
		Fields: map[string]interface{}{
			"monl_id":     nil,
			"recorded_at": nil,
			"kind":        nil,
			"value":       nil,
			"is_latest":   nil,
			"manifest":    nil,
		},
		Where: []string{"id = :id"},
	}.String()
	_, err := e.NamedExec(stmt, m)
	return err
}

func bindStatOpts(opts *StatOpts) (database.SelectBuilder, map[string]interface{}) {
	br := database.SelectBuilder{}
	if opts == nil {
		return br, nil
	}

	br = bindListOpts(opts.ListOpts)
	args := make(map[string]interface{})
	if opts.Kind != "" {
		br.Where = append(br.Where, "kind = :kind")
		args["kind"] = opts.Kind
	}

	if opts.MonlID != "" {
		opts.MonlIDs = append(opts.MonlIDs, opts.MonlID)
	}
	if opts.MonlIDs != nil {
		idKeys := make([]string, 0)
		for i, id := range opts.MonlIDs {
			key := fmt.Sprintf("monl_id%d", i)
			idKeys = append(idKeys, fmt.Sprintf(":%s", key))
			args[key] = id
		}
		br.Where = append(br.Where, "monl_id IN ("+strings.Join(idKeys, ", ")+")")
	}

	if opts.WithLatest {
		br.Where = append(br.Where, "is_latest = 1")
	} else if opts.WithoutLatest {
		br.Where = append(br.Where, "is_latest = 0")
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
