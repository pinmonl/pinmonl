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

func TestPinls(t *testing.T) {
	db, mock, err := dbtest.New()
	assert.Nil(t, err)
	defer db.Close()

	ctx := context.TODO()
	s := NewStore(db)
	pinls := NewPinls(s)

	t.Run("list", testPinlsList(ctx, pinls, mock))
	t.Run("count", testPinlsCount(ctx, pinls, mock))
	t.Run("find", testPinlsFind(ctx, pinls, mock))
	t.Run("create", testPinlsCreate(ctx, pinls, mock))
	t.Run("update", testPinlsUpdate(ctx, pinls, mock))
	t.Run("delete", testPinlsDelete(ctx, pinls, mock))
}

func testPinlsList(ctx context.Context, pinls *Pinls, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			prefix = "SELECT (.+) FROM pinls"
			opts   *PinlOpts
			list   []*model.Pinl
			err    error
		)

		// Test nil opts.
		opts = nil
		mock.ExpectQuery(prefix).
			WillReturnRows(sqlmock.NewRows(pinls.columns()).
				AddRow("pinl-id-1", "user-id-1", "monl-id-1", "http://somewhere.com", "title", "description", "", nil, nil))
		list, err = pinls.List(ctx, opts)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(list))

		// Test filter by user.
		// Test filter by users.
		// Test filter by monls.
	}
}

func testPinlsCount(ctx context.Context, pinls *Pinls, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = regexp.QuoteMeta("SELECT count(*) FROM pinls")
			opts  *PinlOpts
			count int64
			err   error
		)

		opts = &PinlOpts{}
		mock.ExpectQuery(query).
			WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).
				AddRow(1))
		count, err = pinls.Count(ctx, opts)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), count)
	}
}

func testPinlsFind(ctx context.Context, pinls *Pinls, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = "SELECT (.+) FROM pinls WHERE id = \\?"
			id    string
			pinl  *model.Pinl
			err   error
		)

		id = "pinl-id-1"
		mock.ExpectQuery(query).
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows(pinls.columns()).
				AddRow(id, "user-id-1", "monl-id-1", "http://somewhere.com", "title", "description", "", nil, nil))
		pinl, err = pinls.Find(ctx, id)
		assert.Nil(t, err)
		if assert.NotNil(t, pinl) {
			assert.Equal(t, id, pinl.ID)
		}
	}
}

func testPinlsCreate(ctx context.Context, pinls *Pinls, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			pinl *model.Pinl
			err  error
		)

		pinl = &model.Pinl{}
		expectPinlsCreate(mock, pinl)
		err = pinls.Create(ctx, pinl)
		assert.Nil(t, err)
		assert.NotEmpty(t, pinl.ID)
		assert.NotEmpty(t, pinl.CreatedAt)
	}
}

func expectPinlsCreate(mock sqlmock.Sqlmock, pinl *model.Pinl) {
	mock.ExpectExec("INSERT INTO pinls").
		WithArgs(
			sqlmock.AnyArg(),
			pinl.UserID,
			pinl.MonlID,
			pinl.URL,
			pinl.Title,
			pinl.Description,
			pinl.ImageID,
			sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func testPinlsUpdate(ctx context.Context, pinls *Pinls, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			pinl *model.Pinl
			err  error
		)

		pinl = &model.Pinl{ID: "pinl-id-1"}
		expectPinlsUpdate(mock, pinl)
		err = pinls.Update(ctx, pinl)
		assert.Nil(t, err)
		assert.NotEmpty(t, pinl.UpdatedAt)
	}
}

func expectPinlsUpdate(mock sqlmock.Sqlmock, pinl *model.Pinl) {
	mock.ExpectExec("UPDATE pinls (.+) WHERE id = \\?").
		WithArgs(
			pinl.UserID,
			pinl.MonlID,
			pinl.URL,
			pinl.Title,
			pinl.Description,
			pinl.ImageID,
			sqlmock.AnyArg(),
			pinl.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func testPinlsDelete(ctx context.Context, pinls *Pinls, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = regexp.QuoteMeta("DELETE FROM pinls WHERE id = ?")
			id    string
			n     int64
			err   error
		)

		id = "pinl-id-1"
		mock.ExpectExec(query).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))
		n, err = pinls.Delete(ctx, id)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), n)
	}
}
