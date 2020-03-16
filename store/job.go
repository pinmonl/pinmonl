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

// List retrieves jobs by the filter parameters.
func (s *dbJobStore) List(ctx context.Context, opts *JobOpts) ([]model.Job, error) {
	e := s.Queryer(ctx)
	br, args := bindJobOpts(opts)
	rows, err := e.NamedQuery(br.String(), args.Map())
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
	e := s.Queryer(ctx)
	br, _ := bindJobOpts(nil)
	br.Where = []string{"id = :job_id"}
	br.Limit = 1
	rows, err := e.NamedQuery(br.String(), m)
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
	e := s.Execer(ctx)
	br := database.InsertBuilder{
		Into: jobTB,
		Fields: map[string]interface{}{
			"id":           ":job_id",
			"name":         ":job_name",
			"target_id":    ":job_target_id",
			"status":       ":job_status",
			"error":        ":job_error",
			"scheduled_at": ":job_scheduled_at",
			"started_at":   ":job_started_at",
			"created_at":   ":job_created_at",
		},
	}
	_, err := e.NamedExec(br.String(), m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// Update updates job by id.
func (s *dbJobStore) Update(ctx context.Context, m *model.Job) error {
	m2 := *m
	e := s.Execer(ctx)
	br := database.UpdateBuilder{
		From: jobTB,
		Fields: map[string]interface{}{
			"name":         ":job_name",
			"target_id":    ":job_target_id",
			"status":       ":job_status",
			"error":        ":job_error",
			"scheduled_at": ":job_scheduled_at",
			"started_at":   ":job_started_at",
		},
		Where: []string{"id = :job_id"},
	}
	_, err := e.NamedExec(br.String(), m2)
	if err != nil {
		return err
	}
	*m = m2
	return nil
}

// Delete removes job by id.
func (s *dbJobStore) Delete(ctx context.Context, m *model.Job) error {
	e := s.Execer(ctx)
	br := database.DeleteBuilder{
		From:  jobTB,
		Where: []string{"id = :job_id"},
	}
	_, err := e.NamedExec(br.String(), m)
	return err
}

func bindJobOpts(opts *JobOpts) (database.SelectBuilder, database.QueryVars) {
	br := database.SelectBuilder{
		From: jobTB,
		Columns: database.NamespacedColumn(
			[]string{
				"id AS job_id",
				"name AS job_name",
				"target_id AS job_target_id",
				"status AS job_status",
				"error AS job_error",
				"scheduled_at AS job_scheduled_at",
				"started_at AS job_started_at",
				"created_at AS job_created_at",
			},
			jobTB,
		),
	}
	if opts == nil {
		return br, nil
	}

	br = appendListOpts(br, opts.ListOpts)
	args := database.QueryVars{}

	if opts.Status != model.JobStatusEmpty {
		br.Where = append(br.Where, "status = :status")
		args["status"] = opts.Status
	}

	if opts.Name != model.JobEmpty {
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
