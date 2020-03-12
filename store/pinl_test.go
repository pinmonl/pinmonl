package store

import (
	"context"
	"testing"

	"github.com/pinmonl/pinmonl/database/dbtest"
	"github.com/pinmonl/pinmonl/model"
	"github.com/stretchr/testify/assert"
)

func TestPinlStore(t *testing.T) {
	db, err := dbtest.Open()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		dbtest.Reset(db)
		dbtest.Close(db)
	}()

	store := NewStore(db)
	pinls := NewPinlStore(store)
	ctx := context.TODO()
	t.Run("Create", testPinlCreate(ctx, pinls))
	t.Run("List", testPinlList(ctx, pinls))
	t.Run("Find", testPinlFind(ctx, pinls))
	t.Run("Count", testPinlCount(ctx, pinls))
	t.Run("Update", testPinlUpdate(ctx, pinls))
	t.Run("Delete", testPinlDelete(ctx, pinls))
}

func testPinlCreate(ctx context.Context, pinls PinlStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData := []model.Pinl{
			{URL: "http://url.one", Title: "one", Description: "test pinl one", UserID: "user1"},
			{URL: "http://url.two", Title: "two", Description: "test pinl two", UserID: "user2"},
			{URL: "http://url.three", Title: "three", Description: "test pinl three", UserID: "user3"},
		}

		for _, pinl := range testData {
			assert.Nil(t, pinls.Create(ctx, &pinl))
			assert.NotEmpty(t, pinl.ID)
			assert.False(t, pinl.CreatedAt.Time().IsZero())
		}
	}
}

func testPinlList(ctx context.Context, pinls PinlStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, _ := pinls.List(ctx, nil)
		assert.Equal(t, 3, len(testData))

		want := testData[:1]
		got, err := pinls.List(ctx, &PinlOpts{ID: want[0].ID})
		assert.Nil(t, err)
		assert.Equal(t, want, got)

		got, err = pinls.List(ctx, &PinlOpts{UserID: want[0].UserID})
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	}
}

func testPinlFind(ctx context.Context, pinls PinlStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, _ := pinls.List(ctx, nil)

		want := testData[2]
		got := model.Pinl{ID: want.ID}
		err := pinls.Find(ctx, &got)
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	}
}

func testPinlCount(ctx context.Context, pinls PinlStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, _ := pinls.List(ctx, nil)

		want := int64(len(testData))
		got, err := pinls.Count(ctx, nil)
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	}
}

func testPinlUpdate(ctx context.Context, pinls PinlStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, _ := pinls.List(ctx, nil)

		want := testData[1]
		want.Title = "(changed) " + want.Title
		assert.Nil(t, pinls.Update(ctx, &want))
		assert.False(t, want.UpdatedAt.Time().IsZero())

		got := model.Pinl{ID: want.ID}
		pinls.Find(ctx, &got)
		assert.Equal(t, want, got)
	}
}

func testPinlDelete(ctx context.Context, pinls PinlStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, _ := pinls.List(ctx, nil)

		assert.Nil(t, pinls.Delete(ctx, &testData[0]))

		count, _ := pinls.Count(ctx, nil)
		assert.Equal(t, int64(len(testData)-1), count)
	}
}
