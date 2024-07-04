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

package converter

import (
	"errors"
	"net/http"
	"os"
	"os/exec"
)

var converter Converter

func SetConverter(c Converter) {
	converter = c
}

type Converter interface {
	Convert(fijiOutputPath, requestOutputPath, tempMacroPath string) (httpStatusCode int, err error)
}

type ImplementedConverter struct{}

// example command: "./Fiji.app/ImageJ-linux64 --console --memory=2g -macro ./data/docker-convert-image-simple.ijm`"
func (ic ImplementedConverter) Convert(fijiOutputPath, requestOutputPath, tempMacroPath string) (httpStatusCode int, err error) {
	stdOut, err := exec.Command(fijiAppPath, "--console", "--memory=2g", "-macro").Output()
	if err != nil {
		return http.StatusInternalServerError, errors.New("error during Fiji macro execution: " + string(stdOut))
	}

	// file now needs converting from tiff to ome.tiff using bfconvert
	// input will be the output tiff file from fiji, output the original output .ome.tiff
	// example bfconvert command: bfconvert -overwrite  2106-bladder-tma-nimrad.tiff test.ome.tiff
	os.Setenv("BF_MAX_MEM", "4g")
	cmd, err := exec.Command(os.Getenv("BF_TOOLS_CONVERT_PATH"), "-overwrite", "-pyramid-resolutions", "6", "-pyramid-scale", "2", fijiOutputPath, requestOutputPath).CombinedOutput()

	if err != nil {
		return http.StatusInternalServerError, errors.New("error during bfconvert execution: " + string(cmd))
	}

	return http.StatusOK, nil
}
