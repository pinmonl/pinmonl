package queue

import (
	"context"
	"sync"
	"time"

	"github.com/pinmonl/pinmonl/logx"
	"github.com/pinmonl/pinmonl/monler"
	"github.com/pinmonl/pinmonl/pubsub"
	"github.com/pinmonl/pinmonl/store"
)

// SchedulerOpts defines the options of creating scheduler.
type SchedulerOpts struct {
	QueueManager *Manager
	Dispatcher   *Dispatcher
	Pubsub       *pubsub.Server
	Monler       *monler.Repository
	Store        store.Store
	Monls        store.MonlStore
	Monpkgs      store.MonpkgStore
	Pinls        store.PinlStore
	Pkgs         store.PkgStore
	Stats        store.StatStore
}

// Scheduler enqueues job with time interval.
type Scheduler struct {
	qm      *Manager
	dp      *Dispatcher
	ws      *pubsub.Server
	monler  *monler.Repository
	store   store.Store
	monls   store.MonlStore
	monpkgs store.MonpkgStore
	pinls   store.PinlStore
	pkgs    store.PkgStore
	stats   store.StatStore

	pinlJobs map[string]int
	monlJobs map[string]int
}

// NewScheduler creates scheduler.
func NewScheduler(opts *SchedulerOpts) (*Scheduler, error) {
	return &Scheduler{
		qm:      opts.QueueManager,
		dp:      opts.Dispatcher,
		ws:      opts.Pubsub,
		monler:  opts.Monler,
		store:   opts.Store,
		pinls:   opts.Pinls,
		monls:   opts.Monls,
		pkgs:    opts.Pkgs,
		stats:   opts.Stats,
		monpkgs: opts.Monpkgs,

		pinlJobs: make(map[string]int),
		monlJobs: make(map[string]int),
	}, nil
}

// Run starts scheduling.
func (s *Scheduler) Run() error {
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		s.resumeSyncPinls()
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

// processMonler get list of monls which stats are not up-to-date.
func (s *Scheduler) processMonler() error {
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
	wg := &sync.WaitGroup{}
	for i := range monls {
		monl := monls[i]
		wg.Add(1)
		go func() {
			s.dp.FetchMonlerReports(monl)
			wg.Done()
		}()
	}

	wg.Wait()
	return nil
}

func (s *Scheduler) resumeSyncPinls() error {
	ctx := context.TODO()
	pinls, err := s.pinls.List(ctx, &store.PinlOpts{
		EmptyMonlOnly: true,
	})
	if err != nil {
		return err
	}
	if len(pinls) == 0 {
		return nil
	}

	logx.Debug("scheduler: resume SyncPinl processes")
	wg := &sync.WaitGroup{}
	for i := range pinls {
		pinl := pinls[i]
		wg.Add(1)
		go func() {
			s.dp.SyncPinl(pinl)
			wg.Done()
		}()
	}

	wg.Wait()
	return nil
}
