package citadel

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/citadel/citadel"
	"github.com/citadel/citadel/cluster"
	"github.com/citadel/citadel/scheduler"

	"github.com/listhub/please/model"
	"github.com/listhub/please/persistence"
)

type sched struct{}

var myCluster *cluster.Cluster

func (s *sched) ScheduleJob(job model.JobDef) error {

	ports := []*citadel.Port{}
	for _, element := range job.Ports {
		log.Println(element)
		strPort := strings.Split(element, ":")
		toPort, _ := strconv.ParseInt(strPort[0], 10, 0)
		fromPort, _ := strconv.ParseInt(strPort[1], 10, 0)
		port := &citadel.Port{Port: int(fromPort), ContainerPort: int(toPort)}
		ports = append(ports, port)
	}

	image := &citadel.Image{
		Name:        job.Image,
		Memory:      256,
		Cpus:        0.4,
		Type:        "service",
		Environment: job.Environment,
		BindPorts:   ports,
		Args:        strings.Split(job.Command, " "),
	}

	container, err := myCluster.Start(image, true)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("ran container %s\n", container.ID)

	myCluster.ListContainers(false, false, "")

	return nil
}

// New creates an instance of scheduler
func New() model.Scheduler {
	engines := engines()

	var err error
	myCluster, err = cluster.New(scheduler.NewResourceManager(), engines...)
	if err != nil {
		log.Fatal(err)
	}

	if err := myCluster.RegisterScheduler("service", &scheduler.LabelScheduler{}); err != nil {
		log.Fatal(err)
	}

	if err := myCluster.Events(&logHandler{}); err != nil {
		log.Fatal(err)
	}

	return new(sched)
}

type logHandler struct {
}

func (l *logHandler) Handle(e *citadel.Event) error {
	log.Printf("type: %s time: %s image: %s container: %s\n",
		e.Type, e.Time.Format(time.RubyDate), e.Container.Image.Name, e.Container.ID)

	return nil
}

func list() {
	log.Println(myCluster.ListContainers(true, true, ""))
}

func engines() []*citadel.Engine {
	engines := []*citadel.Engine{}
	servers, _ := persistence.Get().GetServers()
	for _, element := range servers {
		address := strings.TrimSuffix(element, "4001") + "2375"
		log.Println(address)
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
