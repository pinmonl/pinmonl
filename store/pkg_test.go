package store

import (
	"context"
	"testing"

	"github.com/pinmonl/pinmonl/database/dbtest"
	"github.com/pinmonl/pinmonl/model"
	"github.com/stretchr/testify/assert"
)

func TestPkgStore(t *testing.T) {
	db, err := dbtest.Open()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		dbtest.Reset(db)
		dbtest.Close(db)
	}()

	store := NewStore(db)
	pkgs := NewPkgStore(store)
	ctx := context.TODO()
	t.Run("Create", testPkgCreate(ctx, pkgs))
	t.Run("List", testPkgList(ctx, pkgs))
	t.Run("Find", testPkgFind(ctx, pkgs))
	t.Run("Update", testPkgUpdate(ctx, pkgs))
	t.Run("Delete", testPkgDelete(ctx, pkgs))
}

func testPkgCreate(ctx context.Context, pkgs PkgStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData := []model.Pkg{
			{URL: "http://url.one", Title: "one", Description: "test pinl one"},
			{URL: "http://url.two", Title: "two", Description: "test pinl two"},
			{URL: "http://url.three", Title: "three", Description: "test pinl three"},
		}

		for _, pkg := range testData {
			assert.Nil(t, pkgs.Create(ctx, &pkg))
			assert.NotEmpty(t, pkg.ID)
			assert.False(t, pkg.CreatedAt.Time().IsZero())
		}
	}
}

func testPkgList(ctx context.Context, pkgs PkgStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, err := pkgs.List(ctx, nil)
		assert.Nil(t, err)
		assert.Equal(t, 3, len(testData))
	}
}

func testPkgFind(ctx context.Context, pkgs PkgStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, _ := pkgs.List(ctx, nil)

		want := testData[1]
		got := model.Pkg{ID: want.ID}
		err := pkgs.Find(ctx, &got)
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	}
}

func testPkgUpdate(ctx context.Context, pkgs PkgStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, _ := pkgs.List(ctx, nil)

		want := testData[2]
		want.Title = "(changed) " + want.Title
		err := pkgs.Update(ctx, &want)
		assert.Nil(t, err)

		got := model.Pkg{ID: want.ID}
		pkgs.Find(ctx, &got)
		assert.Equal(t, want, got)
	}
}

func testPkgDelete(ctx context.Context, pkgs PkgStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, _ := pkgs.List(ctx, nil)

		assert.Nil(t, pkgs.Delete(ctx, &testData[1]))

		testData, _ = pkgs.List(ctx, nil)
		assert.Equal(t, 2, len(testData))
	}
}
