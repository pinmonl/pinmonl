package handleutils

import (
	"context"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/store"
)

func UpdateOrCreateSharetag(ctx context.Context, sharetags *store.Sharetags, shareID, tagID string, data *model.Sharetag) (*model.Sharetag, error) {
	found, err := sharetags.List(ctx, &store.SharetagOpts{
		ShareIDs: []string{shareID},
		TagIDs:   []string{tagID},
	})
	if err != nil {
		return nil, err
	}

	var sharetag *model.Sharetag
	if len(found) > 0 {
		sharetag = found[0]
	} else {
		sharetag = &model.Sharetag{
			ShareID: shareID,
			TagID:   tagID,
		}
	}

	if data != nil {
		sharetag.Kind = data.Kind
		sharetag.Level = data.Level
		sharetag.ParentID = data.ParentID
	}

	if sharetag.ID == "" {
		err = sharetags.Create(ctx, sharetag)
	} else {
		err = sharetags.Update(ctx, sharetag)
	}
	if err != nil {
		return nil, err
	}
	return sharetag, nil
}
