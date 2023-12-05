package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/docker/cli/cli/config/configfile"
	"github.com/docker/cli/cli/config/types"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth/authprovider"
	"github.com/moby/buildkit/util/appcontext"
)

func startBuildKitDaemon(c client.Client) error {

	fmt.Println("BuildKit daemon started")
	// Start the BuildKit daemon
	if err := c.Wait(context.Background()); err != nil {
		return fmt.Errorf("error starting BuildKit daemon: %s", err)
	}

	return nil
}

func main() {
	// Replace this with your BuildKit daemon socket path
	socketPathd := "unix:////var/run/docker/buildkit/buildkitd.sock"

	// Create a context for the BuildKit client
	ctx := appcontext.Context()

	// Create a BuildKit client
	cli, err := client.New(ctx, socketPathd)
	if err != nil {
		log.Panicf("error creating BuildKit client: %s", err)
	}
	// Start BuildKit daemon
	go startBuildKitDaemon(*cli)

	// // Define the LLB state (e.g., using a scratch image and adding a file)
	state := llb.Scratch().File(llb.Copy(llb.Local("context"), "hello", "/hello"))

	dt, err := state.Marshal(
		context.TODO(),
	)
	if err != nil {
		panic(err)
	}
	llb.WriteTo(dt, os.Stdout)

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

	// // Start listening on a socket
	// socketPath := "/tmp/buildkit_daemon.sock" // Replace with your desired socket path
	// l, err := net.Listen("unix", socketPath)
	// if err != nil {
	// 	log.Fatalf("error creating socket listener: %s", err)
	// }
	// defer l.Close()

	// fmt.Printf("BuildKit Daemon listening on %s\n", socketPath)

	// // Accept incoming connections
	// for {
	// 	conn, err := l.Accept()
	// 	if err != nil {
	// 		log.Fatalf("error accepting connection: %s", err)
	// 	}

	// 	// Handle incoming connections concurrently
	// 	go handleConnection(conn)
	// }
}
