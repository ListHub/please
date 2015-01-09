package model

// Scheduler is a wrapper for modules used to interface with the execution of
// jobs
type Scheduler interface {
	ScheduleJob(JobDef) error
}
