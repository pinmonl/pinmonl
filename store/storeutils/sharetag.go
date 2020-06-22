package storeutils

import (
	"context"
	"errors"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/store"
)

// SaveSharetag saves sharetag and validates its relationships.
func SaveSharetag(ctx context.Context, sharetags *store.Sharetags, tags *store.Tags, userID, shareID string, data *model.Sharetag) (*model.Sharetag, error) {
	var (
		sharetag = *data
		isNew    = sharetag.ID == ""
	)

	// Find sharetag by share and tag id.
	found, err := sharetags.List(ctx, &store.SharetagOpts{
		ShareIDs: []string{shareID},
		TagIDs:   []string{sharetag.TagID},
	})
	if err != nil {
		return nil, err
	}
	// Set isNew to false if found.
	if len(found) > 0 && isNew {
		sharetag.ID = found[0].ID
		isNew = false
	}

	// Validate tag.
	tag, err := tags.Find(ctx, sharetag.TagID)
	if err != nil {
		return nil, err
	}
	if tag == nil || tag.UserID != userID {
		return nil, errors.New("tag not found")
	}

	// Validate parent tag.
	if sharetag.ParentID == sharetag.TagID {
		return nil, errors.New("tag and parent cannot be same")
	}
	if sharetag.ParentID != "" {
		parent, err := tags.Find(ctx, sharetag.ParentID)
		if err != nil {
			return nil, err
		}
		if parent == nil || parent.UserID != userID {
			return nil, errors.New("parent not found")
		}
	}

	// Force share id.
	sharetag.ShareID = shareID

	// Save.
	if isNew {
		err = sharetags.Create(ctx, &sharetag)
	} else {
		err = sharetags.Update(ctx, &sharetag)
	}
	if err != nil {
		return nil, err
	}
	return &sharetag, nil
}
