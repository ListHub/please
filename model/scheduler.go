package model

// Scheduler is a wrapper for modules used to store job information
type Scheduler interface {
	ScheduleJob(JobDef) error
	ListContainers() ([]ContainerInfo, error)
}
