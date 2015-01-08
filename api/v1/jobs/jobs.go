package jobs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ListHub/please/model"
	"github.com/ListHub/please/persistence"
	"github.com/zenazn/goji/web"
)

// GetJobs ...
func GetJobs(c web.C, w http.ResponseWriter, r *http.Request) {
	jobs, err := persistence.Load().GetJobs()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to pull jobs from persistence layer: %s", err.Error())
	}

	daters, err := json.Marshal(jobs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to marshal jobs to json: %s", err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Write(daters)
}

// NewJob ...
func NewJob(c web.C, w http.ResponseWriter, r *http.Request) {
	job := new(model.JobDef)
	daters, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unable to read request body: %s", err.Error())
		return
	}
	err = json.Unmarshal(daters, &job)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unable to parse job: %s", err.Error())
		return
	}

	//TODO: Validate job contents

	err = persistence.Load().AddJob(*job)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "unable to persist job: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "sucessfully added job: \n %s", job.ToString())
}

// DeleteJob ...
func DeleteJob(c web.C, w http.ResponseWriter, r *http.Request) {
	err := persistence.Load().DeleteJob(c.URLParams["job_name"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "unable to persist job: %s", err.Error())
	}
	w.WriteHeader(http.StatusOK)
}

// GetJob ...
func GetJob(c web.C, w http.ResponseWriter, r *http.Request) {
	jobName := c.URLParams["job_name"]
	job, err := persistence.Load().GetJob(jobName)
	if err != nil {
		if strings.Contains(err.Error(), "Key not found") {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "No job found with name %s", jobName)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to pull job from persistence layer: %s", err.Error())
	}

	daters, err := json.Marshal(job)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to marshal job to json: %s", err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Write(daters)
}
