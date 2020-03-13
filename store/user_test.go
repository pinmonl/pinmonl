package store

import (
	"context"
	"testing"

	"github.com/pinmonl/pinmonl/database/dbtest"
	"github.com/pinmonl/pinmonl/model"
	"github.com/stretchr/testify/assert"
)

func TestUserStore(t *testing.T) {
	db, err := dbtest.Open()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		dbtest.Reset(db)
		dbtest.Close(db)
	}()

	store := NewStore(db)
	users := NewUserStore(store)
	ctx := context.TODO()
	t.Run("Create", testUserCreate(ctx, users))
	t.Run("List", testUserList(ctx, users))
	t.Run("Find", testUserFind(ctx, users))
	t.Run("FindLogin", testUserFindLogin(ctx, users))
	t.Run("Update", testUserUpdate(ctx, users))
	t.Run("Delete", testUserDelete(ctx, users))
}

func testUserCreate(ctx context.Context, users UserStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData := []model.User{
			{Login: "user1", Name: "User 1"},
			{Login: "user2", Name: "User 2"},
			{Login: "user3", Name: "User 3"},
		}

		for _, user := range testData {
			assert.Nil(t, users.Create(ctx, &user))
			assert.NotEmpty(t, user.ID)
			assert.False(t, user.CreatedAt.Time().IsZero())
		}
	}
}

func testUserList(ctx context.Context, users UserStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, err := users.List(ctx, nil)
		assert.Nil(t, err)
		assert.Equal(t, 3, len(testData))

		want := testData[0:1]
		got, err := users.List(ctx, &UserOpts{Login: want[0].Login})
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	}
}

func testUserFind(ctx context.Context, users UserStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, _ := users.List(ctx, nil)

		want := testData[1]
		got := model.User{ID: want.ID}
		assert.Nil(t, users.Find(ctx, &got))
		assert.Equal(t, want, got)
	}
}

func testUserFindLogin(ctx context.Context, users UserStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, _ := users.List(ctx, nil)

		want := testData[0]
		got := model.User{Login: want.Login}
		assert.Nil(t, users.FindLogin(ctx, &got))
		assert.Equal(t, want, got)
	}
}

func testUserUpdate(ctx context.Context, users UserStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, _ := users.List(ctx, nil)

		want := testData[2]
		want.Name = "user3 new"
		assert.Nil(t, users.Update(ctx, &want))
		assert.False(t, want.UpdatedAt.Time().IsZero())

		got := model.User{ID: want.ID}
		users.Find(ctx, &got)
		assert.Equal(t, want, got)
	}
}

func testUserDelete(ctx context.Context, users UserStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, _ := users.List(ctx, nil)

		assert.Nil(t, users.Delete(ctx, &testData[0]))

		testData, _ = users.List(ctx, nil)
		assert.Equal(t, 2, len(testData))
	}
}
