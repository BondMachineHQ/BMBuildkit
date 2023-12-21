package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
	"github.com/BondMachineHQ/BMBuildkit/pkg/build"
)

func main2() {
	var bmfilePath string
	var imageRef string
	var platformRef string
	flag.StringVar(&bmfilePath, "file", "BMFile", "Path to the BMFile")
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

	var firmwareFilePath string
	var buildPath string
	if bmfile.BitstreamPath == "" && bmfile.Synth != "local" {
		var engine build.SynthEngine

		if bmfile.Vendor == "lattice" {
			engine = &build.Yosys{
				Config: bmfile,
			}
		} else {
			panic(fmt.Errorf("synth engine not available for %s", bmfile.Vendor))
		}

		build.ExecuteEngine(engine)

		firmwareFilePath, err = engine.GetFirmwareFile()
		if err != nil {
			panic(err)
		}

		buildPath, err = engine.GetBuildDir()
		if err != nil {
			panic(err)
		}
	} else if bmfile.BitstreamPath != "" && bmfile.Synth == "local" {
		firmwareFilePath = bmfile.BitstreamPath
		buildPath = filepath.Dir(bmfile.BitstreamPath)
	} else {
		panic(fmt.Errorf("synth engine local as to be set along with a BITSTREAM_PATH"))
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
	}

	platformVariants := make([]*dagger.Container, 0, len(platforms))

	for _, platform := range platforms {
		ctr := client.Container(dagger.ContainerOpts{Platform: platform}).
			//WithLabel("org.opencontainers.image.lattice.ice40", "").
			WithFile("/firmware.bin", project.File(firmwareFilePath))
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

	//bondmachined.PullArtifact(imageRef, platformRef)

}
