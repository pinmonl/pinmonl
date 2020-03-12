package store

import (
	"context"
	"database/sql"
	"testing"

	"github.com/pinmonl/pinmonl/database/dbtest"
	"github.com/pinmonl/pinmonl/model"
	"github.com/stretchr/testify/assert"
)

func TestShareStore(t *testing.T) {
	db, err := dbtest.Open()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		dbtest.Reset(db)
		dbtest.Close(db)
	}()

	store := NewStore(db)
	shares := NewShareStore(store)
	ctx := context.TODO()
	t.Run("Create", testShareCreate(ctx, shares))
	t.Run("List", testShareList(ctx, shares))
	t.Run("Find", testShareFind(ctx, shares))
	t.Run("FindByName", testShareFindByName(ctx, shares))
	t.Run("Count", testShareCount(ctx, shares))
	t.Run("Update", testShareUpdate(ctx, shares))
	t.Run("Delete", testShareDelete(ctx, shares))
}

func testShareCreate(ctx context.Context, shares ShareStore) func(*testing.T) {
	return func(t *testing.T) {
		testData := []model.Share{
			{Name: "Share 1", Description: "Share 1 desc", Readme: "Share 1 readme", UserID: "user1"},
			{Name: "Share 2", Description: "Share 2 desc", Readme: "Share 2 readme", UserID: "user2"},
			{Name: "Share 3", Description: "Share 3 desc", Readme: "Share 3 readme", UserID: "user3"},
		}

		for _, share := range testData {
			assert.Nil(t, shares.Create(ctx, &share))
			assert.NotEmpty(t, share.ID)
			assert.False(t, share.CreatedAt.Time().IsZero())
		}
	}
}

func testShareList(ctx context.Context, shares ShareStore) func(*testing.T) {
	return func(t *testing.T) {
		testData, err := shares.List(ctx, nil)
		assert.Nil(t, err)
		assert.Equal(t, 3, len(testData))

		want := testData[:1]
		got, err := shares.List(ctx, &ShareOpts{UserID: want[0].UserID})
		assert.Nil(t, err)
		assert.Equal(t, want, got)

		got, err = shares.List(ctx, &ShareOpts{Name: want[0].Name})
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	}
}

func testShareFind(ctx context.Context, shares ShareStore) func(*testing.T) {
	return func(t *testing.T) {
		testData, _ := shares.List(ctx, nil)

		want := testData[0]
		got := model.Share{ID: want.ID}
		assert.Nil(t, shares.Find(ctx, &got))
		assert.Equal(t, want, got)
	}
}

func testShareFindByName(ctx context.Context, shares ShareStore) func(*testing.T) {
	return func(t *testing.T) {
		testData, _ := shares.List(ctx, nil)

		want := testData[0]
		got := model.Share{Name: want.Name}
		assert.Equal(t, sql.ErrNoRows, shares.FindByName(ctx, &got))

		got = model.Share{UserID: want.UserID}
		assert.Equal(t, sql.ErrNoRows, shares.FindByName(ctx, &got))

		got = model.Share{Name: want.Name, UserID: want.UserID}
		assert.Nil(t, shares.FindByName(ctx, &got))
		assert.Equal(t, want, got)
	}
}

func testShareCount(ctx context.Context, shares ShareStore) func(*testing.T) {
	return func(t *testing.T) {
		testData, _ := shares.List(ctx, nil)

		got, err := shares.Count(ctx, nil)
		assert.Nil(t, err)
		want := int64(len(testData))
		assert.Equal(t, want, got)
	}
}

func testShareUpdate(ctx context.Context, shares ShareStore) func(*testing.T) {
	return func(t *testing.T) {
		testData, _ := shares.List(ctx, nil)

		want := testData[0]
		want.Name = "(changed) " + want.Name
		assert.Nil(t, shares.Update(ctx, &want))
		assert.False(t, want.UpdatedAt.Time().IsZero())

		got := model.Share{ID: want.ID}
		shares.Find(ctx, &got)
		assert.Equal(t, want, got)
	}
}

func testShareDelete(ctx context.Context, shares ShareStore) func(*testing.T) {
	return func(t *testing.T) {
		testData, _ := shares.List(ctx, nil)

		assert.Nil(t, shares.Delete(ctx, &testData[0]))

		count, _ := shares.Count(ctx, nil)
		assert.Equal(t, int64(2), count)
	}
}
