package model

import "time"

// ReloadJobsHandler should be called whenver a job is added or removed
type ReloadJobsHandler func() error

// Persistence is a wrapper for modules used to store job information
type Persistence interface {
	//If job already exists return error with text "Job already exists"
	AddJob(JobDef) error
	DeleteJob(string) error
	GetJobs() ([]JobDef, error)
	GetJob(string) (JobDef, error)
	SetReloadJobsHandler(ReloadJobsHandler) error
	GetServers() ([]string, error)
	LogContainerStart(jobName, containerID string, time time.Time) error
	LogContainerFinish(jobName, containerID string, time time.Time) error
	GetJobHistory(start, end time.Time) ([]JobRun, error)
}
