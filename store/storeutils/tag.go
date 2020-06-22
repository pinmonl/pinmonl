package storeutils

import (
	"context"
	"errors"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/store"
)

func SaveTag(ctx context.Context, tags *store.Tags, userID string, data *model.Tag) (*model.Tag, error) {
	var (
		tag   = *data
		isNew = tag.ID == ""
	)

	tag.UserID = userID
	if tag.ParentID != "" {
		found, err := tags.List(ctx, &store.TagOpts{
			UserID: userID,
			IDs:    []string{tag.ParentID},
		})
		if err != nil {
			return nil, err
		}
		if len(found) == 0 {
			return nil, errors.New("tag parent ownership inconsistence")
		}
	}

	var err error
	if isNew {
		err = tags.Create(ctx, &tag)
	} else {
		err = tags.Update(ctx, &tag)
	}
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func ReAssociateTags(ctx context.Context, tags *store.Tags, taggables *store.Taggables, target model.Morphable, userID string, tagNames []string) (model.TagList, error) {
	tgList, err := taggables.List(ctx, &store.TaggableOpts{
		Targets: model.MorphableList{target},
	})
	if err != nil {
		return nil, err
	}
	for _, tg := range tgList {
		_, err = taggables.Delete(ctx, tg.ID)
		if err != nil {
			return nil, err
		}
	}

	tList := model.TagList{}
	for _, tn := range tagNames {
		tag := &model.Tag{
			UserID: userID,
			Name:   tn,
		}
		tag, err = tags.FindOrCreate(ctx, tag)
		if err != nil {
			return nil, err
		}
		tList = append(tList, tag)
	}

	for _, t := range tList {
		tg := &model.Taggable{
			TagID:      t.ID,
			TargetID:   target.MorphKey(),
			TargetName: target.MorphName(),
		}
		err = taggables.Create(ctx, tg)
		if err != nil {
			return nil, err
		}
	}

	return tList, nil
}

func GetTags(ctx context.Context, taggables *store.Taggables, targets model.MorphableList) (map[string]model.TagList, error) {
	tgList, err := taggables.ListWithTag(ctx, &store.TaggableOpts{
		Targets: targets,
	})
	if err != nil {
		return nil, err
	}

	return tgList.TagsByTarget(), nil
}
