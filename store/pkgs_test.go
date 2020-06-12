package store

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pinmonl/pinmonl/database/dbtest"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pkgs/pkguri"
	"github.com/stretchr/testify/assert"
)

func TestPkgs(t *testing.T) {
	db, mock, err := dbtest.New()
	assert.Nil(t, err)
	defer db.Close()

	ctx := context.TODO()
	s := NewStore(db)
	pkgs := NewPkgs(s)

	t.Run("list", testPkgsList(ctx, pkgs, mock))
	t.Run("count", testPkgsCount(ctx, pkgs, mock))
	t.Run("find", testPkgsFind(ctx, pkgs, mock))
	t.Run("create", testPkgsCreate(ctx, pkgs, mock))
	t.Run("update", testPkgsUpdate(ctx, pkgs, mock))
	t.Run("delete", testPkgsDelete(ctx, pkgs, mock))
}

func testPkgsList(ctx context.Context, pkgs *Pkgs, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			prefix = "SELECT (.+) FROM pkgs"
			opts   *PkgOpts
			list   []*model.Pkg
			err    error
		)

		// Test nil opts.
		opts = nil
		mock.ExpectQuery(prefix).
			WillReturnRows(sqlmock.NewRows(pkgs.columns()).
				AddRow("pkg-id-1", "https://provider.com/url", "provider", "", "owner/repo", nil, nil))
		list, err = pkgs.List(ctx, opts)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(list))

		// Test filter by provider.
		opts = &PkgOpts{Provider: "provider"}
		mock.ExpectQuery(fmt.Sprintf("%s WHERE provider = ?", prefix)).
			WithArgs(opts.Provider).
			WillReturnRows(sqlmock.NewRows(pkgs.columns()))
		_, err = pkgs.List(ctx, opts)
		assert.Nil(t, err)

		// Test filter by provider host.
		opts = &PkgOpts{ProviderHost: "provider.com"}
		mock.ExpectQuery(fmt.Sprintf("%s WHERE provider_host = ?", prefix)).
			WithArgs(opts.ProviderHost).
			WillReturnRows(sqlmock.NewRows(pkgs.columns()))
		_, err = pkgs.List(ctx, opts)
		assert.Nil(t, err)

		// Test filter by provider uri.
		opts = &PkgOpts{ProviderURI: "uri"}
		mock.ExpectQuery(fmt.Sprintf("%s WHERE provider_uri = ?", prefix)).
			WithArgs(opts.ProviderURI).
			WillReturnRows(sqlmock.NewRows(pkgs.columns()))
		_, err = pkgs.List(ctx, opts)
		assert.Nil(t, err)
	}
}

func testPkgsCount(ctx context.Context, pkgs *Pkgs, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = regexp.QuoteMeta("SELECT count(*) FROM pkgs")
			opts  *PkgOpts
			count int64
			err   error
		)

		opts = &PkgOpts{}
		mock.ExpectQuery(query).
			WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).
				AddRow(1))
		count, err = pkgs.Count(ctx, opts)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), count)
	}
}

func testPkgsFind(ctx context.Context, pkgs *Pkgs, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = "SELECT (.+) FROM pkgs WHERE id = \\?"
			id    string
			pkg   *model.Pkg
			err   error
		)

		id = "pkg-id-1"
		mock.ExpectQuery(query).
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows(pkgs.columns()).
				AddRow(id, "https://provider.com/url", "provider", "", "owner/repo", nil, nil))
		pkg, err = pkgs.Find(ctx, id)
		assert.Nil(t, err)
		if assert.NotNil(t, pkg) {
			assert.Equal(t, id, pkg.ID)
		}
	}
}

func testPkgsFindURI(ctx context.Context, pkgs *Pkgs, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query   = "SELECT (.+) FROM pkgs WHERE provider = \\? AND provider_host = \\? AND provider_uri = \\?"
			prd     string
			prdHost string
			prdURI  string
			pkg     *model.Pkg
			err     error
		)

		prd, prdHost, prdURI = "provider", "provider.com", "owner/repo"
		mock.ExpectQuery(query).
			WithArgs(prd, prdHost, prdURI).
			WillReturnRows(sqlmock.NewRows(pkgs.columns()).
				AddRow("pkg-id-1", "https://provider.com/url", prd, prdHost, prdURI, nil, nil))
		pkg, err = pkgs.FindURI(ctx, &pkguri.PkgURI{
			Provider: prd,
			Host:     prdHost,
			URI:      prdURI,
		})
		assert.Nil(t, err)
		if assert.NotNil(t, pkg) {
			assert.Equal(t, prd, pkg.Provider)
			assert.Equal(t, prdHost, pkg.ProviderHost)
			assert.Equal(t, prdURI, pkg.ProviderURI)
		}
	}
}

func testPkgsCreate(ctx context.Context, pkgs *Pkgs, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			pkg *model.Pkg
			err error
		)

		pkg = &model.Pkg{}
		expectPkgsCreate(mock, pkg)
		err = pkgs.Create(ctx, pkg)
		assert.Nil(t, err)
		assert.NotEmpty(t, pkg.ID)
		assert.NotEmpty(t, pkg.CreatedAt)
	}
}

func expectPkgsCreate(mock sqlmock.Sqlmock, pkg *model.Pkg) {
	mock.ExpectExec("INSERT INTO pkgs").
		WithArgs(
			sqlmock.AnyArg(),
			pkg.URL,
			pkg.Provider,
			pkg.ProviderHost,
			pkg.ProviderURI,
			sqlmock.AnyArg(),
			sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func testPkgsUpdate(ctx context.Context, pkgs *Pkgs, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			pkg *model.Pkg
			err error
		)

		pkg = &model.Pkg{ID: "pkg-id-1"}
		expectPkgsUpdate(mock, pkg)
		err = pkgs.Update(ctx, pkg)
		assert.Nil(t, err)
		assert.NotEmpty(t, pkg.UpdatedAt)
	}
}

func expectPkgsUpdate(mock sqlmock.Sqlmock, pkg *model.Pkg) {
	mock.ExpectExec("UPDATE pkgs (.+) WHERE id = \\?").
		WithArgs(
			pkg.URL,
			pkg.Provider,
			pkg.ProviderHost,
			pkg.ProviderURI,
			sqlmock.AnyArg(),
			pkg.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func testPkgsDelete(ctx context.Context, pkgs *Pkgs, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = regexp.QuoteMeta("DELETE FROM pkgs WHERE id = ?")
			id    string
			n     int64
			err   error
		)

		id = "pkg-id-1"
		mock.ExpectExec(query).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))
		n, err = pkgs.Delete(ctx, id)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), n)
	}
}
