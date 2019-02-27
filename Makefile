.PHONY: all
all: build

export GO111MODULE=on

build:
	go mod download
	go build .