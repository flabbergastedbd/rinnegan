package utils

import (
	"os"
	"strings"
)

var Hostname string

func init() {
	Hostname, _ =  os.Hostname()
	Hostname = strings.Replace(Hostname, ".", "_", -1)
}

