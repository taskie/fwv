.PHONY: build install test fmt coverage dep-ensure dep-graph pre-commit

build:
	go build -v -ldflags "-s -w -X fwv.Revision=$(shell git rev-parse --short HEAD)"
	$(MAKE) -C cmd/fwv build

install:
	go install -v -ldflags "-s -w -X fwv.Revision=$(shell git rev-parse --short HEAD)"
	$(MAKE) -C cmd/fwv install

test:
	go test

fmt:
	find . -name '*.go' | xargs gofmt -w

coverage:
	mkdir -p test/coverage
	go test -coverprofile=test/coverage/cover.out
	go tool cover -html=test/coverage/cover.out -o test/coverage/cover.html

dep-ensure:
	$(MAKE) -C cmd/fwv dep-ensure

dep-graph:
	mkdir -p images
	dep status -dot | dot -Tpng -o images/dependency.png


pre-commit:
	$(MAKE) fmt
	$(MAKE) build
	$(MAKE) coverage
	rm -rf vendor/
	$(MAKE) dep-ensure
	$(MAKE) dep-graph
