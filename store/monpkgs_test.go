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

func TestMonpkgs(t *testing.T) {
	db, mock, err := dbtest.New()
	assert.Nil(t, err)
	defer db.Close()

	ctx := context.TODO()
	s := NewStore(db)
	monpkgs := NewMonpkgs(s)

	t.Run("list", testMonpkgsList(ctx, monpkgs, mock))
	t.Run("count", testMonpkgsCount(ctx, monpkgs, mock))
	t.Run("find", testMonpkgsFind(ctx, monpkgs, mock))
	t.Run("create", testMonpkgsCreate(ctx, monpkgs, mock))
	t.Run("update", testMonpkgsUpdate(ctx, monpkgs, mock))
	t.Run("delete", testMonpkgsDelete(ctx, monpkgs, mock))
}

func testMonpkgsList(ctx context.Context, monpkgs *Monpkgs, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			prefix = "SELECT (.+) FROM monpkgs"
			opts   *MonpkgOpts
			list   []*model.Monpkg
			err    error
		)

		// Test nil opts.
		opts = nil
		mock.ExpectQuery(prefix).
			WillReturnRows(sqlmock.NewRows(monpkgs.columns()).
				AddRow("monpkg-id-1", "monl-id-1", "pkg-id-1", 0))
		list, err = monpkgs.List(ctx, opts)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(list))

		// Test filter by monls.
		opts = &MonpkgOpts{MonlIDs: []string{"monl-id-1", "monl-id-2"}}
		mock.ExpectQuery(fmt.Sprintf(regexp.QuoteMeta("%s WHERE monl_id IN (?,?)"), prefix)).
			WithArgs(opts.MonlIDs[0], opts.MonlIDs[1]).
			WillReturnRows(sqlmock.NewRows(monpkgs.columns()))
		_, err = monpkgs.List(ctx, opts)
		assert.Nil(t, err)

		// Test filter by pkgs.
		opts = &MonpkgOpts{PkgIDs: []string{"pkg-id-1", "pkg-id-2", "pkg-id-3"}}
		mock.ExpectQuery(fmt.Sprintf(regexp.QuoteMeta("%s WHERE pkg_id IN (?,?,?)"), prefix)).
			WithArgs(opts.PkgIDs[0], opts.PkgIDs[1], opts.PkgIDs[2]).
			WillReturnRows(sqlmock.NewRows(monpkgs.columns()))
		_, err = monpkgs.List(ctx, opts)
		assert.Nil(t, err)
	}
}

func testMonpkgsCount(ctx context.Context, monpkgs *Monpkgs, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = regexp.QuoteMeta("SELECT count(*) FROM monpkgs")
			opts  *MonpkgOpts
			count int64
			err   error
		)

		opts = &MonpkgOpts{}
		mock.ExpectQuery(query).
			WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).
				AddRow(1))
		count, err = monpkgs.Count(ctx, opts)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), count)
	}
}

func testMonpkgsFind(ctx context.Context, monpkgs *Monpkgs, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query  = "SELECT (.+) FROM monpkgs WHERE id = \\?"
			id     string
			monpkg *model.Monpkg
			err    error
		)

		id = "monpkg-id-1"
		mock.ExpectQuery(query).
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows(monpkgs.columns()).
				AddRow(id, "monl-id-1", "pkg-id-1", 0))
		monpkg, err = monpkgs.Find(ctx, id)
		assert.Nil(t, err)
		if assert.NotNil(t, monpkg) {
			assert.Equal(t, id, monpkg.ID)
		}
	}
}

func testMonpkgsCreate(ctx context.Context, monpkgs *Monpkgs, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			monpkg *model.Monpkg
			err    error
		)

		monpkg = &model.Monpkg{}
		expectMonpkgsCreate(mock, monpkg)
		err = monpkgs.Create(ctx, monpkg)
		assert.Nil(t, err)
		assert.NotEmpty(t, monpkg.ID)
	}
}

func expectMonpkgsCreate(mock sqlmock.Sqlmock, monpkg *model.Monpkg) {
	mock.ExpectExec("INSERT INTO monpkgs").
		WithArgs(
			sqlmock.AnyArg(),
			monpkg.MonlID,
			monpkg.PkgID,
			monpkg.Kind).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func testMonpkgsUpdate(ctx context.Context, monpkgs *Monpkgs, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			monpkg *model.Monpkg
			err    error
		)

		monpkg = &model.Monpkg{ID: "monpkg-id-1"}
		expectMonpkgsUpdate(mock, monpkg)
		err = monpkgs.Update(ctx, monpkg)
		assert.Nil(t, err)
	}
}

func expectMonpkgsUpdate(mock sqlmock.Sqlmock, monpkg *model.Monpkg) {
	mock.ExpectExec("UPDATE monpkgs (.+) WHERE id = \\?").
		WithArgs(
			monpkg.MonlID,
			monpkg.PkgID,
			monpkg.Kind,
			monpkg.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func testMonpkgsDelete(ctx context.Context, monpkgs *Monpkgs, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = regexp.QuoteMeta("DELETE FROM monpkgs WHERE id = ?")
			id    string
			n     int64
			err   error
		)

		id = "monpkg-id-1"
		mock.ExpectExec(query).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))
		n, err = monpkgs.Delete(ctx, id)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), n)
	}
}
