package job

import (
	"context"
	"time"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/monler"
	"github.com/pinmonl/pinmonl/monler/provider"
	"github.com/pinmonl/pinmonl/pkgs/pkguri"
	"github.com/pinmonl/pinmonl/store"
)

type PkgFromReport struct {
	Report provider.Report
	Pkgs   *store.Pkgs
	Stats  *store.Stats
}

func NewPkgFromReport(report provider.Report, pkgs *store.Pkgs, stats *store.Stats) PkgFromReport {
	return PkgFromReport{
		Report: report,
		Pkgs:   pkgs,
		Stats:  stats,
	}
}

func (p PkgFromReport) String() string {
	return "pkg_from_report"
}

func (p PkgFromReport) Describe() []string {
	return []string{
		p.String(),
		p.Report.String(),
	}
}

func (p PkgFromReport) RunAt() time.Time {
	return time.Time{}
}

func (p PkgFromReport) Run(ctx context.Context) ([]Job, error) {
	defer p.Report.Close()

	pu, err := p.Report.URI()
	if err != nil {
		return nil, err
	}
	pkg, err := p.Pkgs.FindURI(ctx, pu)
	if err != nil {
		return nil, err
	}

	stats, err := p.Report.Stats()
	if err != nil {
		return nil, err
	}
	for _, stat := range stats {
		err := p.saveStat(ctx, pkg.ID, stat)
		if err != nil {
			return nil, err
		}
	}

	for p.Report.Next() {
		tag, err := p.Report.Tag()
		if err != nil {
			return nil, err
		}
		err = p.saveStat(ctx, pkg.ID, tag)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (p PkgFromReport) saveStat(ctx context.Context, pkgID string, stat *model.Stat) error {
	stat.PkgID = pkgID
	if stat.Substats != nil && len(*stat.Substats) > 0 {
		stat.HasChildren = true
	}

	err := p.Stats.Create(ctx, stat)
	if err != nil {
		return err
	}

	if stat.Substats != nil {
		for _, substat := range *stat.Substats {
			substat.ParentID = stat.ID
			err := p.saveStat(ctx, pkgID, substat)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type PkgSelfUpdate struct {
	PkgID string
	Pkgs  *store.Pkgs
	Stats *store.Stats
}

func NewPkgSelfUpdate(pkgID string, pkgs *store.Pkgs, stats *store.Stats) PkgSelfUpdate {
	return PkgSelfUpdate{
		PkgID: pkgID,
		Pkgs:  pkgs,
		Stats: stats,
	}
}

func (p PkgSelfUpdate) String() string {
	return "pkg_self_update"
}

func (p PkgSelfUpdate) Describe() []string {
	return []string{
		p.String(),
		p.PkgID,
	}
}

func (p PkgSelfUpdate) RunAt() time.Time {
	return time.Time{}
}

func (p PkgSelfUpdate) Run(ctx context.Context) ([]Job, error) {
	pkg, err := p.Pkgs.Find(ctx, p.PkgID)
	if err != nil {
		return nil, err
	}

	pu, err := pkg.MarshalPkgURI()
	if err != nil {
		return nil, err
	}

	err = monler.Ping(pu.Provider, pkguri.ToURL(pu).String())
	if err != nil {
		return nil, err
	}

	repo, err := monler.Parse(pu.String())
	if err != nil {
		return nil, err
	}

	report, err := repo.Analyze()
	if err != nil {
		return nil, err
	}
	return []Job{NewPkgFromReport(report, p.Pkgs, p.Stats)}, nil
}

var _ Job = PkgFromReport{}
var _ Job = PkgSelfUpdate{}
