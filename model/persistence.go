package model

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
}
