package store

import (
	"context"
	"testing"

	"github.com/pinmonl/pinmonl/database/dbtest"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/stretchr/testify/assert"
)

func TestSubstatStore(t *testing.T) {
	db, err := dbtest.Open()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		dbtest.Reset(db)
		dbtest.Close(db)
	}()

	mockData := []*model.Substat{
		{Labels: field.Labels{"a1": "v1", "a2": "v2"}},
		{Labels: field.Labels{"a3": "v3", "a4": "v4"}},
		{Labels: field.Labels{"a5": "v5", "a6": "v6"}},
	}

	store := NewStore(db)
	substats := NewSubstatStore(store)
	ctx := context.TODO()
	t.Run("Create", testSubstatCreate(ctx, substats, mockData))
	t.Run("List", testSubstatList(ctx, substats, mockData))
	t.Run("Update", testSubstatUpdate(ctx, substats, mockData))
	t.Run("Delete", testSubstatDelete(ctx, substats, mockData))
}

func testSubstatCreate(ctx context.Context, substats SubstatStore, mockData []*model.Substat) func(t *testing.T) {
	return func(t *testing.T) {
		for _, substat := range mockData {
			assert.Nil(t, substats.Create(ctx, substat))
			assert.NotEmpty(t, substat.ID)
		}
	}
}

func testSubstatList(ctx context.Context, substats SubstatStore, mockData []*model.Substat) func(t *testing.T) {
	return func(t *testing.T) {
		deRef := func(data []*model.Substat) []model.Substat {
			out := make([]model.Substat, len(data))
			for i, ms := range data {
				m := *ms
				out[i] = m
			}
			return out
		}

		want := deRef(mockData)
		got, err := substats.List(ctx, nil)
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	}
}

func testSubstatUpdate(ctx context.Context, substats SubstatStore, mockData []*model.Substat) func(t *testing.T) {
	return func(t *testing.T) {
		want := mockData[0]
		want.Labels["new-a1"] = "new-v"
		assert.Nil(t, substats.Update(ctx, want))
	}
}

func testSubstatDelete(ctx context.Context, substats SubstatStore, mockData []*model.Substat) func(t *testing.T) {
	return func(t *testing.T) {
		del, want := mockData[0], mockData[1:]
		assert.Nil(t, substats.Delete(ctx, del))

		got, _ := substats.List(ctx, nil)
		assert.Equal(t, len(want), len(got))
	}
}
