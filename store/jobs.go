package store

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
)

type Jobs struct {
	*Store
}

type JobOpts struct {
	ListOpts
	Targets    model.MorphableList
	TargetIDs  []string
	TargetName string
	Names      []string
	Status     field.NullValue
	NotEnded   field.NullBool

	joinMonls bool
	joinPkgs  bool
}

func NewJobs(s *Store) *Jobs {
	return &Jobs{s}
}

func (j Jobs) table() string {
	return "jobs"
}

func (j *Jobs) List(ctx context.Context, opts *JobOpts) (model.JobList, error) {
	if opts == nil {
		opts = &JobOpts{}
	}

	qb := j.RunnableBuilder(ctx).
		Select(j.columns()...).From(j.table())
	qb = j.bindOpts(qb, opts)
	qb = addPagination(qb, opts)
	rows, err := qb.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := make([]*model.Job, 0)
	for rows.Next() {
		job, err := j.scan(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, job)
	}
	return list, nil
}

func (j *Jobs) Count(ctx context.Context, opts *JobOpts) (int64, error) {
	if opts == nil {
		opts = &JobOpts{}
	}

	qb := j.RunnableBuilder(ctx).
		Select("count(*)").From(j.table())
	qb = j.bindOpts(qb, opts)
	row := qb.QueryRow()
	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (j *Jobs) Find(ctx context.Context, id string) (*model.Job, error) {
	qb := j.RunnableBuilder(ctx).
		Select(j.columns()...).From(j.table()).
		Where("id = ?", id)
	row := qb.QueryRow()
	job, err := j.scan(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (j Jobs) bindOpts(b squirrel.SelectBuilder, opts *JobOpts) squirrel.SelectBuilder {
	if opts == nil {
		return b
	}

	if len(opts.Targets) > 0 && !opts.Targets.IsMixed() {
		opts.TargetName = opts.Targets.MorphName()
		opts.TargetIDs = opts.Targets.MorphKeys()
	}
	if opts.TargetName != "" {
		b = b.Where("target_name = ?", opts.TargetName)
	}
	if len(opts.TargetIDs) > 0 {
		b = b.Where(squirrel.Eq{"target_id": opts.TargetIDs})
	}

	if len(opts.Names) > 0 {
		b = b.Where(squirrel.Eq{"name": opts.Names})
	}

	if opts.Status.Valid {
		if vs, ok := opts.Status.Value().(model.JobStatus); ok {
			b = b.Where("status = ?", vs)
		}
	}

	if opts.NotEnded.Valid {
		b = b.Where("ended_at IS NULL")
	}

	return b
}

func (j Jobs) columns() []string {
	return []string{
		j.table() + ".id",
		j.table() + ".name",
		j.table() + ".describe",
		j.table() + ".target_id",
		j.table() + ".target_name",
		j.table() + ".status",
		j.table() + ".message",
		j.table() + ".created_at",
		j.table() + ".ended_at",
	}
}

func (j Jobs) scanColumns(job *model.Job) []interface{} {
	return []interface{}{
		&job.ID,
		&job.Name,
		&job.Describe,
		&job.TargetID,
		&job.TargetName,
		&job.Status,
		&job.Message,
		&job.CreatedAt,
		&job.EndedAt,
	}
}

func (j Jobs) scan(row database.RowScanner) (*model.Job, error) {
	var job model.Job
	err := row.Scan(j.scanColumns(&job)...)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (j *Jobs) Create(ctx context.Context, job *model.Job) error {
	job2 := *job
	job2.ID = newID()
	job2.CreatedAt = timestamp()

	qb := j.RunnableBuilder(ctx).
		Insert(j.table()).
		Columns(
			"id",
			"name",
			"describe",
			"target_id",
			"target_name",
			"status",
			"message",
			"created_at",
			"ended_at").
		Values(
			job2.ID,
			job2.Name,
			job2.Describe,
			job2.TargetID,
			job2.TargetName,
			job2.Status,
			job2.Message,
			job2.CreatedAt,
			job2.EndedAt)
	_, err := qb.Exec()
	if err != nil {
		return err
	}
	*job = job2
	return nil
}

func (j *Jobs) Update(ctx context.Context, job *model.Job) error {
	job2 := *job

	qb := j.RunnableBuilder(ctx).
		Update(j.table()).
		Set("name", job2.Name).
		Set("describe", job2.Describe).
		Set("target_id", job2.TargetID).
		Set("target_name", job2.TargetName).
		Set("status", job2.Status).
		Set("message", job2.Message).
		Set("ended_at", job2.EndedAt).
		Where("id = ?", job2.ID)
	_, err := qb.Exec()
	if err != nil {
		return err
	}
	*job = job2
	return nil
}

func (j *Jobs) Delete(ctx context.Context, id string) (int64, error) {
	qb := j.RunnableBuilder(ctx).
		Delete(j.table()).
		Where("id = ?", id)
	res, err := qb.Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (o *JobOpts) JoinMonls() *JobOpts {
	o2 := *o
	o2.joinMonls = true
	return &o2
}

func (o *JobOpts) JoinPkgs() *JobOpts {
	o2 := *o
	o2.joinPkgs = true
	return &o2
}
