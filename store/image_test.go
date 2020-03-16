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

	mockData := []*model.Image{
		{Content: []byte{0x00, 0x02}, Size: 123456},
		{Content: []byte{0x02, 0x03}, Size: 9},
	}

	store := NewStore(db)
	images := NewImageStore(store)
	ctx := context.TODO()
	t.Run("Create", testImageCreate(ctx, images, mockData))
	t.Run("List", testImageList(ctx, images, mockData))
	t.Run("Find", testImageFind(ctx, images, mockData))
	t.Run("Update", testImageUpdate(ctx, images, mockData))
	t.Run("Delete", testImageDelete(ctx, images, mockData))
}

func testImageCreate(ctx context.Context, images ImageStore, mockData []*model.Image) func(*testing.T) {
	return func(t *testing.T) {
		for _, image := range mockData {
			assert.Nil(t, images.Create(ctx, image))
			assert.NotEmpty(t, image.ID)
			assert.False(t, image.CreatedAt.Time().IsZero())
		}
	}
}

func testImageList(ctx context.Context, images ImageStore, mockData []*model.Image) func(*testing.T) {
	return func(t *testing.T) {
		deRef := func(data []*model.Image) []model.Image {
			out := make([]model.Image, len(data))
			for i, mi := range data {
				m := *mi
				out[i] = m
			}
			return out
		}

		want := deRef(mockData)
		got, err := images.List(ctx, nil)
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	}
}

func testImageFind(ctx context.Context, images ImageStore, mockData []*model.Image) func(*testing.T) {
	return func(t *testing.T) {
		want := mockData[0]
		got := model.Image{ID: want.ID}
		assert.Nil(t, images.Find(ctx, &got))
		assert.Equal(t, *want, got)
	}
}

func testImageUpdate(ctx context.Context, images ImageStore, mockData []*model.Image) func(*testing.T) {
	return func(t *testing.T) {
		want := mockData[1]
		want.Content = []byte{0x03, 0x09}
		assert.Nil(t, images.Update(ctx, want))
		assert.False(t, want.UpdatedAt.Time().IsZero())

		got := model.Image{ID: want.ID}
		images.Find(ctx, &got)
		assert.Equal(t, *want, got)
	}
}

func testImageDelete(ctx context.Context, images ImageStore, mockData []*model.Image) func(*testing.T) {
	return func(t *testing.T) {
		del, want := mockData[0], mockData[1:]
		assert.Nil(t, images.Delete(ctx, del))

		got, _ := images.List(ctx, nil)
		assert.Equal(t, len(want), len(got))
	}
}
