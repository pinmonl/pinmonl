package storeutils

import (
	"context"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/monler/provider"
	"github.com/pinmonl/pinmonl/store"
)

func SaveProviderReport(
	ctx context.Context,
	pkgs *store.Pkgs,
	stats *store.Stats,
	report provider.Report,
	force bool,
) (*model.Pkg, model.StatList, error) {
	pkg, err := findOrCreatePkgFromReport(ctx, pkgs, report)
	if err != nil {
		return nil, nil, err
	}
	if !force && !pkg.FetchedAt.Time().IsZero() {
		return pkg, nil, nil
	}

	sList := model.StatList{}
	if rsList, err := saveReportStats(ctx, stats, pkg.ID, report); err == nil {
		sList = append(sList, rsList...)
	} else {
		return nil, nil, err
	}
	if rtList, err := saveReportTags(ctx, stats, pkg.ID, report); err == nil {
		sList = append(sList, rtList...)
	} else {
		return nil, nil, err
	}

	pkg.FetchedAt = field.Now()
	if err := pkgs.Update(ctx, pkg); err != nil {
		return nil, nil, err
	}

	return pkg, sList, nil
}

type statKey struct {
	kind     model.StatKind
	name     string
	value    string
	checksum string
}

func saveReportStats(ctx context.Context, stats *store.Stats, pkgID string, report provider.Report) (model.StatList, error) {
	// Set previous stats, except release kind, to expired.
	if err := updateExpiredStats(ctx, stats, pkgID); err != nil {
		return nil, err
	}

	rsList, err := report.Stats()
	if err != nil {
		return nil, err
	}

	// Save latest stats.
	out := make([]*model.Stat, len(rsList))
	channels := make([]*model.Stat, 0)
	for i, stat := range rsList {
		// Skip for channel stat.
		if stat.Kind == model.ChannelStat {
			channels = append(channels, stat)
			continue
		}

		// Save stat.
		stat2, err := saveStat(ctx, stats, pkgID, stat)
		if err != nil {
			return nil, err
		}
		out[i] = stat2
	}

	if outCh, err := saveChannelStats(ctx, stats, pkgID, channels); err == nil {
		out = append(out, outCh...)
	} else {
		return nil, err
	}

	return out, nil
}

func saveChannelStats(ctx context.Context, stats *store.Stats, pkgID string, channels []*model.Stat) (model.StatList, error) {
	var (
		inserts = make(map[statKey]*model.Stat)
		updates = make(map[statKey]*model.Stat)
		deletes = make(map[statKey]*model.Stat)
		out     = make([]*model.Stat, 0)
	)

	makeKey := func(stat *model.Stat) statKey {
		return statKey{
			value:    stat.Value,
			checksum: stat.Checksum,
		}
	}

	for _, ch := range channels {
		inserts[makeKey(ch)] = ch
	}

	prev, err := stats.List(ctx, &store.StatOpts{
		Kinds:  []model.StatKind{model.ChannelStat},
		PkgIDs: []string{pkgID},
	})
	if err != nil {
		return nil, err
	}
	for _, prevCh := range prev {
		key := makeKey(prevCh)
		if _, has := inserts[key]; has {
			updates[key] = inserts[key]
			updates[key].ID = prevCh.ID
			delete(inserts, key)
		} else {
			deletes[key] = prevCh
		}
	}

	for _, newStat := range inserts {
		if stat, err := saveStat(ctx, stats, pkgID, newStat); err == nil {
			out = append(out, stat)
		} else {
			return nil, err
		}
	}
	for _, oldStat := range updates {
		if err := clearSubstats(ctx, stats, oldStat); err != nil {
			return nil, err
		}
		if stat, err := saveStat(ctx, stats, pkgID, oldStat); err == nil {
			out = append(out, stat)
		} else {
			return nil, err
		}
	}
	for _, delStat := range deletes {
		if err := clearSubstats(ctx, stats, delStat); err != nil {
			return nil, err
		}
		if _, err := stats.Delete(ctx, delStat.ID); err != nil {
			return nil, err
		}
	}
	return out, nil
}

type reportTagKey struct {
	kind  model.StatKind
	value string
}

func saveReportTags(ctx context.Context, stats *store.Stats, pkgID string, report provider.Report) (model.StatList, error) {
	// Get all tags.
	tags := model.StatList{}
	for report.Next() {
		tag, err := report.Tag()
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	// Get previous tags.
	prevTags, err := stats.List(ctx, &store.StatOpts{
		PkgIDs:    []string{pkgID},
		Kinds:     model.ReleaseStatKinds,
		ParentIDs: []string{""},
	})
	if err != nil {
		return nil, err
	}

	prevTagSet := make(map[reportTagKey]*model.Stat)
	for _, tag := range prevTags {
		key := reportTagKey{kind: tag.Kind, value: tag.Value}
		prevTagSet[key] = tag
	}

	latestKeys := make(map[reportTagKey]int)
	for _, tag := range tags.GetLatest() {
		key := reportTagKey{
			kind:  tag.Kind,
			value: tag.Value,
		}
		latestKeys[key]++
	}

	// Save tags.
	out := make([]*model.Stat, len(tags))
	for i := range tags {
		tag := &model.Stat{}
		*tag = *tags[i]
		out[i] = tag

		key := reportTagKey{kind: tag.Kind, value: tag.Value}
		prevTag, has := prevTagSet[key]
		if has {
			tag.ID = prevTag.ID
			// Set is not latest.
			if _, match := latestKeys[key]; !match {
				tag.IsLatest = false
			}
		}
		_, err := saveStat(ctx, stats, pkgID, tag)
		// logrus.Debugf("saving tag %v, %t, %v", tag.Value, tag.IsLatest, err)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func findOrCreatePkgFromReport(ctx context.Context, pkgs *store.Pkgs, report provider.Report) (*model.Pkg, error) {
	pu, err := report.URI()
	if err != nil {
		return nil, err
	}

	// Find by uri.
	pkg, err := pkgs.FindURI(ctx, pu)
	if err != nil {
		return nil, err
	}
	// Create if not found.
	if pkg == nil {
		pkg = &model.Pkg{}
		pkg.UnmarshalPkgURI(pu)
		err := pkgs.Create(ctx, pkg)
		if err != nil {
			return nil, err
		}
	}

	if info, err := report.Pkg(); err == nil && info != nil {
		pkg.Title = info.Title
		pkg.CustomUri = info.CustomUri
		pkg.Description = info.Description
		if err := pkgs.Update(ctx, pkg); err != nil {
			return nil, err
		}
	}

	return pkg, nil
}

func saveStat(ctx context.Context, stats *store.Stats, pkgID string, data *model.Stat) (*model.Stat, error) {
	stat := *data
	stat.PkgID = pkgID
	if stat.Substats != nil && len(*stat.Substats) > 0 {
		stat.HasChildren = true
	}

	// Save.
	var err error
	if stat.ID == "" {
		err = stats.Create(ctx, &stat)
	} else {
		err = stats.Update(ctx, &stat)
	}
	if err != nil {
		return nil, err
	}
	// Save substats.
	if stat.Substats != nil {
		for _, substat := range *stat.Substats {
			substat.ParentID = stat.ID
			_, err := saveStat(ctx, stats, pkgID, substat)
			if err != nil {
				return nil, err
			}
		}
	}
	return &stat, nil
}

func updateExpiredStats(ctx context.Context, stats *store.Stats, pkgID string) error {
	expired, err := stats.List(ctx, &store.StatOpts{
		PkgIDs:        []string{pkgID},
		IsLatest:      field.NewNullBool(true),
		KindsExcluded: model.ReleaseStatKinds,
	})
	if err != nil {
		return err
	}

	for _, ex := range expired {
		ex.IsLatest = false
		if err := stats.Update(ctx, ex); err != nil {
			return err
		}
	}
	return nil
}

func clearSubstats(ctx context.Context, stats *store.Stats, root *model.Stat) error {
	substats, err := stats.List(ctx, &store.StatOpts{
		ParentIDs: []string{root.ID},
	})
	if err != nil {
		return err
	}

	for _, substat := range substats {
		if substat.HasChildren {
			if err := clearSubstats(ctx, stats, substat); err != nil {
				return err
			}
		}
		if _, err := stats.Delete(ctx, substat.ID); err != nil {
			return err
		}
	}
	return nil
}
