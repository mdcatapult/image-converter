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
