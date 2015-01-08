package citadel

import (
	"github.com/listhub/please/model"
	"log"
	"os"
	"testing"
)

func TestCitadel(t *testing.T) {
	os.Chdir("../..")
	cit := New()

	job := model.JobDef{
		Name:    "Test",
		Image:   "ubuntu",
		Command: "echo ooooh",
		Ports:   []string{"9999:9876"},
	}

	cit.ScheduleJob(job)
	log.Println("Hello")
}
