#
# A simple Makefile to easily build, test and run the code
#

.PHONY: default build fmt lint run run_race test clean vet docker_build docker_run docker_clean

APP_NAME := ut4updater

default: build

build:
	go build -o ./bin/${APP_NAME} ./*.go

# http://golang.org/cmd/go/#hdr-Run_gofmt_on_package_sources
fmt:
	go fmt ./...

test:
	go test ./ -v -cover -covermode=count -coverprofile=./coverage.out

clean:
	rm -Rf ./test-resources/installs/004
	rm -Rf ./test-resources/test/
	rm ./bin/*
