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

	mockData := []*model.Pinl{
		{URL: "http://url.one", Title: "one", Description: "test pinl one", UserID: "user1"},
		{URL: "http://url.two", Title: "two", Description: "test pinl two", UserID: "user2"},
		{URL: "http://url.three", Title: "three", Description: "test pinl three", UserID: "user3"},
	}

	store := NewStore(db)
	pinls := NewPinlStore(store)
	ctx := context.TODO()
	t.Run("Create", testPinlCreate(ctx, pinls, mockData))
	t.Run("List", testPinlList(ctx, pinls, mockData))
	t.Run("Find", testPinlFind(ctx, pinls, mockData))
	t.Run("Count", testPinlCount(ctx, pinls, mockData))
	t.Run("Update", testPinlUpdate(ctx, pinls, mockData))
	t.Run("Delete", testPinlDelete(ctx, pinls, mockData))
}

func testPinlCreate(ctx context.Context, pinls PinlStore, mockData []*model.Pinl) func(t *testing.T) {
	return func(t *testing.T) {
		for _, pinl := range mockData {
			assert.Nil(t, pinls.Create(ctx, pinl))
			assert.NotEmpty(t, pinl.ID)
			assert.False(t, pinl.CreatedAt.Time().IsZero())
		}
	}
}

func testPinlList(ctx context.Context, pinls PinlStore, mockData []*model.Pinl) func(t *testing.T) {
	return func(t *testing.T) {
		deRef := func(data []*model.Pinl) []model.Pinl {
			out := make([]model.Pinl, len(data))
			for i, mp := range data {
				m := *mp
				out[i] = m
			}
			return out
		}

		want := deRef(mockData)
		got, err := pinls.List(ctx, nil)
		assert.Nil(t, err)
		assert.Equal(t, want, got)

		want = deRef(mockData[:1])
		got, err = pinls.List(ctx, &PinlOpts{ID: want[0].ID})
		assert.Nil(t, err)
		assert.Equal(t, want, got)

		got, err = pinls.List(ctx, &PinlOpts{UserID: want[0].UserID})
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	}
}

func testPinlFind(ctx context.Context, pinls PinlStore, mockData []*model.Pinl) func(t *testing.T) {
	return func(t *testing.T) {
		want := mockData[2]
		got := model.Pinl{ID: want.ID}
		err := pinls.Find(ctx, &got)
		assert.Nil(t, err)
		assert.Equal(t, *want, got)
	}
}

func testPinlCount(ctx context.Context, pinls PinlStore, mockData []*model.Pinl) func(t *testing.T) {
	return func(t *testing.T) {
		want := int64(len(mockData))
		got, err := pinls.Count(ctx, nil)
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	}
}

func testPinlUpdate(ctx context.Context, pinls PinlStore, mockData []*model.Pinl) func(t *testing.T) {
	return func(t *testing.T) {
		want := mockData[1]
		want.Title = "(changed) " + want.Title
		assert.Nil(t, pinls.Update(ctx, want))
		assert.False(t, want.UpdatedAt.Time().IsZero())

		got := model.Pinl{ID: want.ID}
		pinls.Find(ctx, &got)
		assert.Equal(t, *want, got)
	}
}

func testPinlDelete(ctx context.Context, pinls PinlStore, mockData []*model.Pinl) func(t *testing.T) {
	return func(t *testing.T) {
		del, want := mockData[0], mockData[1:]
		assert.Nil(t, pinls.Delete(ctx, del))

		got, _ := pinls.List(ctx, nil)
		assert.Equal(t, len(want), len(got))
	}
}
