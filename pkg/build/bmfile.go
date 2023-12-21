package build

import "strings"

func ParseBMFile(content []byte) (BMFile, error) {

	var bmfile BMFile

	lines := strings.Split(string(content), "\n")

	for _, l := range lines {
		cmd := strings.Split(l, " ")
		switch cmd[0] {
		case "SYNTH":
			bmfile.Synth = cmd[1]
		case "SYNTH_FILE":
			bmfile.SynthFile = cmd[1]
		case "SYNTH_CMD":
			bmfile.Vendor = strings.Join(cmd[1:], " ")
		case "VENDOR":
			bmfile.Vendor = cmd[1]
		case "MODEL":
			bmfile.Model = cmd[1]
		case "SOURCE":
			bmfile.Source = cmd[1]
		case "BITSTREAM_PATH":
			bmfile.BitstreamPath = cmd[1]
		}
	}

	return bmfile, nil
}

func ExecuteEngine(synthEngine SynthEngine) error {

	err := synthEngine.LoadData()
	if err != nil {
		panic(err)
	}

	err = synthEngine.Synth()
	if err != nil {
		panic(err)
	}

	err = synthEngine.PNR()
	if err != nil {
		panic(err)
	}

	err = synthEngine.Packing()
	if err != nil {
		panic(err)
	}

	return nil
}
