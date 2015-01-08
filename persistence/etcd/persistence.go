package etcd

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/coreos/go-etcd/etcd"
	"github.com/listhub/please/model"
)

type persistence struct {
	reloadJobsHandler model.ReloadJobsHandler
	etcdClient        *etcd.Client
}

func newClient() *etcd.Client {
	return etcd.NewClient([]string{
		model.Config().PersistenceAddress,
	})
}

// AddJob ...
func (p *persistence) AddJob(job model.JobDef) error {
	//TODO: Make sure the job name is a valid etcd name
	jobDaters, err := json.Marshal(job)
	if err != nil {
		return err
	}
	jobStr := string(jobDaters)
	_, err = p.etcdClient.RawSet("/please/jobs/"+job.Name, jobStr, 0)
	return err
}

// DeleteJob ...
func (p *persistence) DeleteJob(jobName string) error {
	//TODO: Make sure the job exists
	_, err := p.etcdClient.RawDelete("/please/jobs/"+jobName, true, true)
	return err
}

// GetJobs ...
func (p *persistence) GetJobs() ([]model.JobDef, error) {
	resp, err := p.etcdClient.Get("/please/jobs/", false, true)
	if err != nil {
		if strings.Contains(err.Error(), "Key not found (/please)") {
			return []model.JobDef{}, nil
		}
		return []model.JobDef{}, err
	}

	jobs := []model.JobDef{}
	for i := 0; i < len(resp.Node.Nodes); i++ {
		job := new(model.JobDef)
		err = json.Unmarshal([]byte(resp.Node.Nodes[i].Value), &job)
		if err != nil {
			return jobs, err
		}
		jobs = append(jobs, *job)
	}

	return jobs, nil
}

// GetJob ..
func (p *persistence) GetJob(jobName string) (model.JobDef, error) {
	resp, err := p.etcdClient.Get("/please/jobs/"+jobName, false, true)
	if err != nil {
		return model.JobDef{}, err
	}

	job := new(model.JobDef)
	err = json.Unmarshal([]byte(resp.Node.Value), &job)
	return *job, err
}

func (p *persistence) SetReloadJobsHandler(handler model.ReloadJobsHandler) error {
	p.reloadJobsHandler = handler
	return nil
}

func (p *persistence) GetServers() ([]string, error) {
	return p.etcdClient.GetCluster(), nil
}

func (p *persistence) setupWatch() {
	log.Println("Setting up etcd watch")

	respChan := make(chan *etcd.Response)
	go func() {
		for {
			resp := <-respChan
			logAction := "nil"
			if resp != nil {
				logAction = resp.Action
			}
			logKey := "nil"
			if resp != nil && resp.Node != nil {
				logKey = resp.Node.Key
			}
			log.Printf("New thing on the watch: action: '%s', node.key: '%s'\n",
				logAction, logKey)

			if p.reloadJobsHandler != nil {
				p.reloadJobsHandler()
			}
		}
	}()
	go p.etcdClient.Watch("/please/jobs/", 0, true, respChan, nil)
	log.Println("watch is setup")
}

// New creates an instance of persistence
func New() model.Persistence {
	p := new(persistence)
	p.etcdClient = newClient()
	p.setupWatch()
	return p
}
