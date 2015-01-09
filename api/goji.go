package api

import (
	"github.com/listhub/please/api/v1/jobs"
	"github.com/listhub/please/api/v1/status"
	"github.com/rs/cors"
	"github.com/zenazn/goji"
)

// ServeAPI ...
func ServeAPI() {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})
	goji.Use(c.Handler)

	goji.Get("/api/v1/jobs", jobs.GetJobs)
	goji.Post("/api/v1/jobs", jobs.NewJob)
	goji.Delete("/api/v1/jobs/:job_name", jobs.DeleteJob)
	goji.Get("/api/v1/jobs/:job_name", jobs.GetJob)
	goji.Get("/api/v1/status", status.GetStatus)

	goji.Serve()
}
