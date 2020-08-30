package runner

import (
	"context"
	"sync"
	"time"

	"github.com/pinmonl/pinmonl/queue"
	"github.com/pinmonl/pinmonl/queue/job"
	"github.com/pinmonl/pinmonl/store"
	"github.com/sirupsen/logrus"
)

type ServerRunner struct {
	Queue  *queue.Manager
	Stores *store.Stores
}

func (s *ServerRunner) Start() error {
	ctx := context.TODO()
	if err := s.bootstrap(ctx); err != nil {
		return err
	}

	if err := s.start(ctx); err != nil {
		return err
	}
	return nil
}

func (s *ServerRunner) bootstrap(ctx context.Context) error {
	logrus.Debugln("runner: bootstrap")
	return nil
}

func (s *ServerRunner) start(ctx context.Context) error {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		s.regularUpdateMonls(ctx)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		s.regularUpdatePkgs(ctx)
		wg.Done()
	}()

	wg.Wait()
	return nil
}

func (s *ServerRunner) regularUpdateMonls(ctx context.Context) error {
	interval := 24 * time.Hour
	ticker := time.NewTicker(interval)
	defer func() {
		ticker.Stop()
	}()
	s.updateMonls(ctx, time.Now().Add(-1*interval))
	for {
		select {
		case <-ticker.C:
			before := time.Now().Add(-1 * interval)
			s.updateMonls(ctx, before)
		}
	}
}

func (s *ServerRunner) updateMonls(ctx context.Context, before time.Time) error {
	logrus.Debugln("runner: start monls update")
	expired, err := s.Stores.Monls.List(ctx, &store.MonlOpts{
		FetchedBefore: before,
	})
	if err != nil {
		return err
	}

	for _, monl := range expired {
		s.Queue.Add(job.NewMonlCrawler(monl.ID))
	}
	logrus.Debugf("runner: %d monls updated", len(expired))
	return nil
}

func (s *ServerRunner) regularUpdatePkgs(ctx context.Context) error {
	interval := 8 * time.Hour
	ticker := time.NewTicker(interval)
	defer func() {
		ticker.Stop()
	}()
	s.updatePkgs(ctx, time.Now().Add(-1*interval))
	for {
		select {
		case <-ticker.C:
			before := time.Now().Add(-1 * interval)
			s.updatePkgs(ctx, before)
		}
	}
}

func (s *ServerRunner) updatePkgs(ctx context.Context, before time.Time) error {
	logrus.Debugf("runner: start pkgs update")
	expired, err := s.Stores.Pkgs.List(ctx, &store.PkgOpts{
		FetchedBefore: before,
	})
	if err != nil {
		return err
	}

	for _, pkg := range expired {
		s.Queue.Add(job.NewPkgCrawler(pkg.ID))
	}
	logrus.Debugf("runner: %d pkgs updated", len(expired))
	return nil
}

var _ Runner = &ServerRunner{}
