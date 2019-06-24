GIT_DESCRIBE ?= $(shell git describe --tags --dirty --broken || git rev-parse --short HEAD)
GO111MODULE  ?= on

export GO111MODULE

build:
	go build -ldflags "-X main.version=$(GIT_DESCRIBE)"

install:
	go install -ldflags "-X main.version=$(GIT_DESCRIBE)"

vendor:
	go mod vendor
	git add vendor && git diff --cached --exit-code > /dev/null || git commit -s -m "Update vendored libs"

.PHONY: build install vendor