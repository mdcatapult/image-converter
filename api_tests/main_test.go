package apitest

import (
	"log"
	"net/http"
	"os"
	"testing"

	"gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/src/converter"

	"gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/src/cropper"
	"gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/src/server"
)

func TestMain(m *testing.M) {
	cropper.SetCropper(mockCropper{})
	converter.SetConverter(mockConverter{})
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

type mockConverter struct{}

func (mc mockConverter) Convert(fijiOutputPath, requestOutputPath, tempMacroPath string) (httpStatusCode int, err error) {
	return http.StatusOK, nil
}
