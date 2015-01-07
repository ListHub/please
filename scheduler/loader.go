package scheduler

import (
	"github.com/listhub/please/please"
	"github.com/listhub/please/scheduler/citadel"
)

// Load persistence from the environment config
func Load() please.Scheduler {
	return citadel.New()
}
