package citadel

import (
	"log"
	"time"

	"github.com/citadel/citadel"

	"github.com/coreos/go-etcd/etcd"
	"github.com/listhub/please/model"
)

type sched struct{}

func (s *sched) ScheduleJob(job model.JobDef) error {
	return nil
}

// New creates an instance of scheduler
func New() model.Scheduler {
	return new(sched)
}

type persistence struct {
}

var etcdClient *etcd.Client

func getEtcdClient() *etcd.Client {
	if etcdClient == nil {
		etcdClient = etcd.NewClient([]string{"http://10.0.10.144:4001", "http://10.0.10.84:4001", "http://10.0.10.248:4001"})
	}
	return etcdClient
}

type logHandler struct {
}

func (l *logHandler) Handle(e *citadel.Event) error {
	log.Printf("type: %s time: %s image: %s container: %s\n",
		e.Type, e.Time.Format(time.RubyDate), e.Container.Image.Name, e.Container.ID)

	return nil
}
