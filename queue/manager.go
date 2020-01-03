package queue

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pinmonl/pinmonl/logx"
	"github.com/pinmonl/pinmonl/model"
	"github.com/pinmonl/pinmonl/model/field"
	"github.com/pinmonl/pinmonl/monl"
	"github.com/pinmonl/pinmonl/store"
)

// ManagerOpts defines the dependencies of manager.
type ManagerOpts struct {
	Parallel int
	Interval time.Duration
	Monl     *monl.Monl

	Store store.Store
	Jobs  store.JobStore
	Pinls store.PinlStore
	Monls store.MonlStore
	Stats store.StatStore
}

// Manager operates the job queue and workers.
type Manager struct {
	sync.Mutex

	interval time.Duration
	parallel int
	jobQueue chan *Job
	workers  []*worker

	store store.Store
	jobs  store.JobStore
}

// NewManager creates manager and workers.
func NewManager(opts ManagerOpts) *Manager {
	m := &Manager{
		interval: opts.Interval,
		parallel: opts.Parallel,
		jobQueue: make(chan *Job, opts.Parallel),
		workers:  make([]*worker, opts.Parallel),

		store: opts.Store,
		jobs:  opts.Jobs,
	}

	for i := range m.workers {
		m.workers[i] = &worker{
			monl:    opts.Monl,
			manager: m,

			store: opts.Store,
			pinls: opts.Pinls,
			monls: opts.Monls,
			stats: opts.Stats,
		}
	}

	return m
}

// Start initiates the goroutines of manager and workers.
func (m *Manager) Start(ctx context.Context) error {
	for _, w := range m.workers {
		go w.run(ctx)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(m.interval):
			if len(m.jobQueue) == 0 {
				logx.Debug("queue manager: updating")
				m.updateQueue(ctx)
			}
		}
	}
	return ctx.Err()
}

// Enqueue pushes job into queue.
func (m *Manager) Enqueue(ctx context.Context, j *model.Job) error {
	m.Lock()
	defer m.Unlock()

	j2 := *j
	j2.Status = model.JobStatusPending
	err := m.jobs.Create(ctx, &j2)
	if err != nil {
		return err
	}
	*j = j2
	return nil
}

// updateQueue updates the job queue by interval.
//
// When there is any pending jobs, its status will change
// from "pending" to "queue" and added to the queue.
func (m *Manager) updateQueue(ctx context.Context) error {
	m.Lock()
	defer m.Unlock()

	jl, err := m.jobs.List(ctx, &store.JobOpts{
		Status:          model.JobStatusPending,
		ScheduledBefore: time.Now(),
		ListOpts:        store.ListOpts{Limit: int64(m.parallel)},
	})
	if err != nil {
		return err
	}
	if len(jl) > 0 {
		ctx, err = m.store.BeginTx(ctx)
		if err != nil {
			logx.Fatalf("queue manager: fails to start tx, err: %s", err)
			return err
		}
		tmpQ := make([]*Job, len(jl))
		for i, j := range jl {
			job, err := m.parseJobAndQueued(ctx, &j)
			if err != nil {
				m.store.Rollback(ctx)
				return err
			}
			tmpQ[i] = job
		}
		m.store.Commit(ctx)
		for _, j := range tmpQ {
			m.jobQueue <- j
		}
		logx.Debugf("queue manager: %d job(s) added to queue", len(jl))
		return nil
	}

	return nil
}

// parseJobAndQueued is shorthand func for wrapping model.Job into Job
// and update status to "queue".
func (m *Manager) parseJobAndQueued(ctx context.Context, job *model.Job) (*Job, error) {
	j := *job
	j.Status = model.JobStatusQueue
	err := m.jobs.Update(ctx, &j)
	if err != nil {
		return nil, err
	}
	*job = j
	return &Job{Detail: j}, nil
}

func (m *Manager) jobStarted(ctx context.Context, job *Job) error {
	m.Lock()
	defer m.Unlock()

	if job.Detail.Status != model.JobStatusQueue {
		return fmt.Errorf("job status does not match")
	}

	jd := job.Detail
	jd.Status = model.JobStatusRunning
	jd.StartedAt = (field.Time)(time.Now())

	err := m.jobs.Update(ctx, &jd)
	if err != nil {
		return err
	}
	job.Detail = jd
	return nil
}

func (m *Manager) jobCompleted(ctx context.Context, job *Job) error {
	m.Lock()
	defer m.Unlock()

	jd := job.Detail
	jd.Status = model.JobStatusNotSet

	err := m.jobs.Update(ctx, &jd)
	if err != nil {
		return err
	}
	job.Detail = jd
	return nil
}

func (m *Manager) jobStopped(ctx context.Context, job *Job) error {
	m.Lock()
	defer m.Unlock()

	jd := job.Detail
	jd.Status = model.JobStatusStopped
	jd.Error = job.Error.Error()

	err := m.jobs.Update(ctx, &jd)
	if err != nil {
		return err
	}
	job.Detail = jd
	return nil
}
