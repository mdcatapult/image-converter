package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/converter"
	"gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/cropper"
)

func Start(port string) error {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
	}))

	router.POST("/convert", converter.Convert)
	router.GET("/crop", cropper.Crop)

	return router.Run(port)
}
