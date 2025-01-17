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
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/go-playground/assert/v2"
	"gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/test_utils"
)

var convertUrl = test_utils.GetUrl("/convert")

func TestBadRequestInvalidJson(t *testing.T) {
	values := map[string]string{"inpert-file": "scooby-dooby", "input-mask-file": "dee-dee", "output-file": "dooby-doo"}
	jsonValue, _ := json.Marshal(values)

	resp, err := http.Post(convertUrl, "json ", bytes.NewBuffer(jsonValue))

	assert.Equal(t, nil, err)

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	responseBody := string(bodyBytes)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "{\"status\":\"Key: 'ConvertRequest.InputFile' Error:Field validation for 'InputFile' failed on the 'required' tag\"}", responseBody)
}

func TestBadRequestIncorrectInputFileFormat(t *testing.T) {
	values := map[string]string{"input-file": "scooby-dooby.xlsx", "input-mask-file": "piglet.tiff", "output-file": "dooby-doo.ome.tiff"}
	jsonValue, _ := json.Marshal(values)

	resp, err := http.Post(convertUrl, "json ", bytes.NewBuffer(jsonValue))

	assert.Equal(t, nil, err)

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	responseBody := string(bodyBytes)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "{\"status\":\"input file extension must be .tiff, input file: scooby-dooby.xlsx\"}", responseBody)
}

func TestBadRequestIncorrectInputMaskFileFormat(t *testing.T) {
	values := map[string]string{"input-file": "scooby-dooby.tiff", "input-mask-file": "bertrand.xlsx", "output-file": "dooby-doo.ome.tiff"}
	jsonValue, _ := json.Marshal(values)

	resp, err := http.Post(convertUrl, "json ", bytes.NewBuffer(jsonValue))

	assert.Equal(t, nil, err)

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	responseBody := string(bodyBytes)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "{\"status\":\"input mask file extension must be .tiff, input mask file: bertrand.xlsx\"}", responseBody)
}

func TestBadRequestIncorrectOutputFileFormat(t *testing.T) {
	values := map[string]string{"input-file": "scooby-dooby.tiff", "input-mask-file": "tiddle.tiff", "output-file": "dooby-doo.text"}
	jsonValue, _ := json.Marshal(values)

	resp, err := http.Post(convertUrl, "json ", bytes.NewBuffer(jsonValue))

	assert.Equal(t, nil, err)

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	responseBody := string(bodyBytes)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "{\"status\":\"output file extension must be ome.tiff, output file: dooby-doo.text\"}", responseBody)
}

func TestBadRequestNonExistentFile(t *testing.T) {
	values := map[string]string{
		"input-file":      "/opt/data/2106xx_Bladder_TMA_NIMRAD-croppyyyy.tiff",
		"input-mask-file": "/opt/data/2106xx_Bladder_TMA_NIMRAD-crop.tiff",
		"output-file":     "/opt/data/converted_file_test.ome.tiff",
	}

	jsonValue, _ := json.Marshal(values)

	resp, err := http.Post(convertUrl, "json ", bytes.NewBuffer(jsonValue))

	assert.Equal(t, nil, err)

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	responseBody := string(bodyBytes)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "{\"status\":\"file: /opt/data/2106xx_Bladder_TMA_NIMRAD-croppyyyy.tiff does not exist\"}", responseBody)
}

func TestFileIsConverted(t *testing.T) {

	values := map[string]string{
		"input-file":      "test-images/raw-image/2106xx_Bladder_TMA_NIMRAD-crop.tiff",
		"input-mask-file": "test-images/raw-image/2106xx_Bladder_TMA_NIMRAD-crop-mask.tiff",
		"output-file":     "test-images/raw-image/converted_file_test.ome.tiff",
	}

	jsonValue, _ := json.Marshal(values)

	resp, err := http.Post(convertUrl, "json ", bytes.NewBuffer(jsonValue))

	assert.Equal(t, nil, err)

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	responseBody := string(bodyBytes)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t,
		"{\"status\":\"conversion from: test-images/raw-image/2106xx_Bladder_TMA_NIMRAD-crop.tiff to: test-images/raw-image/converted_file_test.ome.tiff complete\"}",
		responseBody,
	)
}
