package store

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pinmonl/pinmonl/database/dbtest"
	"github.com/pinmonl/pinmonl/model"
	"github.com/stretchr/testify/assert"
)

func TestImages(t *testing.T) {
	db, mock, err := dbtest.New()
	assert.Nil(t, err)
	defer db.Close()

	ctx := context.TODO()
	s := NewStore(db)
	images := NewImages(s)

	t.Run("list", testImagesList(ctx, images, mock))
	t.Run("count", testImagesCount(ctx, images, mock))
	t.Run("find", testImagesFind(ctx, images, mock))
	t.Run("create", testImagesCreate(ctx, images, mock))
	t.Run("update", testImagesUpdate(ctx, images, mock))
	t.Run("delete", testImagesDelete(ctx, images, mock))
}

func testImagesList(ctx context.Context, images *Images, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			prefix = "SELECT (.+) FROM images"
			opts   *ImageOpts
			list   []*model.Image
			err    error
		)

		// Test nil opts.
		opts = nil
		mock.ExpectQuery(prefix).
			WillReturnRows(sqlmock.NewRows(images.columns()).
				AddRow("image-id-1", "target-1", "target", "", "description", 1, "image/png", nil, nil))
		list, err = images.List(ctx, opts)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(list))

		// Test filter by targets.
		opts = &ImageOpts{Targets: model.MorphableList{
			&model.Pinl{ID: "pinl-id-1"},
			&model.Pinl{ID: "pinl-id-2"},
			&model.Pinl{ID: "pinl-id-3"},
		}}
		mock.ExpectQuery(fmt.Sprintf(regexp.QuoteMeta("%s WHERE target_name = ? AND target_id IN (?,?,?)"), prefix)).
			WithArgs(
				opts.Targets[0].MorphName(),
				opts.Targets[0].MorphKey(),
				opts.Targets[1].MorphKey(),
				opts.Targets[2].MorphKey()).
			WillReturnRows(sqlmock.NewRows(images.columns()))
		_, err = images.List(ctx, opts)
		assert.Nil(t, err)
	}
}

func testImagesCount(ctx context.Context, images *Images, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = regexp.QuoteMeta("SELECT count(*) FROM images")
			opts  *ImageOpts
			count int64
			err   error
		)

		opts = &ImageOpts{}
		mock.ExpectQuery(query).
			WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).
				AddRow(1))
		count, err = images.Count(ctx, opts)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), count)
	}
}

func testImagesFind(ctx context.Context, images *Images, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = "SELECT (.+) FROM images WHERE id = \\?"
			id    string
			image *model.Image
			err   error
		)

		id = "image-id-1"
		mock.ExpectQuery(query).
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows(images.columns()).
				AddRow(id, "target-1", "target", "", "description", 1, "image/png", nil, nil))
		image, err = images.Find(ctx, id)
		assert.Nil(t, err)
		if assert.NotNil(t, image) {
			assert.Equal(t, id, image.ID)
		}
	}
}

func testImagesCreate(ctx context.Context, images *Images, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			image *model.Image
			err   error
		)

		image = &model.Image{}
		expectImagesCreate(mock, image)
		err = images.Create(ctx, image)
		assert.Nil(t, err)
		assert.NotEmpty(t, image.ID)
		assert.NotEmpty(t, image.CreatedAt)
	}
}

func expectImagesCreate(mock sqlmock.Sqlmock, image *model.Image) {
	mock.ExpectExec("INSERT INTO images").
		WithArgs(
			sqlmock.AnyArg(),
			image.TargetID,
			image.TargetName,
			image.Content,
			image.Description,
			image.Size,
			image.ContentType,
			sqlmock.AnyArg(),
			sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func testImagesUpdate(ctx context.Context, images *Images, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			image *model.Image
			err   error
		)

		image = &model.Image{ID: "image-id-1"}
		expectImagesUpdate(mock, image)
		err = images.Update(ctx, image)
		assert.Nil(t, err)
		assert.NotEmpty(t, image.UpdatedAt)
	}
}

func expectImagesUpdate(mock sqlmock.Sqlmock, image *model.Image) {
	mock.ExpectExec("UPDATE images (.+) WHERE id = \\?").
		WithArgs(
			image.TargetID,
			image.TargetName,
			image.Content,
			image.Description,
			image.Size,
			image.ContentType,
			sqlmock.AnyArg(),
			image.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func testImagesDelete(ctx context.Context, images *Images, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = regexp.QuoteMeta("DELETE FROM images WHERE id = ?")
			id    string
			n     int64
			err   error
		)

		id = "image-id-1"
		mock.ExpectExec(query).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))
		n, err = images.Delete(ctx, id)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), n)
	}
}
