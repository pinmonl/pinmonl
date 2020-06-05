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

func TestStats(t *testing.T) {
	db, mock, err := dbtest.New()
	assert.Nil(t, err)
	defer db.Close()

	ctx := context.TODO()
	s := NewStore(db)
	stats := NewStats(s)

	t.Run("list", testStatsList(ctx, stats, mock))
	t.Run("count", testStatsCount(ctx, stats, mock))
	t.Run("find", testStatsFind(ctx, stats, mock))
	t.Run("create", testStatsCreate(ctx, stats, mock))
	t.Run("update", testStatsUpdate(ctx, stats, mock))
	t.Run("delete", testStatsDelete(ctx, stats, mock))
}

func testStatsList(ctx context.Context, stats *Stats, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			prefix = "SELECT (.+) FROM stats"
			opts   *StatOpts
			list   []*model.Stat
			err    error
		)

		// Test nil opts.
		opts = nil
		mock.ExpectQuery(prefix).
			WillReturnRows(sqlmock.NewRows(stats.columns()).
				AddRow("stat-id-1", "pkg-id-1", "", nil, model.TagStat, "v0.1.0", model.StringStat, "checksum", 0, true, false))
		list, err = stats.List(ctx, opts)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(list))

		// Test filter by pkgs.
		opts = &StatOpts{PkgIDs: []string{"pkg-id-1", "pkg-id-2"}}
		mock.ExpectQuery(fmt.Sprintf(regexp.QuoteMeta("%s WHERE pkg_id IN (?,?)"), prefix)).
			WithArgs(opts.PkgIDs[0], opts.PkgIDs[1]).
			WillReturnRows(sqlmock.NewRows(stats.columns()))
		_, err = stats.List(ctx, opts)
		assert.Nil(t, err)

		// Test filter by parents.
		opts = &StatOpts{ParentIDs: []string{"stat-id-1", "stat-id-2"}}
		mock.ExpectQuery(fmt.Sprintf(regexp.QuoteMeta("%s WHERE parent_id IN (?,?)"), prefix)).
			WithArgs(opts.ParentIDs[0], opts.ParentIDs[1]).
			WillReturnRows(sqlmock.NewRows(stats.columns()))
		_, err = stats.List(ctx, opts)
		assert.Nil(t, err)

		// Test filter by kind.
		opts = &StatOpts{Kind: field.NewNullValue(model.TagStat)}
		mock.ExpectQuery(fmt.Sprintf("%s WHERE kind = ?", prefix)).
			WithArgs(opts.Kind.Value()).
			WillReturnRows(sqlmock.NewRows(stats.columns()))
		_, err = stats.List(ctx, opts)
		assert.Nil(t, err)

		// Test filter by latest.
		opts = &StatOpts{IsLatest: field.NewNullBool(true)}
		mock.ExpectQuery(fmt.Sprintf("%s WHERE is_latest = ?", prefix)).
			WithArgs(opts.IsLatest.Value()).
			WillReturnRows(sqlmock.NewRows(stats.columns()))
		_, err = stats.List(ctx, opts)
		assert.Nil(t, err)
	}
}

func testStatsCount(ctx context.Context, stats *Stats, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = regexp.QuoteMeta("SELECT count(*) FROM stats")
			opts  *StatOpts
			count int64
			err   error
		)

		opts = &StatOpts{}
		mock.ExpectQuery(query).
			WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).
				AddRow(1))
		count, err = stats.Count(ctx, opts)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), count)
	}
}

func testStatsFind(ctx context.Context, stats *Stats, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = "SELECT (.+) FROM stats WHERE id = \\? LIMIT 1"
			id    string
			stat  *model.Stat
			err   error
		)

		id = "stat-id-1"
		mock.ExpectQuery(query).
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows(stats.columns()).
				AddRow(id, "pkg-id-1", "", nil, model.TagStat, "v0.1.0", model.StringStat, "checksum", 0, true, false))
		stat, err = stats.Find(ctx, id)
		assert.Nil(t, err)
		if assert.NotNil(t, stat) {
			assert.Equal(t, id, stat.ID)
		}
	}
}

func testStatsCreate(ctx context.Context, stats *Stats, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			stat *model.Stat
			err  error
		)

		stat = &model.Stat{}
		expectStatsCreate(mock, stat)
		err = stats.Create(ctx, stat)
		assert.Nil(t, err)
		assert.NotEmpty(t, stat.ID)
	}
}

func expectStatsCreate(mock sqlmock.Sqlmock, stat *model.Stat) {
	mock.ExpectExec("INSERT INTO stats").
		WithArgs(
			sqlmock.AnyArg(),
			stat.PkgID,
			stat.ParentID,
			stat.RecordedAt,
			stat.Kind,
			stat.Value,
			stat.ValueType,
			stat.Checksum,
			stat.Weight,
			stat.IsLatest,
			stat.HasChildren).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func testStatsUpdate(ctx context.Context, stats *Stats, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			stat *model.Stat
			err  error
		)

		stat = &model.Stat{ID: "stat-id-1"}
		expectStatsUpdate(mock, stat)
		err = stats.Update(ctx, stat)
		assert.Nil(t, err)
	}
}

func expectStatsUpdate(mock sqlmock.Sqlmock, stat *model.Stat) {
	mock.ExpectExec("UPDATE stats (.+) WHERE id = \\? LIMIT 1").
		WithArgs(
			stat.PkgID,
			stat.ParentID,
			stat.RecordedAt,
			stat.Kind,
			stat.Value,
			stat.ValueType,
			stat.Checksum,
			stat.Weight,
			stat.IsLatest,
			stat.HasChildren,
			stat.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func testStatsDelete(ctx context.Context, stats *Stats, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = regexp.QuoteMeta("DELETE FROM stats WHERE id = ? LIMIT 1")
			id    string
			n     int64
			err   error
		)

		id = "stat-id-1"
		mock.ExpectExec(query).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))
		n, err = stats.Delete(ctx, id)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), n)
	}
}
