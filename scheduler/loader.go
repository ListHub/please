package scheduler

import (
	"github.com/listhub/please/model"
	"github.com/listhub/please/scheduler/citadel"
)

// Get persistence from the environment config
func Get() model.Scheduler {
	return citadel.New()
}
