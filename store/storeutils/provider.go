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
	kind  model.StatKind
	name  string
	value string
}

func getStatKey(stat *model.Stat) statKey {
	key := statKey{
		kind: stat.Kind,
		name: stat.Name,
	}

	if stat.Kind == model.FundingStat {
		key.value = stat.Value
	}

	return key
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
	for i, stat := range rsList {
		// Save stat.
		stat2, err := saveStat(ctx, stats, pkgID, stat)
		if err != nil {
			return nil, err
		}
		out[i] = stat2
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
		deref := *tags[i]
		tag := &deref
		out[i] = tag

		key := reportTagKey{kind: tag.Kind, value: tag.Value}
		prevTag, has := prevTagSet[key]
		if has {
			tag = prevTag
			// Skip if the previous tag is not marked as latest.
			if !tag.IsLatest {
				continue
			}
			// Skip if the previous tag is same as the latest one.
			if _, match := latestKeys[key]; match {
				continue
			}
			// Update to unset latest flag.
			tag.IsLatest = false
		}
		_, err := saveStat(ctx, stats, pkgID, tag)
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

	return pkg, nil
}

func saveStat(ctx context.Context, stats *store.Stats, pkgID string, data *model.Stat) (*model.Stat, error) {
	stat := *data
	stat.PkgID = pkgID
	if stat.Substats != nil && len(*stat.Substats) > 0 {
		stat.HasChildren = true
	}

	// Save.
	err := stats.Create(ctx, &stat)
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
