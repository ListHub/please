package persistence

import (
	"github.com/listhub/please/model"
	"github.com/listhub/please/persistence/etcd"
)

var persistence model.Persistence

// Get persistence from the environment config
func Get() model.Persistence {
	if persistence == nil {
		persistence = etcd.New()
	}
	return persistence
}
