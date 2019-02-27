.PHONY: all
all: build

export GO111MODULE=on

clean:
	rm -rf 

build:
	go mod download
	make dist

dist: dist/darwin dist/linux dist/windows

dist/darwin:
	mkdir -p dist/darwin
	GOOS=darwin GOARCH=amd64 go build -o dist/darwin/github-desktop-notifications

dist/linux:
	mkdir -p dist/linux
	GOOS=linux GOARCH=amd64 go build -o dist/linux/github-desktop-notifications

dist/windows:
	mkdir -p dist/windows
	GOOS=windows GOARCH=amd64 go build -o dist/windows/github-desktop-notifications