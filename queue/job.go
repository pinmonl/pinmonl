package queue

import "github.com/pinmonl/pinmonl/model"

// Job defines the communication payload among manager and workers.
type Job struct {
	Detail model.Job
	Error  error
}
