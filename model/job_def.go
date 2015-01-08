package model

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/robfig/cron"
)

// JobDef defines information needed to launch and run a scheduled container
type JobDef struct {
	Name              string            `json:"name"`
	Schedule          string            `json:"schedule"`
	Image             string            `json:"image"`
	Command           string            `json:"command"`
	Ports             []string          `json:"ports"`
	Memory            float64           `json:"memory"`
	CPU               float64           `json:"cpu"`
	Environment       map[string]string `json:"environment"`
	OverrunPrevention bool              `json:"overrun-prevention"`
}

// Validate the jobDef
func (job *JobDef) Validate() []string {

	errorMsgs := []string{}
	if len(job.Name) < 1 {
		errorMsgs = append(errorMsgs, "Job must have a name")
	}
	if len(job.Image) < 1 {
		errorMsgs = append(errorMsgs, "Job must have an image")
	}
	if job.Memory == 0 {
		errorMsgs = append(errorMsgs, "Job must specifiy memory requirements in MB")
	}
	if job.CPU == 0 {
		errorMsgs = append(errorMsgs, "Job must specify CPU requirements")
	}
	if !validateCronSpec(job.Schedule) {
		errorMsgs = append(errorMsgs, "Job must have a valid cron string schedule")
	}

	return errorMsgs
}

// ToString herps the derp
func (job *JobDef) ToString() string {
	d, err := json.Marshal(&job)
	if err != nil {
		return fmt.Sprintf("Error marshalling job: %s", err.Error())
	}
	return string(d)
}

func validateCronSpec(spec string) bool {
	_, err := cron.Parse(spec)
	if err != nil {
		log.Printf("Failure to validate cronspec: %s\n", err.Error())
		return false
	}
	return true
}
