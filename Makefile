GIT_DESCRIBE ?= $(shell git describe --tags --dirty --broken --always || git rev-parse --short HEAD)

build:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-extldflags '-static' -w -s -X main.version=$(GIT_DESCRIBE)"

install:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go install -ldflags "-extldflags '-static' -w -s -X main.version=$(GIT_DESCRIBE)"

.PHONY: build install
