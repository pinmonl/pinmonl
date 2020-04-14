package queue

import (
	"context"
	"sync"
	"time"

	"github.com/pinmonl/pinmonl/logx"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/monler"
	"github.com/pinmonl/pinmonl/pubsub"
	"github.com/pinmonl/pinmonl/store"
)

type SchedulerOpts struct {
	QueueManager *Manager
	Pubsub       *pubsub.Server
	Monler       *monler.Repository
	Store        store.Store
	Pinls        store.PinlStore
	Monls        store.MonlStore
	Pkgs         store.PkgStore
	Stats        store.StatStore
}

type Scheduler struct {
	qm      *Manager
	ws      *pubsub.Server
	monler  *monler.Repository
	txStore store.Store
	pinls   store.PinlStore
	monls   store.MonlStore
	pkgs    store.PkgStore
	stats   store.StatStore

	pinlJobs map[string]int
	monlJobs map[string]int
}

func NewScheduler(opts *SchedulerOpts) (*Scheduler, error) {
	return &Scheduler{
		qm:      opts.QueueManager,
		ws:      opts.Pubsub,
		monler:  opts.Monler,
		txStore: opts.Store,
		pinls:   opts.Pinls,
		monls:   opts.Monls,
		pkgs:    opts.Pkgs,
		stats:   opts.Stats,

		pinlJobs: make(map[string]int),
		monlJobs: make(map[string]int),
	}, nil
}

func (s *Scheduler) Run() error {
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		s.processPinls()
		for {
			select {
			case <-time.After(time.Minute):
				s.processPinls()
			}
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		s.processMonler()
		for {
			select {
			case <-time.After(time.Minute):
				s.processMonler()
			}
		}
		wg.Done()
	}()

	wg.Wait()
	return nil
}

func (s *Scheduler) processPinls() error {
	if len(s.pinlJobs) > 0 {
		return nil
	}

	ctx := context.TODO()
	count, err := s.pinls.Count(ctx, &store.PinlOpts{
		EmptyMonlOnly: true,
	})
	if err != nil {
		return err
	}
	if count == int64(0) {
		return nil
	}

	logx.Debugf("scheduler: process n=%d pinl", count)

	size := int64(50)
	list := make([]model.Pinl, 0)
	for i := int64(0); i < count; i += size {
		pinls, err := s.pinls.List(ctx, &store.PinlOpts{
			EmptyMonlOnly: true,
			ListOpts:      store.ListOpts{Limit: size, Offset: i},
		})
		if err != nil {
			return err
		}
		list = append(list, pinls...)
	}

	for _, pinl := range list {
		s.pinlJobs[pinl.ID]++
		s.qm.Dispatch(s.newProcessPinlJob(pinl.ID))
	}

	return nil
}

func (s *Scheduler) newProcessPinlJob(pinlID string) Job {
	return JobFunc(func() error {
		ctx := context.TODO()
		if c, err := s.txStore.BeginTx(ctx); err != nil {
			return err
		} else {
			ctx = c
		}
		logx.Debug("begin tx of pinl job")

		run := func() error {
			pinl := model.Pinl{ID: pinlID}
			if err := s.pinls.Find(ctx, &pinl); err != nil {
				return err
			}

			url := monler.URLNormalize(pinl.URL)
			found, err := s.monls.List(ctx, &store.MonlOpts{
				URL: url,
			})
			if err != nil {
				return err
			}

			var monl model.Monl
			if len(found) > 0 {
				monl = found[0]
			} else {
				monl = model.Monl{URL: url}
				err = s.monls.Create(ctx, &monl)
				s.qm.Dispatch(s.newProcessMonlerJob(monl.ID))
				if err != nil {
					return err
				}
			}

			pinl.MonlID = monl.ID
			err = s.pinls.Update(ctx, &pinl)
			if err != nil {
				return err
			}

			logx.Debugf("scheduler: pinl \"%s %s\" is associated to %q", pinl.ID, url, pinl.MonlID)
			return nil
		}

		if err := run(); err != nil {
			logx.Debugf("scheduler job: process pinl err %v", err)
			s.txStore.Rollback(ctx)
			return err
		}
		delete(s.pinlJobs, pinlID)
		s.txStore.Commit(ctx)
		return nil
	})
}

func (s *Scheduler) processMonler() error {
	if len(s.monlJobs) > 0 {
		return nil
	}

	ctx := context.TODO()
	before := time.Now().Add(-time.Hour * 24)
	monls, err := s.monls.List(ctx, &store.MonlOpts{
		UpdatedBefore: before,
		ListOpts:      store.ListOpts{Limit: 50},
	})
	if err != nil {
		return err
	}

	logx.Debugf("scheduler: process n=%d monl", len(monls))
	for _, monl := range monls {
		s.monlJobs[monl.ID]++
		s.qm.Dispatch(s.newProcessMonlerJob(monl.ID))
	}

	return nil
}

func (s *Scheduler) newProcessMonlerJob(monlID string) Job {
	return JobFunc(func() error {
		ctx := context.TODO()
		if c, err := s.txStore.BeginTx(ctx); err != nil {
			return err
		} else {
			ctx = c
		}

		run := func() error {
			monl := model.Monl{ID: monlID}
			err := s.monls.Find(ctx, &monl)
			if err != nil {
				return err
			}

			// Gather reports
			reps := make([]monler.Report, 0)
			cred := monler.Credential(nil)
			for _, prd := range s.monler.Providers() {
				if err := s.monler.Ping(prd, monl.URL, cred); err != nil {
					continue
				}
				if rep, err := s.monler.Open(prd, monl.URL, cred); err != nil {
					continue
				} else {
					logx.Debugf("scheduler job: monler report \"%s::%s\"", rep.Provider(), rep.ProviderURI())
					reps = append(reps, rep)
				}
			}

			// Save report to pkg
			for _, rep := range reps {
				if err := rep.Download(); err != nil {
					return err
				}

				pkgs, err := s.pkgs.List(ctx, &store.PkgOpts{
					Provider:    rep.Provider(),
					ProviderURI: rep.ProviderURI(),
					ListOpts:    store.ListOpts{Limit: 1},
				})
				if err != nil {
					return err
				}
				var pkg model.Pkg
				if len(pkgs) > 0 {
					pkg = pkgs[0]
				} else {
					pkg = model.Pkg{
						URL:         rep.URL(),
						Provider:    rep.Provider(),
						ProviderURI: rep.ProviderURI(),
					}
					err := s.pkgs.Create(ctx, &pkg)
					if err != nil {
						return err
					}
				}

				// Save report stats
				if err := saveReportStats(ctx, rep, pkg, s.stats); err != nil {
					return err
				}
				if err := saveReportTags(ctx, rep, pkg, s.stats); err != nil {
					return err
				}
				if err := s.monls.Update(ctx, &monl); err != nil {
					return err
				}
			}
			return nil
		}

		if err := run(); err != nil {
			logx.Debugf("scheduler job: process monler err %v", err)
			s.txStore.Rollback(ctx)
			return err
		}
		delete(s.monlJobs, monlID)
		s.txStore.Commit(ctx)
		return nil
	})
}

func saveReportStats(ctx context.Context, rep monler.Report, pkg model.Pkg, stats store.StatStore) error {
	oldStats, err := stats.List(ctx, &store.StatOpts{
		PkgID:      pkg.ID,
		WithLatest: true,
	})
	if err != nil {
		return err
	}

	for _, rstat := range rep.Stats() {
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

func saveReportTags(ctx context.Context, rep monler.Report, pkg model.Pkg, stats store.StatStore) error {
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
