# Project variables
BINARY_NAME := mcp-for-argo-workflows
BINARY_PATH := bin/$(BINARY_NAME)
MODULE := github.com/Joibel/mcp-for-argo-workflows

# Build variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS := -s -w \
	-X $(MODULE)/internal/version.Version=$(VERSION) \
	-X $(MODULE)/internal/version.Commit=$(COMMIT) \
	-X $(MODULE)/internal/version.BuildTime=$(BUILD_TIME)

# Go variables
GO := go
GOFMT := gofmt
GOIMPORTS := goimports
GOLANGCI_LINT := golangci-lint

# Cross-compilation
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

.PHONY: all build test lint lint-fix fmt vet clean tools help

# Default target
all: fmt vet lint test build

## build: Compile the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -ldflags="$(LDFLAGS)" -o $(BINARY_PATH) ./cmd/$(BINARY_NAME)
	@echo "Built $(BINARY_PATH)"

## test: Run tests with race detection and coverage
test:
	@echo "Running tests..."
	$(GO) test -race -coverprofile=coverage.out -covermode=atomic ./...
	@echo "Coverage report: coverage.out"

## lint: Run golangci-lint
lint:
	@echo "Running linter..."
	$(GOLANGCI_LINT) run ./...

## lint-fix: Run golangci-lint with auto-fix
lint-fix:
	@echo "Running linter with auto-fix..."
	$(GOLANGCI_LINT) run --fix ./...

## fmt: Run gofmt and goimports
fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .
	@if command -v $(GOIMPORTS) >/dev/null 2>&1; then \
		$(GOIMPORTS) -w -local $(MODULE) .; \
	else \
		echo "goimports not installed, skipping import formatting"; \
	fi

## vet: Run go vet
vet:
	@echo "Running go vet..."
	$(GO) vet ./...

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

## tools: Install development dependencies
tools:
	@echo "Installing development tools..."
	$(GO) install golang.org/x/tools/cmd/goimports@latest
	@echo "Note: Install golangci-lint from https://golangci-lint.run/welcome/install/"

## help: Show this help message
help:
	@echo "Available targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /'
