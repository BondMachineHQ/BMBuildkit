package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/BondMachineHQ/BMBuildkit/pkg/build"
	"github.com/docker/cli/cli/config/configfile"
	"github.com/docker/cli/cli/config/types"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth/authprovider"
	"github.com/moby/buildkit/util/appcontext"
)

// "github.com/BondMachineHQ/BMBuildkit/pkg/build"
// "github.com/BondMachineHQ/BMBuildkit/pkg/image"
// "github.com/BondMachineHQ/BMBuildkit/pkg/load"

type buildOpt struct {
	withContainerd bool
	containerd     string
	runc           string
}

func main() {
	var addr string
	var bmfilePath string
	flag.StringVar(&addr, "addr", "", "server socket addr")
	flag.StringVar(&bmfilePath, "file", "BMFile", "server socket addr")
	flag.Parse()

	fileBytes, err := os.ReadFile(bmfilePath)
	if err != nil {
		panic(err)
	}

	bmfile, err := build.ParseBMFile(fileBytes)
	if err != nil {
		panic(err)
	}

	fmt.Println(bmfile)

	var engine build.SynthEngine

	if bmfile.Vendor == "lattice" {
		engine = &build.Yosys{
			Config: bmfile,
		}
	}

	build.ExecuteEngine(engine)

	firmwareFile, err := engine.GetFirmwareFile()
	if err != nil {
		panic(err)
	}

	buildPath, err := engine.GetBuildDir()
	if err != nil {
		panic(err)
	}

	// // Define the LLB state (e.g., using a scratch image and adding a file)
	state := llb.Scratch().File(llb.Copy(llb.Local("context"), firmwareFile, "/firmware.bin"))

	dt, err := state.Marshal(
		context.TODO(),
	)
	if err != nil {
		panic(err)
	}
	llb.WriteTo(dt, os.Stdout)

	ctx := appcontext.Context()

	//cli, err := client.New(ctx, "unix:///run/user/1000/buildkit/buildkitd.sock")
	cli, err := client.New(ctx, addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating client: %v\n", err)
		os.Exit(1)
	}

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
			"context": buildPath,
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

	// Image references
	base := empty.Index
	image1Ref := "docker.io/dciangot/test:v102"
	image2Ref := "docker.io/dciangot/test:v101"

	// // Annotation for image 1
	// image1Annotations := map[string]string{
	// 	"author":  "John Doe",
	// 	"purpose": "Backend Service",
	// }

	// // Annotation for image 2
	// image2Annotations := map[string]string{
	// 	"author":  "Jane Smith",
	// 	"purpose": "Frontend Service",
	// }

	img1ref, err := name.ParseReference(image1Ref)
	if err != nil {
		panic(err)
	}

	img1, err := remote.Get(img1ref)
	if err != nil {
		panic(err)
	}

	img1GOOD, err := img1.Image()
	if err != nil {
		panic(err)
	}

	img2ref, err := name.ParseReference(image2Ref)
	if err != nil {
		panic(err)
	}

	img2, err := remote.Get(img2ref)
	if err != nil {
		panic(err)
	}

	img2GOOD, err := img2.Image()
	if err != nil {
		panic(err)
	}

	// Create a new index image
	idx := mutate.AppendManifests(base,
		mutate.IndexAddendum{Add: img2GOOD, Descriptor: v1.Descriptor{
			Annotations: map[string]string{
				"foo": "bar",
			},
		}},
		mutate.IndexAddendum{Add: img1GOOD},
	)

	// Push the index image to a registry
	dest := "dciangot/index_image:latest"
	tag, err := name.NewTag(dest, name.WeakValidation)
	if err != nil {
		panic(err)
	}

	err = remote.WriteIndex(tag, idx, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Index image %s created and pushed successfully.\n", dest)

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
