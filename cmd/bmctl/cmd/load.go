package cmd

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"

	"github.com/BondMachineHQ/BMBuildkit/pkg/bondmachined"
	"github.com/BondMachineHQ/BMBuildkit/pkg/build"
	"github.com/spf13/cobra"
)

func loadFirmware(cmd *cobra.Command, args []string) error {
	firmwarePath, err := bondmachined.PullArtifact(args[0], args[1])
	if err != nil {
		return err
	}
	deviceID, err := cmd.Flags().GetString("device")
	if err != nil {
		return err
	}

	fmt.Println(firmwarePath)

	customCMD, err := cmd.Flags().GetString("cmd")
	if err != nil {
		return err
	}

	if customCMD != "" {
		// TODO: security checks?
		cmdStr := customCMD + " " + firmwarePath
		cmd := exec.Command("/bin/sh", "-c", cmdStr)
		stderr, _ := cmd.StderrPipe()
		if err := cmd.Start(); err != nil {
			return err
		}

		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	} else {
		vendor := strings.Split(args[1], "/")[0]
		var engine build.SynthEngine
		if vendor == "lattice" {
			engine = &build.Yosys{}
		} else if vendor == "xilinx" {
			engine = &build.Xilinx{}
		} else {
			return fmt.Errorf("synth engine not available for %s", vendor)
		}

		err = engine.LoadFirmware(firmwarePath, deviceID)
		if err != nil {
			return err
		}
	}

	return nil
}

var loadCmd = &cobra.Command{
	Use:   "load [firmaware registry ref] [platform: vendor/board/variant]",
	Short: "Pull and load a firmware",
	Long:  `Pull firmware artifact from a registry and load it into an FPGA board`,
	RunE:  loadFirmware,
	Args:  cobra.MinimumNArgs(2),
}
