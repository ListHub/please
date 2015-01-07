package citadel

import "github.com/listhub/please/please"

type scheduler struct{}

func (s *scheduler) ScheduleJob(job please.JobDef) error {
	return nil
}

// New creates an instance of scheduler
func New() please.Scheduler {
	return new(scheduler)
}
