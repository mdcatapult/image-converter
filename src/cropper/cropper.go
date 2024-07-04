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

package cropper

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

var cropper Cropper

func SetCropper(c Cropper) {
	cropper = c
}

type Cropper interface {
	Crop(cropInstruction, patternFilePath, outputPath string) (croppedImageBytes []byte, err error)
}

type ImplementedCropper struct{}

func (c ImplementedCropper) Crop(cropInstruction, patternFilePath, outputPath string) ([]byte, error) {

	cmd := exec.Command(os.Getenv("BF_TOOLS_CONVERT_PATH"), "-crop", cropInstruction, patternFilePath, outputPath)

	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	return ioutil.ReadFile(outputPath)
}
