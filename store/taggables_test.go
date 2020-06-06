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

func TestTaggables(t *testing.T) {
	db, mock, err := dbtest.New()
	assert.Nil(t, err)
	defer db.Close()

	ctx := context.TODO()
	s := NewStore(db)
	taggables := NewTaggables(s)

	t.Run("list", testTaggablesList(ctx, taggables, mock))
	t.Run("count", testTaggablesCount(ctx, taggables, mock))
	t.Run("find", testTaggablesFind(ctx, taggables, mock))
	t.Run("listWithTags", testTaggablesListWithTags(ctx, taggables, mock))
	t.Run("create", testTaggablesCreate(ctx, taggables, mock))
	t.Run("update", testTaggablesUpdate(ctx, taggables, mock))
	t.Run("delete", testTaggablesDelete(ctx, taggables, mock))
}

func testTaggablesList(ctx context.Context, taggables *Taggables, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			prefix = "SELECT (.+) FROM taggables"
			opts   *TaggableOpts
			list   []*model.Taggable
			err    error
		)

		// Test nil opts.
		opts = nil
		mock.ExpectQuery(prefix).
			WillReturnRows(sqlmock.NewRows(taggables.columns()).
				AddRow("taggable-id-1", "tag-id-1", "target-id-1", "target"))
		list, err = taggables.List(ctx, opts)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(list))

		// Test filter by tags.
		opts = &TaggableOpts{TagIDs: []string{"tag-id-1", "tag-id-2", "tag-id-3"}}
		mock.ExpectQuery(fmt.Sprintf(regexp.QuoteMeta("%s WHERE tag_id IN (?,?,?)"), prefix)).
			WithArgs(opts.TagIDs[0], opts.TagIDs[1], opts.TagIDs[2]).
			WillReturnRows(sqlmock.NewRows(taggables.columns()))
		_, err = taggables.List(ctx, opts)
		assert.Nil(t, err)

		// Test filter by targets.
		opts = &TaggableOpts{Targets: model.MorphableList{
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
			WillReturnRows(sqlmock.NewRows(taggables.columns()))
		_, err = taggables.List(ctx, opts)
		assert.Nil(t, err)
	}
}

func testTaggablesCount(ctx context.Context, taggables *Taggables, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = regexp.QuoteMeta("SELECT count(*) FROM taggables")
			opts  *TaggableOpts
			count int64
			err   error
		)

		opts = &TaggableOpts{}
		mock.ExpectQuery(query).
			WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).
				AddRow(1))
		count, err = taggables.Count(ctx, opts)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), count)
	}
}

func testTaggablesFind(ctx context.Context, taggables *Taggables, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query    = "SELECT (.+) FROM taggables WHERE id = \\?"
			id       string
			taggable *model.Taggable
			err      error
		)

		id = "taggable-id-1"
		mock.ExpectQuery(query).
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows(taggables.columns()).
				AddRow("taggable-id-1", "tag-id-1", "target-id-1", "target"))
		taggable, err = taggables.Find(ctx, id)
		assert.Nil(t, err)
		if assert.NotNil(t, taggable) {
			assert.Equal(t, id, taggable.ID)
		}
	}
}

func testTaggablesListWithTags(ctx context.Context, taggables *Taggables, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			prefix = "SELECT (.+) FROM taggables LEFT JOIN tags ON tags.id = taggables.tag_id"
			opts   *TaggableOpts
			list   []*model.Taggable
			err    error
		)

		opts = &TaggableOpts{Targets: model.MorphableList{
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
			WillReturnRows(sqlmock.NewRows(append((&Tags{}).columns(), "target_id")))
		list, err = taggables.ListWithTags(ctx, opts)
		assert.Nil(t, err)
		assert.Equal(t, 0, len(list))
	}
}

func testTaggablesCreate(ctx context.Context, taggables *Taggables, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			taggable *model.Taggable
			err      error
		)

		taggable = &model.Taggable{}
		expectTaggablesCreate(mock, taggable)
		err = taggables.Create(ctx, taggable)
		assert.Nil(t, err)
		assert.NotEmpty(t, taggable.ID)
	}
}

func expectTaggablesCreate(mock sqlmock.Sqlmock, taggable *model.Taggable) {
	mock.ExpectExec("INSERT INTO taggables").
		WithArgs(
			sqlmock.AnyArg(),
			taggable.TagID,
			taggable.TaggableID,
			taggable.TaggableName).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func testTaggablesUpdate(ctx context.Context, taggables *Taggables, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			taggable *model.Taggable
			err      error
		)

		taggable = &model.Taggable{ID: "taggable-id-1"}
		expectTaggablesUpdate(mock, taggable)
		err = taggables.Update(ctx, taggable)
		assert.Nil(t, err)
	}
}

func expectTaggablesUpdate(mock sqlmock.Sqlmock, taggable *model.Taggable) {
	mock.ExpectExec("UPDATE taggables (.+) WHERE id = \\?").
		WithArgs(
			taggable.TagID,
			taggable.TaggableID,
			taggable.TaggableName,
			taggable.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func testTaggablesDelete(ctx context.Context, taggables *Taggables, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = regexp.QuoteMeta("DELETE FROM taggables WHERE id = ?")
			id    string
			n     int64
			err   error
		)

		id = "taggable-id-1"
		mock.ExpectExec(query).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))
		n, err = taggables.Delete(ctx, id)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), n)
	}
}
