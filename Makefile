.PHONY: build install test test-verbose test-coverage clean

build:
	go build -o note main.go

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
