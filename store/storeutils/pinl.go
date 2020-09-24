package storeutils

import (
	"context"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pkgs/card"
	"github.com/pinmonl/pinmonl/store"
)

func SavePinl(
	ctx context.Context,
	pinls *store.Pinls,
	taggables *store.Taggables,
	pinpkgs *store.Pinpkgs,
	images *store.Images,
	data *model.Pinl,
	withCard bool,
) (*model.Pinl, *model.Image, error) {
	pinl := &model.Pinl{}
	*pinl = *data

	var (
		isNew      = (pinl.ID == "")
		imgbytes   []byte
		hasPinpkgs = (pinl.PkgIDs != nil)
	)

	// Crawls card information by the url.
	if withCard {
		card, err := card.NewCard(pinl.URL)
		if err != nil {
			return nil, nil, err
		}
		pinl.Title = card.Title()
		pinl.Description = card.Description()

		if imgb, err := card.Image(); err == nil && imgb != nil {
			imgbytes = imgb
		}
	}
	pinl.HasPinpkgs = hasPinpkgs

	// Saves pinl.
	var err error
	if isNew {
		err = pinls.Create(ctx, pinl)
	} else {
		err = pinls.Update(ctx, pinl)
	}
	if err != nil {
		return nil, nil, err
	}

	// Associates relationship.
	var image *model.Image
	if withCard && imgbytes != nil {
		image, err = SaveImage(ctx, images, imgbytes, pinl, true)
		if err != nil {
			return nil, nil, err
		}

		pinl.ImageID = image.ID
		if err := pinls.Update(ctx, pinl); err != nil {
			return nil, nil, err
		}
	}

	if tgList, err := savePinlTaggables(ctx, taggables, pinl); err == nil {
		pinl.SetTagPivots(tgList)
	} else {
		return nil, nil, err
	}

	if _, err := savePinlPkgs(ctx, pinpkgs, pinl); err != nil {
		return nil, nil, err
	}

	return pinl, image, nil
}

func savePinlTaggables(ctx context.Context, taggables *store.Taggables, pinl *model.Pinl) (model.TaggableList, error) {
	var (
		tagIDs    = make([]string, 0)
		tagValues = make(map[string]*model.TagPivot)
	)
	if pinl.TagIDs != nil {
		tagIDs = *pinl.TagIDs
	}
	if pinl.TagValues != nil {
		tagValues = *pinl.TagValues
	}
	return SaveTaggables(ctx, taggables, pinl.UserID, pinl, tagIDs, tagValues)
}

func savePinlPkgs(ctx context.Context, pinpkgs *store.Pinpkgs, pinl *model.Pinl) (model.PinpkgList, error) {
	pkgIDs := make([]string, 0)
	if pinl.PkgIDs != nil {
		pkgIDs = *pinl.PkgIDs
	}
	return SavePinpkgs(ctx, pinpkgs, pinl.ID, pkgIDs)
}

func ListPinls(ctx context.Context, pinls *store.Pinls, monpkgs *store.Monpkgs, pinpkgs *store.Pinpkgs, taggables *store.Taggables, opts *store.PinlOpts) (model.PinlList, error) {
	pList, err := pinls.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	tMap, err := GetTaggables(ctx, taggables, pList.Morphables())
	if err != nil {
		return nil, err
	}

	pMap, err := GetPkgs(ctx, monpkgs, pinpkgs, pList)
	if err != nil {
		return nil, err
	}

	pList.SetTagPivots(tMap)
	pList.SetPkgIDs(pMap)
	return pList, nil
}

func GetPinl(ctx context.Context, pinls *store.Pinls, monpkgs *store.Monpkgs, pinpkgs *store.Pinpkgs, taggables *store.Taggables, pinlID string) (*model.Pinl, error) {
	pList, err := ListPinls(ctx, pinls, monpkgs, pinpkgs, taggables, &store.PinlOpts{
		IDs: []string{pinlID},
	})
	if err != nil {
		return nil, err
	}
	if len(pList) == 0 {
		return nil, nil
	}
	return pList[0], nil
}
