package apitest

import (
	"os"
	"testing"

	"gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/src/cropper"
	"gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/src/server"
)

func TestMain(m *testing.M) {
	cropper.SetCropper(mockCropper{})
	go func() {
		server.Start(":8081")
	}()
	os.Exit(m.Run())
}

type mockCropper struct{}

func (mc mockCropper) Crop(cropInstruction, patternFilePath, outputPath string) error {
	f, err := os.Create(outputPath)

	if err != nil {
		return err
	}

	defer f.Close()

	_, err2 := f.WriteString("some image data")

	if err2 != nil {
		return err
	}

	return nil
}
