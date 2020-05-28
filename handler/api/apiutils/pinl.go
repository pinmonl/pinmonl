package apiutils

import (
	"context"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/store"
)

// ListPinlStats fetches pkgs and stats of pinls.
func ListPinlStats(
	ctx context.Context,
	monpkgStore store.MonpkgStore,
	statStore store.StatStore,
	pinls ...model.Pinl,
) (map[string][]model.Pkg, map[string][]model.Stat, error) {
	monlIDs := make([]string, len(pinls))
	for i, p := range pinls {
		monlIDs[i] = p.MonlID
	}
	pkgMap, err := monpkgStore.ListPkgs(ctx, &store.MonpkgOpts{
		MonlIDs: monlIDs,
	})
	if err != nil {
		return nil, nil, err
	}
	pkgIDs := make([]string, 0)
	for _, pkgs := range pkgMap {
		for _, pkg := range pkgs {
			pkgIDs = append(pkgIDs, pkg.ID)
		}
	}
	stats, err := statStore.List(ctx, &store.StatOpts{
		PkgIDs:     pkgIDs,
		WithLatest: true,
	})
	if err != nil {
		return nil, nil, err
	}
	statMap := make(map[string][]model.Stat)
	for _, stat := range stats {
		k := stat.PkgID
		statMap[k] = append(statMap[k], stat)
	}
	return pkgMap, statMap, nil
}
