package bondmachined

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-containerregistry/pkg/crane"
	cranev1 "github.com/google/go-containerregistry/pkg/v1"
	log "github.com/sirupsen/logrus"
)

func PullArtifact(imageName string, platform string) (string, error) {

	splitPlatform := strings.Split(platform, "/")

	vendor := splitPlatform[0]

	board := splitPlatform[1]

	variant := splitPlatform[2]

	pltf := cranev1.Platform{
		OS:           vendor,
		Architecture: board,
		Variant:      variant,
	}

	img, err := crane.Pull(imageName, crane.WithPlatform(&pltf))
	if err != nil {
		return "", err
	}

	fo, err := os.CreateTemp("", "tmpfile-") // in Go version older than 1.17 you can use ioutil.TempFile
	if err != nil {
		return "", err
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	// make a write buffer
	w := bufio.NewWriter(fo)

	crane.Export(img, w)

	return fo.Name(), nil
}

func PullHandler(w http.ResponseWriter, r *http.Request) {

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

	log.Infof("Received request to pull %s", img.ArtifactName)
	log.Infof("Stored artifact at %s", firmwarePath)
	fmt.Fprintf(w, "Pulled artifact.")
}
