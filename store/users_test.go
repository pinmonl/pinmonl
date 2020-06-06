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

func TestUsers(t *testing.T) {
	db, mock, err := dbtest.New()
	assert.Nil(t, err)
	defer db.Close()

	ctx := context.TODO()
	s := NewStore(db)
	users := NewUsers(s)

	t.Run("list", testUsersList(ctx, users, mock))
	t.Run("count", testUsersCount(ctx, users, mock))
	t.Run("find", testUsersFind(ctx, users, mock))
	t.Run("findLogin", testUsersFindLogin(ctx, users, mock))
	t.Run("create", testUsersCreate(ctx, users, mock))
	t.Run("update", testUsersUpdate(ctx, users, mock))
	t.Run("delete", testUsersDelete(ctx, users, mock))
}

func testUsersList(ctx context.Context, users *Users, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			prefix = "SELECT (.+) FROM users"
			opts   *UserOpts
			list   []*model.User
			err    error
		)

		// Test nil opts.
		opts = nil
		mock.ExpectQuery(prefix).
			WillReturnRows(sqlmock.NewRows(users.columns()).
				AddRow("user-id-1", "user1", "pw", "user name 1", "", "hash", model.NormalUser, model.ActiveUser, nil, timestamp(), timestamp()))
		list, err = users.List(ctx, opts)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(list))

		// Test filter by login.
		opts = &UserOpts{Login: "user1"}
		mock.ExpectQuery(fmt.Sprintf("%s WHERE login = ?", prefix)).
			WithArgs(opts.Login).
			WillReturnRows(sqlmock.NewRows(users.columns()))
		_, err = users.List(ctx, opts)
		assert.Nil(t, err)

		// Test filter by status.
		opts = &UserOpts{Status: field.NewNullValue(model.ActiveUser)}
		mock.ExpectQuery(fmt.Sprintf("%s WHERE status = ?", prefix)).
			WithArgs(opts.Status.Value()).
			WillReturnRows(sqlmock.NewRows(users.columns()))
		_, err = users.List(ctx, opts)
		assert.Nil(t, err)

		// Test filter by status.
		opts = &UserOpts{Role: field.NewNullValue(model.NormalUser)}
		mock.ExpectQuery(fmt.Sprintf("%s WHERE role = ?", prefix)).
			WithArgs(opts.Role.Value()).
			WillReturnRows(sqlmock.NewRows(users.columns()))
		_, err = users.List(ctx, opts)
		assert.Nil(t, err)

		// Test pagination.
		opts = &UserOpts{}
		opts.Limit = 10
		opts.Offset = 1
		mock.ExpectQuery(fmt.Sprintf("%s LIMIT %d OFFSET %d", prefix, opts.Limit, opts.Offset)).
			WillReturnRows(sqlmock.NewRows(users.columns()))
		_, err = users.List(ctx, opts)
		assert.Nil(t, err)
	}
}

func testUsersCount(ctx context.Context, users *Users, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = regexp.QuoteMeta("SELECT count(*) FROM users")
			opts  *UserOpts
			count int64
			err   error
		)

		opts = &UserOpts{}
		mock.ExpectQuery(query).
			WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).
				AddRow(1))
		count, err = users.Count(ctx, opts)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), count)
	}
}

func testUsersFind(ctx context.Context, users *Users, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = "SELECT (.+) FROM users WHERE id = \\?"
			id    string
			user  *model.User
			err   error
		)

		id = "user-id-1"
		mock.ExpectQuery(query).
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows(users.columns()).
				AddRow(id, "login", "pw", "name", "", "hash", model.NormalUser, model.ActiveUser, nil, nil, nil))
		user, err = users.Find(ctx, id)
		assert.Nil(t, err)
		if assert.NotNil(t, user) {
			assert.Equal(t, id, user.ID)
		}
	}
}

func testUsersFindLogin(ctx context.Context, users *Users, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = "SELECT (.+) FROM users WHERE login = \\? AND status = \\?"
			login string
			user  *model.User
			err   error
		)

		login = "user1"
		mock.ExpectQuery(query).
			WithArgs(login, model.ActiveUser).
			WillReturnRows(sqlmock.NewRows(users.columns()).
				AddRow("user-id-1", login, "pw", "name", "", "hash", model.NormalUser, model.ActiveUser, nil, nil, nil))
		user, err = users.FindLogin(ctx, login)
		assert.Nil(t, err)
		if assert.NotNil(t, user) {
			assert.Equal(t, login, user.Login)
		}
	}
}

func testUsersCreate(ctx context.Context, users *Users, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			user *model.User
			err  error
		)

		user = &model.User{}
		expectUsersInsert(mock, user)
		err = users.Create(ctx, user)
		assert.Nil(t, err)
		assert.NotEmpty(t, user.ID)
		assert.NotEmpty(t, user.CreatedAt)
	}
}

func expectUsersInsert(mock sqlmock.Sqlmock, user *model.User) {
	mock.ExpectExec("INSERT INTO users").
		WithArgs(
			sqlmock.AnyArg(),
			user.Login,
			user.Password,
			user.Name,
			user.ImageID,
			user.Hash,
			user.Role,
			user.Status,
			user.LastSeen,
			sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func testUsersUpdate(ctx context.Context, users *Users, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			user *model.User
			err  error
		)

		user = &model.User{ID: "user-id-1", Name: "new name"}
		expectUsersUpdate(mock, user)
		err = users.Update(ctx, user)
		assert.Nil(t, err)
		assert.NotEmpty(t, user.UpdatedAt)
	}
}

func expectUsersUpdate(mock sqlmock.Sqlmock, user *model.User) {
	mock.ExpectExec("UPDATE users (.+) WHERE id = \\?").
		WithArgs(
			user.Login,
			user.Password,
			user.Name,
			user.ImageID,
			user.Hash,
			user.Role,
			user.Status,
			user.LastSeen,
			sqlmock.AnyArg(),
			user.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func testUsersDelete(ctx context.Context, users *Users, mock sqlmock.Sqlmock) func(*testing.T) {
	return func(t *testing.T) {
		var (
			query = "DELETE FROM users WHERE id = \\?"
			id    string
			n     int64
			err   error
		)

		id = "user-id-1"
		mock.ExpectExec(query).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))
		n, err = users.Delete(ctx, id)
		assert.Nil(t, err)
		assert.Equal(t, int64(1), n)
	}
}
