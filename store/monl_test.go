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

	store := NewStore(db)
	monls := NewMonlStore(store)
	ctx := context.TODO()
	t.Run("Create", testMonlCreate(ctx, monls))
	t.Run("List", testMonlList(ctx, monls))
	t.Run("Find", testMonlFind(ctx, monls))
	t.Run("Update", testMonlUpdate(ctx, monls))
	t.Run("Delete", testMonlDelete(ctx, monls))
}

func testMonlCreate(ctx context.Context, monls MonlStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData := []model.Monl{
			{URL: "http://url.one", Title: "one", Description: "test pinl one"},
			{URL: "http://url.two", Title: "two", Description: "test pinl two"},
			{URL: "http://url.three", Title: "three", Description: "test pinl three"},
		}

		for _, monl := range testData {
			assert.Nil(t, monls.Create(ctx, &monl))
			assert.NotEmpty(t, monl.ID)
			assert.False(t, monl.CreatedAt.Time().IsZero())
		}
	}
}

func testMonlList(ctx context.Context, monls MonlStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, err := monls.List(ctx, nil)
		assert.Nil(t, err)
		assert.Equal(t, 3, len(testData))

		want := testData[1:2]
		got, err := monls.List(ctx, &MonlOpts{URL: want[0].URL})
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	}
}

func testMonlFind(ctx context.Context, monls MonlStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, _ := monls.List(ctx, nil)

		want := testData[0]
		got := model.Monl{ID: want.ID}
		assert.Nil(t, monls.Find(ctx, &got))
		assert.Equal(t, want, got)
	}
}

func testMonlUpdate(ctx context.Context, monls MonlStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, _ := monls.List(ctx, nil)

		want := testData[0]
		want.Title = "(changed) " + want.Title
		assert.Nil(t, monls.Update(ctx, &want))

		got := model.Monl{ID: want.ID}
		monls.Find(ctx, &got)
		assert.Equal(t, want, got)
	}
}

func testMonlDelete(ctx context.Context, monls MonlStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, _ := monls.List(ctx, nil)

		assert.Nil(t, monls.Delete(ctx, &testData[0]))

		testData, _ = monls.List(ctx, nil)
		assert.Equal(t, 2, len(testData))
	}
}
