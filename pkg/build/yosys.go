package build

import (
	"bufio"
	"fmt"
	"io"
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

func CopyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}

	}

	return
}

func CopyDir(source string, dest string) (err error) {

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		sourcefilepointer := source + "/" + obj.Name()

		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			// create sub-directories - recursively
			err = CopyDir(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// perform copy
			err = CopyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
	return
}

// Uncomment for remote implementation
// func CompressFile(src, dest string) error {

// 	return nil
// }

// func UncompressFile(input []byte, destDir string) error {

// 	return nil
// }

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

	cmdStr := "cd " + t.BuildDir + "; nextpnr-ice40 --package tq144 --hx1k --json " + t.ModuleName + ".json --pcf " + t.ModuleName + ".pcf --asc " + t.ModuleName + ".asc"
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

func (t *Yosys) GetBuildDir() (string, error) {
	return t.BuildDir, nil
}

func (t *Yosys) GetFirmwareFile() (string, error) {
	return "blinky.bin", nil
}
