.PHONY: all build test clean lint vet fmt

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOVET=$(GOCMD) vet
GOFMT=gofmt -w

# Binary names
BINARY_NAME=dsp-service

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

lint:
	golangci-lint run

vet:
	$(GOVET) ./...

fmt:
	$(GOFMT) .

# Run all code quality checks
quality: fmt vet lint

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Install dependencies
deps:
	$(GOGET) -v ./...

# Build and run
run: build
	./$(BINARY_NAME)

# Default target
default: quality test build