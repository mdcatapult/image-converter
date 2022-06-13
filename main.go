package main

import (
	"gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/src/cropper"
	"gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/src/server"
)

func main() {

	cropper.SetCropper(cropper.ImplementedCropper{})
	if err := server.Start(":8080"); err != nil {
		panic(err)
	}
}
