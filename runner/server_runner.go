package runner

import (
	"context"

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
	block := make(chan struct{}, 0)
	<-block
	return nil
}

func (s *ServerRunner) bootstrap(ctx context.Context) error {
	logrus.Debugln("runner: bootstrap")
	if err := s.resumePinlUpdated(ctx); err != nil {
		return err
	}
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

var _ Runner = &ServerRunner{}
