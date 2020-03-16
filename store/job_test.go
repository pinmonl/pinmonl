package store

import (
	"context"
	"testing"

	"github.com/pinmonl/pinmonl/database/dbtest"
	"github.com/pinmonl/pinmonl/model"
	"github.com/stretchr/testify/assert"
)

func TestJobStore(t *testing.T) {
	db, err := dbtest.Open()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		dbtest.Reset(db)
		dbtest.Close(db)
	}()

	mockData := []*model.Job{
		{Name: model.JobPinlCreated, TargetID: "1", Status: model.JobStatusPending},
		{Name: model.JobPinlUpdated, TargetID: "2", Status: model.JobStatusPending},
	}

	store := NewStore(db)
	jobs := NewJobStore(store)
	ctx := context.TODO()
	t.Run("Create", testJobCreate(ctx, jobs, mockData))
	t.Run("List", testJobList(ctx, jobs, mockData))
	t.Run("Find", testJobFind(ctx, jobs, mockData))
	t.Run("Update", testJobUpdate(ctx, jobs, mockData))
	t.Run("Delete", testJobDelete(ctx, jobs, mockData))
}

func testJobCreate(ctx context.Context, jobs JobStore, mockData []*model.Job) func(t *testing.T) {
	return func(t *testing.T) {
		for _, job := range mockData {
			assert.Nil(t, jobs.Create(ctx, job))
			assert.NotEmpty(t, job.ID)
			assert.False(t, job.CreatedAt.Time().IsZero())
		}
	}
}

func testJobList(ctx context.Context, jobs JobStore, mockData []*model.Job) func(t *testing.T) {
	return func(t *testing.T) {
		deRef := func(data []*model.Job) []model.Job {
			out := make([]model.Job, len(data))
			for i, mj := range data {
				m := *mj
				out[i] = m
			}
			return out
		}

		want := deRef(mockData)
		got, err := jobs.List(ctx, nil)
		assert.Nil(t, err)
		assert.Equal(t, want, got)

		want = deRef(mockData)
		got, err = jobs.List(ctx, &JobOpts{Status: model.JobStatusPending})
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	}
}

func testJobFind(ctx context.Context, jobs JobStore, mockData []*model.Job) func(t *testing.T) {
	return func(t *testing.T) {
		want := mockData[0]
		got := model.Job{ID: want.ID}
		assert.Nil(t, jobs.Find(ctx, &got))
		assert.Equal(t, *want, got)
	}
}

func testJobUpdate(ctx context.Context, jobs JobStore, mockData []*model.Job) func(t *testing.T) {
	return func(t *testing.T) {
		want := mockData[0]
		want.Status = model.JobStatusQueue
		assert.Nil(t, jobs.Update(ctx, want))

		got := model.Job{ID: want.ID}
		jobs.Find(ctx, &got)
		assert.Equal(t, *want, got)
	}
}

func testJobDelete(ctx context.Context, jobs JobStore, mockData []*model.Job) func(t *testing.T) {
	return func(t *testing.T) {
		del, want := mockData[0], mockData[1:]
		assert.Nil(t, jobs.Delete(ctx, del))

		got, _ := jobs.List(ctx, nil)
		assert.Equal(t, len(want), len(got))
	}
}
