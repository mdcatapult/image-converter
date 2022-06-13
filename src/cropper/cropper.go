package cropper

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

var cropper Cropper

func SetCropper(c Cropper) {
	cropper = c
}

type Cropper interface {
	Crop(cropInstruction, patternFilePath, outputPath string) (croppedImageBytes []byte, err error)
}

type ImplementedCropper struct{}

func (c ImplementedCropper) Crop(cropInstruction, patternFilePath, outputPath string) ([]byte, error) {

	cmd := exec.Command(os.Getenv("BF_TOOLS_CONVERT_PATH"), "-crop", cropInstruction, patternFilePath, outputPath)

	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	return ioutil.ReadFile(outputPath)
}
