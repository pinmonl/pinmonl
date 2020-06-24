package queue

import (
	"context"

	"github.com/pinmonl/pinmonl/exchange"
	"github.com/pinmonl/pinmonl/queue/job"
	"github.com/pinmonl/pinmonl/store"
	"github.com/sirupsen/logrus"
)

type Scheduler struct {
	Queue    *Manager
	Exchange *exchange.Manager
	Stores   *store.Stores
}

func (s *Scheduler) Start() error {
	ctx := context.TODO()
	if err := s.bootstrap(ctx); err != nil {
		return err
	}
	block := make(chan struct{}, 0)
	<-block
	return nil
}

func (s *Scheduler) bootstrap(ctx context.Context) error {
	logrus.Debugln("scheduler: bootstrap")
	if err := s.resumePinlUpdated(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Scheduler) resumePinlUpdated(ctx context.Context) error {
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

func (s *Scheduler) keepExchangeAlive(ctx context.Context) error {
	return nil
}
