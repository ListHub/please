package cron

import (
	"errors"
	"log"

	"github.com/listhub/please/model"
	"github.com/listhub/please/persistence"
	"github.com/listhub/please/scheduler"
	"github.com/robfig/cron"
)

// cronic is a wrapper for cron system
type cronic struct {
	currentCron *cron.Cron
}

var cronicSingleton *cronic

func cronicInstance() *cronic {
	if cronicSingleton == nil {
		cronicSingleton = new(cronic)
	}
	return cronicSingleton
}

// Start runing inside the singleton
func Start() {
	cronicInstance().StartCron()
}

// StartCron will populate the cron system with entries from the persistence
// layer and then start it up
func (cronic *cronic) StartCron() {
	persistence.Get().SetReloadJobsHandler(cronic.reloadJobsHandler)

	cronic.currentCron = cron.New()
	loadExistingJobs(cronic.currentCron)
	cronic.currentCron.Start()
}

func (cronic *cronic) reloadJobsHandler() error {
	log.Printf("reloading jobs")
	nextCron := cron.New()
	err := loadExistingJobs(nextCron)
	if err != nil {
		log.Printf("Unable to reload jobs: " + err.Error())
		return err
	}

	nextCron.Start()
	cronic.currentCron.Stop()
	cronic.currentCron = nextCron

	return nil
}

func loadExistingJobs(c *cron.Cron) error {
	jobs, err := persistence.Get().GetJobs()
	if err != nil {
		return errors.New("Unable to load jobs from persistence layer: " + err.Error())
	}
	for i := 0; i < len(jobs); i++ {
		job := jobs[i]
		err := c.AddFunc(job.Schedule, getRunJobFn(job))
		if err != nil {
			log.Printf("Error loading job '%s': %s\n", job.Name, err.Error())
		} else {
			log.Printf("Loaded job '%s' with schedule '%s'\n", job.Name, job.Schedule)
		}
	}
	return nil
}

func getRunJobFn(job model.JobDef) func() {
	return func() {
		log.Printf("scheduling job: %s\n", job.Name)
		err := scheduler.Get().ScheduleJob(job)
		if err != nil {
			log.Printf("unexpected error when scheduling job: %s\n", err.Error())
		}
	}
}
