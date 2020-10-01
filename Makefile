GIT_DESCRIBE ?= $(shell git describe --tags --dirty --broken --always || git rev-parse --short HEAD)
GO111MODULE  ?= on

export GO111MODULE

build:
	CGO_ENABLED=0 go build -ldflags "-extldflags '-static' -w -s -X main.version=$(GIT_DESCRIBE)"

install:
	CGO_ENABLED=0 go install -ldflags "-extldflags '-static' -w -s -X main.version=$(GIT_DESCRIBE)"

.PHONY: build install