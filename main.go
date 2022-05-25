package main

import "gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/server"

func main() {
	if err := server.Start(":8080"); err != nil {
		panic(err)
	}
}
