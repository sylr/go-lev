GIT_DESCRIBE ?= $(shell git describe --tags --dirty --broken --always || git rev-parse --short HEAD)
GO111MODULE  ?= on

export GO111MODULE

build:
	go build -ldflags "-X main.version=$(GIT_DESCRIBE)"

install:
	go install -ldflags "-X main.version=$(GIT_DESCRIBE)"

.PHONY: build install