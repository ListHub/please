package persistence

import (
	"github.com/listhub/please/persistence/etcd"
	"github.com/listhub/please/please"
)

// Load persistence from the environment config
func Load() please.Persistence {
	return etcd.New()
}
