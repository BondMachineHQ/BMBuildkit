package build

import (
	"bufio"
	"fmt"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

type Xilinx struct {
	Config     BMFile
	BuildDir   string
	ModuleName string
}

// LoadData (only when we are going to implemente remote build)uncompress the source folder sent by the BMctl to the bondmachined
func (t *Xilinx) LoadData() error {

	return fmt.Errorf("NOT IMPLEMENTED YET")
}

func (t *Xilinx) Synth() error {

	return fmt.Errorf("NOT IMPLEMENTED YET")
}

func (t *Xilinx) PNR() error {

	return fmt.Errorf("NOT IMPLEMENTED YET")
}

func (t *Xilinx) Packing() error {

	return fmt.Errorf("NOT IMPLEMENTED YET")
}

func (t *Xilinx) LoadFirmware(firmwarePath string, deviceID string) error {
	cmdStr := "xbutil program --device " + deviceID + " --user " + firmwarePath
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

func (t *Xilinx) GetBuildDir() (string, error) {
	return t.BuildDir, nil
}

func (t *Xilinx) GetFirmwareFile() (string, error) {
	return t.ModuleName + ".bin", nil
}
