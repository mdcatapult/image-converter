package test_utils

import (
	"fmt"
	"os"
	"strings"
)

type URL struct {
	Url string
}

func GetUrl(path string) string {
	return fmt.Sprintf("http://%s%s", getHostNameAndPort(), path)
}

func (u URL) WithParam(param string) URL {
	if !strings.Contains(u.Url, "?") {
		u.Url = u.Url + "?"
	} else {
		u.Url = u.Url + "&"
	}

	return URL{fmt.Sprintf("%s%s", u.Url, param)}
}

// if running in CI, gets the hostname and port from an env var, else uses localhost and the port mapping defined
// in the local docker-compose file
func getHostNameAndPort() string {
	hostnameAndPort := os.Getenv("HOSTNAME_FROM_CI")

	if hostnameAndPort == "" {
		return "localhost:8081"
	}
	return hostnameAndPort
}