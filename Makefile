.PHONY: linux_agent osx_agent windows_agent clean

BASE_DIR:=$(shell pwd)
GOARCH:=amd64

linux_agent: export GOOS := linux
linux_agent:
	go build -o build/agent ./agent

osx_agent: export GOOS := darwin
osx_agent:
	go build -o build/agent ./agent

windows_agent: export GOOS := windows
windows_agent:
	go build -o build/agent ./agent

clean:
	rm build/agent
