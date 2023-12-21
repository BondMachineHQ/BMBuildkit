package build

import (
	"fmt"

	"github.com/docker/cli/cli/config/configfile"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth/authprovider"
)

func BuildState(bitstreamPath string) (llb.State, error) {

	return llb.State{}, nil
}

func AuthConfig(sess *session.Session) (session.Attachable, error) {

	return authprovider.NewDockerAuthProvider(&configfile.ConfigFile{}), nil
}

func BuildImages(images BMFile, sess *session.Session) error {

	return fmt.Errorf("NOT IMPLEMENTED YET")
}

func BuildIndex(images []BMImage, sess *session.Session) error {

	return fmt.Errorf("NOT IMPLEMENTED YET")
}

func PushToRegistry(imageRef string) error {

	return fmt.Errorf("NOT IMPLEMENTED YET")
}
