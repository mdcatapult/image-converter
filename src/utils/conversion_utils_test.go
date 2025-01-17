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
	"io/ioutil"
	"os"
	"testing"

	"github.com/go-playground/assert/v2"
	"gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/src/model"
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
		InputFile:     "some-input-directory/test.tiff",
		InputMaskFile: "some-input-directory/test-mask.tiff",
		OutputFile:    "some-output-directory/test.ome.tiff",
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

	expectedResult := `open("some-input-directory/test.tiff");
		run("Split Channels");
		open("some-input-directory/test-mask.tiff");
		run("Split Channels");
		run("Merge Channels...", "c1=[test.tiff (red)] c2=[test.tiff (green)] c3=[test.tiff (blue)] c4=[test-mask.tiff (red)] c5=[test-mask.tiff (green)] c6=[test-mask.tiff (blue)] create");
		run("16-bit");
		saveAs("Tiff", "some-output-directory/test.ome.tiff");`

	assert.Equal(t, expectedResult, fileAsString)
}
