package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/docker/cli/cli/config/configfile"
	"github.com/docker/cli/cli/config/types"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth/authprovider"
	"github.com/moby/buildkit/util/appcontext"
)

type buildOpt struct {
	withContainerd bool
	containerd     string
	runc           string
}

func main() {
	var opt buildOpt
	flag.BoolVar(&opt.withContainerd, "with-containerd", true, "enable containerd worker")
	flag.StringVar(&opt.containerd, "containerd", "v1.7.2", "containerd version")
	flag.StringVar(&opt.runc, "runc", "v1.1.7", "runc version")
	flag.Parse()

	// // Define the LLB state (e.g., using a scratch image and adding a file)
	state := llb.Scratch().File(llb.Copy(llb.Local("context"), "hello", "/hello"))

	dt, err := state.Marshal(
		context.TODO(),
	)
	if err != nil {
		panic(err)
	}
	llb.WriteTo(dt, os.Stdout)

	ctx := appcontext.Context()
	cli, err := client.New(ctx, "unix:///run/user/1000/buildkit/buildkitd.sock")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating client: %v\n", err)
		os.Exit(1)
	}

	localDir := "/home/dciangot/git/BMBuildkit"

	// Set up the session with credential helper
	sess, err := session.NewSession(ctx, "my-session", "")
	if err != nil {
		panic(err)
	}

	// Attach an auth provider to handle registry credentials
	authProvider := authprovider.NewDockerAuthProvider(&configfile.ConfigFile{
		AuthConfigs: map[string]types.AuthConfig{
			"https://index.docker.io/v1/": {
				Username: "dciangot",
				Password: "Vatice44",
			},
		},
		//CredentialsStore: "pass",
	})
	sess.Allow(authProvider)

	// Solve the LLB definition
	resp, err := cli.Solve(ctx, dt, client.SolveOpt{
		Exports: []client.ExportEntry{
			{
				Type: "image",
				Attrs: map[string]string{
					"name": "docker.io/dciangot/test:v102",
					"push": "true",
					"annotation.org.opencontainers.image.xilinx.alveo": "",
					"annotation.org.opencontainers.image.title":        "easda",
				},
			},
		},

		LocalDirs: map[string]string{
			"context": localDir,
		},
		SharedSession: sess,
	}, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error solving: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(resp.ExporterResponse)
	// Extract the image digest
	dgst := resp.ExporterResponse["containerimage.digest"]
	fmt.Printf("Built image with digest %s\n", dgst)

	// dockerCLI, err := dockerClient.NewClientWithOpts(dockerClient.FromEnv, dockerClient.WithAPIVersionNegotiation())
	// if err != nil {
	// 	fmt.Println("Error creating Docker client:", err)
	// 	os.Exit(1)
	// }

	// // Define the annotation key and value you want to filter by
	// annotationKey := "annotation.org.opencontainers.image.title"
	// annotationValue := "easda"

	// // Create a filter map
	// filters := filters.NewArgs(filters.KeyValuePair{
	// 	Key:   "label",
	// 	Value: fmt.Sprintf("%s=%s", annotationKey, annotationValue),
	// })

	// // List images with the applied filter
	// images, err := dockerCLI.ImageList(ctx, apiTypes.ImageListOptions{Filters: filters})
	// if err != nil {
	// 	fmt.Println("Error listing images:", err)
	// 	os.Exit(1)
	// }

	// fmt.Println("Found images:", images)
	// // Print out the filtered images
	// for _, image := range images {
	// 	fmt.Println("Found image:", image.ID)
	// }
}
