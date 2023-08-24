package cropper

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const (
	ErrMissingXParam = "must supply x coordinate"
	ErrXParamType    = "x must be integer"

	ErrMissingYParam = "must supply y coordinate"
	ErrYParamType    = "y must be integer"

	ErrMissingCropSizeParam = "must supply crop-size param"
	ErrCropSizeParamType    = "crop size must be integer"

	ErrMissingExperimentDir = "must supply experiment-directory"
	ErrOutOfBounds          = "coords out of bounds"
)

func Crop(c *gin.Context) {

	x, y, experimentDir, cropSize, err := getParams(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	patternFilePath, outputPath := getPaths(experimentDir)

	startX := x - cropSize/2
	startY := y - cropSize/2

	if err := validateCoords(startX, startY, cropSize, patternFilePath, readImageMetadata); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	cropInstruction := fmt.Sprintf("%v,%v,%v,%v", startX, startY, cropSize, cropSize)

	croppedImageBytes, err := cropper.Crop(cropInstruction, patternFilePath, outputPath)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, errors.Wrap(err, fmt.Sprintf("crop instruction was %s", cropInstruction)))
		return
	}

	if err := os.Remove(outputPath); err != nil {
		fmt.Println("error cleaning up file:", err)
	}

	c.Data(http.StatusOK, "application/octet-stream", croppedImageBytes)
}

func getParams(c *gin.Context) (x, y int64, experimentDir string, cropSize int64, err error) {

	xParam, ok := c.GetQuery("x")
	if !ok {
		return 0, 0, "", 0, errors.New(ErrMissingXParam)
	}

	x, err = strconv.ParseInt(xParam, 10, 64)
	if err != nil {
		return 0, 0, "", 0, errors.New(ErrXParamType)
	}

	yParam, ok := c.GetQuery("y")
	if !ok {
		return 0, 0, "", 0, errors.New(ErrMissingYParam)
	}

	y, err = strconv.ParseInt(yParam, 10, 64)
	if err != nil {
		return 0, 0, "", 0, errors.New(ErrYParamType)
	}

	experimentDir, ok = c.GetQuery("experiment-directory")
	if !ok {
		return 0, 0, "", 0, errors.New(ErrMissingExperimentDir)
	}

	cropSizeParam, ok := c.GetQuery("crop-size")
	if !ok {
		return 0, 0, "", 0, errors.New(ErrMissingCropSizeParam)
	}

	cropSize, err = strconv.ParseInt(cropSizeParam, 10, 64)
	if err != nil {
		return 0, 0, "", 0, errors.New(ErrCropSizeParamType)
	}
	return
}

func getPaths(experimentDir string) (patternFilePath, outputPath string) {

	dataMountPath := os.Getenv("DSP_ATLAS_DATA")

	println(dataMountPath)

	patternFilePath = fmt.Sprintf("%v/%v/raw-image/channels.pattern", dataMountPath, experimentDir)

	println(patternFilePath)

	outputPath = GetCroppedImageName(experimentDir)

	return
}

func GetCroppedImageName(experimentDir string) string {
	return experimentDir + "-cropped.tiff"
}

func readImageMetadata(patternFilePath string) (string, error) {
	out, err := exec.Command(os.Getenv("BF_TOOLS_INFO_PATH"), "-nopix", patternFilePath).Output()

	return string(out), errors.Wrap(err, "Couldn't read info about the raw image to validate coordinates!")
}

func validateCoords(x, y, croppedImageSize int64, patternFilePath string, imageMetadataReader func(string) (string, error)) error {

	if x < 0 || y < 0 {
		return errors.New(ErrOutOfBounds)
	}

	rawImgData, err := imageMetadataReader(patternFilePath)
	if err != nil {
		// failing this means we can't properly validate the user's given coordinates.
		// If the coords are invalid, the error will be caught later by bftools but the message won't be as nice.
		// This isn't a reason abort the process, so don't return error here.
		return nil
	}

	heightSplit := strings.Split(rawImgData, "Height = ")
	heightStr := strings.Split(heightSplit[1], "\n")[0]
	rawImageHeight, _ := strconv.ParseInt(heightStr, 10, 64)

	widthSplit := strings.Split(rawImgData, "Width = ")
	widthStr := strings.Split(widthSplit[1], "\n")[0]
	rawImageWidth, _ := strconv.ParseInt(widthStr, 10, 64)

	if x+croppedImageSize > rawImageWidth || y+croppedImageSize > rawImageHeight {
		return errors.New(ErrOutOfBounds)
	}

	return nil

}
