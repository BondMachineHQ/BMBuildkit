package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"dagger.io/dagger"
	"github.com/BondMachineHQ/BMBuildkit/pkg/build"
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
	var imageRef string
	var platformRef string
	flag.StringVar(&addr, "addr", "", "server socket addr")
	flag.StringVar(&bmfilePath, "file", "BMFile", "server socket addr")
	flag.StringVar(&imageRef, "target", "", "image name")
	flag.StringVar(&platformRef, "platform", "lattice/ice40/yosys", "platform name [prod/type/board]")
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

	// create Dagger client
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	project := client.Host().Directory(buildPath)

	var platforms = []dagger.Platform{
		dagger.Platform(platformRef),
		"lattice/ice40-6k/yosys",
		// "linux/arm64", // a.k.a. aarch64
		// "linux/s390x", // a.k.a. IBM S/390
	}

	platformVariants := make([]*dagger.Container, 0, len(platforms))

	for _, platform := range platforms {
		ctr := client.Container(dagger.ContainerOpts{Platform: platform}).
			WithLabel("org.opencontainers.image.lattice.ice40", "").
			WithFile("/firmware.bin", project.File(firmwareFile))
		platformVariants = append(platformVariants, ctr)
	}

	imageDigest, err := client.
		Container().
		Publish(ctx, imageRef, dagger.ContainerPublishOpts{
			PlatformVariants: platformVariants,
			// Some registries may require explicit use of docker mediatypes
			// rather than the default OCI mediatypes
			// MediaTypes: dagger.Dockermediatypes,
		})
	if err != nil {
		panic(err)
	}
	fmt.Println("Pushed multi-platform image w/ digest: ", imageDigest)

}
