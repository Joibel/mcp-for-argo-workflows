# Project variables
BINARY_NAME := mcp-for-argo-workflows
BINARY_PATH := bin/$(BINARY_NAME)
MODULE := github.com/Joibel/mcp-for-argo-workflows
DIST_DIR := dist

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

# Platforms for cross-compilation
PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

.PHONY: all build test lint lint-fix fmt vet clean tools help \
	build-all build-darwin-amd64 build-darwin-arm64 \
	build-linux-amd64 build-linux-arm64 build-windows-amd64 \
	checksums dist-clean

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
	$(GO) test -race -coverprofile=coverage.out -covermode=atomic ./internal/...
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

# =============================================================================
# Cross-compilation targets
# =============================================================================

## build-all: Build binaries for all platforms
build-all: dist-clean build-darwin-amd64 build-darwin-arm64 build-linux-amd64 build-linux-arm64 build-windows-amd64 checksums
	@echo "All platform builds complete. Binaries in $(DIST_DIR)/"

## build-darwin-amd64: Build for macOS Intel
build-darwin-amd64:
	@echo "Building for darwin/amd64..."
	@mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build -ldflags="$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/$(BINARY_NAME)

## build-darwin-arm64: Build for macOS Apple Silicon
build-darwin-arm64:
	@echo "Building for darwin/arm64..."
	@mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO) build -ldflags="$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/$(BINARY_NAME)

## build-linux-amd64: Build for Linux x86_64
build-linux-amd64:
	@echo "Building for linux/amd64..."
	@mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -ldflags="$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/$(BINARY_NAME)

## build-linux-arm64: Build for Linux ARM64
build-linux-arm64:
	@echo "Building for linux/arm64..."
	@mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build -ldflags="$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/$(BINARY_NAME)

## build-windows-amd64: Build for Windows x86_64
build-windows-amd64:
	@echo "Building for windows/amd64..."
	@mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -ldflags="$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/$(BINARY_NAME)

## checksums: Generate SHA256 checksums for all binaries
checksums:
	@echo "Generating checksums..."
	@cd $(DIST_DIR) && sha256sum $(BINARY_NAME)-* > checksums.txt
	@echo "Checksums written to $(DIST_DIR)/checksums.txt"

## dist-clean: Remove distribution artifacts
dist-clean:
	@echo "Cleaning dist directory..."
	@rm -rf $(DIST_DIR)/
