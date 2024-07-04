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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/src/cropper"
	"gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/test_utils"
)

var cropUrl = test_utils.URL{Url: test_utils.GetUrl("/crop")}

func TestBadRequestMissingX(t *testing.T) {

	res, err := http.Get(cropUrl.WithParam("y=2000").WithParam("experiment-directory=hello").Url)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, cropper.ErrMissingXParam, getResponseError(res))
}

func TestBadRequestMissingY(t *testing.T) {

	res, err := http.Get(cropUrl.WithParam("x=2000").WithParam("experiment-directory=hello").Url)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, cropper.ErrMissingYParam, getResponseError(res))
}

func TestBadRequestMissingExperimentDir(t *testing.T) {

	res, err := http.Get(cropUrl.WithParam("y=2000").WithParam("x=2000").WithParam("crop-size=100").Url)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, cropper.ErrMissingExperimentDir, getResponseError(res))
}

func TestBadRequestXType(t *testing.T) {

	res, err := http.Get(cropUrl.WithParam("y=2000").WithParam("x=hello").WithParam("crop-size=100").WithParam("experiment-directory=hello").Url)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, cropper.ErrXParamType, getResponseError(res))
}

func TestBadRequestYType(t *testing.T) {

	res, err := http.Get(cropUrl.WithParam("x=2000").WithParam("y=hello").WithParam("crop-size=100").WithParam("experiment-directory=hello").Url)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, cropper.ErrYParamType, getResponseError(res))
}

func TestCropFile(t *testing.T) {
	experimentDir := "test-images"
	os.Setenv("BF_TOOLS_INFO_PATH", "/opt/bftools/showinf")
	os.Setenv("BF_TOOLS_CONVERT_PATH", "/opt/bftools/bfconvert")

	res, err := http.Get(cropUrl.WithParam("y=50").WithParam("x=50").WithParam("crop-size=100").WithParam(fmt.Sprintf("experiment-directory=%s", experimentDir)).Url)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	imageBytes, err := ioutil.ReadAll(res.Body)

	assert.Nil(t, err)

	// assert that we got some image back as bytes
	assert.True(t, len(imageBytes) > 0)

}

func getResponseError(resp *http.Response) string {
	bytes, _ := ioutil.ReadAll(resp.Body)

	var responseWithError struct {
		Error string `json:"error"`
	}
	_ = json.Unmarshal(bytes, &responseWithError)

	return responseWithError.Error

}
