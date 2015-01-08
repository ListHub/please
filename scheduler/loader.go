package scheduler

import (
	"github.com/listhub/please/model"
	"github.com/listhub/please/scheduler/citadel"
)

var scheduler model.Scheduler

// Get persistence from the environment config
func Get() model.Scheduler {
	if scheduler == nil {
		scheduler = citadel.New()
	}
	return scheduler
}
