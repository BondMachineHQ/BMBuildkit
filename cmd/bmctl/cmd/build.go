package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"dagger.io/dagger"
	"github.com/BondMachineHQ/BMBuildkit/pkg/build"
	"github.com/spf13/cobra"
)

func buildFirmware(cmd *cobra.Command, args []string) error {
	bmfileName, err := cmd.Flags().GetString("bmfile")
	if err != nil {
		return err
	}
	bmfilePath := args[0] + "/" + bmfileName

	imageRef, err := cmd.Flags().GetString("target")
	if err != nil {
		return err
	}

	platformRef, err := cmd.Flags().GetString("platform")
	if err != nil {
		return err
	}

	fileBytes, err := os.ReadFile(bmfilePath)
	if err != nil {
		return err
	}

	bmfile, err := build.ParseBMFile(fileBytes)
	if err != nil {
		return err
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
			return fmt.Errorf("synth engine not available for %s", bmfile.Vendor)
		}

		build.ExecuteEngine(engine)

		firmwareFilePath, err = engine.GetFirmwareFile()
		if err != nil {
			return err
		}

		buildPath, err = engine.GetBuildDir()
		if err != nil {
			return err
		}
	} else if bmfile.BitstreamPath != "" && bmfile.Synth == "local" {
		firmwareFilePath = bmfile.BitstreamPath
		buildPath = filepath.Dir(bmfile.BitstreamPath)
	} else {
		return fmt.Errorf("synth engine local as to be set along with a BITSTREAM_PATH")
	}

	// create Dagger client
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
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
		return err
	}
	fmt.Println("Pushed multi-platform image w/ digest: ", imageDigest)

	return nil
}

var buildCmd = &cobra.Command{
	Use:   "build [context path (pwd default)]",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	RunE:  buildFirmware,
	Args:  cobra.MinimumNArgs(1),
}
