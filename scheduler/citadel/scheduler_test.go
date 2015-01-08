package citadel

import (
	"github.com/listhub/please/model"
	"log"
	"testing"
)

func TestCitadel(t *testing.T) {
	cit := New()

	job := model.JobDef{
		Name:    "Test",
		Image:   "busybox",
		Command: "sleep 1000",
		Ports:   []string{"9999:9876"},
	}

	cit.ScheduleJob(job)
	log.Println("Hello")
}
