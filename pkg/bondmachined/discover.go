package bondmachined

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func DiscoveryHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		log.Errorf("Unsupported method %s", http.MethodPost)
		fmt.Fprintf(w, "Failed to read the pulling request")
		return
	}

	log.Infof("Discovery not supported yet")
	fmt.Fprintf(w, "Discovery not supported yet")
}
