package job

import (
	"context"
	"time"

	"github.com/pinmonl/pinmonl/monler"
	"github.com/pinmonl/pinmonl/monler/provider"
	"github.com/pinmonl/pinmonl/store"
)

type MonlCreated struct {
	MonlID string
	Monls  *store.Monls
	Pkgs   *store.Pkgs
	Stats  *store.Stats
}

func NewMonlCreated(monlID string, monls *store.Monls, pkgs *store.Pkgs, stats *store.Stats) MonlCreated {
	return MonlCreated{
		MonlID: monlID,
		Monls:  monls,
		Pkgs:   pkgs,
		Stats:  stats,
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
		job := NewPkgFromReport(reports[i], m.Pkgs, m.Stats)
		jobs = append(jobs, job)
	}
	return jobs, nil
}

var _ Job = MonlCreated{}
