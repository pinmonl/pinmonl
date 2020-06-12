package job

import (
	"context"
	"time"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/monler"
	"github.com/pinmonl/pinmonl/monler/provider"
	"github.com/pinmonl/pinmonl/store"
)

type MonlCreated struct {
	MonlID  string
	Monls   *store.Monls
	Pkgs    *store.Pkgs
	Stats   *store.Stats
	Monpkgs *store.Monpkgs
}

func NewMonlCreated(monlID string, monls *store.Monls, pkgs *store.Pkgs, stats *store.Stats, monpkgs *store.Monpkgs) MonlCreated {
	return MonlCreated{
		MonlID:  monlID,
		Monls:   monls,
		Pkgs:    pkgs,
		Stats:   stats,
		Monpkgs: monpkgs,
	}
}

func (m MonlCreated) String() string {
	return "monl_created"
}

func (m MonlCreated) Describe() []string {
	return []string{
		m.String(),
		m.MonlID,
	}
}

func (m MonlCreated) RunAt() time.Time {
	return time.Time{}
}

func (m MonlCreated) Run(ctx context.Context) ([]Job, error) {
	monl, err := m.Monls.Find(ctx, m.MonlID)
	if err != nil {
		return nil, err
	}

	repos, err := monler.Guess(monl.URL)
	if err != nil {
		return nil, err
	}

	reports := make([]provider.Report, 0)
	for i := range repos {
		repo := repos[i]
		report, err := repo.Analyze()
		if err != nil {
			continue
		}
		reports = append(reports, report)
		repo.Close()
	}

	jobs := make([]Job, len(reports))
	for i := range reports {
		report := reports[i]
		job := NewPkgFromReport(report, m.Pkgs, m.Stats, m.onPkgCompleted)
		jobs[i] = job
	}
	return jobs, nil
}

func (m MonlCreated) onPkgCompleted(ctx context.Context, pkg *model.Pkg) error {
	monl, err := m.Monls.Find(ctx, m.MonlID)
	if err != nil {
		return err
	}

	_, err = m.Monpkgs.FindOrCreate(ctx, &model.Monpkg{
		MonlID: monl.ID,
		PkgID:  pkg.ID,
	})
	return err
}

var _ Job = MonlCreated{}
