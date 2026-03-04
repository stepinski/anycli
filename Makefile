# anycli Makefile
#
# WHY A MAKEFILE?
# Go has `go build`, `go test` etc. but a Makefile gives you:
#   1. Muscle memory — `make test` works in every Go project you'll ever touch
#   2. Documented shortcuts — new contributors immediately see what commands exist
#   3. Dependency chains — `make release` can depend on `make test` passing first
#   4. CI parity — CI runs the same `make` commands as you do locally
#
# PHONY targets are targets that don't produce a file with that name.
# Without .PHONY, if a file called "test" existed, `make test` would do nothing.
.PHONY: build test lint clean install dev help

# Variables
BINARY     = anycli
VERSION   ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT    ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE      ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS    = -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Default target — what runs when you just type `make`
# Convention: make the default target print help
.DEFAULT_GOAL := help

## build: compile the binary for the current platform
build:
	go build $(LDFLAGS) -o $(BINARY) .

## install: install to $GOPATH/bin (makes `anycli` available system-wide)
install:
	go install $(LDFLAGS) .

## test: run all tests with race detector
# -race detects data races in concurrent code — always use it
# -count=1 disables test caching — ensures tests actually run
# ./... means all packages in the module
test:
	go test -race -count=1 -timeout 30s ./...

## test-verbose: run tests with full output
test-verbose:
	go test -race -count=1 -v -timeout 30s ./...

## test-cover: run tests and open coverage report in browser
test-cover:
	go test -race -count=1 -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html   # macOS; use `xdg-open` on Linux

## lint: run golangci-lint (install: brew install golangci-lint)
lint:
	golangci-lint run ./...

## clean: remove build artifacts
clean:
	rm -f $(BINARY) coverage.out coverage.html
	go clean -testcache

## dev: build and run with example args (useful during development)
dev: build
	./$(BINARY) --help

## tidy: tidy and verify go modules
tidy:
	go mod tidy
	go mod verify

## help: print this help message
# This works by reading lines starting with ##
# It's a self-documenting Makefile pattern — very popular in open source
help:
	@echo "anycli — terminal client for AnythingLLM"
	@echo ""
	@echo "Usage:"
	@awk 'BEGIN {FS = ":.*##"; printf ""} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
