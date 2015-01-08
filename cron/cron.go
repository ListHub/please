package cron

import (
	"fmt"
	"log"

	"github.com/listhub/please/model"
	"github.com/listhub/please/persistence"
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
	jobs, err := persistence.Get().GetJobs()
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

func getRunJobFn(job model.JobDef) func() {
	return func() {
		log.Printf("scheduling job: %s\n", job.Name)
		scheduler.Get().ScheduleJob(job)
	}
}
