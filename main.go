package main

import (
	"github.com/listhub/please/api"
	"github.com/listhub/please/cron"
)

func main() {

	cron.Start()

	api.ServeAPI()
}
