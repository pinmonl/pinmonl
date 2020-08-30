package storeutils

import (
	"context"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/store"
)

func LoadPkgsLatestStats(ctx context.Context, stats *store.Stats, pkgs model.PkgList) error {
	sList, err := stats.List(ctx, &store.StatOpts{
		PkgIDs:   pkgs.Keys(),
		IsLatest: field.NewNullBool(true),
	})
	if err != nil {
		return err
	}

	for i := range pkgs {
		pkgStats := sList.GetPkgID(pkgs[i].ID)
		pkgs[i].SetStats(pkgStats)
	}
	return nil
}

func GetPkgs(ctx context.Context, monpkgs *store.Monpkgs, pinpkgs *store.Pinpkgs, pList model.PinlList) (map[string]model.PkgList, error) {
	var (
		monlIDs = make([]string, 0)
		pinlIDs = make([]string, 0)
		monlMap = make(map[string][]string)
	)
	for _, p := range pList {
		if p.HasPinpkgs {
			pinlIDs = append(pinlIDs, p.ID)
		} else {
			monlIDs = append(monlIDs, p.MonlID)
		}
		monlMap[p.MonlID] = append(monlMap[p.MonlID], p.ID)
	}

	mpList, err := monpkgs.ListWithPkg(ctx, &store.MonpkgOpts{
		MonlIDs: monlIDs,
	})
	if err != nil {
		return nil, err
	}

	ppList, err := pinpkgs.ListWithPkg(ctx, &store.PinpkgOpts{
		PinlIDs: pinlIDs,
	})
	if err != nil {
		return nil, err
	}

	out := make(map[string]model.PkgList)
	for _, mp := range mpList {
		for _, k := range monlMap[mp.MonlID] {
			out[k] = append(out[k], mp.Pkg)
		}
	}
	for _, pp := range ppList {
		k := pp.PinlID
		out[k] = append(out[k], pp.Pkg)
	}
	return out, nil
}
