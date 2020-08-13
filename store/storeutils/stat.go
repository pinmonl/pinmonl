package storeutils

import (
	"context"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/store"
)

func ListStatTree(ctx context.Context, stats *store.Stats, rootStats []*model.Stat) (model.StatList, error) {
	var children []*model.Stat

	parent := model.StatList(rootStats).GetHasChildren()
	if len(parent) > 0 {
		children2, err := stats.List(ctx, &store.StatOpts{
			ParentIDs: parent.Keys(),
		})
		if err != nil {
			return nil, err
		}
		children = children2
	}

	if len(children) > 0 {
		children2, err := ListStatTree(ctx, stats, children)
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
