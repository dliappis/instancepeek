BINDIR := "build"
CWD = $(shell pwd)

build: setup
	CGO_ENABLED=0 gox -osarch="linux/amd64 darwin/amd64 windows/amd64" -output="$(BINDIR)/{{.Dir}}_{{.OS}}_{{.Arch}}" ./

setup:
	go get github.com/mitchellh/gox

test:
	go test $(CWD)/providers/aws

.DEFAULT: build

.PHONY: build setup test
