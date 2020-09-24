package storeutils

import (
	"context"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/store"
)

func SavePinpkgs(ctx context.Context, pinpkgs *store.Pinpkgs, pinlID string, pkgIDs []string) (model.PinpkgList, error) {
	origPpList, err := pinpkgs.List(ctx, &store.PinpkgOpts{
		PinlIDs: []string{pinlID},
	})
	if err != nil {
		return nil, err
	}

	var (
		deletes = make(map[string]*model.Pinpkg)
		saves   = make(map[string]*model.Pinpkg)
	)
	for _, pp := range origPpList {
		deletes[pp.PkgID] = pp
	}
	for _, pkgID := range pkgIDs {
		pp := &model.Pinpkg{
			PinlID: pinlID,
			PkgID:  pkgID,
		}
		if match, has := deletes[pkgID]; has {
			pp.ID = match.ID
			delete(deletes, pkgID)
		}
		saves[pkgID] = pp
	}

	out := make([]*model.Pinpkg, 0)
	for _, savePp := range saves {
		var err error
		if savePp.ID == "" {
			err = pinpkgs.Create(ctx, savePp)
		} else {
			err = pinpkgs.Update(ctx, savePp)
		}
		if err != nil {
			return nil, err
		}
		out = append(out, savePp)
	}
	for _, delPp := range deletes {
		if _, err := pinpkgs.Delete(ctx, delPp.ID); err != nil {
			return nil, err
		}
	}

	return out, nil
}
