package store

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pinmonl/pinmonl/database/dbtest"
	"github.com/pinmonl/pinmonl/model"
	"github.com/stretchr/testify/assert"
)

func TestShares(t *testing.T) {
	db, mock, err := dbtest.New()
	assert.Nil(t, err)
	defer db.Close()

	ctx := context.TODO()
	s := NewStore(db)
	shares := NewShares(s)

	t.Run("list", testSharesList(ctx, shares, mock))
	t.Run("count", testSharesCount(ctx, shares, mock))
	t.Run("find", testSharesFind(ctx, shares, mock))
	t.Run("create", testSharesCreate(ctx, shares, mock))
	t.Run("update", testSharesUpdate(ctx, shares, mock))
	t.Run("delete", testSharesDelete(ctx, shares, mock))
}

func testSharesList(ctx context.Context, shares *Shares, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			prefix = "SELECT (.+) FROM shares"
			opts   *ShareOpts
			list   []*model.Share
			err    error
		)

		// Test nil opts.
		opts = nil
		mock.ExpectQuery(prefix).
			WillReturnRows(sqlmock.NewRows(shares.columns()).
				AddRow("share-id-1", "user-id-1", "user/share", "share name", "description", "", model.PublishedShare, nil, nil))
		list, err = shares.List(ctx, opts)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(list))
	}
}

func testSharesCount(ctx context.Context, shares *Shares, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = regexp.QuoteMeta("SELECT count(*) FROM shares")
			opts  *ShareOpts
			count int64
			err   error
		)

		opts = &ShareOpts{}
		mock.ExpectQuery(query).
			WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).
				AddRow(1))
		count, err = shares.Count(ctx, opts)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), count)
	}
}

func testSharesFind(ctx context.Context, shares *Shares, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = "SELECT (.+) FROM shares WHERE id = \\?"
			id    string
			share *model.Share
			err   error
		)

		id = "share-id-1"
		mock.ExpectQuery(query).
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows(shares.columns()).
				AddRow(id, "user-id-1", "user/share", "share name", "description", "", model.PublishedShare, nil, nil))
		share, err = shares.Find(ctx, id)
		assert.Nil(t, err)
		if assert.NotNil(t, share) {
			assert.Equal(t, id, share.ID)
		}
	}
}

func testSharesCreate(ctx context.Context, shares *Shares, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			share *model.Share
			err   error
		)

		share = &model.Share{}
		expectSharesCreate(mock, share)
		err = shares.Create(ctx, share)
		assert.Nil(t, err)
		assert.NotEmpty(t, share.ID)
		assert.NotEmpty(t, share.CreatedAt)
	}
}

func expectSharesCreate(mock sqlmock.Sqlmock, share *model.Share) {
	mock.ExpectExec("INSERT INTO shares").
		WithArgs(
			sqlmock.AnyArg(),
			share.UserID,
			share.Slug,
			share.Name,
			share.Description,
			share.ImageID,
			share.Status,
			sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func testSharesUpdate(ctx context.Context, shares *Shares, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			share *model.Share
			err   error
		)

		share = &model.Share{ID: "share-id-1"}
		expectSharesUpdate(mock, share)
		err = shares.Update(ctx, share)
		assert.Nil(t, err)
		assert.NotEmpty(t, share.UpdatedAt)
	}
}

func expectSharesUpdate(mock sqlmock.Sqlmock, share *model.Share) {
	mock.ExpectExec("UPDATE shares (.+) WHERE id = \\?").
		WithArgs(
			share.UserID,
			share.Slug,
			share.Name,
			share.Description,
			share.ImageID,
			share.Status,
			sqlmock.AnyArg(),
			share.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func testSharesDelete(ctx context.Context, shares *Shares, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = regexp.QuoteMeta("DELETE FROM shares WHERE id = ?")
			id    string
			n     int64
			err   error
		)

		id = "share-id-1"
		mock.ExpectExec(query).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))
		n, err = shares.Delete(ctx, id)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), n)
	}
}
