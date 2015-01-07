package main

import (
	"fmt"

	"github.com/listhub/please/api"
	"github.com/robfig/cron"
)

func main() {

	c := cron.New()
	c.AddFunc("5 * * * * *", func() { fmt.Println("Every five minutes") })
	c.Start()

	api.ServeAPI()
}
