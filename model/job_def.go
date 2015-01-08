package model

import (
	"encoding/json"
	"fmt"
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

// ToString herps the derp
func (job *JobDef) ToString() string {
	d, err := json.Marshal(&job)
	if err != nil {
		return fmt.Sprintf("Error marshalling job: %s", err.Error())
	}
	return string(d)
}
