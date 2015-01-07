package etcd

import (
	"encoding/json"

	"github.com/coreos/go-etcd/etcd"
	"github.com/listhub/please/please"
)

type persistence struct {
}

var etcdClient *etcd.Client

func getEtcdClient() *etcd.Client {
	if etcdClient == nil {
		etcdClient = etcd.NewClient([]string{"http://10.0.10.144:4001", "http://10.0.10.84:4001", "http://10.0.10.248:4001"})
	}
	return etcdClient
}

// AddJob ...
func (p *persistence) AddJob(job please.JobDef) error {
	client := getEtcdClient()
	//TODO: Make sure the job doesn't already exist
	//TODO: Make sure the job name is a valid etcd name
	jobDaters, err := json.Marshal(job)
	if err != nil {
		return err
	}
	jobStr := string(jobDaters)
	_, err = client.RawSet("/please/jobs/"+job.Name, jobStr, 0)
	return err
}

// DeleteJob ...
func (p *persistence) DeleteJob(jobName string) error {
	client := getEtcdClient()
	//TODO: Make sure the job exists
	_, err := client.RawDelete("/please/jobs/"+jobName, true, true)
	return err
}

// GetJobs ...
func (p *persistence) GetJobs() ([]please.JobDef, error) {
	client := getEtcdClient()
	resp, err := client.Get("/please/jobs/", false, true)
	if err != nil {
		return []please.JobDef{}, err
	}

	jobs := []please.JobDef{}
	for i := 0; i < len(resp.Node.Nodes); i++ {
		job := new(please.JobDef)
		err = json.Unmarshal([]byte(resp.Node.Nodes[i].Value), &job)
		if err != nil {
			return jobs, err
		}
		jobs = append(jobs, *job)
	}

	return jobs, nil
}

// GetJob ..
func (p *persistence) GetJob(jobName string) (please.JobDef, error) {
	client := getEtcdClient()
	resp, err := client.Get("/please/jobs/"+jobName, false, true)
	if err != nil {
		return please.JobDef{}, err
	}

	job := new(please.JobDef)
	err = json.Unmarshal([]byte(resp.Node.Value), &job)
	return *job, err
}

// New creates an instance of persistence
func New() please.Persistence {
	return new(persistence)
}
