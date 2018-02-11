DATE=$(shell date -u '+%Y-%m-%d %H:%M:%S')
COMMIT=$(shell git log --format=%h -1)
VERSION=main.version=${TRAVIS_BUILD_NUMBER} ${COMMIT} ${DATE} $@
COMPILE_FLAGS=-ldflags="-X '${VERSION}'"

build:
	@go build ${COMPILE_FLAGS}

amd64:
	@GOOS=linux GOARCH=amd64 go build ${COMPILE_FLAGS}

darwin:
	@GOOS=darwin go build ${COMPILE_FLAGS}

dep:
	@dep ensure

test:
	@go test .

clean:
	@rm -f bwpool_exporter

all: clean test build
