.PHONY: build install test test-verbose test-coverage clean

# Get version from git tag or use "dev" for local builds
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

build:
	go build -ldflags "-X github.com/nathan-nicholson/note/internal/version.Version=$(VERSION)" -o note main.go

install:
	go install

test:
	go test ./...

test-verbose:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean:
	rm -f note coverage.out coverage.html
