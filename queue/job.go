package queue

// Job defines the communication payload among manager and workers.
type Job interface {
	Run() error
}

// JobFunc provides a simple wrapper for Job.
type JobFunc func() error

// Run implements Job interface.
func (j JobFunc) Run() error { return j() }
