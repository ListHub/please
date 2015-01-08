package status

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/listhub/please/scheduler"
	"github.com/zenazn/goji/web"
)

// GetStatus returns info about running and previously runnign containers
func GetStatus(c web.C, w http.ResponseWriter, r *http.Request) {
	containers, err := scheduler.Get().ListContainers()

	daters, err := json.Marshal(containers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to marshal status to json: %s", err.Error())
		log.Printf("Unable to marshal status to json: %s\n", err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Write(daters)
}
