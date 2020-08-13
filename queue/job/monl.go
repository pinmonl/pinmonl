package job

import (
	"context"
	"time"

	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/monler"
	"github.com/pinmonl/pinmonl/monler/provider"
	"github.com/pinmonl/pinmonl/pkgs/monlutils"
	"github.com/pinmonl/pinmonl/pkgs/pkgdata"
	"github.com/pinmonl/pinmonl/store"
	"github.com/pinmonl/pinmonl/store/storeutils"
	"github.com/sirupsen/logrus"
)

// MonlCrawler defines the job when monl is created.
//
// It finds the monler reports by the url and saves into
// database.
type MonlCrawler struct {
	MonlID   string
	monl     *model.Monl
	reports  map[string]provider.Report
	parentID string
	derived  []string
}

func NewMonlCrawler(monlID string) *MonlCrawler {
	return &MonlCrawler{MonlID: monlID}
}

func (m *MonlCrawler) WithParentID(parentID string) *MonlCrawler {
	m.parentID = parentID
	return m
}

func (m *MonlCrawler) String() string {
	return "monl_crawler"
}

func (m *MonlCrawler) Describe() []string {
	return []string{
		m.String(),
		m.MonlID,
	}
}

func (m *MonlCrawler) Target() model.Morphable {
	return model.Monl{ID: m.MonlID}
}

func (m *MonlCrawler) RunAt() time.Time {
	return time.Time{}
}

func (m *MonlCrawler) PreRun(ctx context.Context) error {
	stores := StoresFrom(ctx)

	// Fetch data of monl.
	monl, err := stores.Monls.Find(ctx, m.MonlID)
	if err != nil {
		return err
	}
	m.monl = monl

	// Convert to repos by the url.
	var repos []provider.Repo
	if monlutils.IsHttp(monl.URL) {
		repos2, err := monler.GuessWithout([]string{pkgdata.GitProvider}, monl.URL)
		if err != nil {
			return err
		}
		repos = repos2
	} else if monler.IsMonler(monl.URL) {
		repo, err := monler.Parse(monl.URL)
		if err != nil {
			return err
		}
		repos = append(repos, repo)
	}

	// Treat as plain website if empty.
	if len(repos) == 0 {
		repo, err := monler.Open(pkgdata.WebsiteProvider, monl.URL)
		if err != nil {
			return err
		}
		repos = append(repos, repo)
	}

	// Get reports.
	reports := make(map[string]provider.Report, 0)
	derived := make([]string, 0)
	for i := range repos {
		repo := repos[i]

		report, err := repo.Analyze()
		if err != nil {
			logrus.Debugf("job monl: %s analyze err(%v)", report, err)
			continue
		}
		rpu, err := report.URI()
		if err != nil {
			continue
		}
		reports[rpu.String()] = report

		// Derived only when there is no parent
		// to avoid recursive derive.
		if m.parentID == "" {
			urls, err := repo.Derived()
			if err != nil {
				continue
			}
			derived = append(derived, urls...)
		}

		repo.Close()
	}

	m.reports = reports
	m.derived = derived
	return nil
}

func (m *MonlCrawler) Run(ctx context.Context) ([]Job, error) {
	stores := StoresFrom(ctx)

	jobs := make([]Job, 0)
	for uri := range m.reports {
		report := m.reports[uri]
		logrus.Debugf("job monl: %q start", report)
		defer report.Close()

		// Create one if monl of the report uri does not exist.
		if m.monl.URL != uri {
			directMonl, isNew, err := storeutils.FindOrCreateMonl(ctx, stores.Monls, uri)
			if err != nil {
				return nil, err
			}
			logrus.Debugf("job monl: %q create direct uri", report)
			if isNew {
				jobs = append(jobs, NewMonlCrawler(directMonl.ID))
			}
		}

		// Save monler report.
		pkg, _, err := storeutils.SaveProviderReport(ctx, stores.Pkgs, stores.Stats, report, false)
		if err != nil {
			return nil, err
		}

		// Find or create the relation between monl and pkg.
		monpkg, err := stores.Monpkgs.FindOrCreate(ctx, &model.Monpkg{
			MonlID: m.MonlID,
			PkgID:  pkg.ID,
		})
		if err != nil {
			return nil, err
		}
		// Set to direct relation.
		monpkg.Kind = model.MonpkgDirect
		err = stores.Monpkgs.Update(ctx, monpkg)
		if err != nil {
			return nil, err
		}

		// Save for derived relation.
		if m.parentID != "" {
			// Search for relation.
			found, err := stores.Monpkgs.List(ctx, &store.MonpkgOpts{
				MonlIDs: []string{m.parentID},
				PkgIDs:  []string{pkg.ID},
			})
			if err != nil {
				return nil, err
			}
			if len(found) > 0 {
				continue
			}

			// Insert as derived if not existed.
			rel := &model.Monpkg{
				MonlID: m.parentID,
				PkgID:  pkg.ID,
				Kind:   model.MonpkgDerived,
			}
			err = stores.Monpkgs.Create(ctx, rel)
			if err != nil {
				return nil, err
			}
		}
		logrus.Debugf("job monl: %q end", report)
	}

	for i := range m.derived {
		u, err := monlutils.NormalizeURL(m.derived[i])
		if err != nil {
			continue
		}
		uri := u.String()

		if _, skip := m.reports[uri]; skip {
			continue
		}

		derivedMonl, isNew, err := storeutils.FindOrCreateMonl(ctx, stores.Monls, uri)
		if err != nil {
			return nil, err
		}
		if isNew {
			jobs = append(jobs, NewMonlCrawler(derivedMonl.ID).WithParentID(m.monl.ID))
		}
	}

	m.monl.FetchedAt = field.Now()
	err := stores.Monls.Update(ctx, m.monl)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

var _ Job = &MonlCrawler{}
