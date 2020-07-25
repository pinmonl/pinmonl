package storeutils

import (
	"context"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pkgs/card"
	"github.com/pinmonl/pinmonl/store"
)

func SavePinl(ctx context.Context, pinls *store.Pinls, images *store.Images, data *model.Pinl, withCard bool) (*model.Pinl, *model.Image, error) {
	var (
		pinl     = *data
		isNew    = (pinl.ID == "")
		imgbytes []byte
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

	// Saves pinl.
	var err error
	if isNew {
		err = pinls.Create(ctx, &pinl)
	} else {
		err = pinls.Update(ctx, &pinl)
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
		if err := pinls.Update(ctx, &pinl); err != nil {
			return nil, nil, err
		}
	}

	return &pinl, image, nil
}

func ListPinlsWithLatestStats(ctx context.Context, pinls *store.Pinls, monpkgs *store.Monpkgs, stats *store.Stats, taggables *store.Taggables, opts *store.PinlOpts) (model.PinlList, error) {
	pList, err := pinls.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	mpList, err := monpkgs.ListWithPkg(ctx, &store.MonpkgOpts{
		MonlIDs: pList.MonlKeys(),
	})
	if err != nil {
		return nil, err
	}

	sList, err := stats.List(ctx, &store.StatOpts{
		PkgIDs: mpList.Pkgs().Keys(),
	})
	if err != nil {
		return nil, err
	}
	if sList2, err := ListStatTree(ctx, stats, sList); err == nil {
		sList = sList2
	} else {
		return nil, err
	}

	tMap, err := GetTags(ctx, taggables, pList.Morphables())
	if err != nil {
		return nil, err
	}

	mpList.Pkgs().SetStats(sList)
	pList.SetPkgs(mpList.PkgsByMonl())
	pList.SetTagNames(tMap)

	return pList, nil
}

func PinlWithLatestStats(ctx context.Context, pinls *store.Pinls, monpkgs *store.Monpkgs, stats *store.Stats, taggables *store.Taggables, pinlID string) (*model.Pinl, error) {
	pList, err := ListPinlsWithLatestStats(ctx, pinls, monpkgs, stats, taggables, &store.PinlOpts{
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
