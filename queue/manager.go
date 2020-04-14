package queue

import (
	"sync"
	"time"

	"github.com/pinmonl/pinmonl/logx"
)

// ManagerOpts defines the parameters for creating Manager.
type ManagerOpts struct {
	MaxJob    int
	MaxWorker int
}

// Manager operates the job queue and workers.
type Manager struct {
	workers []*worker
	pool    chan chan Job
	queue   chan Job
}

// NewManager creates Manager.
func NewManager(o *ManagerOpts) (*Manager, error) {
	m := &Manager{
		queue:   make(chan Job, o.MaxJob),
		pool:    make(chan chan Job, o.MaxWorker),
		workers: make([]*worker, o.MaxWorker),
	}
	for i := range m.workers {
		w := newWorker(m.pool)
		m.workers[i] = w
	}
	return m, nil
}

// Start starts queue running service.
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
	logx.Debugf("queue: started and %d workers are running", len(m.workers))
	wg.Wait()
	return nil
}

func (m *Manager) start() error {
	for {
		select {
		case job := <-m.queue:
			w := <-m.pool
			w <- job
		}
	}
	return nil
}

// Dispatch adds job to queue.
func (m *Manager) Dispatch(job Job) error {
	return m.dispatch(job, 0)
}

// DispatchAfter adds job to queue after the given duration.
func (m *Manager) DispatchAfter(job Job, after time.Duration) error {
	return m.dispatch(job, after)
}

// DispatchAt adds job to queue at the given time.
func (m *Manager) DispatchAt(job Job, at time.Time) error {
	return m.dispatch(job, time.Until(at))
}

func (m *Manager) dispatch(job Job, after time.Duration) error {
	go func() {
		logx.Debugf("queue: one job is going to enqueue after %q", after)
		time.Sleep(after)
		m.queue <- job
		logx.Debugf("queue: one job is added to queue")
	}()
	return nil
}
