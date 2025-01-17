/*
 * Copyright 2024 Medicines Discovery Catapult
 * Licensed under the Apache License, Version 2.0 (the "Licence");
 * you may not use this file except in compliance with the Licence.
 * You may obtain a copy of the Licence at
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the Licence is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the Licence for the specific language governing permissions and
 * limitations under the Licence.
 */

package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/src/model"
)

// create a temporary macro file to use with fiji, return the temp filename if successful
func CreateTempMacroFile(request model.ConvertRequest, tempDir string) (*os.File, error) {

	err := os.MkdirAll(tempDir, os.ModeTemporary)
	if err != nil {
		return nil, errors.New("error creating temp file in directory: " + tempDir)
	}

	tmpFilePath := fmt.Sprintf("%s/macro.ijm", tempDir)
	tempMacroFile, err := os.Create(tmpFilePath)
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
		InputFile:         request.InputFile,
		InputFilename:     filepath.Base(request.InputFile),
		InputMaskFile:     request.InputMaskFile,
		InputMaskFilename: filepath.Base(request.InputMaskFile),
		OutputFile:        request.OutputFile,
	}

	templateString :=
		`open("{{.InputFile}}");
		run("Split Channels");
		open("{{.InputMaskFile }}");
		run("Split Channels");
		run("Merge Channels...", "c1=[{{.InputFilename}} (red)] c2=[{{.InputFilename}} (green)] c3=[{{.InputFilename}} (blue)] c4=[{{.InputMaskFilename}} (red)] c5=[{{.InputMaskFilename}} (green)] c6=[{{.InputMaskFilename}} (blue)] create");
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
