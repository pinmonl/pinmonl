package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
)

// JobOpts defines the parameters for job filtering.
type JobOpts struct {
	ListOpts
	Name            model.JobName
	Status          model.JobStatus
	ScheduledAfter  time.Time
	ScheduledBefore time.Time
}

// JobStore defines the services of job.
type JobStore interface {
	List(context.Context, *JobOpts) ([]model.Job, error)
	Find(context.Context, *model.Job) error
	Create(context.Context, *model.Job) error
	Update(context.Context, *model.Job) error
	Delete(context.Context, *model.Job) error
}

type dbJobStore struct {
	Store
}

// NewJobStore creates job store.
func NewJobStore(s Store) JobStore {
	return &dbJobStore{s}
}

var jobCols = []string{"id", "name", "target_id", "status", "error", "scheduled_at", "started_at", "created_at"}

// List retrieves jobs by the filter parameters.
func (s *dbJobStore) List(ctx context.Context, opts *JobOpts) ([]model.Job, error) {
	e := s.Exter(ctx)
	br, args := bindJobOpts(opts)
	br.Columns = jobCols
	br.From = jobTB
	stmt := br.String()
	rows, err := e.NamedQuery(stmt, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Job
	for rows.Next() {
		var m model.Job
		err = rows.StructScan(&m)
		if err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

// Find retrieves job by id.
func (s *dbJobStore) Find(ctx context.Context, m *model.Job) error {
	e := s.Exter(ctx)
	stmt := database.SelectBuilder{
		From:    jobTB,
		Columns: jobCols,
		Where:   []string{"id = :id"},
		Limit:   1,
	}.String()
	rows, err := e.NamedQuery(stmt, m)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return sql.ErrNoRows
	}
	var m2 model.Job
	err = rows.StructScan(&m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// Create inserts the field of job with generated id.
func (s *dbJobStore) Create(ctx context.Context, m *model.Job) error {
	m2 := *m
	m2.ID = newUID()
	m2.CreatedAt = timestamp()
	e := s.Exter(ctx)
	stmt := database.InsertBuilder{
		Into: jobTB,
		Fields: map[string]interface{}{
			"id":           nil,
			"name":         nil,
			"target_id":    nil,
			"status":       nil,
			"error":        nil,
			"scheduled_at": nil,
			"started_at":   nil,
			"created_at":   nil,
		},
	}.String()
	_, err := e.NamedExec(stmt, m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// Update updates job by id.
func (s *dbJobStore) Update(ctx context.Context, m *model.Job) error {
	m2 := *m
	e := s.Exter(ctx)
	stmt := database.UpdateBuilder{
		From: jobTB,
		Fields: map[string]interface{}{
			"name":         nil,
			"target_id":    nil,
			"status":       nil,
			"error":        nil,
			"scheduled_at": nil,
			"started_at":   nil,
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

// Delete removes job by id.
func (s *dbJobStore) Delete(ctx context.Context, m *model.Job) error {
	e := s.Exter(ctx)
	stmt := database.DeleteBuilder{
		From:  jobTB,
		Where: []string{"id = :id"},
	}.String()
	_, err := e.NamedExec(stmt, m)
	return err
}

func bindJobOpts(opts *JobOpts) (database.SelectBuilder, map[string]interface{}) {
	br := database.SelectBuilder{}
	if opts == nil {
		return br, nil
	}

	br = bindListOpts(opts.ListOpts)
	args := make(map[string]interface{})

	if opts.Status != model.JobStatusNotSet {
		br.Where = append(br.Where, "status = :status")
		args["status"] = opts.Status
	}

	if opts.Name != model.JobNotSet {
		br.Where = append(br.Where, "name = :name")
		args["name"] = opts.Name
	}

	if !opts.ScheduledAfter.IsZero() {
		br.Where = append(br.Where, "scheduled_at >= :scheduled_after")
		args["scheduled_after"] = opts.ScheduledAfter
	}
	if !opts.ScheduledBefore.IsZero() {
		br.Where = append(br.Where, "(scheduled_at <= :scheduled_before OR scheduled_at IS NULL)")
		args["scheduled_before"] = opts.ScheduledBefore
	}

	return br, args
}
