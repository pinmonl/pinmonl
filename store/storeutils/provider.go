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
) (*model.Pkg, model.StatList, error) {
	pkg, err := findOrCreatePkgFromReport(ctx, pkgs, report)
	if err != nil {
		return nil, nil, err
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

func saveReportStats(ctx context.Context, stats *store.Stats, pkgID string, report provider.Report) (model.StatList, error) {
	rsList, err := report.Stats()
	if err != nil {
		return nil, err
	}

	out := make([]*model.Stat, len(rsList))
	for i := range rsList {
		stat := *rsList[i]
		// If stat is latest, update previous latest stats as expired.
		if stat.IsLatest {
			err := updateExpiredStats(ctx, stats, pkgID, &stat)
			if err != nil {
				return nil, err
			}
		}
		// Save stat.
		s2, err := saveStat(ctx, stats, pkgID, &stat)
		if err != nil {
			return nil, err
		}
		out[i] = s2
	}
	return out, nil
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
		PkgIDs: []string{pkgID},
		Kind:   field.NewNullValue(model.TagStat),
	})
	if err != nil {
		return nil, err
	}

	var latest *model.Stat
	if ts := tags.GetLatest(); len(ts) > 0 {
		latest = ts[0]
	}

	// Save tags.
	out := make([]*model.Stat, len(tags))
	for i := range tags {
		tag := *tags[i]
		preVals := prevTags.GetValue(tag.Value)
		// If tag is already existed.
		if len(preVals) > 0 {
			tag = *preVals[0]
			out[i] = &tag
			// Skip if the previous tag is not marked as latest.
			if !tag.IsLatest {
				continue
			}
			// Skip if the previous tag is same as the latest one.
			if tag.Value == latest.Value {
				continue
			}
			// Update to unset latest flag.
			tag.IsLatest = false
		}
		_, err := saveStat(ctx, stats, pkgID, &tag)
		if err != nil {
			return nil, err
		}
		out[i] = &tag
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

func updateExpiredStats(ctx context.Context, stats *store.Stats, pkgID string, stat *model.Stat) error {
	expired, err := stats.List(ctx, &store.StatOpts{
		PkgIDs:   []string{pkgID},
		Kind:     field.NewNullValue(stat.Kind),
		IsLatest: field.NewNullBool(true),
	})
	if err != nil {
		return err
	}

	for _, ex := range expired {
		ex.IsLatest = false
		err := stats.Update(ctx, ex)
		if err != nil {
			return err
		}
	}
	return nil
}
