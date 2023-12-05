package build

import (
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

func BuildImages(images []BMImage, sess *session.Session) error {

	return nil
}

func BuildIndex(images []BMImage, sess *session.Session) error {

	return nil
}

func PushToRegistry(imageRef string) error {

	return nil
}
