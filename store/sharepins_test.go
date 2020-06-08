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

func TestSharepins(t *testing.T) {
	db, mock, err := dbtest.New()
	assert.Nil(t, err)
	defer db.Close()

	ctx := context.TODO()
	s := NewStore(db)
	sharepins := NewSharepins(s)

	t.Run("list", testSharepinsList(ctx, sharepins, mock))
	t.Run("count", testSharepinsCount(ctx, sharepins, mock))
	t.Run("find", testSharepinsFind(ctx, sharepins, mock))
	t.Run("create", testSharepinsCreate(ctx, sharepins, mock))
	t.Run("update", testSharepinsUpdate(ctx, sharepins, mock))
	t.Run("delete", testSharepinsDelete(ctx, sharepins, mock))
}

func testSharepinsList(ctx context.Context, sharepins *Sharepins, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			prefix = "SELECT (.+) FROM sharepins"
			opts   *SharepinOpts
			list   []*model.Sharepin
			err    error
		)

		// Test nil opts.
		opts = nil
		mock.ExpectQuery(prefix).
			WillReturnRows(sqlmock.NewRows(sharepins.columns()).
				AddRow("sharepin-id-1", "share-id-1", "pinl-id-1", model.PublishedShare))
		list, err = sharepins.List(ctx, opts)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(list))

		// Test filter by shares.
		opts = &SharepinOpts{ShareIDs: []string{"share-id-1", "share-id-2"}}
		mock.ExpectQuery(fmt.Sprintf(regexp.QuoteMeta("%s WHERE share_id IN (?,?)"), prefix)).
			WithArgs(opts.ShareIDs[0], opts.ShareIDs[1]).
			WillReturnRows(sqlmock.NewRows(sharepins.columns()))
		_, err = sharepins.List(ctx, opts)
		assert.Nil(t, err)

		// Test filter by pinls.
		opts = &SharepinOpts{PinlIDs: []string{"pinl-id-1", "pinl-id-2", "pinl-id-3"}}
		mock.ExpectQuery(fmt.Sprintf(regexp.QuoteMeta("%s WHERE pinl_id IN (?,?,?)"), prefix)).
			WithArgs(opts.PinlIDs[0], opts.PinlIDs[1], opts.PinlIDs[2]).
			WillReturnRows(sqlmock.NewRows(sharepins.columns()))
		_, err = sharepins.List(ctx, opts)
		assert.Nil(t, err)
	}
}

func testSharepinsCount(ctx context.Context, sharepins *Sharepins, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = regexp.QuoteMeta("SELECT count(*) FROM sharepins")
			opts  *SharepinOpts
			count int64
			err   error
		)

		opts = &SharepinOpts{}
		mock.ExpectQuery(query).
			WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).
				AddRow(1))
		count, err = sharepins.Count(ctx, opts)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), count)
	}
}

func testSharepinsFind(ctx context.Context, sharepins *Sharepins, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query    = "SELECT (.+) FROM sharepins WHERE id = \\?"
			id       string
			sharepin *model.Sharepin
			err      error
		)

		id = "sharepin-id-1"
		mock.ExpectQuery(query).
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows(sharepins.columns()).
				AddRow("sharepin-id-1", "share-id-1", "pinl-id-1", model.PublishedShare))
		sharepin, err = sharepins.Find(ctx, id)
		assert.Nil(t, err)
		if assert.NotNil(t, sharepin) {
			assert.Equal(t, id, sharepin.ID)
		}
	}
}

func testSharepinsCreate(ctx context.Context, sharepins *Sharepins, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			sharepin *model.Sharepin
			err      error
		)

		sharepin = &model.Sharepin{}
		expectSharepinsCreate(mock, sharepin)
		err = sharepins.Create(ctx, sharepin)
		assert.Nil(t, err)
		assert.NotEmpty(t, sharepin.ID)
	}
}

func expectSharepinsCreate(mock sqlmock.Sqlmock, sharepin *model.Sharepin) {
	mock.ExpectExec("INSERT INTO sharepins").
		WithArgs(
			sqlmock.AnyArg(),
			sharepin.ShareID,
			sharepin.PinlID,
			sharepin.Status).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func testSharepinsUpdate(ctx context.Context, sharepins *Sharepins, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			sharepin *model.Sharepin
			err      error
		)

		sharepin = &model.Sharepin{ID: "sharepin-id-1"}
		expectSharepinsUpdate(mock, sharepin)
		err = sharepins.Update(ctx, sharepin)
		assert.Nil(t, err)
	}
}

func expectSharepinsUpdate(mock sqlmock.Sqlmock, sharepin *model.Sharepin) {
	mock.ExpectExec("UPDATE sharepins (.+) WHERE id = \\?").
		WithArgs(
			sharepin.ShareID,
			sharepin.PinlID,
			sharepin.Status,
			sharepin.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func testSharepinsDelete(ctx context.Context, sharepins *Sharepins, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = regexp.QuoteMeta("DELETE FROM sharepins WHERE id = ?")
			id    string
			n     int64
			err   error
		)

		id = "sharepin-id-1"
		mock.ExpectExec(query).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))
		n, err = sharepins.Delete(ctx, id)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), n)
	}
}
