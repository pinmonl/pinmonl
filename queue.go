package main

import (
	"github.com/pinmonl/pinmonl/config"
	"github.com/pinmonl/pinmonl/logx"
	"github.com/pinmonl/pinmonl/monler"
	"github.com/pinmonl/pinmonl/pubsub"
	"github.com/pinmonl/pinmonl/queue"
)

func initQueueManager(cfg *config.Config) *queue.Manager {
	qm, err := queue.NewManager(&queue.ManagerOpts{
		MaxJob:    cfg.Queue.MaxJob,
		MaxWorker: cfg.Queue.MaxWorker,
	})
	if err != nil {
		logx.Panic(err)
	}
	return qm
}

func initQueueScheduler(
	qm *queue.Manager,
	dp *queue.Dispatcher,
	ws *pubsub.Server,
	ml *monler.Repository,
	ss stores,
) *queue.Scheduler {
	sched, err := queue.NewScheduler(&queue.SchedulerOpts{
		QueueManager: qm,
		Dispatcher:   dp,
		Pubsub:       ws,
		Monler:       ml,
		Store:        ss.store,
		Pinls:        ss.pinls,
		Monls:        ss.monls,
		Pkgs:         ss.pkgs,
		Stats:        ss.stats,
		Monpkgs:      ss.monpkgs,
	})
	if err != nil {
		logx.Panic(err)
	}
	return sched
}

func initQueueDispatcher(
	qm *queue.Manager,
	ws *pubsub.Server,
	ml *monler.Repository,
	ss stores,
) *queue.Dispatcher {
	dp, err := queue.NewDispatcher(&queue.DispatcherOpts{
		QueueManager: qm,
		Pubsub:       ws,
		Monler:       ml,
		Store:        ss.store,
		Monls:        ss.monls,
		Monpkgs:      ss.monpkgs,
		Pinls:        ss.pinls,
		Pkgs:         ss.pkgs,
		Stats:        ss.stats,
	})
	if err != nil {
		logx.Panic(err)
	}
	return dp
}
