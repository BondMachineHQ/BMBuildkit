package bondmachined

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/BondMachineHQ/BMBuildkit/pkg/build"
	log "github.com/sirupsen/logrus"
)

// LoadHandler pull and loads the firmware from the registry to the board
func LoadHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		log.Errorf("Unsupported method %s", http.MethodPost)
		fmt.Fprintf(w, "Failed to read the pulling request")
		return
	}

	var bodyBytes []byte
	_, err := r.Body.Read(bodyBytes)
	if err != nil {
		log.Errorf("failed to read the request: %s", err)
		fmt.Fprintf(w, "Failed to read the pulling request")
		return
	}

	var img PullArtifactRequest
	err = json.Unmarshal(bodyBytes, &img)
	if err != nil {
		log.Errorf("Invalid pull artifact request: %s", err)
		fmt.Fprintf(w, "Invalid pull artifact request")
		return
	}

	firmwarePath, err := PullArtifact(img.ArtifactName, img.BoardModel)
	if err != nil {
		log.Errorf("Error during pulling of the image: %s", err)
		fmt.Fprintf(w, "Error during pulling of the image")
		return
	}

	var engine *build.Yosys
	if img.BoardModel == "lattice" {
		engine = &build.Yosys{}
	} else {
		log.Errorf("synth engine not available for %s", img.BoardModel)
		fmt.Fprintf(w, "synth engine not available")
		return
	}

	err = engine.LoadFirmware(firmwarePath, "")
	if err != nil {
		log.Errorf("Error during loading firmware: %s", err)
		fmt.Fprintf(w, "Error during loading firmware")
		return
	}
	//log.Infof("Discovery not supported yet")
	//fmt.Fprintf(w, "Discovery not supported yet")
}
