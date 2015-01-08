package model

// NewJobEventHandler is used to be notified when new jobs are added/updated to the
// persistence layer
type NewJobEventHandler func(job JobDef) error

// DeleteJobEventHandler is used to be notified when existing jobs are removed
// from the persistence layer
type DeleteJobEventHandler func(jobName string) error

// Persistence is a wrapper for modules used to store job information
type Persistence interface {
	//If job already exists return error with text "Job already exists"
	AddJob(JobDef) error
	DeleteJob(string) error
	GetJobs() ([]JobDef, error)
	GetJob(string) (JobDef, error)
	SetNewJobEventHandler(NewJobEventHandler) error
	SetDeleteJobEventHandler(DeleteJobEventHandler) error
}
