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

func TestTags(t *testing.T) {
	db, mock, err := dbtest.New()
	assert.Nil(t, err)
	defer db.Close()

	ctx := context.TODO()
	s := NewStore(db)
	tags := NewTags(s)

	t.Run("list", testTagsList(ctx, tags, mock))
	t.Run("count", testTagsCount(ctx, tags, mock))
	t.Run("find", testTagsFind(ctx, tags, mock))
	t.Run("create", testTagsCreate(ctx, tags, mock))
	t.Run("update", testTagsUpdate(ctx, tags, mock))
	t.Run("delete", testTagsDelete(ctx, tags, mock))
}

func testTagsList(ctx context.Context, tags *Tags, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			prefix = "SELECT (.+) FROM tags"
			opts   *TagOpts
			list   []*model.Tag
			err    error
		)

		// Test nil opts.
		opts = nil
		mock.ExpectQuery(prefix).
			WillReturnRows(sqlmock.NewRows(tags.columns()).
				AddRow("tag-id-1", "name", "user-id-1", "", 0, "#colorhex", "#bgcolorhex", false, nil, nil))
		list, err = tags.List(ctx, opts)
		assert.Nil(t, err)
		assert.NotNil(t, list)

		// Test filter by user.
		opts = &TagOpts{UserID: "user-id-1"}
		mock.ExpectQuery(fmt.Sprintf("%s WHERE user_id = ?", prefix)).
			WithArgs(opts.UserID).
			WillReturnRows(sqlmock.NewRows(tags.columns()))
		_, err = tags.List(ctx, opts)
		assert.Nil(t, err)

		// Test filter by name.
		opts = &TagOpts{Name: "tag name"}
		mock.ExpectQuery(fmt.Sprintf("%s WHERE name = ?", prefix)).
			WithArgs(opts.Name).
			WillReturnRows(sqlmock.NewRows(tags.columns()))
		_, err = tags.List(ctx, opts)
		assert.Nil(t, err)

		// Test filter by name pattern.
		opts = &TagOpts{NamePattern: "%tag%"}
		mock.ExpectQuery(fmt.Sprintf("%s WHERE name LIKE ?", prefix)).
			WithArgs(opts.NamePattern).
			WillReturnRows(sqlmock.NewRows(tags.columns()))
		_, err = tags.List(ctx, opts)
		assert.Nil(t, err)

		// Test filter by parents.
		opts = &TagOpts{ParentIDs: []string{"parent-id-1", "parent-id-2"}}
		mock.ExpectQuery(fmt.Sprintf(regexp.QuoteMeta("%s WHERE parent_id IN (?,?)"), prefix)).
			WithArgs(opts.ParentIDs[0], opts.ParentIDs[1]).
			WillReturnRows(sqlmock.NewRows(tags.columns()))
		_, err = tags.List(ctx, opts)
		assert.Nil(t, err)

		// Test filter by level.
		opts = &TagOpts{Level: field.NewNullInt64(0)}
		mock.ExpectQuery(fmt.Sprintf(regexp.QuoteMeta("%s WHERE level = ?"), prefix)).
			WithArgs(opts.Level.Value()).
			WillReturnRows(sqlmock.NewRows(tags.columns()))
		_, err = tags.List(ctx, opts)
		assert.Nil(t, err)
	}
}

func testTagsCount(ctx context.Context, tags *Tags, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = regexp.QuoteMeta("SELECT count(*) FROM tags")
			opts  *TagOpts
			count int64
			err   error
		)

		opts = &TagOpts{}
		mock.ExpectQuery(query).
			WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).
				AddRow(1))
		count, err = tags.Count(ctx, opts)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), count)
	}
}

func testTagsFind(ctx context.Context, tags *Tags, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = "SELECT (.+) FROM tags WHERE id = \\? LIMIT 1"
			id    string
			tag   *model.Tag
			err   error
		)

		id = "tag-id-1"
		mock.ExpectQuery(query).
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows(tags.columns()).
				AddRow(id, "name", "user-id-1", "", 0, "#colorhex", "#bgcolorhex", false, nil, nil))
		tag, err = tags.Find(ctx, id)
		assert.Nil(t, err)
		if assert.NotNil(t, tag) {
			assert.Equal(t, id, tag.ID)
		}
	}
}

func testTagsCreate(ctx context.Context, tags *Tags, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			tag *model.Tag
			err error
		)

		tag = &model.Tag{}
		expectTagsCreate(mock, tag)
		err = tags.Create(ctx, tag)
		assert.Nil(t, err)
		assert.NotEmpty(t, tag.ID)
		assert.NotEmpty(t, tag.CreatedAt)
	}
}

func expectTagsCreate(mock sqlmock.Sqlmock, tag *model.Tag) {
	mock.ExpectExec("INSERT INTO tags").
		WithArgs(
			sqlmock.AnyArg(),
			tag.Name,
			tag.UserID,
			tag.ParentID,
			tag.Level,
			tag.Color,
			tag.BgColor,
			tag.HasChildren,
			sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func testTagsUpdate(ctx context.Context, tags *Tags, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			tag *model.Tag
			err error
		)

		tag = &model.Tag{ID: "tag-id-1"}
		expectTagsUpdate(mock, tag)
		err = tags.Update(ctx, tag)
		assert.Nil(t, err)
		assert.NotEmpty(t, tag.UpdatedAt)
	}
}

func expectTagsUpdate(mock sqlmock.Sqlmock, tag *model.Tag) {
	mock.ExpectExec("UPDATE tags (.+) WHERE id = \\? LIMIT 1").
		WithArgs(
			tag.Name,
			tag.UserID,
			tag.ParentID,
			tag.Level,
			tag.Color,
			tag.BgColor,
			tag.HasChildren,
			sqlmock.AnyArg(),
			tag.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func testTagsDelete(ctx context.Context, tags *Tags, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			id  string
			n   int64
			err error
		)

		id = "tag-id-1"
		mock.ExpectExec("DELETE FROM tags WHERE id = \\? LIMIT 1").
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))
		n, err = tags.Delete(ctx, id)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), n)
	}
}
