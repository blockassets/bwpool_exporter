BINARY=bwpool_exporter
BINARY_LINUX=$(BINARY)-linux-amd64
BINARY_DARWIN=$(BINARY)-darwin-amd64

.DEFAULT_GOAL:=all

DATE=$(shell date -u '+%Y-%m-%d %H:%M:%S')
COMMIT=$(shell git log --format=%h -1)

build: VERSION=main.version=$(COMMIT) $(DATE)
build: COMPILE_FLAGS=-o $(BINARY) -ldflags="-X '$(VERSION)'"
build:
	go build $(COMPILE_FLAGS)

darwin: GOOS=darwin
darwin: GOARCH=amd64
darwin: VERSION=main.version=$(TRAVIS_BUILD_NUMBER) $(COMMIT) $(DATE) $(GOOS) $(GOARCH)
darwin: COMPILE_FLAGS=-o $(BINARY_DARWIN) -ldflags="-s -w -X '$(VERSION)'" # -s -w makes binary size smaller
darwin:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(COMPILE_FLAGS)

linux: GOOS=linux
linux: GOARCH=amd64
linux: VERSION=main.version=$(TRAVIS_BUILD_NUMBER) $(COMMIT) $(DATE) $(GOOS) $(GOARCH)
linux: COMPILE_FLAGS=-o $(BINARY_LINUX) -ldflags="-s -w -X '$(VERSION)'" # -s -w makes binary size smaller
linux:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(COMPILE_FLAGS)

gzip: darwin linux
	gzip -9 $(BINARY_LINUX)
	gzip -9 $(BINARY_DARWIN)

test:
	@go test ./...

dep:
	@dep ensure

clean:
	@rm -f $(BINARY) $(BINARY_LINUX) $(BINARY_DARWIN) $(BINARY_LINUX).gz $(BINARY_DARWIN).gz

all: clean test build
