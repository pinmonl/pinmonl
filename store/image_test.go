package store

import (
	"context"
	"testing"

	"github.com/pinmonl/pinmonl/database/dbtest"
	"github.com/pinmonl/pinmonl/model"
	"github.com/stretchr/testify/assert"
)

func TestImageStore(t *testing.T) {
	db, err := dbtest.Open()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		dbtest.Reset(db)
		dbtest.Close(db)
	}()

	store := NewStore(db)
	images := NewImageStore(store)
	ctx := context.TODO()
	t.Run("Create", testImageCreate(ctx, images))
	t.Run("List", testImageList(ctx, images))
	t.Run("Find", testImageFind(ctx, images))
	t.Run("Update", testImageUpdate(ctx, images))
	t.Run("Delete", testImageDelete(ctx, images))
}

func testImageCreate(ctx context.Context, images ImageStore) func(*testing.T) {
	return func(t *testing.T) {
		testData := []model.Image{
			{Content: []byte{0x00, 0x02}, Size: 123456},
			{Content: []byte{0x02, 0x03}, Size: 9},
		}

		for _, image := range testData {
			assert.Nil(t, images.Create(ctx, &image))
			assert.NotEmpty(t, image.ID)
			assert.False(t, image.CreatedAt.Time().IsZero())
		}
	}
}

func testImageList(ctx context.Context, images ImageStore) func(*testing.T) {
	return func(t *testing.T) {
		testData, err := images.List(ctx, nil)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(testData))
	}
}

func testImageFind(ctx context.Context, images ImageStore) func(*testing.T) {
	return func(t *testing.T) {
		testData, _ := images.List(ctx, nil)

		want := testData[0]
		got := model.Image{ID: want.ID}
		assert.Nil(t, images.Find(ctx, &got))
		assert.Equal(t, want, got)
	}
}

func testImageUpdate(ctx context.Context, images ImageStore) func(*testing.T) {
	return func(t *testing.T) {
		testData, _ := images.List(ctx, nil)

		want := testData[1]
		want.Content = []byte{0x03, 0x09}
		assert.Nil(t, images.Update(ctx, &want))
		assert.False(t, want.UpdatedAt.Time().IsZero())

		got := model.Image{ID: want.ID}
		images.Find(ctx, &got)
		assert.Equal(t, want, got)
	}
}

func testImageDelete(ctx context.Context, images ImageStore) func(*testing.T) {
	return func(t *testing.T) {
		testData, _ := images.List(ctx, nil)

		assert.Nil(t, images.Delete(ctx, &testData[0]))

		testData, _ = images.List(ctx, nil)
		assert.Equal(t, 1, len(testData))
	}
}
