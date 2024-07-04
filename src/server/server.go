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

package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/src/converter"
	"gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/src/cropper"
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
