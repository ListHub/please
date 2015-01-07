package jobs

import (
	"fmt"
	"net/http"

	"github.com/zenazn/goji/web"
)

// GetJobs ...
func GetJobs(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Getting all jobs")
}

// NewJob ...
func NewJob(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Creating new job")
}

// DeleteJob ...
func DeleteJob(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Deleting job %s!", c.URLParams["job_name"])
}

// ReplaceJob ...
func ReplaceJob(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Replacing job %s!", c.URLParams["job_name"])
}

// GetJob ...
func GetJob(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Getting job %s!", c.URLParams["job_name"])
}
