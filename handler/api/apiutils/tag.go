package apiutils

import (
	"context"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/store"
)

// FindOrCreateTagsByName handles the tags searching and creation.
func FindOrCreateTagsByName(
	ctx context.Context,
	tags store.TagStore,
	owner model.User,
	tagNames []string,
) ([]model.Tag, error) {
	if len(tagNames) == 0 {
		return []model.Tag{}, nil
	}

	exists, err := tags.List(ctx, &store.TagOpts{
		Names:  tagNames,
		UserID: owner.ID,
	})
	if err != nil {
		return nil, err
	}

	var ts []model.Tag
	for _, name := range tagNames {
		var t model.Tag
		found := false
		for _, et := range exists {
			if et.Name == name {
				t = et
				found = true
				break
			}
		}
		if !found {
			t = model.Tag{Name: name, UserID: owner.ID}
			err = tags.Create(ctx, &t)
			if err != nil {
				return nil, err
			}
		}
		ts = append(ts, t)
	}
	return ts, nil
}
