.PHONY: all test lint coverage build clean

# Variables
VERSION := $(shell git describe --tags --always --dirty)
BUILD_FLAGS := -ldflags="-X 'github.com/yourusername/module.Version=$(VERSION)'"

all: lint test build

build:
    go build $(BUILD_FLAGS) ./...

test:
    go test -v -race ./...
    go test -v ./test/integration/...

test-short:
    go test -v -short ./...

lint:
    golangci-lint run

coverage:
    go test -coverprofile=coverage.txt -covermode=atomic ./...
    go tool cover -html=coverage.txt

clean:
    go clean
    rm -f coverage.txt

# Run examples
example-basic:
    go run ./examples/basic

example-advanced:
    go run ./examples/advanced