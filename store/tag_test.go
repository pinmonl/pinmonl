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

	store := NewStore(db)
	tags := NewTagStore(store)
	ctx := context.TODO()
	t.Run("Create", testTagCreate(ctx, tags))
	t.Run("List", testTagList(ctx, tags))
	t.Run("Find", testTagFind(ctx, tags))
	t.Run("FindByName", testTagFindByName(ctx, tags))
	t.Run("Count", testTagCount(ctx, tags))
	t.Run("Update", testTagUpdate(ctx, tags))
	t.Run("Delete", testTagDelete(ctx, tags))
}

func testTagCreate(ctx context.Context, tags TagStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData := []model.Tag{
			{Name: "tag1", UserID: "user-test-id1", Level: 0, Sort: 0},
			{Name: "tag2", UserID: "user-test-id1", Level: 0, Sort: 0},
			{Name: "tag3", UserID: "user-test-id2", Level: 0, Sort: 1},
			{Name: "tag4", UserID: "user-test-id2", Level: 0, Sort: 0},
			{Name: "tag5", UserID: "user-test-id3", Level: 0, Sort: 1},
			{Name: "tag6", UserID: "user-test-id3", Level: 0, Sort: 0},
		}

		for _, tag := range testData {
			assert.Nil(t, tags.Create(ctx, &tag))
			assert.NotEmpty(t, tag.ID)
			assert.False(t, tag.CreatedAt.Time().IsZero())
		}
	}
}

func testTagList(ctx context.Context, tags TagStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, err := tags.List(ctx, nil)
		assert.Nil(t, err)
		assert.Equal(t, 6, len(testData))

		want := testData[0:2]
		got, err := tags.List(ctx, &TagOpts{UserID: want[0].UserID})
		assert.Nil(t, err)
		assert.Equal(t, want, got)

		want = testData[2:4]
		opts := TagOpts{}
		for _, t := range want {
			opts.Names = append(opts.Names, t.Name)
		}
		got, err = tags.List(ctx, &opts)
		assert.Nil(t, err)
		assert.Equal(t, want, want)
	}
}

func testTagFind(ctx context.Context, tags TagStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, _ := tags.List(ctx, nil)

		want := testData[0]
		got := model.Tag{ID: want.ID}
		assert.Nil(t, tags.Find(ctx, &got))
		assert.Equal(t, want, got)
	}
}

func testTagFindByName(ctx context.Context, tags TagStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, _ := tags.List(ctx, nil)

		want := testData[0]
		got := model.Tag{UserID: want.UserID, Name: want.Name}
		assert.Nil(t, tags.FindByName(ctx, &got))
		assert.Equal(t, want, got)
	}
}

func testTagCount(ctx context.Context, tags TagStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, _ := tags.List(ctx, nil)

		want := int64(len(testData))
		got, err := tags.Count(ctx, nil)
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	}
}

func testTagUpdate(ctx context.Context, tags TagStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, _ := tags.List(ctx, nil)

		want := testData[1]
		want.Name = "(changed)" + want.Name
		assert.Nil(t, tags.Update(ctx, &want))

		got := model.Tag{ID: want.ID}
		tags.Find(ctx, &got)
		assert.Equal(t, want, got)
	}
}

func testTagDelete(ctx context.Context, tags TagStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, _ := tags.List(ctx, nil)

		assert.Nil(t, tags.Delete(ctx, &testData[0]))

		want := int64(len(testData) - 1)
		got, _ := tags.Count(ctx, nil)
		assert.Equal(t, want, got)
	}
}
