package apitest

import (
	"os"
	"testing"

	"gitlab.mdcatapult.io/informatics/software-engineering/mdc-minerva-image-converter/server"
)

func TestMain(m *testing.M) {
	go func() {
		server.Start(":8081")
	}()
	os.Exit(m.Run())
}
