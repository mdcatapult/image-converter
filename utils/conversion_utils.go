package utils

import (
	"bytes"
	"errors"
	"gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/model"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// create a temporary macro file to use with fiji, return the temp filename if successful
func CreateTempMacroFile(request model.ConvertRequest, tempDir string) (*os.File, error) {
	tempMacroFile, err := ioutil.TempFile(tempDir, "macro*.ijm")
	if err != nil {
		return nil, errors.New("error creating temp file in directory: " + tempDir)
	}

	macroString, err := createFijiMacroString(request)
	if err != nil {
		return nil, err
	}

	_, err = tempMacroFile.WriteString(macroString)
	if err != nil {
		return nil, err
	}

	return tempMacroFile, nil
}

func createFijiMacroString(request model.ConvertRequest) (string, error) {

	requestInputFilenames := model.ConvertRequestForFijiMacro{
		InputFile:     request.InputFile,
		InputFilename: filepath.Base(request.InputFile),
		InputMaskFile: request.InputMaskFile,
		InputMaskFilename: filepath.Base(request.InputMaskFile),
		OutputFile: request.OutputFile,
	}

	templateString :=
		`open("{{.InputMaskFile }}");
		run("Split Channels");
		open("{{.InputFile}}");
		run("Split Channels");
		run("Merge Channels...", "c1=[{{.InputFilename}} (red)] c2={{.InputFilename}} (green)] c3=[{{.InputFilename}} (blue)] c4=[{.InputMaskFilename}} (red)] c5=[{{.InputMaskFilename}} (green)] c6=[{{.InputMaskFilename}} (blue)] create");
		selectWindow("Composite");
		run("16-bit");
		saveAs("Tiff", "{{.OutputFile}}");`

	macroTemplate, err := template.New("_").Parse(templateString)

	if err != nil {
		return "", err
	}

	var templateBuffer bytes.Buffer
	if err = macroTemplate.Execute(&templateBuffer, requestInputFilenames); err != nil {
		return "", err
	}

	log.Println(templateBuffer.String())

	return templateBuffer.String(), nil
}

// removes the specified extension if it exists, else returns the original filename
func StripFileExtension(name, extensionToRemove string) string {
	index := strings.LastIndex(name, extensionToRemove)

	if index == -1 {
		return name
	}

	return name[:index]
}
