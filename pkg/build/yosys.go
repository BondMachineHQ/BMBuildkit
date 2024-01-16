package build

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Yosys struct {
	Config     BMFile
	BuildDir   string
	ModuleName string
}

// LoadData (only when we are going to implemente remote build)uncompress the source folder sent by the BMctl to the bondmachined
func (t *Yosys) LoadData() error {

	dir, err := os.MkdirTemp("/tmp", ".yosys-build")
	if err != nil {
		panic(err)
	}

	t.BuildDir = dir

	t.ModuleName = strings.Split(t.Config.SynthFile, ".v")[0]

	fmt.Println(t.ModuleName)

	source := t.Config.Source
	dest := dir

	err = CopyDir(source, dest)
	if err != nil {
		panic(err)
	}

	return nil
}

func (t *Yosys) Synth() error {

	fmt.Println(t.Config.SynthCMD)

	if t.Config.SynthCMD == "" {
		cmdStr := "cd " + t.BuildDir + "; yosys -p 'synth_ice40 -top " + t.ModuleName + " -json " + t.ModuleName + ".json' " + t.Config.SynthFile
		cmd := exec.Command("/bin/sh", "-c", cmdStr)
		stderr, _ := cmd.StderrPipe()
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}

		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}

	return nil
}

func (t *Yosys) PNR() error {

	cmdStr := "cd " + t.BuildDir + "; nextpnr-ice40 --hx1k --json " + t.ModuleName + ".json --pcf " + t.ModuleName + ".pcf --asc " + t.ModuleName + ".asc"
	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	return nil
}

func (t *Yosys) Packing() error {

	cmdStr := "cd " + t.BuildDir + "; icepack " + t.ModuleName + ".asc " + t.ModuleName + ".bin"
	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	return nil
}

func (t *Yosys) LoadFirmware(firmwarePath string, deviceID string) error {
	cmdStr := "tar -xvf "+ firmwarePath + " && iceprog firmware.bin"
	cmd := exec.Command("/bin/sh", "-c", cmdStr)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	return nil
}

func (t *Yosys) GetBuildDir() (string, error) {
	return t.BuildDir, nil
}

func (t *Yosys) GetFirmwareFile() (string, error) {
	return t.ModuleName + ".bin", nil
}
