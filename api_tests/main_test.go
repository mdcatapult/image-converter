package apitest

import (
	"log"
	"os"
	"testing"

	"gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/src/cropper"
	"gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/src/server"
)

func TestMain(m *testing.M) {
	cropper.SetCropper(mockCropper{})
	go func() {
		if err := server.Start(":8081"); err != nil {
			log.Fatal(err)
		}
	}()
	os.Exit(m.Run())
}

type mockCropper struct{}

func (mc mockCropper) Crop(cropInstruction, patternFilePath, outputPath string) ([]byte, error) {
	bytes := []byte("some image data")
	return bytes, nil
}
