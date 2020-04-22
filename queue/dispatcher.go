package queue

import (
	"context"
	"fmt"
	"sync"

	"github.com/pinmonl/pinmonl/logx"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/monler"
	"github.com/pinmonl/pinmonl/monler/github"
	"github.com/pinmonl/pinmonl/pubsub"
	"github.com/pinmonl/pinmonl/store"
)

// DispatcherOpts defines the options of creating dispatcher.
type DispatcherOpts struct {
	QueueManager *Manager
	Pubsub       *pubsub.Server
	Monler       *monler.Repository
	Store        store.Store
	Monls        store.MonlStore
	Monpkgs      store.MonpkgStore
	Pinls        store.PinlStore
	Pkgs         store.PkgStore
	Stats        store.StatStore
}

// Dispatcher handles request of job.
type Dispatcher struct {
	qm      *Manager
	ws      *pubsub.Server
	monler  *monler.Repository
	store   store.Store
	monls   store.MonlStore
	monpkgs store.MonpkgStore
	pinls   store.PinlStore
	pkgs    store.PkgStore
	stats   store.StatStore
}

// NewDispatcher creates dispatcher.
func NewDispatcher(opts *DispatcherOpts) (*Dispatcher, error) {
	return &Dispatcher{
		qm:      opts.QueueManager,
		ws:      opts.Pubsub,
		monler:  opts.Monler,
		store:   opts.Store,
		monls:   opts.Monls,
		monpkgs: opts.Monpkgs,
		pinls:   opts.Pinls,
		pkgs:    opts.Pkgs,
		stats:   opts.Stats,
	}, nil
}

// txFunc wraps database transaction.
func (d *Dispatcher) txFunc(ctx context.Context, fn func(context.Context) error) error {
	ctx2, err := d.store.BeginTx(ctx)
	if err != nil {
		return err
	}
	ctx = ctx2

	err = fn(ctx)
	if err != nil {
		if err := d.store.Rollback(ctx); err != nil {
			return err
		}
		return err
	}
	return d.store.Commit(ctx)
}

// SyncPinl updates monl_id of pinl and create monl if not exist.
func (d *Dispatcher) SyncPinl(pinl model.Pinl) error {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	job := JobFunc(func() error {
		defer wg.Done()
		ctx := context.TODO()

		url := monler.URLNormalize(pinl.URL)
		found, err := d.monls.List(ctx, &store.MonlOpts{
			URL: url,
		})
		if err != nil {
			return err
		}

		err = d.txFunc(ctx, func(ctx context.Context) error {
			var monl model.Monl
			if len(found) > 0 {
				monl = found[0]
			} else {
				monl = model.Monl{URL: url}
				err = d.monls.Create(ctx, &monl)
				if err != nil {
					return err
				}
				go d.FetchMonlerReports(monl)
			}

			pinl.MonlID = monl.ID
			err = d.pinls.Update(ctx, &pinl)
			if err != nil {
				return err
			}

			logx.Debugf("dispatcher: %s(%s) is associated to %s", pinl.ID, url, pinl.MonlID)
			return nil
		})

		if err != nil {
			logx.Debugf("dispatcher: err (SyncPinl) %v", err)
			return err
		}
		return nil
	})

	err := d.qm.Dispatch(job)
	wg.Wait()
	return err
}

// FetchMonlerReports downloads and saves monler reports.
func (d *Dispatcher) FetchMonlerReports(monl model.Monl) error {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	job := JobFunc(func() error {
		defer wg.Done()
		ctx := context.TODO()

		// Gather reports
		rs := reports{}
		cred := monler.Credential(nil)
		for _, prd := range d.monler.Providers() {
			if err := d.monler.Ping(prd, monl.URL, cred); err != nil {
				continue
			}
			r, err := d.monler.Open(prd, monl.URL, cred)
			if err != nil {
				continue
			}
			logx.Debugf("dispatcher: monler report %s", textMarshalReport(r))
			rs = append(rs, r)
		}

		// Save report to pkg
		for _, report := range rs {
			err := d.saveMonlerReport(ctx, monl, report, cred, rs)
			if err != nil {
				return err
			}
		}

		return nil
	})

	err := d.qm.Dispatch(job)
	wg.Wait()
	return err
}

// reports provides shorthand functions to walk through each report.
type reports []monler.Report

// has reports whether the report is existed.
func (rs reports) has(report monler.Report) bool {
	for _, r := range rs {
		if r.Provider() == report.Provider() &&
			r.ProviderURI() == report.ProviderURI() {
			return true
		}
	}
	return false
}

// saveMonlerReport saves report and its derived reports into database.
func (d *Dispatcher) saveMonlerReport(ctx context.Context, monl model.Monl, report monler.Report, cred monler.Credential, previousReports reports) error {
	logx.Debugf("dispatcher: %s is downloading", textMarshalReport(report))
	if err := report.Download(); err != nil {
		return err
	}

	pkgs, err := d.pkgs.List(ctx, &store.PkgOpts{
		Provider:    report.Provider(),
		ProviderURI: report.ProviderURI(),
		ListOpts:    store.ListOpts{Limit: 1},
	})
	if err != nil {
		return err
	}

	err = d.txFunc(ctx, func(ctx context.Context) error {
		logx.Debugf("dispatcher: %s begins transaction", textMarshalReport(report))
		var pkg model.Pkg
		if len(pkgs) > 0 {
			pkg = pkgs[0]
		} else {
			pkg = model.Pkg{
				URL:         report.URL(),
				Provider:    report.Provider(),
				ProviderURI: report.ProviderURI(),
			}
			err := d.pkgs.Create(ctx, &pkg)
			if err != nil {
				return err
			}
		}

		err := associateMonpkg(ctx, d.monpkgs, monl, pkg)
		if err != nil {
			return err
		}

		// Save report stats
		if err := saveReportStats(ctx, d.stats, report, pkg); err != nil {
			return err
		}
		if err := saveReportTags(ctx, d.stats, report, pkg); err != nil {
			return err
		}
		if err := d.monls.Update(ctx, &monl); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logx.Debugf("dispatcher: err (FetchMonlerReports) %v", err)
		return err
	}

	derived, err := report.Derived(d.monler, cred)
	if err == nil {
		previous := append(previousReports, report)
		current := append(previous, derived...)

		for _, dReport := range derived {
			// Skip for recusive derived from Github.
			if dReport.Provider() == github.Name && report.Provider() == github.Name {
				continue
			}
			// Skip derived report if already exists.
			if previous.has(dReport) {
				logx.Debugf("dispatcher: skipping derived report %s", textMarshalReport(dReport))
				continue
			}
			logx.Debugf("dispatcher: running derived report %s", textMarshalReport(dReport))
			err := d.saveMonlerReport(ctx, monl, dReport, cred, current)
			if err != nil {
				//
			}
		}
	}

	if err := report.Close(); err != nil {
		return err
	}

	logx.Debugf("dispatcher: done saving of %s", textMarshalReport(report))

	return nil
}

// associateMonpkg associates monl and pkg.
func associateMonpkg(ctx context.Context, monpkgs store.MonpkgStore, monl model.Monl, pkg model.Pkg) error {
	found, err := monpkgs.List(ctx, &store.MonpkgOpts{
		MonlIDs: []string{monl.ID},
		PkgIDs:  []string{pkg.ID},
	})
	if err != nil {
		return err
	}
	if len(found) > 0 {
		return nil
	}
	return monpkgs.Associate(ctx, monl, pkg)
}

// saveReportStats saves report stats.
func saveReportStats(ctx context.Context, stats store.StatStore, report monler.Report, pkg model.Pkg) error {
	oldStats, err := stats.List(ctx, &store.StatOpts{
		PkgID:      pkg.ID,
		WithLatest: true,
	})
	if err != nil {
		return err
	}

	for _, rstat := range report.Stats() {
		stat, err := parseStat(rstat, pkg)
		if err != nil {
			return err
		}

		for _, old := range model.StatList(oldStats).FindKind(string(rstat.Kind)) {
			old.IsLatest = false
			err = stats.Update(ctx, &old)
			if err != nil {
				return err
			}
		}
		stat.IsLatest = true
		err = stats.Create(ctx, stat)
		if err != nil {
			return err
		}
	}
	return nil
}

// saveReportTags saves report tags.
func saveReportTags(ctx context.Context, stats store.StatStore, report monler.Report, pkg model.Pkg) error {
	opts := &store.StatOpts{
		PkgID: pkg.ID,
		Kind:  string(monler.KindTag),
	}
	if report.LatestTag() != nil {
		opts.After = report.LatestTag().RecordedAt.Time()
	}
	oldTags, err := stats.List(ctx, opts)
	if err != nil {
		return err
	}

	for report.Next() {
		tag := report.Tag()
		found := model.StatList(oldTags).FindValue(tag.Value)
		if len(found) > 0 {
			break
		}

		stat, err := parseStat(tag, pkg)
		if err != nil {
			return err
		}
		if ltag := report.LatestTag(); ltag != nil && tag.Value == ltag.Value{
			stat.IsLatest = true
		}
		err = stats.Create(ctx, stat)
		if err != nil {
			return err
		}
	}
	return nil
}

// parseStat parses monler.Stat to model.Stat.
func parseStat(stat *monler.Stat, pkg model.Pkg) (*model.Stat, error) {
	s := model.Stat{
		PkgID:      pkg.ID,
		RecordedAt: stat.RecordedAt,
		Kind:       string(stat.Kind),
		Value:      stat.Value,
		Digest:     stat.Digest,
		Labels:     stat.Labels,
	}
	return &s, nil
}

// textMarshalReport marshals report into plain text.
func textMarshalReport(r monler.Report) string {
	return fmt.Sprintf("%s://%s", r.Provider(), r.ProviderURI())
}
