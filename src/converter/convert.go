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
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/src/model"
	"gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/src/utils"
)

const (
	tempDirectory    = "/opt/temp"
	fijiAppPath      = "/opt/fiji/Fiji.app/ImageJ-linux64"
	fijiFileSuffix   = "-from-fiji"
	tiffExtension    = ".tiff"
	omeTiffExtension = ".ome.tiff"
)

func Convert(c *gin.Context) {

	var convertRequest model.ConvertRequest

	if err := c.ShouldBindBodyWith(&convertRequest, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": err.Error()})
		return
	}

	correctInputFileType := strings.HasSuffix(convertRequest.InputFile, tiffExtension)
	if !correctInputFileType {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"status": "input file extension must be .tiff, input file: " + convertRequest.InputFile},
		)
		return
	}

	correctInputMaskFileType := strings.HasSuffix(convertRequest.InputMaskFile, tiffExtension)
	if !correctInputMaskFileType {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"status": "input mask file extension must be .tiff, input mask file: " + convertRequest.InputMaskFile},
		)
		return
	}

	correctOutputFileType := strings.HasSuffix(convertRequest.OutputFile, omeTiffExtension)
	if !correctOutputFileType {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"status": "output file extension must be ome.tiff, output file: " + convertRequest.OutputFile},
		)
		return
	}

	log.Println("starting file conversion from: " + convertRequest.InputFile + " to:" + convertRequest.OutputFile)

	httpStatusCode, err := convertFile(convertRequest)
	if err != nil {
		c.AbortWithStatusJSON(httpStatusCode, gin.H{"status": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "conversion from: " + convertRequest.InputFile + " to: " + convertRequest.OutputFile + " complete"})

	tiffFilenameToRemove := utils.StripFileExtension(convertRequest.OutputFile, omeTiffExtension) + fijiFileSuffix + tiffExtension
	err = os.Remove(tiffFilenameToRemove)

	if err != nil {
		log.Println("error removing file: " + err.Error())
	}

	log.Println("removed intermediate file from fiji: " + tiffFilenameToRemove)
}

// from an image conversion request, check the input file can be loaded
// make a temporary macro file, create a macro string from the request, then write to the temp file
// run the fiji conversion command to create a tiff file
// run the bfconvert tool to create the final .ome.tiff file
func convertFile(request model.ConvertRequest) (HttpStatusCode int, e error) {

	if _, err := os.Stat(request.InputFile); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return http.StatusBadRequest, errors.New("file: " + request.InputFile + " does not exist")
		}

		return http.StatusInternalServerError, errors.New("error opening input file: " + request.InputFile)
	}

	if _, err := os.Stat(request.InputMaskFile); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return http.StatusBadRequest, errors.New("mask file: " + request.InputMaskFile + " does not exist")
		}

		return http.StatusInternalServerError, errors.New("error opening input mask file: " + request.InputFile)
	}

	// before making the temp macro file, make a copy of the request with an output file renamed from tiff to .ome.tiff
	// as this will be used in the macro, and later to refer to the output tiff file generated from fiji with bfconvert
	fijiRequest := model.ConvertRequest{
		InputFile:     request.InputFile,
		InputMaskFile: request.InputMaskFile,
		OutputFile:    utils.StripFileExtension(request.OutputFile, ".ome.tiff") + fijiFileSuffix + tiffExtension,
	}

	tempMacroFile, err := utils.CreateTempMacroFile(fijiRequest, tempDirectory)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer os.Remove(tempMacroFile.Name())

	return converter.Convert(fijiRequest.OutputFile, request.OutputFile, tempMacroFile.Name())

}
