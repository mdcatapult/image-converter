package apitest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/go-playground/assert/v2"
	"gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/cropper"
	"gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/test_utils"
)

var cropUrl = test_utils.URL{test_utils.GetUrl("/crop")}

func TestBadRequestMissingX(t *testing.T) {

	res, err := http.Get(cropUrl.WithParam("y=2000").WithParam("experiment-directory=hello").Url)

	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, cropper.ErrMissingXParam, getResponseError(res))
}

func TestBadRequestMissingY(t *testing.T) {

	res, err := http.Get(cropUrl.WithParam("x=2000").WithParam("experiment-directory=hello").Url)

	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, cropper.ErrMissingYParam, getResponseError(res))
}

func TestBadRequestMissingExperimentDir(t *testing.T) {

	res, err := http.Get(cropUrl.WithParam("y=2000").WithParam("x=2000").WithParam("crop-size=100").Url)

	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, cropper.ErrMissingExperimentDir, getResponseError(res))
}

func TestBadRequestXType(t *testing.T) {

	res, err := http.Get(cropUrl.WithParam("y=2000").WithParam("x=hello").WithParam("crop-size=100").WithParam("experiment-directory=hello").Url)

	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, cropper.ErrXParamType, getResponseError(res))
}

func TestBadRequestYType(t *testing.T) {

	res, err := http.Get(cropUrl.WithParam("x=2000").WithParam("y=hello").WithParam("crop-size=100").WithParam("experiment-directory=hello").Url)

	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, cropper.ErrYParamType, getResponseError(res))
}

func TestCropFile(t *testing.T) {
	experimentDir := "test-images/crop-test-data"
	dspMountPath := "/"
	os.Setenv("DSP_MNT_PATH", dspMountPath)
	os.Setenv("BF_TOOLS_INFO_PATH", "/opt/bftools/showinf")
	os.Setenv("BF_TOOLS_CONVERT_PATH", "/opt/bftools/bfconvert")

	os.MkdirAll(fmt.Sprintf("%s/%s", dspMountPath, experimentDir), os.ModeTemporary)

	res, err := http.Get(cropUrl.WithParam("y=50").WithParam("x=50").WithParam("crop-size=100").WithParam(fmt.Sprintf("experiment-directory=%s", experimentDir)).Url)

	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	imageFilePath := fmt.Sprintf("%s/%s/%s", os.Getenv("DSP_MNT_PATH"), experimentDir, cropper.GetCroppedImageName(experimentDir))
	_, err = os.Stat(imageFilePath)
	assert.Equal(t, nil, err)

}

func getResponseError(resp *http.Response) string {
	bytes, _ := io.ReadAll(resp.Body)

	var responseWithError struct {
		Error string `json:"error"`
	}
	json.Unmarshal(bytes, &responseWithError)

	return responseWithError.Error

}
