package queue

import (
	"context"
	"sync"
	"time"

	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/queue/job"
	"github.com/sirupsen/logrus"
)

type Manager struct {
	queue     chan job.Job
	workers   []worker
	readyPool chan chan *workerJob
	txer      database.Txer
}

func NewManager(txer database.Txer, maxJob, workerN int) (*Manager, error) {
	readyPool := make(chan chan *workerJob, workerN)

	workers := make([]worker, workerN)
	for i := range workers {
		workers[i] = newWorker(readyPool)
	}

	return &Manager{
		txer:      txer,
		queue:     make(chan job.Job, maxJob),
		readyPool: readyPool,
		workers:   workers,
	}, nil
}

func (m *Manager) Start() error {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		m.start()
		wg.Done()
	}()

	for _, w := range m.workers {
		wg.Add(1)
		go func() {
			w.start()
			wg.Done()
		}()
	}

	wg.Wait()
	return nil
}

func (m *Manager) start() error {
	for {
		select {
		case job := <-m.queue:
			w := <-m.readyPool
			w <- newWorkerJob(m, job)
		}
	}
	return nil
}

func (m *Manager) Add(job job.Job) error {
	go func() {
		if time.Now().Before(job.RunAt()) {
			d := time.Until(job.RunAt())
			logrus.Debugf("queue: job %s will be added after %s", job.Describe(), d)
			<-time.After(d)
		}
		logrus.Debugf("queue: added job %s", job.Describe())
		m.queue <- job
	}()
	return nil
}

// worker is the job runner of the queue.
type worker struct {
	readyPool chan chan *workerJob
	assigned  chan *workerJob
}

func newWorker(readyPool chan chan *workerJob) worker {
	return worker{
		readyPool: readyPool,
		assigned:  make(chan *workerJob),
	}
}

func (w worker) start() error {
	for {
		w.readyPool <- w.assigned
		select {
		case job := <-w.assigned:
			ctx := context.Background()
			job.Run(ctx)
		}
	}
	return nil
}

// workerJob handles the returned data after job is done.
type workerJob struct {
	mgr *Manager
	job job.Job
}

func newWorkerJob(mgr *Manager, job job.Job) *workerJob {
	return &workerJob{
		mgr: mgr,
		job: job,
	}
}

func (w *workerJob) Run(ctx context.Context) error {
	return w.run(ctx)
}

func (w *workerJob) run(ctx context.Context) error {
	err := w.mgr.txer.TxFunc(ctx, func(ctx context.Context) bool {
		jobs, err := w.job.Run(ctx)
		if err != nil {
			logrus.Debugf("queue: job %s done with err(%s)", w.job.Describe(), err)
			return false
		}
		logrus.Debugf("queue: job %s done", w.job.Describe())

		for _, job := range jobs {
			w.mgr.Add(job)
		}
		return true
	})
	return err
}
