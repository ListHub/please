package etcd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/coreos/go-etcd/etcd"
	"github.com/listhub/please/model"
)

const historyTTL uint64 = 30 * 24 * 60 * 60
const jobHistoryKey string = "/please/history/jobs/"

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

func (p *persistence) LogContainerStart(jobName, containerID string, startTime time.Time) error {
	containerKey := fmt.Sprintf("/please/history/containers/%s", containerID)
	p.etcdClient.SetDir(containerKey, historyTTL)

	timeStr := startTime.UTC().Format(time.RFC3339)
	p.etcdClient.Set(containerKey+"/start", timeStr, 0)
	p.etcdClient.Set(containerKey+"/job", jobName, 0)

	jobKey := fmt.Sprintf("/please/history/jobs/%s/%s", jobName, timeStr)
	p.etcdClient.SetDir(jobKey, historyTTL)
	p.etcdClient.Set(jobKey+"/status", string(model.JobRunStatusActive), 0)
	p.etcdClient.Set(jobKey+"/container-id", containerID, 0)

	return nil
}

func (p *persistence) lookupStartTimeForContainer(containerID, containerKey string) (string, error) {
	startTimeResp, err := p.etcdClient.Get(containerKey+"/start", false, false)
	if err != nil || startTimeResp == nil || startTimeResp.Node == nil {
		log.Printf("Unable to retrieve start time for "+
			"container-id %s: %s", containerID, err.Error())
		return "", err
	}
	return startTimeResp.Node.Value, nil
}

func (p *persistence) LogContainerFinish(jobName, containerID string, endTime time.Time) error {
	containerKey := fmt.Sprintf("/please/history/containers/%s", containerID)
	endTimeStr := endTime.UTC().Format(time.RFC3339)
	p.etcdClient.Set(containerKey+"/finish", endTimeStr, 0)

	startTimeStr, err := p.lookupStartTimeForContainer(containerID, containerKey)
	if err != nil {
		return err
	}

	jobKey := fmt.Sprintf("/please/history/jobs/%s/%s", jobName, startTimeStr)
	p.etcdClient.Set(jobKey+"/status", string(model.JobRunStatusFinished), 0)
	p.etcdClient.Set(jobKey+"/finish", endTimeStr, 0)

	return nil
}

func (p *persistence) GetJobHistory(start, end time.Time) ([]model.JobRun, error) {
	result := []model.JobRun{}
	resp, err := p.etcdClient.Get(jobHistoryKey, false, true)
	if err != nil || resp == nil || resp.Node == nil {
		return result, errors.New("Unable to pull job history from etcd:" + err.Error())
	}

	for _, jobNode := range resp.Node.Nodes {
		jobName := strings.Replace(jobNode.Key, jobHistoryKey, "", 1)

		for _, runNode := range jobNode.Nodes {
			jobRun := parseJobRun(runNode, jobName)
			result = append(result, jobRun)
		}
	}

	return result, nil
}

func parseJobRun(runNode *etcd.Node, jobName string) model.JobRun {
	jobKey := jobHistoryKey + jobName + "/"
	jobStartStr := strings.Replace(runNode.Key, jobKey, "", 1)

	jobRun := model.JobRun{}
	jobRun.JobName = jobName
	jobRun.Start, _ = time.Parse(time.RFC3339, jobStartStr)
	for _, runAttrNode := range runNode.Nodes {
		if strings.HasSuffix(runAttrNode.Key, "status") {
			jobRun.Status = model.JobRunStatus(runAttrNode.Value)
		}
		if strings.HasSuffix(runAttrNode.Key, "finish") {
			jobRun.Finish, _ = time.Parse(time.RFC3339, runAttrNode.Value)
		}
		if strings.HasSuffix(runAttrNode.Key, "container-id") {
			jobRun.ContainerID = runAttrNode.Value
		}
	}
	return jobRun
}

func (p *persistence) setupWatch() {
	log.Println("Setting up etcd watch")

	respChan := make(chan *etcd.Response)
	go func() {
		for {
			resp := <-respChan
			if resp == nil {
				continue
			}
			logKey := "nil"
			if resp.Node != nil {
				logKey = resp.Node.Key
			}
			log.Printf("New thing on the watch: action: '%s', node.key: '%s'\n",
				resp.Action, logKey)

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
