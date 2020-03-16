package store

import (
	"context"
	"testing"

	"github.com/pinmonl/pinmonl/database/dbtest"
	"github.com/pinmonl/pinmonl/model"
	"github.com/stretchr/testify/assert"
)

func TestSharetagStore(t *testing.T) {
	db, err := dbtest.Open()
	if err != nil {
		t.Error(err)
	}
	defer func() {
		dbtest.Reset(db)
		dbtest.Close(db)
	}()

	mockShares := []*model.Share{
		{Name: "share1"},
		{Name: "share2"},
		{Name: "share3"},
	}
	mockTags := []*model.Tag{
		{Name: "tag1", UserID: "user-test-id1", Level: 0, Sort: 0},
		{Name: "tag2", UserID: "user-test-id1", Level: 0, Sort: 0},
		{Name: "tag3", UserID: "user-test-id2", Level: 0, Sort: 1},
		{Name: "tag4", UserID: "user-test-id2", Level: 0, Sort: 0},
		{Name: "tag5", UserID: "user-test-id3", Level: 0, Sort: 1},
		{Name: "tag6", UserID: "user-test-id3", Level: 0, Sort: 0},
	}
	mockData := []*model.Sharetag{
		{Share: mockShares[0], Tag: mockTags[0], Kind: model.SharetagKindMust},
		{Share: mockShares[0], Tag: mockTags[1], Kind: model.SharetagKindAny},
		{Share: mockShares[1], Tag: mockTags[2], Kind: model.SharetagKindMust},
		{Share: mockShares[1], Tag: mockTags[3], Kind: model.SharetagKindAny},
		{Share: mockShares[2], Tag: mockTags[4], Kind: model.SharetagKindMust},
		{Share: mockShares[2], Tag: mockTags[5], Kind: model.SharetagKindAny},
	}

	store := NewStore(db)
	shares := NewShareStore(store)
	tags := NewTagStore(store)
	sharetags := NewSharetagStore(store)
	ctx := context.TODO()

	for _, share := range mockShares {
		shares.Create(ctx, share)
	}
	for _, tag := range mockTags {
		tags.Create(ctx, tag)
	}

	t.Run("Create", testSharetagCreate(ctx, sharetags, mockData))
	t.Run("List", testSharetagList(ctx, sharetags, mockData))
	t.Run("ListTags", testSharetagListTags(ctx, sharetags, mockData))
	t.Run("ListShares", testSharetagListShares(ctx, sharetags, mockData))
	t.Run("ClearByKind", testSharetagClearByKind(ctx, sharetags, mockData))
	t.Run("AssocTag", testSharetagAssocTag(ctx, sharetags, mockData))
	t.Run("AssocTags", testSharetagAssocTags(ctx, sharetags, mockData))
	t.Run("ReAssocTags", testSharetagReAssocTags(ctx, sharetags, mockData))
	t.Run("Delete", testSharetagDelete(ctx, sharetags, mockData))
	t.Run("DissocTag", testSharetagDissocTag(ctx, sharetags, mockData))
	t.Run("DissocTags", testSharetagDissocTags(ctx, sharetags, mockData))
}

func testSharetagCreate(ctx context.Context, sharetags SharetagStore, mockData []*model.Sharetag) func(*testing.T) {
	return func(t *testing.T) {
		for _, sharetag := range mockData {
			sharetag.ShareID = sharetag.Share.ID
			sharetag.TagID = sharetag.Tag.ID
			assert.Nil(t, sharetags.Create(ctx, sharetag))
		}
	}
}

func testSharetagList(ctx context.Context, sharetags SharetagStore, mockData []*model.Sharetag) func(*testing.T) {
	return func(t *testing.T) {
		deRef := func(data []*model.Sharetag) []model.Sharetag {
			out := make([]model.Sharetag, len(data))
			for i, mst := range data {
				m := *mst
				m.Share = nil
				m.Tag = nil
				out[i] = m
			}
			return out
		}

		want := deRef(mockData)
		got, err := sharetags.List(ctx, nil)
		assert.Nil(t, err)
		assert.Equal(t, want, got)

		want = deRef(mockData[:2])
		got, err = sharetags.List(ctx, &SharetagOpts{ShareID: want[0].ShareID})
		assert.Nil(t, err)
		assert.Equal(t, want, got)

		want = deRef(mockData[:1])
		got, err = sharetags.List(ctx, &SharetagOpts{TagID: want[0].TagID})
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	}
}

func testSharetagListTags(ctx context.Context, sharetags SharetagStore, mockData []*model.Sharetag) func(*testing.T) {
	return func(t *testing.T) {
		deRef := func(data []*model.Sharetag) map[string][]model.Tag {
			out := map[string][]model.Tag{}
			for _, mst := range data {
				k := mst.ShareID
				out[k] = append(out[k], *mst.Tag)
			}
			return out
		}

		want := deRef(mockData)
		got, err := sharetags.ListTags(ctx, nil)
		assert.Nil(t, err)
		assert.Equal(t, want, got)

		want = deRef(mockData[2:4])
		got, err = sharetags.ListTags(ctx, &SharetagOpts{ShareID: mockData[2].ShareID})
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	}
}

func testSharetagListShares(ctx context.Context, sharetags SharetagStore, mockData []*model.Sharetag) func(*testing.T) {
	return func(t *testing.T) {
		deRef := func(data []*model.Sharetag) map[string][]model.Share {
			out := map[string][]model.Share{}
			for _, mst := range data {
				k := mst.TagID
				out[k] = append(out[k], *mst.Share)
			}
			return out
		}

		want := deRef(mockData)
		got, err := sharetags.ListShares(ctx, nil)
		assert.Nil(t, err)
		assert.Equal(t, want, got)

		want = deRef(mockData[4:5])
		got, err = sharetags.ListShares(ctx, &SharetagOpts{TagID: mockData[4].TagID})
		assert.Nil(t, err)
		assert.Equal(t, want, got)
	}
}

func testSharetagClearByKind(ctx context.Context, sharetags SharetagStore, mockData []*model.Sharetag) func(*testing.T) {
	return func(t *testing.T) {
		share := *(mockData[0].Share)
		want := ([]model.Tag)(nil)
		assert.Nil(t, sharetags.ClearByKind(ctx, share, model.SharetagKindMust))
		check, _ := sharetags.ListTags(ctx, &SharetagOpts{ShareID: share.ID, Kind: model.SharetagKindMust})
		got := check[share.ID]
		assert.Equal(t, want, got)

		want = []model.Tag{*(mockData[1].Tag)}
		check, _ = sharetags.ListTags(ctx, &SharetagOpts{ShareID: share.ID, Kind: model.SharetagKindAny})
		got = check[share.ID]
		assert.Equal(t, want, got)

		want = nil
		assert.Nil(t, sharetags.ClearByKind(ctx, share, model.SharetagKindAny))
		check, _ = sharetags.ListTags(ctx, &SharetagOpts{ShareID: share.ID, Kind: model.SharetagKindAny})
		got = check[share.ID]
		assert.Equal(t, want, got)

		sharetags.Create(ctx, mockData[0])
		sharetags.Create(ctx, mockData[1])
	}
}

func testSharetagAssocTag(ctx context.Context, sharetags SharetagStore, mockData []*model.Sharetag) func(*testing.T) {
	return func(t *testing.T) {
		share := *(mockData[0].Share)
		want := []model.Tag{*(mockData[1].Tag), *(mockData[2].Tag)}
		assert.Nil(t, sharetags.AssocTag(ctx, share, model.SharetagKindAny, *(mockData[2].Tag)))
		check, _ := sharetags.ListTags(ctx, &SharetagOpts{ShareID: share.ID, Kind: model.SharetagKindAny})
		got := check[share.ID]
		assert.Equal(t, want, got)

		sharetags.ClearByKind(ctx, share, model.SharetagKindMust)
		sharetags.ClearByKind(ctx, share, model.SharetagKindAny)
		sharetags.Create(ctx, mockData[0])
		sharetags.Create(ctx, mockData[1])
	}
}

func testSharetagAssocTags(ctx context.Context, sharetags SharetagStore, mockData []*model.Sharetag) func(*testing.T) {
	return func(t *testing.T) {
		share := *(mockData[0].Share)
		tags := []model.Tag{*(mockData[3].Tag)}
		want := append([]model.Tag{*(mockData[0].Tag)}, tags...)
		assert.Nil(t, sharetags.AssocTags(ctx, share, model.SharetagKindMust, tags))
		check, _ := sharetags.ListTags(ctx, &SharetagOpts{ShareID: share.ID, Kind: model.SharetagKindMust})
		got := check[share.ID]
		assert.Equal(t, want, got)

		sharetags.ClearByKind(ctx, share, model.SharetagKindMust)
		sharetags.ClearByKind(ctx, share, model.SharetagKindAny)
		sharetags.Create(ctx, mockData[0])
		sharetags.Create(ctx, mockData[1])
	}
}

func testSharetagReAssocTags(ctx context.Context, sharetags SharetagStore, mockData []*model.Sharetag) func(*testing.T) {
	return func(t *testing.T) {
		share := *(mockData[0].Share)
		want := []model.Tag{*(mockData[1].Tag), *(mockData[3].Tag)}
		assert.Nil(t, sharetags.ReAssocTags(ctx, share, model.SharetagKindMust, want))
		check, _ := sharetags.ListTags(ctx, &SharetagOpts{ShareID: share.ID, Kind: model.SharetagKindMust})
		got := check[share.ID]
		assert.Equal(t, want, got)

		sharetags.ClearByKind(ctx, share, model.SharetagKindMust)
		sharetags.ClearByKind(ctx, share, model.SharetagKindAny)
		sharetags.Create(ctx, mockData[0])
		sharetags.Create(ctx, mockData[1])
	}
}

func testSharetagDelete(ctx context.Context, sharetags SharetagStore, mockData []*model.Sharetag) func(*testing.T) {
	return func(t *testing.T) {
		share := *(mockData[0].Share)
		tag := *(mockData[0].Tag)
		want := ([]model.Tag)(nil)
		assert.Nil(t, sharetags.Delete(ctx, &model.Sharetag{ShareID: share.ID, TagID: tag.ID}))
		check, _ := sharetags.ListTags(ctx, &SharetagOpts{ShareID: share.ID, Kind: model.SharetagKindMust})
		got := check[share.ID]
		assert.Equal(t, want, got)

		sharetags.ClearByKind(ctx, share, model.SharetagKindMust)
		sharetags.ClearByKind(ctx, share, model.SharetagKindAny)
		sharetags.Create(ctx, mockData[0])
		sharetags.Create(ctx, mockData[1])
	}
}

func testSharetagDissocTag(ctx context.Context, sharetags SharetagStore, mockData []*model.Sharetag) func(*testing.T) {
	return func(t *testing.T) {
		share := *(mockData[0].Share)
		tag := *(mockData[0].Tag)
		want := ([]model.Tag)(nil)
		assert.Nil(t, sharetags.DissocTag(ctx, share, model.SharetagKindMust, tag))
		check, _ := sharetags.ListTags(ctx, &SharetagOpts{ShareID: share.ID, Kind: model.SharetagKindMust})
		got := check[share.ID]
		assert.Equal(t, want, got)

		sharetags.ClearByKind(ctx, share, model.SharetagKindMust)
		sharetags.ClearByKind(ctx, share, model.SharetagKindAny)
		sharetags.Create(ctx, mockData[0])
		sharetags.Create(ctx, mockData[1])
	}
}

func testSharetagDissocTags(ctx context.Context, sharetags SharetagStore, mockData []*model.Sharetag) func(*testing.T) {
	return func(t *testing.T) {
		share := *(mockData[0].Share)
		tags := []model.Tag{*(mockData[0].Tag)}
		want := ([]model.Tag)(nil)
		assert.Nil(t, sharetags.DissocTags(ctx, share, model.SharetagKindMust, tags))
		check, _ := sharetags.ListTags(ctx, &SharetagOpts{ShareID: share.ID, Kind: model.SharetagKindMust})
		got := check[share.ID]
		assert.Equal(t, want, got)

		sharetags.ClearByKind(ctx, share, model.SharetagKindMust)
		sharetags.ClearByKind(ctx, share, model.SharetagKindAny)
		sharetags.Create(ctx, mockData[0])
		sharetags.Create(ctx, mockData[1])
	}
}
