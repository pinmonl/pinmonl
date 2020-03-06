package model

import (
	"database/sql/driver"

	"github.com/pinmonl/pinmonl/model/field"
)

// Job defines the queue task.
type Job struct {
	ID          string     `json:"id"          db:"id"`
	Name        JobName    `json:"name"        db:"name"`
	TargetID    string     `json:"targetId"    db:"target_id"`
	Status      JobStatus  `json:"status"      db:"status"`
	Error       string     `json:"error"       db:"error"`
	ScheduledAt field.Time `json:"scheduledAt" db:"scheduled_at"`
	StartedAt   field.Time `json:"startedAt"   db:"started_at"`
	CreatedAt   field.Time `json:"createdAt"   db:"created_at"`
}

// JobStatus is type of job status
type JobStatus int

// Job status
const (
	// JobStatusNotSet indicates zero value of job status.
	JobStatusNotSet JobStatus = iota
	// JobStatusPending indicates the job is pending.
	JobStatusPending
	// JobStatusQueue indicates the job is added to queue.
	JobStatusQueue
	// JobStatusRunning indicates the job is running.
	JobStatusRunning
	// JobStatusStopped indicates the job is stopped.
	JobStatusStopped
)

var jobStatusSeqs = []string{
	"",
	"pending",
	"queue",
	"running",
	"stopped",
}

// Scan implements sql.Scanner interface.
func (js *JobStatus) Scan(value interface{}) error {
	var vs string
	switch value.(type) {
	case []uint8:
		vs = string(value.([]uint8))
	case string:
		vs = value.(string)
	}
	for i, s := range jobStatusSeqs {
		if s == vs {
			*js = JobStatus(i)
			break
		}
	}
	return nil
}

// String returns the code of job status.
func (js JobStatus) String() string {
	vi := int(js)
	for i, s := range jobStatusSeqs {
		if i == vi {
			return s
		}
	}
	return ""
}

// Value implements driver.Value interface.
func (js JobStatus) Value() (driver.Value, error) {
	return js.String(), nil
}

// JobName is type of job name
type JobName int

// Job name
const (
	// JobNotSet indicates zero value of job name.
	JobNotSet JobName = iota
	// JobPinlCreated defines the job after pinl created.
	JobPinlCreated
	// JobPinlUpdated defines the job after pinl updated.
	JobPinlUpdated
	// JobPkgCreated defines the job after pkg created.
	JobPkgCreated
	// JobPkgRegularUpdate defines the job of pkg regular update.
	JobPkgRegularUpdate
)

var jobNameSeqs = []string{
	"",
	"pinl.created",
	"pinl.updated",
	"pkg.created",
	"pkg.regular_update",
}

// Scan implements sql.Scanner interface.
func (jn *JobName) Scan(value interface{}) error {
	var vs string
	switch value.(type) {
	case []uint8:
		vs = string(value.([]uint8))
	case string:
		vs = value.(string)
	}
	for i, s := range jobNameSeqs {
		if s == vs {
			*jn = JobName(i)
			break
		}
	}
	return nil
}

// String returns the code of job name.
func (jn JobName) String() string {
	vi := int(jn)
	for i, s := range jobNameSeqs {
		if i == vi {
			return s
		}
	}
	return ""
}

// Value implements driver.Value interface.
func (jn JobName) Value() (driver.Value, error) {
	return jn.String(), nil
}
