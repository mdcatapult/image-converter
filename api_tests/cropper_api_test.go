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
	json.Unmarshal(bytes, &responseWithError)

	return responseWithError.Error

}
