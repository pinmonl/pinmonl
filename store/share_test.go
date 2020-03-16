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

	mockData := []*model.Share{
		{Name: "Share 1", Description: "Share 1 desc", Readme: "Share 1 readme", UserID: "user1"},
		{Name: "Share 2", Description: "Share 2 desc", Readme: "Share 2 readme", UserID: "user2"},
		{Name: "Share 3", Description: "Share 3 desc", Readme: "Share 3 readme", UserID: "user3"},
	}

	store := NewStore(db)
	shares := NewShareStore(store)
	ctx := context.TODO()
	t.Run("Create", testShareCreate(ctx, shares, mockData))
	t.Run("List", testShareList(ctx, shares, mockData))
	t.Run("Find", testShareFind(ctx, shares, mockData))
	t.Run("FindByName", testShareFindByName(ctx, shares, mockData))
	t.Run("Count", testShareCount(ctx, shares, mockData))
	t.Run("Update", testShareUpdate(ctx, shares, mockData))
	t.Run("Delete", testShareDelete(ctx, shares, mockData))
}

func testShareCreate(ctx context.Context, shares ShareStore, mockData []*model.Share) func(*testing.T) {
	return func(t *testing.T) {
		for _, share := range mockData {
			assert.Nil(t, shares.Create(ctx, share))
			assert.NotEmpty(t, share.ID)
			assert.False(t, share.CreatedAt.Time().IsZero())
		}
	}
}

func testShareList(ctx context.Context, shares ShareStore, mockData []*model.Share) func(*testing.T) {
	return func(t *testing.T) {
		deRef := func(data []*model.Share) []model.Share {
			out := make([]model.Share, len(data))
			for i, ms := range data {
				m := *ms
				out[i] = m
			}
			return out
		}

		want := deRef(mockData)
		got, err := shares.List(ctx, nil)
		assert.Nil(t, err)
		assert.Equal(t, want, got)

		want = deRef(mockData[:1])
		got, err = shares.List(ctx, &ShareOpts{UserID: want[0].UserID})
		assert.Nil(t, err)
		assert.Equal(t, want, got)

		got, err = shares.List(ctx, &ShareOpts{Name: want[0].Name})
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	}
}

func testShareFind(ctx context.Context, shares ShareStore, mockData []*model.Share) func(*testing.T) {
	return func(t *testing.T) {
		want := mockData[0]
		got := model.Share{ID: want.ID}
		assert.Nil(t, shares.Find(ctx, &got))
		assert.Equal(t, *want, got)
	}
}

func testShareFindByName(ctx context.Context, shares ShareStore, mockData []*model.Share) func(*testing.T) {
	return func(t *testing.T) {
		want := mockData[0]
		got := model.Share{Name: want.Name}
		assert.Equal(t, sql.ErrNoRows, shares.FindByName(ctx, &got))

		got = model.Share{UserID: want.UserID}
		assert.Equal(t, sql.ErrNoRows, shares.FindByName(ctx, &got))

		got = model.Share{Name: want.Name, UserID: want.UserID}
		assert.Nil(t, shares.FindByName(ctx, &got))
		assert.Equal(t, *want, got)
	}
}

func testShareCount(ctx context.Context, shares ShareStore, mockData []*model.Share) func(*testing.T) {
	return func(t *testing.T) {
		want := int64(len(mockData))
		got, err := shares.Count(ctx, nil)
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	}
}

func testShareUpdate(ctx context.Context, shares ShareStore, mockData []*model.Share) func(*testing.T) {
	return func(t *testing.T) {
		want := mockData[0]
		want.Name = "(changed) " + want.Name
		assert.Nil(t, shares.Update(ctx, want))
		assert.False(t, want.UpdatedAt.Time().IsZero())

		got := model.Share{ID: want.ID}
		shares.Find(ctx, &got)
		assert.Equal(t, *want, got)
	}
}

func testShareDelete(ctx context.Context, shares ShareStore, mockData []*model.Share) func(*testing.T) {
	return func(t *testing.T) {
		del, want := mockData[0], mockData[1:]
		assert.Nil(t, shares.Delete(ctx, del))

		got, _ := shares.List(ctx, nil)
		assert.Equal(t, len(want), len(got))
	}
}
