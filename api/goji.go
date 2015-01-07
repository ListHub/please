package api

import (
	"github.com/listhub/please/api/v1/jobs"
	"github.com/listhub/please/please"
	"github.com/zenazn/goji"
)

// ServeAPI ...
func ServeAPI(cfg *please.Config) {
	goji.Get("/api/v1/jobs", jobs.GetJobs)
	goji.Post("/api/v1/jobs", jobs.NewJob)
	goji.Delete("/api/v1/jobs/:job_name", jobs.DeleteJob)
	goji.Post("/api/v1/jobs/:job_name", jobs.ReplaceJob)
	goji.Get("/api/v1/jobs/:job_name", jobs.GetJob)

	goji.Serve()
}
