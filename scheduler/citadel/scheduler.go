package citadel

import (
	"errors"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/citadel/citadel"
	"github.com/citadel/citadel/cluster"
	"github.com/citadel/citadel/scheduler"

	"github.com/listhub/please/model"
	"github.com/listhub/please/persistence"
)

type sched struct {
	myCluster *cluster.Cluster
}

func containerName(job model.JobDef) string {
	return job.Name + "_" + time.Now().UTC().Format("2006-01-02T150405Z")
}

func parseJobNameFromContainerName(containerName string) string {
	re, _ := regexp.Compile("_[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{6}Z")
	return strings.Replace(re.ReplaceAllString(containerName, ""), "/", "", 1)
}

func (s *sched) ScheduleJob(job model.JobDef) error {

	image := &citadel.Image{
		ContainerName: containerName(job),
		Name:          job.Image,
		Memory:        job.Memory,
		Cpus:          job.CPU,
		Type:          "service",
		Environment:   job.Environment,
		BindPorts:     parsePorts(job.Ports),
		Args:          strings.Split(job.Command, " "),
	}

	container, err := s.myCluster.Start(image, true)
	if err != nil {
		return err
	}

	log.Printf("ran job %s in container %s\n", job.Name, container.ID)

	return nil
}

func (s *sched) init() error {
	engines := engines()

	var err error
	s.myCluster, err = cluster.New(scheduler.NewResourceManager(), engines...)
	if err != nil {
		return errors.New("Unable to instantiate cluster: " + err.Error())
	}

	err = s.myCluster.RegisterScheduler("service", &scheduler.LabelScheduler{})
	if err != nil {
		return errors.New("Unable to register scheduler: " + err.Error())
	}

	err = s.myCluster.Events(s)
	if err != nil {
		return errors.New("Unable to register for events: " + err.Error())
	}

	return nil
}

// New creates an instance of scheduler
func New() model.Scheduler {

	s := new(sched)
	err := s.init()
	if err != nil {
		log.Printf("Error initializing scheduler: %ss", err.Error())
	}

	return s
}

func (s *sched) Handle(e *citadel.Event) error {
	jobName := parseJobNameFromContainerName(e.Container.Name)
	if e.Type == "start" {
		persistence.Get().LogContainerStart(jobName, e.Container.ID, e.Time)
	}

	if e.Type == "die" {
		persistence.Get().LogContainerFinish(jobName, e.Container.ID, e.Time)
	}

	log.Printf("citadel event - job: %s type: %s time: %s name: %s container: %s\n",
		jobName, e.Type, e.Time.Format(time.RubyDate), e.Container.Name, e.Container.ID)

	return nil
}

func (s *sched) list() {
	log.Println(s.myCluster.ListContainers(true, true, ""))
}

func parsePorts(jobPorts []string) []*citadel.Port {
	ports := []*citadel.Port{}
	for _, element := range jobPorts {
		strPort := strings.Split(element, ":")
		toPort, _ := strconv.ParseInt(strPort[0], 10, 0)
		fromPort, _ := strconv.ParseInt(strPort[1], 10, 0)
		port := &citadel.Port{Port: int(fromPort), ContainerPort: int(toPort)}
		ports = append(ports, port)
	}
	return ports
}

func engines() []*citadel.Engine {
	engines := []*citadel.Engine{}
	servers, _ := persistence.Get().GetServers()
	for _, element := range servers {
		address := strings.TrimSuffix(element, "4001") + "2375"
		engine := &citadel.Engine{
			ID:     address,
			Addr:   address,
			Memory: 2048,
			Cpus:   4,
			Labels: []string{"local"},
		}

		engines = append(engines, engine)

		if err := engine.Connect(nil); err != nil {
			log.Fatal(err)
		}
	}
	return engines
}
