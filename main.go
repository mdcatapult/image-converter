package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/model"
	"gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/utils"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	tempDirectory    = "/opt/temp"
	fijiAppPath      = "/opt/fiji/Fiji.app/ImageJ-linux64"
	bfconvertAppPath = "/opt/bftools/bfconvert"
	fijiFileSuffix   = "-from-fiji"
	tiffExtension    = ".tiff"
	omeTiffExtension = ".ome.tiff"
)

func main() {
	router := gin.Default()

	router.POST("/convert", convertImage)

	err := router.Run()
	if err != nil {
		panic(err)
	}
}

func convertImage(c *gin.Context) {
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

	// before making the temp macro file, make a copy of the request with an output file renamed from tiff to .ome.tiff
	// as this will be used in the macro, and later to refer to the output tiff file generated from fiji with bfconvert
	fijiRequest := model.ConvertRequest{
		InputFile:  request.InputFile,
		OutputFile: utils.StripFileExtension(request.OutputFile, ".ome.tiff") + fijiFileSuffix + tiffExtension,
	}

	tempMacroFile, err := utils.CreateTempMacroFile(fijiRequest, tempDirectory)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer os.Remove(tempMacroFile.Name())

	//example command: "./Fiji.app/ImageJ-linux64 --console --memory=2g -macro ./data/docker-convert-image-simple.ijm`"
	stdOut, err := exec.Command(fijiAppPath, "--console", "--memory=2g", "-macro", tempMacroFile.Name()).Output()
	if err != nil {
		return http.StatusInternalServerError, errors.New("error during Fiji macro execution: " + string(stdOut))
	}

	// file now needs converting from tiff to ome.tiff using bfconvert
	// input will be the output tiff file from fiji, output the original output .ome.tiff
	// example bfconvert command: bfconvert -overwrite  2106-bladder-tma-nimrad.tiff test.ome.tiff
	stdOut, err = exec.Command(bfconvertAppPath, "-overwrite", fijiRequest.OutputFile, request.OutputFile).Output()
	if err != nil {
		return http.StatusInternalServerError, errors.New("error during bfconvert execution: " + string(stdOut))
	}

	return http.StatusOK, nil
}
