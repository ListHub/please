package main

import (
	"fmt"

	"github.com/listhub/please/api"
	"github.com/listhub/please/persistence/etcd"
	"github.com/listhub/please/please"
	"github.com/robfig/cron"
)

func main() {

	c := cron.New()
	c.AddFunc("5 * * * * *", func() { fmt.Println("Every five minutes") })
	c.Start()

	api.ServeAPI(LoadConfig())
}

// LoadConfig loads the config from the enviroment
func LoadConfig() *please.Config {
	return &please.Config{
		Persistence: etcd.New(),
	}
}
