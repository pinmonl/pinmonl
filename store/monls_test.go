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

func TestMonls(t *testing.T) {
	db, mock, err := dbtest.New()
	assert.Nil(t, err)
	defer db.Close()

	ctx := context.TODO()
	s := NewStore(db)
	monls := NewMonls(s)

	t.Run("list", testMonlsList(ctx, monls, mock))
	t.Run("count", testMonlsCount(ctx, monls, mock))
	t.Run("find", testMonlsFind(ctx, monls, mock))
	t.Run("create", testMonlsCreate(ctx, monls, mock))
	t.Run("update", testMonlsUpdate(ctx, monls, mock))
	t.Run("delete", testMonlsDelete(ctx, monls, mock))
}

func testMonlsList(ctx context.Context, monls *Monls, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			prefix = "SELECT (.+) FROM monls"
			opts   *MonlOpts
			list   []*model.Monl
			err    error
		)

		// Test nil opts.
		opts = nil
		mock.ExpectQuery(prefix).
			WillReturnRows(sqlmock.NewRows(monls.columns()).
				AddRow("monl-id-1", "https://somewhere.com", nil, nil))
		list, err = monls.List(ctx, opts)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(list))

		// Test filter by url.
		opts = &MonlOpts{URL: "http://somewhere.com"}
		mock.ExpectQuery(fmt.Sprintf("%s WHERE url = ?", prefix)).
			WithArgs(opts.URL).
			WillReturnRows(sqlmock.NewRows(monls.columns()))
		_, err = monls.List(ctx, opts)
		assert.Nil(t, err)
	}
}

func testMonlsCount(ctx context.Context, monls *Monls, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = regexp.QuoteMeta("SELECT count(*) FROM monls")
			opts  *MonlOpts
			count int64
			err   error
		)

		opts = &MonlOpts{}
		mock.ExpectQuery(query).
			WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).
				AddRow(1))
		count, err = monls.Count(ctx, opts)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), count)
	}
}

func testMonlsFind(ctx context.Context, monls *Monls, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = "SELECT (.+) FROM monls WHERE id = \\?"
			id    string
			monl  *model.Monl
			err   error
		)

		id = "monl-id-1"
		mock.ExpectQuery(query).
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows(monls.columns()).
				AddRow(id, "https://provider.com/url", nil, nil))
		monl, err = monls.Find(ctx, id)
		assert.Nil(t, err)
		if assert.NotNil(t, monl) {
			assert.Equal(t, id, monl.ID)
		}
	}
}

func testMonlsCreate(ctx context.Context, monls *Monls, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			monl *model.Monl
			err  error
		)

		monl = &model.Monl{}
		expectMonlsCreate(mock, monl)
		err = monls.Create(ctx, monl)
		assert.Nil(t, err)
		assert.NotEmpty(t, monl.ID)
		assert.NotEmpty(t, monl.CreatedAt)
	}
}

func expectMonlsCreate(mock sqlmock.Sqlmock, monl *model.Monl) {
	mock.ExpectExec("INSERT INTO monls").
		WithArgs(
			sqlmock.AnyArg(),
			monl.URL,
			sqlmock.AnyArg(),
			sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func testMonlsUpdate(ctx context.Context, monls *Monls, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			monl *model.Monl
			err  error
		)

		monl = &model.Monl{ID: "monl-id-1"}
		expectMonlsUpdate(mock, monl)
		err = monls.Update(ctx, monl)
		assert.Nil(t, err)
		assert.NotEmpty(t, monl.UpdatedAt)
	}
}

func expectMonlsUpdate(mock sqlmock.Sqlmock, monl *model.Monl) {
	mock.ExpectExec("UPDATE monls (.+) WHERE id = \\?").
		WithArgs(
			monl.URL,
			sqlmock.AnyArg(),
			monl.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func testMonlsDelete(ctx context.Context, monls *Monls, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = regexp.QuoteMeta("DELETE FROM monls WHERE id = ?")
			id    string
			n     int64
			err   error
		)

		id = "monl-id-1"
		mock.ExpectExec(query).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))
		n, err = monls.Delete(ctx, id)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), n)
	}
}
