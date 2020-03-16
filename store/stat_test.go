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

	mockData := []*model.Stat{
		{Kind: "tag", Value: "v1.0.0", RecordedAt: field.Time(time.Now().Round(time.Second).UTC()), IsLatest: false, Labels: field.Labels{"a1": "v1"}},
		{Kind: "tag", Value: "v2.0.0", RecordedAt: field.Time(time.Now().Round(time.Second).UTC()), IsLatest: true, Labels: field.Labels{"a2": "v2"}},
	}

	store := NewStore(db)
	stats := NewStatStore(store)
	ctx := context.TODO()
	t.Run("Create", testStatCreate(ctx, stats, mockData))
	t.Run("List", testStatList(ctx, stats, mockData))
	t.Run("Update", testStatUpdate(ctx, stats, mockData))
}

func testStatCreate(ctx context.Context, stats StatStore, mockData []*model.Stat) func(t *testing.T) {
	return func(t *testing.T) {
		for _, stat := range mockData {
			assert.Nil(t, stats.Create(ctx, stat))
			assert.NotEmpty(t, stat.ID)
		}
	}
}

func testStatList(ctx context.Context, stats StatStore, mockData []*model.Stat) func(t *testing.T) {
	return func(t *testing.T) {
		deRef := func(data []*model.Stat) []model.Stat {
			out := make([]model.Stat, len(data))
			for i, ms := range data {
				m := *ms
				out[i] = m
			}
			return out
		}

		want := deRef(mockData)
		got, err := stats.List(ctx, &StatOpts{Kind: "tag"})
		assert.Nil(t, err)
		assert.Equal(t, want, got)

		want = deRef(mockData[:1])
		got, err = stats.List(ctx, &StatOpts{WithoutLatest: true})
		assert.Nil(t, err)
		assert.Equal(t, want, got)

		want = deRef(mockData[1:2])
		got, err = stats.List(ctx, &StatOpts{WithLatest: true})
		assert.Nil(t, err)
		assert.Equal(t, want, got)

		want = deRef(mockData)
		got, err = stats.List(ctx, &StatOpts{Before: time.Now().Add(time.Minute)})
		assert.Nil(t, err)
		assert.Equal(t, want, got)

		want = nil
		got, err = stats.List(ctx, &StatOpts{After: time.Now().Add(time.Minute)})
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	}
}

func testStatUpdate(ctx context.Context, stats StatStore, mockData []*model.Stat) func(t *testing.T) {
	return func(t *testing.T) {
		want := mockData[0]
		want.Value = "v1.1.0"
		assert.Nil(t, stats.Update(ctx, want))
	}
}
