package store

import (
	"context"
	"testing"

	"github.com/pinmonl/pinmonl/database/dbtest"
	"github.com/pinmonl/pinmonl/model"
	"github.com/stretchr/testify/assert"
)

func TestTagStore(t *testing.T) {
	db, err := dbtest.Open()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		dbtest.Reset(db)
		dbtest.Close(db)
	}()

	mockData := []*model.Tag{
		{Name: "tag1", UserID: "user-test-id1", Level: 0, Sort: 0},
		{Name: "tag2", UserID: "user-test-id1", Level: 0, Sort: 0},
		{Name: "tag3", UserID: "user-test-id2", Level: 0, Sort: 1},
		{Name: "tag4", UserID: "user-test-id2", Level: 0, Sort: 0},
		{Name: "tag5", UserID: "user-test-id3", Level: 0, Sort: 1},
		{Name: "tag6", UserID: "user-test-id3", Level: 0, Sort: 0},
	}

	store := NewStore(db)
	tags := NewTagStore(store)
	ctx := context.TODO()
	t.Run("Create", testTagCreate(ctx, tags, mockData))
	t.Run("List", testTagList(ctx, tags, mockData))
	t.Run("Find", testTagFind(ctx, tags, mockData))
	t.Run("FindByName", testTagFindByName(ctx, tags, mockData))
	t.Run("Count", testTagCount(ctx, tags, mockData))
	t.Run("Update", testTagUpdate(ctx, tags, mockData))
	t.Run("Delete", testTagDelete(ctx, tags, mockData))
}

func testTagCreate(ctx context.Context, tags TagStore, mockData []*model.Tag) func(t *testing.T) {
	return func(t *testing.T) {
		for _, tag := range mockData {
			assert.Nil(t, tags.Create(ctx, tag))
			assert.NotEmpty(t, tag.ID)
			assert.False(t, tag.CreatedAt.Time().IsZero())
		}
	}
}

func testTagList(ctx context.Context, tags TagStore, mockData []*model.Tag) func(t *testing.T) {
	return func(t *testing.T) {
		deRef := func(data []*model.Tag) []model.Tag {
			out := make([]model.Tag, len(data))
			for i, mt := range data {
				m := *mt
				out[i] = m
			}
			return out
		}

		want := deRef(mockData)
		got, err := tags.List(ctx, nil)
		assert.Nil(t, err)
		assert.Equal(t, want, got)

		want = deRef(mockData[0:2])
		got, err = tags.List(ctx, &TagOpts{UserID: want[0].UserID})
		assert.Nil(t, err)
		assert.Equal(t, want, got)

		want = deRef(mockData[2:4])
		opts := TagOpts{}
		for _, t := range want {
			opts.Names = append(opts.Names, t.Name)
		}
		got, err = tags.List(ctx, &opts)
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	}
}

func testTagFind(ctx context.Context, tags TagStore, mockData []*model.Tag) func(t *testing.T) {
	return func(t *testing.T) {
		want := mockData[0]
		got := model.Tag{ID: want.ID}
		assert.Nil(t, tags.Find(ctx, &got))
		assert.Equal(t, *want, got)
	}
}

func testTagFindByName(ctx context.Context, tags TagStore, mockData []*model.Tag) func(t *testing.T) {
	return func(t *testing.T) {
		want := mockData[0]
		got := model.Tag{UserID: want.UserID, Name: want.Name}
		assert.Nil(t, tags.FindByName(ctx, &got))
		assert.Equal(t, *want, got)
	}
}

func testTagCount(ctx context.Context, tags TagStore, mockData []*model.Tag) func(t *testing.T) {
	return func(t *testing.T) {
		want := int64(len(mockData))
		got, err := tags.Count(ctx, nil)
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	}
}

func testTagUpdate(ctx context.Context, tags TagStore, mockData []*model.Tag) func(t *testing.T) {
	return func(t *testing.T) {
		want := mockData[1]
		want.Name = "(changed)" + want.Name
		assert.Nil(t, tags.Update(ctx, want))

		got := model.Tag{ID: want.ID}
		tags.Find(ctx, &got)
		assert.Equal(t, *want, got)
	}
}

func testTagDelete(ctx context.Context, tags TagStore, mockData []*model.Tag) func(t *testing.T) {
	return func(t *testing.T) {
		del, want := mockData[0], mockData[1:]
		assert.Nil(t, tags.Delete(ctx, del))

		got, _ := tags.List(ctx, nil)
		assert.Equal(t, len(want), len(got))
	}
}
