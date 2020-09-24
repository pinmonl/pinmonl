package storeutils

import (
	"context"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/store"
)

func SaveTaggables(
	ctx context.Context,
	taggables *store.Taggables,
	userID string,
	target model.Morphable,
	tagIDs []string,
	pivotValues map[string]*model.TagPivot,
) (model.TaggableList, error) {
	origTgList, err := taggables.List(ctx, &store.TaggableOpts{
		TargetIDs:  []string{target.MorphKey()},
		TargetName: target.MorphName(),
	})
	if err != nil {
		return nil, err
	}

	var (
		deletes = make(map[string]*model.Taggable)
		saves   = make(map[string]*model.Taggable)
	)
	for _, tg := range origTgList {
		deletes[tg.TagID] = tg
	}
	for _, tagID := range tagIDs {
		tg := &model.Taggable{
			TagID:      tagID,
			TargetID:   target.MorphKey(),
			TargetName: target.MorphName(),
		}
		if match, has := deletes[tagID]; has {
			tg.ID = match.ID
			delete(deletes, tagID)
		}
		if pivot, hasPivot := pivotValues[tagID]; hasPivot {
			tg.Value = pivot.Value
			tg.ValuePrefix = pivot.Prefix
			tg.ValueSuffix = pivot.Suffix
		}
		saves[tagID] = tg
	}

	out := make([]*model.Taggable, 0)
	for _, saveTg := range saves {
		var err error
		if saveTg.ID == "" {
			err = taggables.Create(ctx, saveTg)
		} else {
			err = taggables.Update(ctx, saveTg)
		}
		if err != nil {
			return nil, err
		}
		out = append(out, saveTg)
	}
	for _, delTg := range deletes {
		if _, err := taggables.Delete(ctx, delTg.ID); err != nil {
			return nil, err
		}
	}

	return out, nil
}

func GetTaggables(ctx context.Context, taggables *store.Taggables, targets model.MorphableList) (map[string]model.TaggableList, error) {
	tgList, err := taggables.List(ctx, &store.TaggableOpts{
		Targets: targets,
	})
	if err != nil {
		return nil, err
	}

	return tgList.ByTarget(), nil
}
