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

	store := NewStore(db)
	substats := NewSubstatStore(store)
	ctx := context.TODO()
	t.Run("Create", testSubstatCreate(ctx, substats))
	t.Run("List", testSubstatList(ctx, substats))
	t.Run("Update", testSubstatUpdate(ctx, substats))
	t.Run("Delete", testSubstatDelete(ctx, substats))
}

func testSubstatCreate(ctx context.Context, substats SubstatStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData := []model.Substat{
			{Labels: field.Labels{"a1": "v1", "a2": "v2"}},
			{Labels: field.Labels{"a3": "v3", "a4": "v4"}},
			{Labels: field.Labels{"a5": "v5", "a6": "v6"}},
		}

		for _, substat := range testData {
			assert.Nil(t, substats.Create(ctx, &substat))
			assert.NotEmpty(t, substat.ID)
		}
	}
}

func testSubstatList(ctx context.Context, substats SubstatStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, err := substats.List(ctx, nil)
		assert.Nil(t, err)
		assert.Equal(t, 3, len(testData))
	}
}

func testSubstatUpdate(ctx context.Context, substats SubstatStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, _ := substats.List(ctx, nil)

		want := testData[0]
		want.Labels["new-a1"] = "new-v"
		assert.Nil(t, substats.Update(ctx, &want))
	}
}

func testSubstatDelete(ctx context.Context, substats SubstatStore) func(t *testing.T) {
	return func(t *testing.T) {
		testData, _ := substats.List(ctx, nil)

		assert.Nil(t, substats.Delete(ctx, &testData[0]))

		testData, _ = substats.List(ctx, nil)
		assert.Equal(t, 2, len(testData))
	}
}
