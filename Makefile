.PHONY: build install test test-verbose test-coverage clean

# Get version from git tags, or use dev + short commit hash
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -X github.com/nathan-nicholson/note/internal/version.Version=$(VERSION)

build:
	go build -ldflags "$(LDFLAGS)" -o note main.go

install:
	go install -ldflags "$(LDFLAGS)"

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
