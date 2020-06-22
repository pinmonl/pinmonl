package store

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pinmonl/pinmonl/database/dbtest"
	"github.com/pinmonl/pinmonl/model"
	"github.com/stretchr/testify/assert"
)

func TestJobs(t *testing.T) {
	db, mock, err := dbtest.New()
	assert.Nil(t, err)
	defer db.Close()

	ctx := context.TODO()
	s := NewStore(db)
	jobs := NewJobs(s)

	t.Run("list", testJobsList(ctx, jobs, mock))
	t.Run("count", testJobsCount(ctx, jobs, mock))
	t.Run("find", testJobsFind(ctx, jobs, mock))
	t.Run("create", testJobsCreate(ctx, jobs, mock))
	t.Run("update", testJobsUpdate(ctx, jobs, mock))
	t.Run("delete", testJobsDelete(ctx, jobs, mock))
}

func testJobsList(ctx context.Context, jobs *Jobs, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			prefix = "SELECT (.+) FROM jobs"
			opts   *JobOpts
			list   []*model.Job
			err    error
		)

		// Test nil opts.
		opts = nil
		mock.ExpectQuery(prefix).
			WillReturnRows(sqlmock.NewRows(jobs.columns()).
				AddRow("job-id-1", "target-1", "target", "", "description", 1, "job/png", nil, nil))
		list, err = jobs.List(ctx, opts)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(list))

		// Test filter by targets.
		opts = &JobOpts{Targets: model.MorphableList{
			&model.Pinl{ID: "pinl-id-1"},
			&model.Pinl{ID: "pinl-id-2"},
			&model.Pinl{ID: "pinl-id-3"},
		}}
		mock.ExpectQuery(fmt.Sprintf(regexp.QuoteMeta("%s WHERE target_name = ? AND target_id IN (?,?,?)"), prefix)).
			WithArgs(
				opts.Targets[0].MorphName(),
				opts.Targets[0].MorphKey(),
				opts.Targets[1].MorphKey(),
				opts.Targets[2].MorphKey()).
			WillReturnRows(sqlmock.NewRows(jobs.columns()))
		_, err = jobs.List(ctx, opts)
		assert.Nil(t, err)
	}
}

func testJobsCount(ctx context.Context, jobs *Jobs, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = regexp.QuoteMeta("SELECT count(*) FROM jobs")
			opts  *JobOpts
			count int64
			err   error
		)

		opts = &JobOpts{}
		mock.ExpectQuery(query).
			WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).
				AddRow(1))
		count, err = jobs.Count(ctx, opts)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), count)
	}
}

func testJobsFind(ctx context.Context, jobs *Jobs, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = "SELECT (.+) FROM jobs WHERE id = \\?"
			id    string
			job   *model.Job
			err   error
		)

		id = "job-id-1"
		mock.ExpectQuery(query).
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows(jobs.columns()).
				AddRow(id, "target-1", "target", "", "description", 1, "job/png", nil, nil))
		job, err = jobs.Find(ctx, id)
		assert.Nil(t, err)
		if assert.NotNil(t, job) {
			assert.Equal(t, id, job.ID)
		}
	}
}

func testJobsCreate(ctx context.Context, jobs *Jobs, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			job *model.Job
			err error
		)

		job = &model.Job{}
		expectJobsCreate(mock, job)
		err = jobs.Create(ctx, job)
		assert.Nil(t, err)
		assert.NotEmpty(t, job.ID)
		assert.NotEmpty(t, job.CreatedAt)
	}
}

func expectJobsCreate(mock sqlmock.Sqlmock, job *model.Job) {
	mock.ExpectExec("INSERT INTO jobs").
		WithArgs(
			sqlmock.AnyArg(),
			job.Name,
			job.Describe,
			job.TargetID,
			job.TargetName,
			job.Status,
			job.Message,
			sqlmock.AnyArg(),
			job.EndedAt).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func testJobsUpdate(ctx context.Context, jobs *Jobs, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			job *model.Job
			err error
		)

		job = &model.Job{ID: "job-id-1"}
		expectJobsUpdate(mock, job)
		err = jobs.Update(ctx, job)
		assert.Nil(t, err)
	}
}

func expectJobsUpdate(mock sqlmock.Sqlmock, job *model.Job) {
	mock.ExpectExec("UPDATE jobs (.+) WHERE id = \\?").
		WithArgs(
			job.Name,
			job.Describe,
			job.TargetID,
			job.TargetName,
			job.Status,
			job.Message,
			job.EndedAt,
			job.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func testJobsDelete(ctx context.Context, jobs *Jobs, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = regexp.QuoteMeta("DELETE FROM jobs WHERE id = ?")
			id    string
			n     int64
			err   error
		)

		id = "job-id-1"
		mock.ExpectExec(query).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))
		n, err = jobs.Delete(ctx, id)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), n)
	}
}
