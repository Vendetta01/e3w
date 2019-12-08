
.PHONY: all build dep dep-build

all: build

dep-build: dep build

build:
	export GOPATH="~/go/"
	CGO_ENABLED=0 go build -o e3w

dep:
	export GOPATH="~/go/"
	rm -rf vendor/*
	dep ensure


