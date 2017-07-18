#
# A simple Makefile to easily build, test and run the code
#

.PHONY: default build fmt lint run run_race test clean vet docker_build docker_run docker_clean

APP_NAME := ut4updater

default: build

build:
	go build -o ./bin/${APP_NAME} ./src/ut4updater/ut4updater.go

# http://golang.org/cmd/go/#hdr-Run_gofmt_on_package_sources
fmt:
	go fmt ./src/...

test:
	go test ./src/... -v -cover

clean:
	rm ./bin/*
