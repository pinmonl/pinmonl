package model

import "github.com/pinmonl/pinmonl/model/field"

type Job struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Describe   string     `json:"describe"`
	TargetID   string     `json:"targetId"`
	TargetName string     `json:"targetName"`
	Status     JobStatus  `json:"status"`
	Message    string     `json:"message"`
	CreatedAt  field.Time `json:"createdAt"`
	EndedAt    field.Time `json:"endedAt"`
}

type JobStatus int

const (
	JobInProgress JobStatus = iota
	JobCompleted
	JobFailed
)

type JobList []*Job
