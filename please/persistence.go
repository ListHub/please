package please

// Persistence is a wrapper for modules used to store job information
type Persistence interface {
	AddJob(JobDef) error
	DeleteJob(string) error
	GetJobs() ([]JobDef, error)
	GetJob(string) (JobDef, error)
}
