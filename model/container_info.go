package model

import "time"

// ContainerStatus type used in const block below
type ContainerStatus int

// Container status definitions
const (
	ContainerStatusActive ContainerStatus = iota
	ContainerStatusFinished
)

// ContainerInfo is a wrapper around a job and the details about it's execution
type ContainerInfo struct {
	Job    *JobDef         `json:"job"`
	Status ContainerStatus `json:"status"`
	Start  time.Time       `json:"start"`
	Finish time.Time       `json:"finish"`
}
