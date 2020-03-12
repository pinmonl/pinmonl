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

	store := NewStore(db)
	jobs := NewJobStore(store)
	ctx := context.TODO()
	t.Run("Create", testJobCreate(ctx, jobs))
	t.Run("List", testJobList(ctx, jobs))
	t.Run("Find", testJobFind(ctx, jobs))
	t.Run("Update", testJobUpdate(ctx, jobs))
	t.Run("Delete", testJobDelete(ctx, jobs))
}

func testJobCreate(ctx context.Context, jobs JobStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData := []model.Job{
			{Name: model.JobPinlCreated, TargetID: "1", Status: model.JobStatusPending},
			{Name: model.JobPinlUpdated, TargetID: "2", Status: model.JobStatusPending},
		}

		for _, job := range testData {
			assert.Nil(t, jobs.Create(ctx, &job))
			assert.NotEmpty(t, job.ID)
			assert.False(t, job.CreatedAt.Time().IsZero())
		}
	}
}

func testJobList(ctx context.Context, jobs JobStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, err := jobs.List(ctx, nil)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(testData))

		want := testData
		got, err := jobs.List(ctx, &JobOpts{Status: model.JobStatusPending})
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	}
}

func testJobFind(ctx context.Context, jobs JobStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, _ := jobs.List(ctx, nil)

		want := testData[0]
		got := model.Job{ID: want.ID}
		assert.Nil(t, jobs.Find(ctx, &got))
		assert.Equal(t, want, got)
	}
}

func testJobUpdate(ctx context.Context, jobs JobStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, _ := jobs.List(ctx, nil)

		want := testData[0]
		want.Status = model.JobStatusQueue
		assert.Nil(t, jobs.Update(ctx, &want))

		got := model.Job{ID: want.ID}
		jobs.Find(ctx, &got)
		assert.Equal(t, want, got)
	}
}

func testJobDelete(ctx context.Context, jobs JobStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, _ := jobs.List(ctx, nil)

		assert.Nil(t, jobs.Delete(ctx, &testData[0]))

		testData, _ = jobs.List(ctx, nil)
		assert.Equal(t, 1, len(testData))
	}
}
