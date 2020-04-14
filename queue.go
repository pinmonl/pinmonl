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
		MaxJob:    4,
		MaxWorker: 4,
	})
	if err != nil {
		logx.Panic(err)
	}
	return qm
}

func initQueueScheduler(
	qm *queue.Manager,
	ws *pubsub.Server,
	ml *monler.Repository,
	ss stores,
) *queue.Scheduler {
	sched, err := queue.NewScheduler(&queue.SchedulerOpts{
		QueueManager: qm,
		Pubsub:       ws,
		Monler:       ml,
		Store:        ss.store,
		Pinls:        ss.pinls,
		Monls:        ss.monls,
		Pkgs:         ss.pkgs,
		Stats:        ss.stats,
	})
	if err != nil {
		logx.Panic(err)
	}
	return sched
}
