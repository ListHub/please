package cron

import (
	"fmt"

	"github.com/listhub/please/persistence"
	"github.com/listhub/please/please"
	"github.com/listhub/please/scheduler"
	"github.com/robfig/cron"
)

// StartCron will populate the cron system with entries from the persistence
// layer and then start it up
func StartCron() {
	c := cron.New()
	loadExistingJobs(c)
	c.Start()
}

func loadExistingJobs(c *cron.Cron) {
	jobs, err := persistence.Load().GetJobs()
	if err != nil {
		panic("Unable to load jobs from persistence layer")
	}
	for i := 0; i < len(jobs); i++ {
		job := jobs[i]
		err := c.AddFunc(job.Schedule, getRunJobFn(job))
		if err != nil {
			fmt.Printf("Error loading job '%s': %s\n", job.Name, err.Error())
		}
	}
}

func getRunJobFn(job please.JobDef) func() {
	return func() {
		fmt.Printf("scheduling job: %s\n", job.Name)
		scheduler.Load().ScheduleJob(job)
	}
}