package utils

import (
	"github.com/go-playground/assert/v2"
	"gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/model"
	"io/ioutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	// make a temp directory to store temp macro files
	// delete directory after test run
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	tempDir := workingDir + "/tmp"
	err = os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	exitCode := m.Run()

	err = os.RemoveAll(tempDir)
	if err != nil {
		panic(err)
	}

	os.Exit(exitCode)
}

func TestMacroFileIsCorrect(t *testing.T) {

	request := model.ConvertRequest{
		InputFile:  "some-input-directory/test.tiff",
		InputMaskFile: "some-input-directory/test-mask.tiff",
		OutputFile: "some-output-directory/test.ome.tiff",
	}

	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	tempDir := workingDir + "/tmp"

	file, err := CreateTempMacroFile(request, tempDir)

	assert.Equal(t, nil, err)

	fileBytes, err := ioutil.ReadFile(file.Name())
	assert.Equal(t, nil, err)

	fileAsString := string(fileBytes)

	expectedResult := `open("some-input-directory/test-mask.tiff");
		run("Split Channels");
		open("some-input-directory/test.tiff");
		run("Split Channels");
		run("Merge Channels...", "c1=[test.tiff (red)] c2=test.tiff (green)] c3=[test.tiff (blue)] c4=[{.InputMaskFilename}} (red)] c5=[test-mask.tiff (green)] c6=[test-mask.tiff (blue)] create");
		selectWindow("Composite");
		run("16-bit");
		saveAs("Tiff", "some-output-directory/test.ome.tiff");`

	assert.Equal(t, expectedResult, fileAsString)
}
