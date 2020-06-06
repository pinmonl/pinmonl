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

func TestPinpkgs(t *testing.T) {
	db, mock, err := dbtest.New()
	assert.Nil(t, err)
	defer db.Close()

	ctx := context.TODO()
	s := NewStore(db)
	pinpkgs := NewPinpkgs(s)

	t.Run("list", testPinpkgsList(ctx, pinpkgs, mock))
	t.Run("count", testPinpkgsCount(ctx, pinpkgs, mock))
	t.Run("find", testPinpkgsFind(ctx, pinpkgs, mock))
	t.Run("create", testPinpkgsCreate(ctx, pinpkgs, mock))
	t.Run("update", testPinpkgsUpdate(ctx, pinpkgs, mock))
	t.Run("delete", testPinpkgsDelete(ctx, pinpkgs, mock))
}

func testPinpkgsList(ctx context.Context, pinpkgs *Pinpkgs, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			prefix = "SELECT (.+) FROM pinpkgs"
			opts   *PinpkgOpts
			list   []*model.Pinpkg
			err    error
		)

		// Test nil opts.
		opts = nil
		mock.ExpectQuery(prefix).
			WillReturnRows(sqlmock.NewRows(pinpkgs.columns()).
				AddRow("pinpkg-id-1", "pinl-id-1", "pkg-id-1"))
		list, err = pinpkgs.List(ctx, opts)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(list))

		// Test filter by pinls.
		opts = &PinpkgOpts{PinlIDs: []string{"pinl-id-1"}}
		mock.ExpectQuery(fmt.Sprintf(regexp.QuoteMeta("%s WHERE pinl_id IN (?)"), prefix)).
			WithArgs(opts.PinlIDs[0]).
			WillReturnRows(sqlmock.NewRows(pinpkgs.columns()))
		_, err = pinpkgs.List(ctx, opts)
		assert.Nil(t, err)

		// Test filter by pkgs.
		opts = &PinpkgOpts{PkgIDs: []string{"pkg-id-1", "pkg-id-2"}}
		mock.ExpectQuery(fmt.Sprintf(regexp.QuoteMeta("%s WHERE pkg_id IN (?,?)"), prefix)).
			WithArgs(opts.PkgIDs[0], opts.PkgIDs[1]).
			WillReturnRows(sqlmock.NewRows(pinpkgs.columns()))
		_, err = pinpkgs.List(ctx, opts)
		assert.Nil(t, err)
	}
}

func testPinpkgsCount(ctx context.Context, pinpkgs *Pinpkgs, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = regexp.QuoteMeta("SELECT count(*) FROM pinpkgs")
			opts  *PinpkgOpts
			count int64
			err   error
		)

		opts = &PinpkgOpts{}
		mock.ExpectQuery(query).
			WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).
				AddRow(1))
		count, err = pinpkgs.Count(ctx, opts)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), count)
	}
}

func testPinpkgsFind(ctx context.Context, pinpkgs *Pinpkgs, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query  = "SELECT (.+) FROM pinpkgs WHERE id = \\?"
			id     string
			pinpkg *model.Pinpkg
			err    error
		)

		id = "pinpkg-id-1"
		mock.ExpectQuery(query).
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows(pinpkgs.columns()).
				AddRow(id, "pinl-id-1", "pkg-id-1"))
		pinpkg, err = pinpkgs.Find(ctx, id)
		assert.Nil(t, err)
		if assert.NotNil(t, pinpkg) {
			assert.Equal(t, id, pinpkg.ID)
		}
	}
}

func testPinpkgsCreate(ctx context.Context, pinpkgs *Pinpkgs, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			pinpkg *model.Pinpkg
			err    error
		)

		pinpkg = &model.Pinpkg{}
		expectPinpkgsCreate(mock, pinpkg)
		err = pinpkgs.Create(ctx, pinpkg)
		assert.Nil(t, err)
		assert.NotEmpty(t, pinpkg.ID)
	}
}

func expectPinpkgsCreate(mock sqlmock.Sqlmock, pinpkg *model.Pinpkg) {
	mock.ExpectExec("INSERT INTO pinpkgs").
		WithArgs(
			sqlmock.AnyArg(),
			pinpkg.PinlID,
			pinpkg.PkgID).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func testPinpkgsUpdate(ctx context.Context, pinpkgs *Pinpkgs, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			pinpkg *model.Pinpkg
			err    error
		)

		pinpkg = &model.Pinpkg{ID: "pinpkg-id-1"}
		expectPinpkgsUpdate(mock, pinpkg)
		err = pinpkgs.Update(ctx, pinpkg)
		assert.Nil(t, err)
	}
}

func expectPinpkgsUpdate(mock sqlmock.Sqlmock, pinpkg *model.Pinpkg) {
	mock.ExpectExec("UPDATE pinpkgs (.+) WHERE id = \\?").
		WithArgs(
			pinpkg.PinlID,
			pinpkg.PkgID,
			pinpkg.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func testPinpkgsDelete(ctx context.Context, pinpkgs *Pinpkgs, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = regexp.QuoteMeta("DELETE FROM pinpkgs WHERE id = ?")
			id    string
			n     int64
			err   error
		)

		id = "pinpkg-id-1"
		mock.ExpectExec(query).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))
		n, err = pinpkgs.Delete(ctx, id)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), n)
	}
}
