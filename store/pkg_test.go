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

	mockData := []*model.Pkg{
		{Vendor: "test", VendorURI: "uri1", URL: "http://url.one", Title: "one", Description: "test pinl one"},
		{Vendor: "test", VendorURI: "uri2", URL: "http://url.two", Title: "two", Description: "test pinl two"},
		{Vendor: "test", VendorURI: "uri3", URL: "http://url.three", Title: "three", Description: "test pinl three"},
	}

	store := NewStore(db)
	pkgs := NewPkgStore(store)
	ctx := context.TODO()
	t.Run("Create", testPkgCreate(ctx, pkgs, mockData))
	t.Run("List", testPkgList(ctx, pkgs, mockData))
	t.Run("Find", testPkgFind(ctx, pkgs, mockData))
	t.Run("Update", testPkgUpdate(ctx, pkgs, mockData))
	t.Run("Delete", testPkgDelete(ctx, pkgs, mockData))
}

func testPkgCreate(ctx context.Context, pkgs PkgStore, mockData []*model.Pkg) func(t *testing.T) {
	return func(t *testing.T) {
		for _, pkg := range mockData {
			assert.Nil(t, pkgs.Create(ctx, pkg))
			assert.NotEmpty(t, pkg.ID)
			assert.False(t, pkg.CreatedAt.Time().IsZero())
		}
	}
}

func testPkgList(ctx context.Context, pkgs PkgStore, mockData []*model.Pkg) func(t *testing.T) {
	return func(t *testing.T) {
		deRef := func(data []*model.Pkg) []model.Pkg {
			out := make([]model.Pkg, len(data))
			for i, mp := range data {
				m := *mp
				out[i] = m
			}
			return out
		}

		want := deRef(mockData)
		got, err := pkgs.List(ctx, nil)
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	}
}

func testPkgFind(ctx context.Context, pkgs PkgStore, mockData []*model.Pkg) func(t *testing.T) {
	return func(t *testing.T) {
		want := mockData[1]
		got := model.Pkg{ID: want.ID}
		err := pkgs.Find(ctx, &got)
		assert.Nil(t, err)
		assert.Equal(t, *want, got)
	}
}

func testPkgUpdate(ctx context.Context, pkgs PkgStore, mockData []*model.Pkg) func(t *testing.T) {
	return func(t *testing.T) {
		want := mockData[2]
		want.Title = "(changed) " + want.Title
		err := pkgs.Update(ctx, want)
		assert.Nil(t, err)

		got := model.Pkg{ID: want.ID}
		pkgs.Find(ctx, &got)
		assert.Equal(t, *want, got)
	}
}

func testPkgDelete(ctx context.Context, pkgs PkgStore, mockData []*model.Pkg) func(t *testing.T) {
	return func(t *testing.T) {
		del, want := mockData[0], mockData[1:]
		assert.Nil(t, pkgs.Delete(ctx, del))

		got, _ := pkgs.List(ctx, nil)
		assert.Equal(t, len(want), len(got))
	}
}
