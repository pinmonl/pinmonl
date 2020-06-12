package handleutils

import (
	"context"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/store"
)

func ListStatsOfPkgs(ctx context.Context, stats *store.Stats, pkgs []*model.Pkg) (model.StatList, error) {
	return listStatsOfPkgs(ctx, stats, pkgs, field.NullBool{})
}

func ListLatestStatsOfPkgs(ctx context.Context, stats *store.Stats, pkgs []*model.Pkg) (model.StatList, error) {
	return listStatsOfPkgs(ctx, stats, pkgs, field.NewNullBool(true))
}

func listStatsOfPkgs(ctx context.Context, stats *store.Stats, pkgs []*model.Pkg, isLatest field.NullBool) (model.StatList, error) {
	opts := &store.StatOpts{
		PkgIDs:    model.PkgList(pkgs).Keys(),
		ParentIDs: []string{""},
		IsLatest:  isLatest,
	}

	sList, err := stats.List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return ListStatsTree(ctx, stats, sList)
}

func ListStatsTree(ctx context.Context, stats *store.Stats, rootStats []*model.Stat) (model.StatList, error) {
	children, err := stats.List(ctx, &store.StatOpts{
		ParentIDs: model.StatList(rootStats).Keys(),
	})
	if err != nil {
		return nil, err
	}

	if len(children) > 0 {
		children2, err := ListStatsTree(ctx, stats, children)
		if err != nil {
			return nil, err
		}
		children = children2
	}

	out := make([]*model.Stat, len(rootStats))
	for i := range rootStats {
		stat := *rootStats[i]
		subList := model.StatList(children).GetParentID(stat.ID)
		subList = append(model.StatList{}, subList...)
		stat.Substats = &subList
		out[i] = &stat
	}

	return out, nil
}
