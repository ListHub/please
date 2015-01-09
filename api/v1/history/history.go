package history

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/listhub/please/persistence"
	"github.com/zenazn/goji/web"
)

// GetContainerHistory returns info about running and previously runnign containers
func GetContainerHistory(c web.C, w http.ResponseWriter, r *http.Request) {
	//TODO: parse these times from the request
	start := time.Now().UTC().AddDate(0, -1, 0)
	end := time.Now().UTC()
	containers, err := persistence.Get().GetJobHistory(start, end)

	daters, err := json.Marshal(containers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unable to marshal status to json: %s", err.Error())
		log.Printf("Unable to marshal status to json: %s\n", err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Write(daters)
}
