package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-playground/assert/v2"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

var convertUrl = fmt.Sprintf("http://%s/convert", getHostNameAndPort())

// if running in CI, gets the hostname and port from an env var, else uses localhost and the port mapping defined
// in the local docker-compose file
func getHostNameAndPort() string {
	hostnameAndPort := os.Getenv("HOSTNAME_FROM_CI")

	if hostnameAndPort == "" {
		return "localhost:8081"
	}
	return hostnameAndPort
}


func TestBadRequestInvalidJson(t *testing.T) {
	values := map[string]string{"inpert-file": "scooby-dooby", "output-file": "dooby-doo"}
	jsonValue, _ := json.Marshal(values)

	resp, err := http.Post(convertUrl, "json ", bytes.NewBuffer(jsonValue))

	assert.Equal(t, nil, err)

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	responseBody := string(bodyBytes)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "{\"status\":\"Key: 'ConvertRequest.InputFile' Error:Field validation for 'InputFile' failed on the 'required' tag\"}", responseBody)
}

func TestBadRequestIncorrectInputFormat(t *testing.T) {
	values := map[string]string{"input-file": "scooby-dooby.xlsx", "output-file": "dooby-doo.ome.tiff"}
	jsonValue, _ := json.Marshal(values)

	resp, err := http.Post(convertUrl, "json ", bytes.NewBuffer(jsonValue))

	assert.Equal(t, nil, err)

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	responseBody := string(bodyBytes)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "{\"status\":\"input file extension must be .tiff, input file: scooby-dooby.xlsx\"}", responseBody)
}

func TestBadRequestIncorrectOutputFormat(t *testing.T) {
	values := map[string]string{"input-file": "scooby-dooby.tiff", "output-file": "dooby-doo.text"}
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
		"input-file":  "/opt/data/2106xx_Bladder_TMA_NIMRAD-croppyyyy.tiff",
		"output-file": "/opt/data/converted_file_test.ome.tiff",
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
		"input-file":  "/opt/data/2106xx_Bladder_TMA_NIMRAD-crop.tiff",
		"output-file": "/opt/data/converted_file_test.ome.tiff",
	}

	jsonValue, _ := json.Marshal(values)

	resp, err := http.Post(convertUrl, "json ", bytes.NewBuffer(jsonValue))

	assert.Equal(t, nil, err)

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	responseBody := string(bodyBytes)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t,
		"{\"status\":\"conversion from: /opt/data/2106xx_Bladder_TMA_NIMRAD-crop.tiff to: /opt/data/converted_file_test.ome.tiff complete\"}",
		responseBody,
	)
}

// TODO maybe test the hash of the file?