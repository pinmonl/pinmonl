package model

import (
	"github.com/pinmonl/pinmonl/model/field"
)

// Job defines the queue task.
type Job struct {
	ID          string     `json:"id"          db:"job_id"`
	Name        JobName    `json:"name"        db:"job_name"`
	TargetID    string     `json:"targetId"    db:"job_target_id"`
	Status      JobStatus  `json:"status"      db:"job_status"`
	Error       string     `json:"error"       db:"job_error"`
	ScheduledAt field.Time `json:"scheduledAt" db:"job_scheduled_at"`
	StartedAt   field.Time `json:"startedAt"   db:"job_started_at"`
	CreatedAt   field.Time `json:"createdAt"   db:"job_created_at"`
}

// JobStatus is type of job status.
type JobStatus int

// Job status
const (
	// JobStatusEmpty indicates zero value of job status.
	JobStatusEmpty JobStatus = iota
	// JobStatusPending indicates the job is pending.
	JobStatusPending
	// JobStatusQueue indicates the job is added to queue.
	JobStatusQueue
	// JobStatusRunning indicates the job is running.
	JobStatusRunning
	// JobStatusStopped indicates the job is stopped.
	JobStatusStopped
)

// JobName is type of job name.
type JobName int

// Job name
const (
	// JobEmpty indicates zero value of job name.
	JobEmpty JobName = iota
	// JobPinlCreated defines the job after pinl created.
	JobPinlCreated
	// JobPinlUpdated defines the job after pinl updated.
	JobPinlUpdated
	// JobPkgCreated defines the job after pkg created.
	JobPkgCreated
	// JobPkgRegularUpdate defines the job of pkg regular update.
	JobPkgRegularUpdate
)
