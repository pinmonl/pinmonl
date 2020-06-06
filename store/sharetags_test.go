package store

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pinmonl/pinmonl/database/dbtest"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/stretchr/testify/assert"
)

func TestSharetags(t *testing.T) {
	db, mock, err := dbtest.New()
	assert.Nil(t, err)
	defer db.Close()

	ctx := context.TODO()
	s := NewStore(db)
	sharetags := NewSharetags(s)

	t.Run("list", testSharetagsList(ctx, sharetags, mock))
	t.Run("count", testSharetagsCount(ctx, sharetags, mock))
	t.Run("find", testSharetagsFind(ctx, sharetags, mock))
	t.Run("create", testSharetagsCreate(ctx, sharetags, mock))
	t.Run("update", testSharetagsUpdate(ctx, sharetags, mock))
	t.Run("delete", testSharetagsDelete(ctx, sharetags, mock))
}

func testSharetagsList(ctx context.Context, sharetags *Sharetags, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			prefix = "SELECT (.+) FROM sharetags"
			opts   *SharetagOpts
			list   []*model.Sharetag
			err    error
		)

		// Test nil opts.
		opts = nil
		mock.ExpectQuery(prefix).
			WillReturnRows(sqlmock.NewRows(sharetags.columns()).
				AddRow("sharetag-id-1", "share-id-1", "tag-id-1", model.SharetagMust, "", 0, false))
		list, err = sharetags.List(ctx, opts)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(list))

		// Test filter by kind.
		opts = &SharetagOpts{Kind: field.NewNullValue(model.SharetagMust)}
		mock.ExpectQuery(fmt.Sprintf("%s WHERE kind = ?", prefix)).
			WithArgs(opts.Kind.Value()).
			WillReturnRows(sqlmock.NewRows(sharetags.columns()))
		_, err = sharetags.List(ctx, opts)
		assert.Nil(t, err)
	}
}

func testSharetagsCount(ctx context.Context, sharetags *Sharetags, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = regexp.QuoteMeta("SELECT count(*) FROM sharetags")
			opts  *SharetagOpts
			count int64
			err   error
		)

		opts = &SharetagOpts{}
		mock.ExpectQuery(query).
			WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).
				AddRow(1))
		count, err = sharetags.Count(ctx, opts)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), count)
	}
}

func testSharetagsFind(ctx context.Context, sharetags *Sharetags, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query    = "SELECT (.+) FROM sharetags WHERE id = \\?"
			id       string
			sharetag *model.Sharetag
			err      error
		)

		id = "sharetag-id-1"
		mock.ExpectQuery(query).
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows(sharetags.columns()).
				AddRow("sharetag-id-1", "share-id-1", "tag-id-1", model.SharetagMust, "", 0, false))
		sharetag, err = sharetags.Find(ctx, id)
		assert.Nil(t, err)
		if assert.NotNil(t, sharetag) {
			assert.Equal(t, id, sharetag.ID)
		}
	}
}

func testSharetagsCreate(ctx context.Context, sharetags *Sharetags, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			sharetag *model.Sharetag
			err      error
		)

		sharetag = &model.Sharetag{}
		expectSharetagsCreate(mock, sharetag)
		err = sharetags.Create(ctx, sharetag)
		assert.Nil(t, err)
		assert.NotEmpty(t, sharetag.ID)
	}
}

func expectSharetagsCreate(mock sqlmock.Sqlmock, sharetag *model.Sharetag) {
	mock.ExpectExec("INSERT INTO sharetags").
		WithArgs(
			sqlmock.AnyArg(),
			sharetag.ShareID,
			sharetag.TagID,
			sharetag.Kind,
			sharetag.ParentID,
			sharetag.Level,
			sharetag.HasChildren).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func testSharetagsUpdate(ctx context.Context, sharetags *Sharetags, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			sharetag *model.Sharetag
			err      error
		)

		sharetag = &model.Sharetag{ID: "sharetag-id-1"}
		expectSharetagsUpdate(mock, sharetag)
		err = sharetags.Update(ctx, sharetag)
		assert.Nil(t, err)
	}
}

func expectSharetagsUpdate(mock sqlmock.Sqlmock, sharetag *model.Sharetag) {
	mock.ExpectExec("UPDATE sharetags (.+) WHERE id = \\?").
		WithArgs(
			sharetag.ShareID,
			sharetag.TagID,
			sharetag.Kind,
			sharetag.ParentID,
			sharetag.Level,
			sharetag.HasChildren,
			sharetag.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func testSharetagsDelete(ctx context.Context, sharetags *Sharetags, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = regexp.QuoteMeta("DELETE FROM sharetags WHERE id = ?")
			id    string
			n     int64
			err   error
		)

		id = "sharetag-id-1"
		mock.ExpectExec(query).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))
		n, err = sharetags.Delete(ctx, id)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), n)
	}
}
