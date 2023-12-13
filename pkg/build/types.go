package build

import "github.com/moby/buildkit/client/llb"

type BMImage struct {
	ImageName string
	State     llb.State
}

type BMFile struct {
	Synth     string
	SynthCMD  string
	SynthFile string
	Vendor    string
	Model     string
	Source    string
}

type SynthEngine interface {
	LoadData() error
	Synth() error
	PNR() error
	Packing() error
	GetFirmwareFile() (string, error)
	GetBuildDir() (string, error)
}
