package utils

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func CreateIfNotExists(path string) string {
	_, err := os.Stat(path)
	//Create directory if doesn't exist
	if os.IsNotExist(err) {
		os.MkdirAll(path, os.ModeDir)
	}
	return path
}

func GetBinDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func GetLogDir() string {
	logDir := filepath.Join(GetBinDir(), "logs")
	return CreateIfNotExists(logDir)
}

func ReadFile(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}
