package citadel

import (
	"log"
	"strings"

	"github.com/citadel/citadel"
	"github.com/citadel/citadel/cluster"
	"github.com/citadel/citadel/scheduler"
)

func main() {

	client := getEtcdClient()

	engines := []*citadel.Engine{}

	for _, element := range client.GetCluster() {
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

	my_cluster, err := cluster.New(scheduler.NewResourceManager(), engines...)
	if err != nil {
		log.Fatal(err)
	}
	defer my_cluster.Close()

	if err := my_cluster.RegisterScheduler("service", &scheduler.LabelScheduler{}); err != nil {
		log.Fatal(err)
	}

	if err := my_cluster.Events(&logHandler{}); err != nil {
		log.Fatal(err)
	}

	log.Println(my_cluster.ListContainers(true, true, ""))

	image := &citadel.Image{
		Name:   "crosbymichael/redis",
		Memory: 256,
		Cpus:   0.4,
		Type:   "service",
	}

	container, err := my_cluster.Start(image, false)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("ran container %s\n", container.ID)

	containers := my_cluster.ListContainers(false, false, "")

	c1 := containers[0]

	if err := my_cluster.Kill(c1, 9); err != nil {
		log.Fatal(err)
	}

	if err := my_cluster.Remove(c1); err != nil {
		log.Fatal(err)
	}
}
