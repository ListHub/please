package etcd

import "github.com/listhub/please/please"

type persistence struct {
}

// AddJob ...
func (p *persistence) AddJob(job please.JobDef) error {
	return nil
}

// DeleteJob ...
func (p *persistence) DeleteJob(jobName string) error {
	return nil
}

// GetJobs ...
func (p *persistence) GetJobs() ([]please.JobDef, error) {
	return nil, nil
}

// GetJob ..
func (p *persistence) GetJob(jobName string) (please.JobDef, error) {
	return please.JobDef{}, nil
}

// New creates an instance of persistence
func New() please.Persistence {
	return new(persistence)
}
