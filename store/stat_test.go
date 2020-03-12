package store

import (
	"context"
	"testing"
	"time"

	"github.com/pinmonl/pinmonl/database/dbtest"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/stretchr/testify/assert"
)

func TestStatStore(t *testing.T) {
	db, err := dbtest.Open()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		dbtest.Reset(db)
		dbtest.Close(db)
	}()

	store := NewStore(db)
	stats := NewStatStore(store)
	ctx := context.TODO()
	t.Run("Create", testStatCreate(ctx, stats))
	t.Run("List", testStatList(ctx, stats))
	t.Run("Update", testStatUpdate(ctx, stats))
}

func testStatCreate(ctx context.Context, stats StatStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData := []model.Stat{
			{Kind: "tag", Value: "v1.0.0", RecordedAt: field.Time(time.Now()), IsLatest: false},
			{Kind: "tag", Value: "v2.0.0", RecordedAt: field.Time(time.Now()), IsLatest: true},
		}

		for _, stat := range testData {
			assert.Nil(t, stats.Create(ctx, &stat))
			assert.NotEmpty(t, stat.ID)
		}
	}
}

func testStatList(ctx context.Context, stats StatStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, err := stats.List(ctx, &StatOpts{Kind: "tag"})
		assert.Nil(t, err)
		assert.Equal(t, 2, len(testData))

		want := testData[0:1]
		got, err := stats.List(ctx, &StatOpts{WithoutLatest: true})
		assert.Nil(t, err)
		assert.Equal(t, want, got)

		want = testData[1:2]
		got, err = stats.List(ctx, &StatOpts{WithLatest: true})
		assert.Nil(t, err)
		assert.Equal(t, want, got)

		want = testData
		got, err = stats.List(ctx, &StatOpts{Before: time.Now()})
		assert.Nil(t, err)
		assert.Equal(t, want, got)

		want = nil
		got, err = stats.List(ctx, &StatOpts{After: time.Now()})
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	}
}

func testStatUpdate(ctx context.Context, stats StatStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, _ := stats.List(ctx, nil)

		want := testData[0]
		want.Value = "v1.1.0"
		assert.Nil(t, stats.Update(ctx, &want))
	}
}
