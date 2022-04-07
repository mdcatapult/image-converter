package utils

import (
	"bytes"
	"errors"
	"gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/model"
	"io/ioutil"
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

	filename := filepath.Base(request.InputFile)

	templateString := `open("{{.InputFile }}");
		run("Split Channels");
		run("Merge Channels...", ` + "\"c1=[" + filename + " (red)] c2=[" + filename + " (green)] c3=[" + filename + " (blue)] create\");" +
		`
		run("16-bit");
		saveAs("Tiff", "{{.OutputFile}}");`

	macroTemplate, err := template.New("_").Parse(templateString)

	if err != nil {
		return "", err
	}

	var templateBuffer bytes.Buffer
	if err = macroTemplate.Execute(&templateBuffer, request); err != nil {
		return "", err
	}

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
