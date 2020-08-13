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
	if err := s.resumePinlUpdated(ctx); err != nil {
		return err
	}
	return nil
}

func (s *ServerRunner) start(ctx context.Context) error {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		s.regularUpdatePkgs(ctx)
		wg.Done()
	}()

	wg.Wait()
	return nil
}

func (s *ServerRunner) resumePinlUpdated(ctx context.Context) error {
	pList, err := s.Stores.Pinls.List(ctx, &store.PinlOpts{
		MonlIDs: []string{""},
	})
	if err != nil {
		return err
	}

	for i := range pList {
		j := job.NewPinlUpdated(pList[i].ID)
		s.Queue.Add(j)
	}
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
