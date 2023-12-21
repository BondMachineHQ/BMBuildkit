package build

import "github.com/moby/buildkit/client/llb"

type BMImage struct {
	ImageName string
	State     llb.State
}

type BMFile struct {
	Synth         string
	SynthCMD      string
	SynthFile     string
	BitstreamPath string
	Vendor        string
	Model         string
	Source        string
}

type SynthEngine interface {
	LoadData() error
	Synth() error
	PNR() error
	Packing() error
	LoadFirmware(firmwarePath string) error
	GetFirmwareFile() (string, error)
	GetBuildDir() (string, error)
}
