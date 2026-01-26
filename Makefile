SHELL := /bin/bash
APP := depo-cleaner
SCRIPT := scripts/build.sh
DIST := dist
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

.PHONY: help build clean tidy fmt test

default: help

help:
	@echo "Targets:"
	@echo "  build       Build all targets via scripts/build.sh"
	@echo "  clean       Remove $(DIST)/ artifacts"
	@echo "  tidy        Run go mod tidy"
	@echo "  fmt         Format Go code"
	@echo "  test        Run unit tests"

build: tidy clean 
	@echo "Building $(APP) (VERSION=$(VERSION))..."
	@VERSION=$(VERSION) $(SCRIPT)

clean:
	@echo "Cleaning previous builds..."
	@rm -rf $(DIST)

tidy:
	go mod tidy

fmt:
	go fmt ./...

test:
	go test ./...
