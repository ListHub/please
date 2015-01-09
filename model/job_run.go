package model

import "time"

// JobRunStatus type used in const block below
type JobRunStatus string

// Container status definitions
const (
	JobRunStatusActive   JobRunStatus = "active"
	JobRunStatusFinished              = "finished"
)

// JobRun is a wrapper around a job and the details about it's execution
type JobRun struct {
	JobName     string       `json:"job-name"`
	Status      JobRunStatus `json:"status"`
	ContainerID string       `json:"container-id"`
	Start       time.Time    `json:"start"`
	Finish      time.Time    `json:"finish,omitempty"`
}
