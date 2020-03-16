package store

import (
	"context"
	"testing"

	"github.com/pinmonl/pinmonl/database/dbtest"
	"github.com/pinmonl/pinmonl/model"
	"github.com/stretchr/testify/assert"
)

func TestTaggableStore(t *testing.T) {
	db, err := dbtest.Open()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		dbtest.Reset(db)
		dbtest.Close(db)
	}()

	mockTags := []*model.Tag{
		{Name: "tag1", UserID: "user-test-id1", Level: 0, Sort: 0},
		{Name: "tag2", UserID: "user-test-id1", Level: 0, Sort: 0},
		{Name: "tag3", UserID: "user-test-id2", Level: 0, Sort: 1},
		{Name: "tag4", UserID: "user-test-id2", Level: 0, Sort: 0},
		{Name: "tag5", UserID: "user-test-id3", Level: 0, Sort: 1},
		{Name: "tag6", UserID: "user-test-id3", Level: 0, Sort: 0},
	}
	mockPinls := []*model.Pinl{
		{Title: "pinl1", URL: "http://url.one"},
		{Title: "pinl2", URL: "http://url.two"},
		{Title: "pinl3", URL: "http://url.three"},
		{Title: "pinl4", URL: "http://url.four"},
		{Title: "pinl5", URL: "http://url.five"},
		{Title: "pinl6", URL: "http://url.six"},
	}
	mockData := []*model.Taggable{
		{Tag: mockTags[0], Pinl: mockPinls[0]},
		{Tag: mockTags[1], Pinl: mockPinls[1]},
		{Tag: mockTags[2], Pinl: mockPinls[2]},
		{Tag: mockTags[3], Pinl: mockPinls[3]},
		{Tag: mockTags[4], Pinl: mockPinls[4]},
		{Tag: mockTags[5], Pinl: mockPinls[5]},
	}

	store := NewStore(db)
	tags := NewTagStore(store)
	pinls := NewPinlStore(store)
	taggables := NewTaggableStore(store)
	ctx := context.TODO()

	for _, pinl := range mockPinls {
		pinls.Create(ctx, pinl)
	}
	for _, tag := range mockTags {
		tags.Create(ctx, tag)
	}

	t.Run("Create", testTaggableCreate(ctx, taggables, mockData))
	t.Run("List", testTaggableList(ctx, taggables, mockData))
	t.Run("ListTags", testTaggableListTags(ctx, taggables, mockData))
	t.Run("ClearTags", testTaggableClearTags(ctx, taggables, mockData))
	t.Run("AssocTag", testTaggableAssocTag(ctx, taggables, mockData))
	t.Run("AssocTags", testTaggableAssocTags(ctx, taggables, mockData))
	t.Run("ReAssocTags", testTaggableReAssocTags(ctx, taggables, mockData))
	t.Run("Delete", testTaggableDelete(ctx, taggables, mockData))
	t.Run("DissocTag", testTaggableDissocTag(ctx, taggables, mockData))
	t.Run("DissocTags", testTaggableDissocTags(ctx, taggables, mockData))
}

func testTaggableCreate(ctx context.Context, taggables TaggableStore, mockData []*model.Taggable) func(*testing.T) {
	return func(t *testing.T) {
		for _, taggable := range mockData {
			taggable.TagID = taggable.Tag.ID
			taggable.TargetID = taggable.Pinl.MorphKey()
			taggable.TargetName = taggable.Pinl.MorphName()
			assert.Nil(t, taggables.Create(ctx, taggable))
		}
	}
}

func testTaggableList(ctx context.Context, taggables TaggableStore, mockData []*model.Taggable) func(*testing.T) {
	return func(t *testing.T) {
		deRef := func(data []*model.Taggable) []model.Taggable {
			out := make([]model.Taggable, len(data))
			for i, mt := range data {
				m := *mt
				m.Tag = nil
				m.Pinl = nil
				out[i] = m
			}
			return out
		}

		want := deRef(mockData)
		got, err := taggables.List(ctx, nil)
		assert.Nil(t, err)
		assert.Equal(t, want, got)

		want = deRef(mockData[:1])
		got, err = taggables.List(ctx, &TaggableOpts{Target: model.Pinl{ID: want[0].TargetID}})
		assert.Nil(t, err)
		assert.Equal(t, want, got)

		want = deRef(mockData[3:4])
		got, err = taggables.List(ctx, &TaggableOpts{Tags: []model.Tag{{ID: want[0].TagID}}})
		assert.Nil(t, err)
		assert.Equal(t, want, got)

		got, err = taggables.List(ctx, &TaggableOpts{TagIDs: []string{want[0].TagID}})
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	}
}

func testTaggableListTags(ctx context.Context, taggables TaggableStore, mockData []*model.Taggable) func(*testing.T) {
	return func(t *testing.T) {
		deRef := func(data []*model.Taggable) map[string][]model.Tag {
			out := map[string][]model.Tag{}
			for _, mt := range data {
				m := *mt
				k := m.TargetID
				out[k] = append(out[k], *m.Tag)
			}
			return out
		}

		want := deRef(mockData)
		got, err := taggables.ListTags(ctx, nil)
		assert.Nil(t, err)
		assert.Equal(t, want, got)

		want = deRef(mockData[:1])
		got, err = taggables.ListTags(ctx, &TaggableOpts{Target: model.Pinl{ID: mockData[0].TargetID}})
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	}
}

func testTaggableClearTags(ctx context.Context, taggables TaggableStore, mockData []*model.Taggable) func(*testing.T) {
	return func(t *testing.T) {
		pinl := *(mockData[0].Pinl)
		want := ([]model.Tag)(nil)
		assert.Nil(t, taggables.ClearTags(ctx, pinl))
		check, _ := taggables.ListTags(ctx, &TaggableOpts{Target: pinl})
		got := check[pinl.ID]
		assert.Equal(t, want, got)

		taggables.ClearTags(ctx, pinl)
		taggables.Create(ctx, mockData[0])
	}
}

func testTaggableAssocTag(ctx context.Context, taggables TaggableStore, mockData []*model.Taggable) func(*testing.T) {
	return func(t *testing.T) {
		pinl := *(mockData[0].Pinl)
		tag := *(mockData[1].Tag)
		want := []model.Tag{*(mockData[0].Tag), tag}
		assert.Nil(t, taggables.AssocTag(ctx, pinl, tag))
		check, _ := taggables.ListTags(ctx, &TaggableOpts{Target: pinl})
		got := check[pinl.ID]
		assert.Equal(t, want, got)

		taggables.ClearTags(ctx, pinl)
		taggables.Create(ctx, mockData[0])
	}
}

func testTaggableAssocTags(ctx context.Context, taggables TaggableStore, mockData []*model.Taggable) func(*testing.T) {
	return func(t *testing.T) {
		pinl := *(mockData[0].Pinl)
		tags := []model.Tag{*(mockData[2].Tag), *(mockData[3].Tag)}
		want := append([]model.Tag{*(mockData[0].Tag)}, tags...)
		assert.Nil(t, taggables.AssocTags(ctx, pinl, tags))
		check, _ := taggables.ListTags(ctx, &TaggableOpts{Target: pinl})
		got := check[pinl.ID]
		assert.Equal(t, want, got)

		taggables.ClearTags(ctx, pinl)
		taggables.Create(ctx, mockData[0])
	}
}

func testTaggableReAssocTags(ctx context.Context, taggables TaggableStore, mockData []*model.Taggable) func(*testing.T) {
	return func(t *testing.T) {
		pinl := *(mockData[0].Pinl)
		want := []model.Tag{*(mockData[2].Tag), *(mockData[3].Tag)}
		assert.Nil(t, taggables.ReAssocTags(ctx, pinl, want))
		check, _ := taggables.ListTags(ctx, &TaggableOpts{Target: pinl})
		got := check[pinl.ID]
		assert.Equal(t, want, got)

		taggables.ClearTags(ctx, pinl)
		taggables.Create(ctx, mockData[0])
	}
}

func testTaggableDelete(ctx context.Context, taggables TaggableStore, mockData []*model.Taggable) func(*testing.T) {
	return func(t *testing.T) {
		pinl := *(mockData[0].Pinl)
		tag := *(mockData[0].Tag)
		want := ([]model.Tag)(nil)
		assert.Nil(t, taggables.Delete(ctx, &model.Taggable{TagID: tag.ID, TargetID: pinl.MorphKey(), TargetName: pinl.MorphName()}))
		check, _ := taggables.ListTags(ctx, &TaggableOpts{Target: pinl})
		got := check[pinl.ID]
		assert.Equal(t, want, got)

		taggables.ClearTags(ctx, pinl)
		taggables.Create(ctx, mockData[0])
	}
}

func testTaggableDissocTag(ctx context.Context, taggables TaggableStore, mockData []*model.Taggable) func(*testing.T) {
	return func(t *testing.T) {
		pinl := *(mockData[0].Pinl)
		tag := *(mockData[0].Tag)
		want := ([]model.Tag)(nil)
		assert.Nil(t, taggables.DissocTag(ctx, pinl, tag))
		check, _ := taggables.ListTags(ctx, &TaggableOpts{Target: pinl})
		got := check[pinl.ID]
		assert.Equal(t, want, got)

		taggables.ClearTags(ctx, pinl)
		taggables.Create(ctx, mockData[0])
	}
}

func testTaggableDissocTags(ctx context.Context, taggables TaggableStore, mockData []*model.Taggable) func(*testing.T) {
	return func(t *testing.T) {
		pinl := *(mockData[0].Pinl)
		tags := []model.Tag{*(mockData[2].Tag), *(mockData[3].Tag)}
		taggables.AssocTags(ctx, pinl, tags)
		want := []model.Tag{*(mockData[0].Tag)}
		assert.Nil(t, taggables.DissocTags(ctx, pinl, tags))
		check, _ := taggables.ListTags(ctx, &TaggableOpts{Target: pinl})
		got := check[pinl.ID]
		assert.Equal(t, want, got)

		taggables.ClearTags(ctx, pinl)
		taggables.Create(ctx, mockData[0])
	}
}
