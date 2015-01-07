package please

import "fmt"

// JobDef defines information needed to launch and run a scheduled container
type JobDef struct {
	Name        string            `json:"name"`
	Schedule    string            `json:"schedule"`
	Image       string            `json:"image"`
	Command     string            `json:"command"`
	Ports       []string          `json:"ports"`
	Environment map[string]string `json:"environment"`
}

// ToString herps the derp
func (job *JobDef) ToString() string {
	return fmt.Sprintf(
		"\tname: %s\n"+
			"\tschedule: %s\n"+
			"\timage: %s\n"+
			"\tcommand: %s\n", job.Name, job.Schedule, job.Image, job.Command)
}
