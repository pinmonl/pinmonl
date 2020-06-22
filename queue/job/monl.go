package job

import (
	"context"
	"time"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/monler"
	"github.com/pinmonl/pinmonl/monler/provider"
	"github.com/pinmonl/pinmonl/pkgs/pkguri"
	"github.com/pinmonl/pinmonl/store/storeutils"
)

// MonlCreated defines the job when monl is created.
//
// It finds the monler reports by the url and saves into
// database.
type MonlCreated struct {
	MonlID  string
	reports []provider.Report
}

func NewMonlCreated(monlID string) *MonlCreated {
	return &MonlCreated{
		MonlID: monlID,
	}
}

func (m *MonlCreated) String() string {
	return "monl_created"
}

func (m *MonlCreated) Describe() []string {
	return []string{
		m.String(),
		m.MonlID,
	}
}

func (m *MonlCreated) Target() model.Morphable {
	return model.Monl{ID: m.MonlID}
}

func (m *MonlCreated) RunAt() time.Time {
	return time.Time{}
}

func (m *MonlCreated) PreRun(ctx context.Context) error {
	stores := StoresFrom(ctx)
	monl, err := stores.Monls.Find(ctx, m.MonlID)
	if err != nil {
		return err
	}

	repos, err := monler.GuessWithout([]string{pkguri.GitProvider}, monl.URL)
	if err != nil {
		return err
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

	m.reports = reports
	return nil
}

func (m *MonlCreated) Run(ctx context.Context) ([]Job, error) {
	stores := StoresFrom(ctx)
	for i := range m.reports {
		report := m.reports[i]
		defer report.Close()

		pkg, _, err := storeutils.SaveProviderReport(ctx, stores.Pkgs, stores.Stats, report)
		if err != nil {
			return nil, err
		}

		_, err = stores.Monpkgs.FindOrCreate(ctx, &model.Monpkg{
			MonlID: m.MonlID,
			PkgID:  pkg.ID,
		})
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

var _ Job = &MonlCreated{}
