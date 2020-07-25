package storeutils

import (
	"context"
	"strings"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/store"
	"github.com/sirupsen/logrus"
)

func SaveTag(ctx context.Context, tags *store.Tags, userID string, data *model.Tag) (*model.Tag, error) {
	var (
		tag   = *data
		isNew = tag.ID == ""
	)

	tag.UserID = userID
	tag.Level = strings.Count(tag.Name, "/")
	if tag.Level > 0 {
		splits := strings.Split(tag.Name, "/")
		parentName := strings.Join(splits[:len(splits)-1], "/")
		parent, err := tags.FindName(ctx, userID, parentName)
		if err != nil {
			return nil, err
		}
		// If parent not found
		// OR parent is self (occurs when rename with same prefix).
		if parent == nil || parent.ID == tag.ID {
			parent = &model.Tag{Name: parentName}
		}

		parent.HasChildren = true
		if p2, err := SaveTag(ctx, tags, userID, parent); err == nil {
			parent = p2
		} else {
			return nil, err
		}

		tag.ParentID = parent.ID
	} else {
		tag.ParentID = ""
	}

	if !isNew {
		if err := rebuildTagHierarchy(ctx, tags, &tag); err != nil {
			return nil, err
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

func rebuildTagHierarchy(ctx context.Context, tags *store.Tags, tag *model.Tag) error {
	orig, err := tags.Find(ctx, tag.ID)
	if err != nil {
		return err
	}

	// Replaces prefix of children to new path.
	if orig.HasChildren {
		_, err := replaceTagsPrefix(ctx, tags, orig.Name, tag.Name)
		if err != nil {
			return err
		}
	}

	// Updates parent HasChildren flag.
	if orig.ParentID != tag.ParentID {
		count, err := tags.Count(ctx, &store.TagOpts{
			ParentIDs: []string{orig.ParentID},
		})
		if err != nil {
			return err
		}

		// count <= 1 to exclude this tag.
		if count <= 1 {
			parent, err := tags.Find(ctx, orig.ParentID)
			if err != nil {
				return err
			}
			parent.HasChildren = false
			err = tags.Update(ctx, parent)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func replaceTagsPrefix(ctx context.Context, tags *store.Tags, fromPrefix, toPrefix string) (int64, error) {
	tList, err := tags.List(ctx, &store.TagOpts{
		NamePattern: fromPrefix + "/%",
	})
	logrus.Debugln(tList, err)
	if err != nil {
		return 0, err
	}
	count := int64(0)
	for _, t := range tList {
		t.Name = toPrefix + strings.TrimPrefix(t.Name, fromPrefix)
		t.Level = strings.Count(t.Name, "/")
		err := tags.Update(ctx, t)
		if err != nil {
			return 0, err
		}
		count++
	}
	return count, nil
}

func ReAssociateTags(ctx context.Context, tags *store.Tags, taggables *store.Taggables, target model.Morphable, userID string, tagNames []string) (model.TagList, error) {
	_, err := taggables.DeleteByTarget(ctx, target)
	if err != nil {
		return nil, err
	}

	tList := model.TagList{}
	for _, tagName := range tagNames {
		tag, err := tags.FindName(ctx, userID, tagName)
		if err != nil {
			return nil, err
		}
		if tag == nil {
			if t2, err := SaveTag(ctx, tags, userID, &model.Tag{Name: tagName}); err == nil {
				tag = t2
			} else {
				return nil, err
			}
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

func DeleteTag(ctx context.Context, tags *store.Tags, taggables *store.Taggables, data *model.Tag) (int64, error) {
	tag := *data
	tag.ParentID = ""

	if err := rebuildTagHierarchy(ctx, tags, &tag); err != nil {
		return 0, err
	}

	// Get all children with same prefix.
	children, err := tags.List(ctx, &store.TagOpts{
		NamePattern: tag.Name + "/%",
	})
	if err != nil {
		return 0, err
	}

	// Append to del array.
	dels := make([]*model.Tag, 1+len(children))
	dels[0] = &tag
	for i := range children {
		dels[i+1] = children[i]
	}

	// Delete tags and theirs relation.
	for _, del := range dels {
		if _, err := tags.Delete(ctx, del.ID); err != nil {
			return 0, err
		}

		if _, err := taggables.DeleteByTag(ctx, del); err != nil {
			return 0, err
		}
	}

	// Update parent

	return int64(len(dels)), nil
}
