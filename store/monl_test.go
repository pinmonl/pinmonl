package store

import (
	"context"
	"testing"

	"github.com/pinmonl/pinmonl/database/dbtest"
	"github.com/pinmonl/pinmonl/model"
	"github.com/stretchr/testify/assert"
)

func TestMonlStore(t *testing.T) {
	db, err := dbtest.Open()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		dbtest.Reset(db)
		dbtest.Close(db)
	}()

	mockData := []*model.Monl{
		{URL: "http://url.one", Title: "one", Description: "test pinl one"},
		{URL: "http://url.two", Title: "two", Description: "test pinl two"},
		{URL: "http://url.three", Title: "three", Description: "test pinl three"},
	}

	store := NewStore(db)
	monls := NewMonlStore(store)
	ctx := context.TODO()
	t.Run("Create", testMonlCreate(ctx, monls, mockData))
	t.Run("List", testMonlList(ctx, monls, mockData))
	t.Run("Find", testMonlFind(ctx, monls, mockData))
	t.Run("Update", testMonlUpdate(ctx, monls, mockData))
	t.Run("Delete", testMonlDelete(ctx, monls, mockData))
}

func testMonlCreate(ctx context.Context, monls MonlStore, mockData []*model.Monl) func(t *testing.T) {
	return func(t *testing.T) {
		for _, monl := range mockData {
			assert.Nil(t, monls.Create(ctx, monl))
			assert.NotEmpty(t, monl.ID)
			assert.False(t, monl.CreatedAt.Time().IsZero())
		}
	}
}

func testMonlList(ctx context.Context, monls MonlStore, mockData []*model.Monl) func(t *testing.T) {
	return func(t *testing.T) {
		deRef := func(data []*model.Monl) []model.Monl {
			out := make([]model.Monl, len(data))
			for i, mm := range data {
				m := *mm
				out[i] = m
			}
			return out
		}

		want := deRef(mockData)
		got, err := monls.List(ctx, nil)
		assert.Nil(t, err)
		assert.Equal(t, want, got)

		want = deRef(mockData[1:2])
		got, err = monls.List(ctx, &MonlOpts{URL: want[0].URL})
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	}
}

func testMonlFind(ctx context.Context, monls MonlStore, mockData []*model.Monl) func(t *testing.T) {
	return func(t *testing.T) {
		want := mockData[0]
		got := model.Monl{ID: want.ID}
		assert.Nil(t, monls.Find(ctx, &got))
		assert.Equal(t, *want, got)
	}
}

func testMonlUpdate(ctx context.Context, monls MonlStore, mockData []*model.Monl) func(t *testing.T) {
	return func(t *testing.T) {
		want := mockData[0]
		want.Title = "(changed) " + want.Title
		assert.Nil(t, monls.Update(ctx, want))

		got := model.Monl{ID: want.ID}
		monls.Find(ctx, &got)
		assert.Equal(t, *want, got)
	}
}

func testMonlDelete(ctx context.Context, monls MonlStore, mockData []*model.Monl) func(t *testing.T) {
	return func(t *testing.T) {
		del, want := mockData[0], mockData[1:]
		assert.Nil(t, monls.Delete(ctx, del))

		got, _ := monls.List(ctx, nil)
		assert.Equal(t, len(want), len(got))
	}
}
