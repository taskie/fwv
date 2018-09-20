.PHONY: build install

build:
	go build -v -ldflags "-s -w -X main.revision=$(shell git rev-parse --short HEAD)"

install:
	go install -v -ldflags "-s -w -X main.revision=$(shell git rev-parse --short HEAD)"

dep-ensure:
	dep ensure
	dep status
