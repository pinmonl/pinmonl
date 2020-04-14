package queue

type worker struct {
	pool     chan chan Job
	assigned chan Job
}

func newWorker(pool chan chan Job) *worker {
	return &worker{
		pool:     pool,
		assigned: make(chan Job),
	}
}

func (w *worker) start() error {
	for {
		w.pool <- w.assigned
		select {
		case job := <-w.assigned:
			job.Run()
		}
	}
	return nil
}

func (w *worker) stop() error {
	return nil
}
