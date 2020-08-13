package queue

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/pinmonl/pinmonl/database"
	"github.com/pinmonl/pinmonl/exchange"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/pubsub"
	"github.com/pinmonl/pinmonl/queue/job"
	"github.com/pinmonl/pinmonl/store"
	"github.com/sirupsen/logrus"
)

type Manager struct {
	*sync.Mutex
	queue     chan *workerJob
	workers   []worker
	readyPool chan chan *workerJob
	txer      database.Txer
	errchs    map[string][]chan error

	stores   *store.Stores
	exchange *exchange.Manager
	hub      pubsub.Pubsuber
}

func NewManager(txer database.Txer, maxJob, workerN int) (*Manager, error) {
	readyPool := make(chan chan *workerJob, workerN)

	workers := make([]worker, workerN)
	for i := range workers {
		workers[i] = newWorker(readyPool)
	}

	return &Manager{
		Mutex:     &sync.Mutex{},
		txer:      txer,
		queue:     make(chan *workerJob, maxJob),
		readyPool: readyPool,
		workers:   workers,
		errchs:    make(map[string][]chan error),
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
			w <- job
		}
	}
	return nil
}

func (m *Manager) Add(job job.Job) <-chan error {
	ctx := context.TODO()
	has, cherr := m.enqueue(job)
	if has {
		return cherr
	}

	wjob := newWorkerJob(m, job)
	if err := wjob.startRecord(ctx); err != nil {
		cherr <- err
		return cherr
	}

	go func() {
		if time.Now().Before(job.RunAt()) {
			d := time.Until(job.RunAt())
			logrus.Debugf("queue: job %s will be added after %s", job.Describe(), d)
			<-time.After(d)
		}
		logrus.Debugf("queue: added job %s", job.Describe())
		m.queue <- wjob

		// Wait until job is completed.
		err := <-wjob.done
		if qerr := m.dequeue(job, err); qerr != nil {
			logrus.Debugf("queue: dequeue err(%v)", qerr)
		}
	}()
	return cherr
}

func (m *Manager) enqueue(job job.Job) (bool, chan error) {
	m.Lock()
	defer m.Unlock()

	ch := make(chan error)
	k := m.jobKey(job)
	has := false
	if chs, ok := m.errchs[k]; ok && len(chs) > 0 {
		has = true
	}

	m.errchs[k] = append(m.errchs[k], ch)
	return has, ch
}

func (m *Manager) dequeue(job job.Job, err error) error {
	m.Lock()
	defer m.Unlock()
	k := m.jobKey(job)
	if _, ok := m.errchs[k]; !ok {
		return fmt.Errorf("queue: job dequeue not found %s", k)
	}

	for i := range m.errchs[k] {
		chs := m.errchs[k][i]
		defer func() {
			close(chs)
			logrus.Debugf("queue: dequeue channel closed %s", k)
		}()
		select {
		case chs <- err:
		default:
		}
	}
	delete(m.errchs, k)
	return nil
}

func (m *Manager) jobKey(job job.Job) string {
	return strings.Join(job.Describe(), "::")
}

func (m *Manager) Stores(stores *store.Stores) *Manager {
	m.stores = stores
	return m
}

func (m *Manager) ExchangeManager(exm *exchange.Manager) *Manager {
	m.exchange = exm
	return m
}

func (m *Manager) Pubsuber(hub pubsub.Pubsuber) *Manager {
	m.hub = hub
	return m
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
			ctx := context.TODO()
			job.Run(ctx)
		}
	}
	return nil
}

// workerJob handles the returned data after job is done.
type workerJob struct {
	mgr      *Manager
	txer     database.Txer
	stores   *store.Stores
	exchange *exchange.Manager
	hub      pubsub.Pubsuber
	job      job.Job
	record   *model.Job
	done     chan error
}

func newWorkerJob(mgr *Manager, job job.Job) *workerJob {
	return &workerJob{
		mgr:      mgr,
		txer:     mgr.txer,
		stores:   mgr.stores,
		exchange: mgr.exchange,
		hub:      mgr.hub,
		job:      job,
		done:     make(chan error, 1),
	}
}

func (w *workerJob) Run(ctx context.Context) error {
	return w.run(ctx)
}

func (w *workerJob) run(ctx context.Context) error {
	var (
		nexts  []job.Job
		outerr error
	)

	// Inject services
	ctx = job.WithStores(ctx, w.stores)
	ctx = job.WithExchangeManager(ctx, w.exchange)
	ctx = job.WithPubsuber(ctx, w.hub)

	if err := w.job.PreRun(ctx); err != nil {
		w.endRecord(ctx, model.JobFailed, err.Error())
		logrus.Debugf("queue: job err with (%s)", err)
		return err
	}

	txerr := w.txer.TxFunc(ctx, func(ctx context.Context) bool {
		jobs, err := w.job.Run(ctx)
		if err != nil {
			logrus.Debugf("queue: job %s done with err(%s)", w.job.Describe(), err)
			outerr = err
			return false
		}
		nexts = jobs
		logrus.Debugf("queue: job %s done", w.job.Describe())
		return true
	})
	if txerr != nil {
		logrus.Debugf("queue: job %s tx err(%s)", w.job.Describe(), txerr)
		return txerr
	}

	for _, job := range nexts {
		w.mgr.Add(job)
	}

	var (
		status  = model.JobCompleted
		message = ""
	)
	if outerr != nil {
		status, message = model.JobFailed, outerr.Error()
	}
	if err := w.endRecord(ctx, status, message); err != nil {
		return err
	}
	return nil
}

func (w *workerJob) startRecord(ctx context.Context) error {
	record := &model.Job{
		Name:     w.job.String(),
		Describe: strings.Join(w.job.Describe(), ", "),
	}

	if target := w.job.Target(); target != nil {
		record.TargetID = target.MorphKey()
		record.TargetName = target.MorphName()
	}

	// err := w.stores.Jobs.Create(ctx, record)
	// if err != nil {
	// 	return err
	// }
	w.record = record
	return nil
}

func (w *workerJob) endRecord(ctx context.Context, status model.JobStatus, message string) error {
	if w.record == nil {
		return nil
	}

	record := *w.record
	// record.Status = status
	// record.Message = message
	// record.EndedAt = field.Now()
	// err := w.stores.Jobs.Update(ctx, &record)
	// if err != nil {
	// 	w.done <- err
	// 	return err
	// }
	w.done <- nil
	w.record = &record
	return nil
}
