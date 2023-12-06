package build

type Yosys struct {
	source string
}

func CompressFile(src, dest string) error {

	return nil
}

func UncompressFile(input []byte, destDir string) error {

	return nil
}

// LoadData uncompress the source folder sent by the BMctl to the bondmachined
func (*Yosys) LoadData(sourceData []byte) (string, error) {

	return "", nil
}

func (*Yosys) Synth() error {

	return nil
}

func (*Yosys) PNR() error {

	return nil
}

func (*Yosys) Packing() error {

	return nil
}
