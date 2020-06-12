package handleutils

import (
	"context"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/store"
)

func UpdateOrCreateTag(ctx context.Context, tags *store.Tags, userID, name string, data *model.Tag) (*model.Tag, error) {
	found, err := tags.List(ctx, &store.TagOpts{
		UserID: userID,
		Name:   name,
	})
	if err != nil {
		return nil, err
	}

	var tag *model.Tag
	if len(found) > 0 {
		tag = found[0]
	} else {
		tag = &model.Tag{
			UserID: userID,
			Name:   name,
		}
	}

	if data != nil {
		tag.Level = data.Level
		tag.ParentID = data.ParentID
	}

	if tag.ID == "" {
		err = tags.Create(ctx, tag)
	} else {
		err = tags.Update(ctx, tag)
	}
	if err != nil {
		return nil, err
	}
	return tag, nil
}
