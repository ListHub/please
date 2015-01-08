package persistence

import (
	"github.com/ListHub/please/model"
	"github.com/ListHub/please/persistence/etcd"
)

// Load persistence from the environment config
func Load() model.Persistence {
	return etcd.New()
}
