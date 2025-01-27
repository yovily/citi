.PHONY: test clean build coverage lint help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=auth-service
PKG_LIST=$(shell go list ./... | grep -v /vendor/)

# Main targets
all: test build

help: ## Display available commands
	@echo "Available commands:"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	$(GOBUILD) -o $(BINARY_NAME) -v

clean: ## Remove binary and test coverage files
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.out

test: ## Run tests
	$(GOTEST) $(PKG_LIST)

auth:
	go test -v -cover ./pkg/auth/...

ldap:
	go test -v -cover ./pkg/ldap/...

resolver:
	go test -v -cover ./pkg/resolver/...

handler:
	go test -v -cover ./internal/handler/...

platform:
	go test -v -cover ./internal/platform/...

test-verbose: ## Run tests with verbose output
	$(GOTEST) -v $(PKG_LIST)

coverage: ## Run tests with coverage
	$(GOTEST) -cover $(PKG_LIST)

coverage-html: ## Generate HTML coverage report
	$(GOTEST) -coverprofile=coverage.out ./pkg/auth/...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	open coverage.html  # For macOS
	# xdg-open coverage.html  # Uncomment for Linux
	# start coverage.html     # Uncomment for Windows

deps: ## Download dependencies
	$(GOMOD) download

tidy: ## Tidy up module dependencies
	$(GOMOD) tidy

verify: ## Verify dependencies
	$(GOMOD) verify

lint: ## Run linter
	golangci-lint run

# Run all quality checks
check: deps verify lint test coverage ## Run all quality checks